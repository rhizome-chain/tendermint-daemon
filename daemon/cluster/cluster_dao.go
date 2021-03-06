package cluster

import (
	"github.com/tendermint/tendermint/libs/log"
	
	"github.com/rhizome-chain/tendermint-daemon/daemon/common"
	"github.com/rhizome-chain/tendermint-daemon/types"
)

const (
	PathMember    = "member"
	PathHeartbeat = "heartbeat"
	PathLeader    = "leader"
)

// clusterDao kv store model for job
type clusterDao struct {
	config common.DaemonConfig
	logger log.Logger
	client types.Client
}

var _ Repository = (*clusterDao)(nil)

func NewRepository(config common.DaemonConfig, logger log.Logger, client types.Client) Repository {
	dao := &clusterDao{config: config, logger: logger, client: client}
	return dao
}

func (dao *clusterDao) PutMember(member *Member) (err error) {
	bytes, err := types.BasicCdc.MarshalBinaryBare(member)
	
	if err != nil {
		dao.logger.Error("PutMember marshal", err)
		return err
	}
	
	dao.logger.Info("PutMember :", "member", member)
	
	msg := types.NewTxMsg(types.TxSetSync, common.SpaceDaemon, PathMember, member.NodeID, bytes)
	
	return dao.client.BroadcastTxSync(msg)
}

func (dao *clusterDao) GetMember(nodeID string) (member *Member, err error) {
	msg := types.NewViewMsgOne(common.SpaceDaemon, PathMember, nodeID)
	
	member = &Member{}
	err = dao.client.GetObject(msg, member)
	return member, err
}

func (dao *clusterDao) HasMember(nodeID string) (ok bool) {
	msg := types.NewViewMsgHas(common.SpaceDaemon, PathMember, nodeID)
	ok, err := dao.client.Has(msg)
	
	if err != nil {
		dao.logger.Error("HasMember ", err)
	}
	return ok
}

// PutLeader set leader
func (dao *clusterDao) PutLeader(leader string) (err error) {
	msg := types.NewTxMsg(types.TxSetSync, common.SpaceDaemon, PathLeader, "", []byte(leader))
	// fmt.Println(" -------- PutLeader :", leader)
	return dao.client.BroadcastTxSync(msg)
}

// GetLeader get leader id
func (dao *clusterDao) GetLeader() (leader string, err error) {
	msg := types.NewViewMsgOne(common.SpaceDaemon, PathLeader, "")
	data, err := dao.client.Query(msg)
	return string(data), err
}

func (dao *clusterDao) GetAllMembers() (members []*Member, err error) {
	msg := types.NewViewMsgMany(common.SpaceDaemon, PathMember, "", "")
	
	members = []*Member{}
	
	err = dao.client.GetMany(msg, func(key []byte, value []byte) bool {
		member := Member{}
		err = dao.client.UnmarshalObject(value, &member)
		if err != nil {
			dao.logger.Error("GetAllMembers unmarshal member : ", err)
		} else {
			members = append(members, &member)
		}
		
		return true
	})
	
	return members, err
}

func (dao *clusterDao) GetAllMemberIDs() (memberIDs []string, err error) {
	msg := types.NewViewMsgKeys(common.SpaceDaemon, PathMember, "", "")
	
	memberIDs, err = dao.client.GetKeys(msg)
	
	return memberIDs, err
}

func (dao *clusterDao) PutHeartbeat(nodeID string, blockHeight int64) (err error) {
	bytes, err := types.BasicCdc.MarshalBinaryBare(blockHeight)
	
	if err != nil {
		dao.logger.Error("PutHeartbeat : Member : ", err)
		return err
	}
	
	msg := types.NewTxMsg(types.TxSetSync, common.SpaceDaemon, PathHeartbeat, nodeID, bytes)
	
	return dao.client.BroadcastTxSync(msg)
}

func (dao *clusterDao) GetHeartbeats(handler func(nodeid string, blockHeight int64)) (err error) {
	msg := types.NewViewMsgMany(common.SpaceDaemon, PathHeartbeat, "", "")
	err = dao.client.GetMany(msg, func(key []byte, value []byte) bool {
		var heartbeat int64
		nodeid := string(key)
		err = dao.client.UnmarshalObject(value, &heartbeat)
		
		if err != nil {
			dao.logger.Error("GetHeartbeats unmarshal heartbeat", "key", string(key),
				"value", value, "err", err)
		} else {
			handler(nodeid, heartbeat)
		}
		
		return true
	})
	
	return err
}
