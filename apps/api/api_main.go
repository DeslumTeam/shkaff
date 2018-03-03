package api

import (
	"fmt"
	"shkaff/drivers/maindb"
	"shkaff/drivers/stat"
	"shkaff/internal/logger"
	"shkaff/internal/options"

	"github.com/gin-gonic/gin"
	logging "github.com/op/go-logging"
)

type API struct {
	cfg    *options.ShkaffConfig
	report *stat.StatDB
	router *gin.Engine
	psql   *maindb.PSQL
	log    *logging.Logger
}

func InitAPI() (api *API) {
	gin.SetMode(gin.ReleaseMode)
	api = &API{
		cfg:    options.InitControlConfig(),
		router: gin.Default(),
		report: stat.InitStat(),
		psql:   maindb.InitPSQL(),
		log:    logger.GetLogs("API"),
	}
	v1 := api.router.Group("/api/v1")
	//CRUD Operations with Users
	{
		v1.POST("/CreateUser", api.createUser)
		v1.POST("/UpdateUser/:UserID", api.updateUser)
		v1.GET("/GetUser/:UserID", api.getUser)
		v1.DELETE("/DeleteUser/:UserID", api.deleteUser)
	}
	//CRUD Operations with DatabaseSettings
	{
		v1.POST("/CreateDatabase", api.createDatabase)
		v1.POST("/UpdateDatabase/:DatabaseID", api.updateDatabase)
		v1.GET("/GetDatabase/:DatabaseID", api.getDatabase)
		v1.DELETE("/DeleteDatabase/:DatabaseID", api.deleteDatabase)
	}
	//CRUD Operations with Tasks
	{
		v1.POST("/CreateTask", api.createTask)
		v1.POST("/UpdateTask/:TaskID", api.updateTask)
		v1.GET("/GetTask/:TaskID", api.getTask)
		v1.DELETE("/DeleteTask/:TaskID", api.deleteTask)
	}
	//Statistic
	{
		v1.GET("/GetStat/:TaskID", api.getTaskStat)
	}
	return
}

func (api *API) Run() {
	api.log.Info("Start API")
	uri := fmt.Sprintf("%s:%d", api.cfg.SHKAFF_UI_HOST, api.cfg.SHKAFF_UI_PORT)
	err := api.router.Run(uri)
	if err != nil {
		api.log.Fatal(err)
	}
	return
}

func (api *API) Stop() {
	api.log.Info("Stop Dashboard")
}
