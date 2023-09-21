package api

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"xll-job/orm"
	"xll-job/orm/do"
	"xll-job/scheduler/core"
	"xll-job/scheduler/handle"
	"xll-job/utils"
	"xll-job/web/dto"
)

// JobInfoApi /*todo 除了添加的方法其他方法都差个更新缓存以及替换cron操作*/
type JobInfoApi struct {
	router *gin.RouterGroup
}

func NewJobInfoApi(router *gin.RouterGroup) *JobInfoApi {
	return &JobInfoApi{
		router: router,
	}
}

func (jobInfoApi *JobInfoApi) Router() {
	jobInfoApi.router.POST("/add", jobInfoApi.Add)
	jobInfoApi.router.POST("/update", jobInfoApi.update)
	jobInfoApi.router.GET("/stop", jobInfoApi.stop)
	jobInfoApi.router.GET("/start", jobInfoApi.start)
	jobInfoApi.router.GET("/delete", jobInfoApi.delete)
}

func (jobInfoApi *JobInfoApi) update(ctx *gin.Context) {
	var jobDto dto.JobInfoDto
	ctx.BindJSON(&jobDto)
	if jobDto.Id == 0 {
		ctx.JSON(500, gin.Error{
			Meta: "修改任务需要绑定id",
		})
		ctx.Done()
		return
	}
	//七位cron转六位
	jobDto.Cron = utils.RemoveYearField(jobDto.Cron)
	var jobInfoDo do.JobInfoDo
	copier.Copy(&jobInfoDo, &jobDto)
	orm.DB.Updates(&jobInfoDo)

	ctx.JSON(200, map[string]string{
		"msg": "ok",
	})
}

func (jobInfoApi *JobInfoApi) Add(ctx *gin.Context) {
	var jobDto dto.JobInfoDto
	ctx.BindJSON(&jobDto)

	//七位cron转六位
	jobDto.Cron = utils.RemoveYearField(jobDto.Cron)

	var jobInfoDo do.JobInfoDo
	copier.Copy(&jobInfoDo, &jobDto)
	jobInfoDo.Enable = false
	orm.DB.Create(&jobInfoDo)
	//添加缓存
	manager := handle.JobManagerMap[jobDto.ManageId]
	scheduler, _ := core.NewScheduler(jobInfoDo.Cron, jobInfoDo.JobHandler, manager, false)
	manager.Schedulers[jobInfoDo.Id] = scheduler
	ctx.JSON(200, map[string]interface{}{
		"message": "ok",
	})
}

func (jobInfoApi *JobInfoApi) stop(context *gin.Context) {

}

func (jobInfoApi *JobInfoApi) start(context *gin.Context) {

}

func (jobInfoApi *JobInfoApi) delete(context *gin.Context) {

}
