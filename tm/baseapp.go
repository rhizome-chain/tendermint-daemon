package tm

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	
	"github.com/tendermint/tendermint/abci/example/code"
	abcitypes "github.com/tendermint/tendermint/abci/types"
	cfg "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/libs/log"
	tmtypes "github.com/tendermint/tendermint/types"
	"github.com/tendermint/tendermint/version"
	dbm "github.com/tendermint/tm-db"
	
	"github.com/rhizome-chain/tendermint-daemon/tm/events"
	"github.com/rhizome-chain/tendermint-daemon/tm/store"
	"github.com/rhizome-chain/tendermint-daemon/tm/tmcom"
	"github.com/rhizome-chain/tendermint-daemon/types"
)

var (
	ProtocolVersion version.Protocol = 0x1
)

const (
	ValidatorSetChangePrefix string = "val:"
)

var (
	stateKey = []byte("stateKey")
)

type State struct {
	db      dbm.DB
	Size    int64  `json:"size"`
	Height  int64  `json:"height"`
	AppHash []byte `json:"app_hash"`
}

func loadState(db dbm.DB) State {
	var state State
	state.db = db
	stateBytes, err := db.Get(stateKey)
	if err != nil {
		panic(err)
	}
	if len(stateBytes) == 0 {
		return state
	}
	err = json.Unmarshal(stateBytes, &state)
	if err != nil {
		panic(err)
	}
	return state
}

func saveState(state State) {
	stateBytes, err := json.Marshal(state)
	if err != nil {
		panic(err)
	}
	state.db.Set(stateKey, stateBytes)
}

type BaseApplication struct {
	abcitypes.BaseApplication
	config *cfg.Config
	logger log.Logger
	state  State
	// validator set
	ValUpdates         []abcitypes.ValidatorUpdate
	valAddrToPubKeyMap map[string]abcitypes.PubKey
	spaces             map[string]*store.Registry
}

var _ abcitypes.Application = (*BaseApplication)(nil)
var _ types.SpaceRegistry = (*BaseApplication)(nil)

func NewBaseApplication(config *cfg.Config, logger log.Logger) (bapp *BaseApplication) {
	bapp = &BaseApplication{
		config:             config,
		logger:             logger,
		ValUpdates:         []abcitypes.ValidatorUpdate{},
		valAddrToPubKeyMap: make(map[string]abcitypes.PubKey),
		spaces:             make(map[string]*store.Registry),
	}
	
	registry := bapp.registerSpace(tmcom.SpaceDaemonState)
	
	bapp.state = loadState(registry.DB())
	
	return bapp
}

func (app *BaseApplication) RegisterSpace(name string) {
	app.registerSpace(name)
}

func (app *BaseApplication) RegisterSpaceIfNotExist(name string) {
	if _, ok := app.spaces[name]; !ok {
		app.registerSpace(name)
	}
}

func (app *BaseApplication) registerSpace(name string) *store.Registry {
	if reg, ok := app.spaces[name]; ok {
		app.logger.Error("[ERROR] Register Space '" + name + "' already exists.")
		return reg
	}
	
	db, err := dbm.NewGoLevelDB(name, app.config.DBDir())
	if err != nil {
		panic("[ERROR] Register Space '" + name + "' : " + err.Error())
	}
	
	storeRegistry := store.NewRegistry(db)
	app.spaces[name] = storeRegistry
	return storeRegistry
}

func (app *BaseApplication) getSpace(name string) *store.Registry {
	storeRegistry, ok := app.spaces[name]
	if !ok {
		panic(fmt.Sprintf("DB Space[%s] is not registered.", name))
	}
	return storeRegistry
}

func (app *BaseApplication) getSpaceStoreAny(space string, path string) *store.Store {
	storeRegistry, ok := app.spaces[space]
	if !ok {
		storeRegistry = app.registerSpace(space)
		app.logger.Info(fmt.Sprintf("[WARN] Force to make Space[%s]",space))
	}
	
	return storeRegistry.GetOrMakeStore(path)
}

func (app *BaseApplication) IncreaseTxSize() {
	app.state.Size++
}

func (app *BaseApplication) Info(req abcitypes.RequestInfo) abcitypes.ResponseInfo {
	res := abcitypes.ResponseInfo{
		Data:             fmt.Sprintf("{\"size\":%v}", app.state.Size),
		Version:          version.ABCIVersion,
		AppVersion:       ProtocolVersion.Uint64(),
		LastBlockHeight:  app.state.Height,
		LastBlockAppHash: app.state.AppHash,
	}
	return res
}

func (app *BaseApplication) Commit() abcitypes.ResponseCommit {
	// Using a memdb - just return the big endian size of the db
	appHash := make([]byte, 8)
	binary.PutVarint(appHash, app.state.Size)
	app.state.AppHash = appHash
	app.state.Height++
	saveState(app.state)
	events.PublishBlockEvent(events.CommitEvent{Height: app.state.Height, Size: app.state.Size, AppHash: app.state.AppHash})
	return abcitypes.ResponseCommit{Data: appHash}
}

func (app *BaseApplication) InitChain(req abcitypes.RequestInitChain) abcitypes.ResponseInitChain {
	for _, v := range req.Validators {
		r := app.updateValidator(v)
		if r.IsErr() {
			app.logger.Error("Error updating validators", "r", r)
		}
	}
	return abcitypes.ResponseInitChain{}
}

func (app *BaseApplication) BeginBlock(req abcitypes.RequestBeginBlock) abcitypes.ResponseBeginBlock {
	// reset valset changes
	app.ValUpdates = make([]abcitypes.ValidatorUpdate, 0)
	
	for _, ev := range req.ByzantineValidators {
		if ev.Type == tmtypes.ABCIEvidenceTypeDuplicateVote {
			// decrease voting power by 1
			if ev.TotalVotingPower == 0 {
				continue
			}
			app.updateValidator(abcitypes.ValidatorUpdate{
				PubKey: app.valAddrToPubKeyMap[string(ev.Validator.Address)],
				Power:  ev.TotalVotingPower - 1,
			})
		}
	}
	
	events.PublishBlockEvent(events.BeginBlockEvent{Height: req.Header.Height, Time: req.Header.Time})
	return abcitypes.ResponseBeginBlock{}
}

func (app *BaseApplication) EndBlock(req abcitypes.RequestEndBlock) abcitypes.ResponseEndBlock {
	events.PublishBlockEvent(events.EndBlockEvent{Height: req.GetHeight(), Size: req.Size()})
	return abcitypes.ResponseEndBlock{ValidatorUpdates: app.ValUpdates}
}

func (app *BaseApplication) isValidatorTx(tx []byte) bool {
	return strings.HasPrefix(string(tx), ValidatorSetChangePrefix)
}

// format is "val:pubkey!power"
// pubkey is a base64-encoded 32-byte ed25519 key
func (app *BaseApplication) execValidatorTx(tx []byte) abcitypes.ResponseDeliverTx {
	tx = tx[len(ValidatorSetChangePrefix):]
	
	// get the pubkey and power
	pubKeyAndPower := strings.Split(string(tx), "!")
	if len(pubKeyAndPower) != 2 {
		return abcitypes.ResponseDeliverTx{
			Code: code.CodeTypeEncodingError,
			Log:  fmt.Sprintf("Expected 'pubkey!power'. Got %v", pubKeyAndPower)}
	}
	pubkeyS, powerS := pubKeyAndPower[0], pubKeyAndPower[1]
	
	// decode the pubkey
	pubkey, err := base64.StdEncoding.DecodeString(pubkeyS)
	if err != nil {
		return abcitypes.ResponseDeliverTx{
			Code: code.CodeTypeEncodingError,
			Log:  fmt.Sprintf("Pubkey (%s) is invalid base64", pubkeyS)}
	}
	
	// decode the power
	power, err := strconv.ParseInt(powerS, 10, 64)
	if err != nil {
		return abcitypes.ResponseDeliverTx{
			Code: code.CodeTypeEncodingError,
			Log:  fmt.Sprintf("Power (%s) is not an int", powerS)}
	}
	
	// update
	return app.updateValidator(abcitypes.Ed25519ValidatorUpdate(pubkey, power))
}

func (app *BaseApplication) updateValidator(v abcitypes.ValidatorUpdate) abcitypes.ResponseDeliverTx {
	key := []byte("val:" + string(v.PubKey.Data))
	
	pubkey := ed25519.PubKeyEd25519{}
	copy(pubkey[:], v.PubKey.Data)
	
	if v.Power == 0 {
		// remove validator
		hasKey, err := app.state.db.Has(key)
		if err != nil {
			panic(err)
		}
		if !hasKey {
			pubStr := base64.StdEncoding.EncodeToString(v.PubKey.Data)
			return abcitypes.ResponseDeliverTx{
				Code: code.CodeTypeUnauthorized,
				Log:  fmt.Sprintf("Cannot remove non-existent validator %s", pubStr)}
		}
		app.state.db.Delete(key)
		delete(app.valAddrToPubKeyMap, string(pubkey.Address()))
	} else {
		// add or update validator
		value := bytes.NewBuffer(make([]byte, 0))
		if err := abcitypes.WriteMessage(&v, value); err != nil {
			return abcitypes.ResponseDeliverTx{
				Code: code.CodeTypeEncodingError,
				Log:  fmt.Sprintf("Error encoding validator: %v", err)}
		}
		app.state.db.Set(key, value.Bytes())
		app.valAddrToPubKeyMap[string(pubkey.Address())] = v.PubKey
	}
	
	// we only update the changes array if we successfully updated the tree
	app.ValUpdates = append(app.ValUpdates, v)
	
	return abcitypes.ResponseDeliverTx{Code: code.CodeTypeOK}
}
