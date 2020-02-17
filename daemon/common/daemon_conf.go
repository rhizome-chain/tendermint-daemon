package common

const (
	SpaceDaemon = "daemon"
)

type DaemonConfig struct {
	ChainID  string
	NodeID   string
	NodeName string
	// HeartbeatInterval Heartbeat Interval
	//HeartbeatInterval uint
	
	// CheckHeartbeatInterval Heartbeat check Interval
	// CheckHeartbeatInterval uint
	
	// AliveThreasholdSecond Heartbeat time Threshold
	AliveThresholdSeconds uint
}


func NewDaemonConfig() *DaemonConfig {
	conf := &DaemonConfig{
		AliveThresholdSeconds: uint(2),
	}
	return conf
}