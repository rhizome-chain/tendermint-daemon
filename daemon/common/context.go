package common

import (
	client "github.com/rhizome-chain/tendermint-daemon/tm/client"
	"github.com/rhizome-chain/tendermint-daemon/types"
	cfg "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/node"
	"time"
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
}

type DefaultContext struct {
	tmNode       *node.Node
	tmCfg        *cfg.Config
	config       DaemonConfig
	logger       log.Logger
	client       types.Client
	clusterState ClusterState
}

func NewContext(tmCfg *cfg.Config, logger log.Logger, tmNode *node.Node, config DaemonConfig) Context {
	client := client.NewClient(tmCfg, logger, types.BasicCdc)
	return &DefaultContext{
		tmNode: tmNode,
		tmCfg:  tmCfg,
		config: config,
		logger: logger,
		client: client,
	}
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
