package api

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
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
		log:    logger.GetLogs("Dashboard"),
	}
	currDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	api.router.LoadHTMLGlob(path.Join(currDir, "static", "html", "*"))
	api.router.Static("css", path.Join(currDir, "static", "css"))
	api.router.Static("js", path.Join(currDir, "static", "js"))
	api.router.Static("plugins", path.Join(currDir, "static", "plugins"))
	api.router.Static("img", path.Join(currDir, "static", "img"))

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
	ui := api.router.Group("ui")
	{
		ui.GET("/get_users", api.getUsers)
		ui.GET("/get_tasks", api.getTasks)
		ui.GET("/get_databases", api.getDatabases)
		ui.GET("/activate_task", api.changeTaskStatus)
		ui.GET("/deactivate_task", api.changeTaskStatus)
	}
	page := api.router.Group("/")
	{
		page.GET("/", api.General)
		page.GET("/dashboard", api.Dashboard)
		page.GET("/databases", api.Databases)
		page.GET("/tasks/active", api.ActiveTasks)
		page.GET("/tasks/unactive", api.UnactiveTasks)
	}
	return
}

func (api *API) Run() {
	api.log.Info("Start Dashboard")
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
