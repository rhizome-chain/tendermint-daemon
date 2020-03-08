package common

import (
	"time"
	
	cfg "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/node"
	
	"github.com/rhizome-chain/tendermint-daemon/tm/client"
	"github.com/rhizome-chain/tendermint-daemon/types"
)

type ClusterState interface {
	IsLeader() bool
	GetAliveMemberIDs() []string
}

type Context interface {
	log.Logger
	LastBlockTime() time.Time
	LastBlockIndex() int64
	GetNodeID() string
	GetClient() types.Client
	GetTMConfig() *cfg.Config
	GetConfig() DaemonConfig
	GetClusterState() ClusterState
	SetClusterState(state ClusterState)
	GetPrivateValidatorPubKey() crypto.PubKey
	CheckValidator() (ok bool)
	IsValidator() (ok bool)
}

type DefaultContext struct {
	tmNode          *node.Node
	tmCfg           *cfg.Config
	config          DaemonConfig
	logger          log.Logger
	client          types.Client
	clusterState    ClusterState
	validatorPubKey crypto.PubKey
	isValidator     bool
}

func NewContext(tmCfg *cfg.Config, logger log.Logger, tmNode *node.Node, config DaemonConfig) Context {
	client := client.NewClient(tmCfg, logger, types.BasicCdc)
	ctx := &DefaultContext{
		tmNode:          tmNode,
		tmCfg:           tmCfg,
		config:          config,
		logger:          logger,
		client:          client,
		validatorPubKey: tmNode.PrivValidator().GetPubKey(),
	}
	ctx.CheckValidator()
	return ctx
}

var _ Context = (*DefaultContext)(nil)

func (ctx *DefaultContext) GetClusterState() ClusterState {
	return ctx.clusterState
}

func (ctx *DefaultContext) SetClusterState(state ClusterState) {
	ctx.clusterState = state
}

func (ctx *DefaultContext) LastBlockTime() time.Time {
	return ctx.tmNode.ConsensusState().GetState().LastBlockTime
}

func (ctx *DefaultContext) LastBlockIndex() int64 {
	return ctx.tmNode.ConsensusState().GetState().LastBlockHeight
}
func (ctx *DefaultContext) GetNodeID() string {
	return ctx.config.NodeID
}
func (ctx *DefaultContext) GetClient() types.Client {
	return ctx.client
}

func (ctx *DefaultContext) GetTMConfig() *cfg.Config {
	return ctx.tmCfg
}

func (ctx *DefaultContext) GetConfig() DaemonConfig {
	return ctx.config
}

func (ctx *DefaultContext) GetPrivateValidatorPubKey() crypto.PubKey {
	return ctx.validatorPubKey
}

func (ctx *DefaultContext) CheckValidator() (ok bool) {
	valKey := ctx.GetPrivateValidatorPubKey()
	validators := ctx.GetClient().GetValidators()
	for _, val := range validators {
		if val.PubKey == valKey {
			ok = true
			break
		}
	}
	ctx.isValidator = ok
	return ok
}

func (ctx *DefaultContext) IsValidator() (ok bool) {
	return ctx.isValidator
}

// Info log info
func (ctx *DefaultContext) Info(msg string, keyvals ...interface{}) {
	ctx.logger.Info(msg, keyvals...)
}

// Info log info
func (ctx *DefaultContext) Debug(msg string, keyvals ...interface{}) {
	ctx.logger.Debug(msg, keyvals...)
}

// Info log info
func (ctx *DefaultContext) Error(msg string, keyvals ...interface{}) {
	ctx.logger.Error(msg, keyvals...)
}

// Info log info
func (ctx *DefaultContext) With(keyvals ...interface{}) log.Logger {
	return ctx.logger.With(keyvals...)
}
