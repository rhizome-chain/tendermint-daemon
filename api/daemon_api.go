package api

import (
	"encoding/json"
	
	"github.com/rhizome-chain/tendermint-daemon/daemon"
	"github.com/rhizome-chain/tendermint-daemon/daemon/job"
	
	"net/http"
	
	"github.com/gin-gonic/gin"
)

// DaemonAPI ..
type DaemonAPI struct {
	daemon *daemon.Daemon
}

func NewDaemonAPI(daemon *daemon.Daemon) (api *DaemonAPI) {
	api = &DaemonAPI{
		daemon: daemon,
	}
	return api
}

func (api *DaemonAPI) RelativePath() string {
	return "daemon"
}

func (api *DaemonAPI) SetHandlers(group *gin.RouterGroup) {
	group.GET("info/cluster", api.getInfoCluster)
	group.GET("info/node", api.getInfoNode)
	group.GET("info/config", api.getInfoConfig)
	
	group.GET("jobs", api.getJobs)
	group.POST("job/add/factory/:factory/jobid/:jobid", api.addJob)
	group.DELETE("job/:jobid", api.removeJob)
}


func (api *DaemonAPI) getInfoConfig(context *gin.Context) {
	config := api.daemon.GetDaemonConfig()
	
	bytes, err := json.Marshal(config)
	
	if err != nil {
		context.Status(http.StatusBadRequest)
		context.Writer.WriteString(err.Error())
		context.Writer.Flush()
		return
	}
	
	context.Writer.Write(bytes)
	context.Writer.Flush()
}

func (api *DaemonAPI) getInfoNode(context *gin.Context) {
	config := api.daemon.GetTMConfig()
	
	bytes, err := json.Marshal(config)
	
	if err != nil {
		context.Status(http.StatusBadRequest)
		context.Writer.WriteString(err.Error())
		context.Writer.Flush()
		return
	}
	
	context.Writer.Write(bytes)
	context.Writer.Flush()
}


func (api *DaemonAPI) getInfoCluster(context *gin.Context) {
	cluster := api.daemon.GetCluster()
	
	bytes, err := json.Marshal(cluster)
	
	if err != nil {
		context.Status(http.StatusBadRequest)
		context.Writer.WriteString(err.Error())
		context.Writer.Flush()
		return
	}
	
	context.Writer.Write(bytes)
	context.Writer.Flush()
}


func (api *DaemonAPI) getJobs(context *gin.Context) {
	memberJobs, err := api.daemon.GetJobRepository().GetAllMemberJobIDs()
	if err != nil {
		context.Status(http.StatusBadRequest)
		context.Writer.WriteString(err.Error())
		context.Writer.Flush()
		return
	}
	
	bytes, err := json.Marshal(memberJobs)
	
	if err != nil {
		context.Status(http.StatusBadRequest)
		context.Writer.WriteString(err.Error())
		context.Writer.Flush()
		return
	}
	
	context.Writer.Write(bytes)
	context.Writer.Flush()
}


func (api *DaemonAPI) addJob(context *gin.Context) {
	data, err := context.GetRawData()
	if err != nil {
		context.Status(http.StatusBadRequest)
		context.Writer.WriteString(err.Error())
		context.Writer.Flush()
		return
	}
	
	factory := context.Param("factory")
	jobID := context.Param("jobid")
	
	var j job.Job
	
	if len(jobID) > 0 {
		j = job.NewWithID(factory, jobID, data)
	} else {
		j = job.New(factory, data)
	}
	
	err = api.daemon.GetJobRepository().PutJob(j)
	
	if err != nil {
		context.Status(http.StatusInternalServerError)
		context.Writer.WriteString(err.Error())
		context.Writer.Flush()
		return
	}
	
	context.Writer.Write([]byte(j.ID))
	context.Writer.Flush()
}

func (api *DaemonAPI) removeJob(context *gin.Context) {
	jobID := context.Param("jobid")
	
	err := api.daemon.GetJobRepository().RemoveJob(string(jobID))
	if err != nil {
		context.Status(http.StatusInternalServerError)
		context.Writer.WriteString(err.Error())
		context.Writer.Flush()
		return
	}
	context.Writer.WriteString("ok")
	context.Writer.Flush()
}
