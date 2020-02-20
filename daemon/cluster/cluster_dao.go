package cluster

import (
	"time"
	
	"github.com/tendermint/tendermint/libs/log"
	
	"github.com/rhizome-chain/tendermint-daemon/daemon/common"
	"github.com/rhizome-chain/tendermint-daemon/types"
)

const (
	PathMember    = "member"
	PathHeartbeat = "heartbeat"
	PathLeader = "leader"
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
	
	msg := types.NewTxMsg(types.TxSet, common.SpaceDaemon, PathMember, member.NodeID, bytes)
	
	return dao.client.BroadcastTxSync(msg)
}

func (dao *clusterDao) GetMember(nodeID string) (member Member, err error) {
	msg := types.NewViewMsgOne(common.SpaceDaemon, PathMember, nodeID)
	
	member = Member{}
	err = dao.client.GetObject(msg, &member)
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
	msg := types.NewTxMsg(types.TxSet, common.SpaceDaemon, PathLeader, "", []byte(leader))
	//fmt.Println(" -------- PutLeader :", leader)
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


func (dao *clusterDao) PutHeartbeat(nodeID string) (err error) {
	bytes, err := types.BasicCdc.MarshalBinaryBare(time.Now())
	
	if err != nil {
		dao.logger.Error("PutHeartbeat : Member : ", err)
		return err
	}
	
	msg := types.NewTxMsg(types.TxSet, common.SpaceDaemon, PathHeartbeat, nodeID, bytes)
	
	return dao.client.BroadcastTxAsync(msg)
}

func (dao *clusterDao) GetHeartbeats(handler func(nodeid string, time time.Time)) (err error) {
	msg := types.NewViewMsgMany(common.SpaceDaemon, PathHeartbeat, "", "")
	err = dao.client.GetMany(msg, func(key []byte, value []byte) bool {
		time := time.Time{}
		nodeid := string(key)
		err = dao.client.UnmarshalObject(value, &time)
		
		if err != nil {
			dao.logger.Error("GetHeartbeats unmarshal time", "key", string(key),
				"value", string(value), err)
		} else {
			handler(nodeid, time)
		}
		
		return true
	})
	
	return err
}
