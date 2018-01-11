package dashboard

import (
	"fmt"
	"log"
	"net/http"
	"shkaff/config"
	"shkaff/drivers/maindb"
	"shkaff/drivers/stat"
	"strconv"

	"github.com/gin-gonic/gin"
)

// type PostTask struct {
// 	TaskName    string    `json:"task_name" db:"task_name"`
// 	Verb        int       `json:"verb" db:"verb"`
// 	ThreadCount int       `json:"thread_count" db:"thread_count"`
// 	Gzip        bool      `json:"gzip" db:"gzip"`
// 	Ipv6        bool      `json:"ipv6" db:"ipv6"`
// 	Host        string    `json:"host" db:"host"`
// 	Port        int       `json:"port" db:"port"`
// 	StartTime   time.Time `json:"start_time" db:"start_time"`
// 	DBUser      string    `json:"db_user" db:"db_user"`
// 	DBPassword  string    `json:"db_password" db:"db_password"`
// 	Database    string    `json:"database"`
// 	Sheet       string    `json:"sheet"`
// }

type API struct {
	cfg    *config.ShkaffConfig
	report *stat.StatDB
	router *gin.Engine
	psql   *maindb.PSQL
}

func InitAPI() (api *API) {
	api = &API{
		cfg:    config.InitControlConfig(),
		router: gin.Default(),
		report: stat.InitStat(),
		psql:   maindb.InitPSQL(),
	}
	v1 := api.router.Group("/api/v1")
	//CRUD Operation with Tasks
	{
		v1.POST("/CreateTask", api.createTask)
		v1.POST("/UpdateTask/:TaskID", api.updateTask)
		v1.GET("/GetTask/:TaskID", api.getTask)
		v1.DELETE("/DeleteTask/:TaskID", api.deleteTask)
	}
	{
		v1.POST("/CreateDatabase", api.createTask)
		v1.POST("/UpdateDatabase/:DatabaseID", api.updateDatabase)
		v1.GET("/GetDatabase/:DatabaseID", api.getDatabase)
		v1.DELETE("/DeleteDatabase/:DatabaseID", api.deleteDatabase)
	}
	//Statistic
	{
		v1.GET("/GetStat/:TaskID", api.getTaskStat)
	}

	return
}

func (api *API) Run() {
	log.Println("Start Dashboard")
	uri := fmt.Sprintf("%s:%d", api.cfg.SHKAFF_UI_HOST, api.cfg.SHKAFF_UI_PORT)
	err := api.router.Run(uri)
	if err != nil {
		log.Fatalln(err)
	}
	return
}

func (api *API) createTask(c *gin.Context) {
	setStrings, err := api.checkTaskParameters(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Error": err.Error()})
		return
	}
	_, err = api.psql.CreateTask(setStrings)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Error": err.Error()})
		return
	}
	taskName := setStrings["task_name"].(string)
	if taskName == "" {
		c.JSON(http.StatusNotFound, gin.H{"Error": err.Error()})
		return
	}
	task, err := api.psql.GetTaskByName(taskName)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusNotFound, gin.H{"Error": "TaskID not found"})
		return
	}
	c.JSON(http.StatusOK, task)
	return
}

func (api *API) updateTask(c *gin.Context) {
	taskID := c.Param("TaskID")
	taskIDInt, err := strconv.Atoi(taskID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Error": err.Error()})
		return
	}
	_, err = api.psql.GetTask(taskIDInt, true)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Error": "TaskID not found"})
		return
	}
	setStrings, err := api.checkTaskParameters(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Error": err.Error()})
		return
	}
	_, err = api.psql.UpdateTask(taskIDInt, setStrings)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Status": "OK"})
	return
}

func (api *API) getTask(c *gin.Context) {
	taskID := c.Param("TaskID")
	taskIDInt, err := strconv.Atoi(taskID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Error": "Bad taskID"})
		return
	}
	task, err := api.psql.GetTask(taskIDInt, false)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusNotFound, gin.H{"Error": "TaskID not found"})
		return
	}
	c.JSON(http.StatusOK, task)
	return
}

func (api *API) deleteTask(c *gin.Context) {
	taskID := c.Param("TaskID")
	taskIDInt, err := strconv.Atoi(taskID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Error": "Bad taskID"})
		return
	}
	_, err = api.psql.GetTask(taskIDInt, true)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Error": "TaskID not found"})
		return
	}
	_, err = api.psql.DeleteTask(taskIDInt)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Error": "TaskID not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Result": "Success"})
	return
}

func (api *API) getTaskStat(c *gin.Context) {
	taskID := c.Param("TaskID")
	_, err := strconv.Atoi(taskID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Error": "Bad taskID"})
		return
	}
	task, err := api.report.StandartStatSelect()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, task)
	return

}

func (api *API) createDatabase(c *gin.Context) {
	setStrings, err := api.checkDatabaseParameters(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Error": err.Error()})
		return
	}
	_, err = api.psql.CreateDatabase(setStrings)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, "OK")
	return
}

func (api *API) getDatabase(c *gin.Context) {
	DatabaseID := c.Param("DatabaseID")
	DatabaseIDInt, err := strconv.Atoi(DatabaseID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Error": "Bad DatabaseID"})
		return
	}
	database, err := api.psql.GetDatabase(DatabaseIDInt)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusNotFound, gin.H{"Error": "DatabaseID not found"})
		return
	}
	c.JSON(http.StatusOK, database)
	return
}
func (api *API) updateDatabase(c *gin.Context) {
	databaseID := c.Param("DatabaseID")
	databaseIDInt, err := strconv.Atoi(databaseID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Error": err.Error()})
		return
	}
	_, err = api.psql.GetDatabase(databaseIDInt)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Error": "DatabaseID not found"})
		return
	}
	setStrings, err := api.checkDatabaseParameters(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Error": err.Error()})
		return
	}
	_, err = api.psql.UpdateDatabase(databaseIDInt, setStrings)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Status": "OK"})
	return
}

func (api *API) deleteDatabase(c *gin.Context) {
	databaseID := c.Param("DatabaseID")
	databaseIDInt, err := strconv.Atoi(databaseID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Error": "Bad DatabaseID"})
		return
	}
	_, err = api.psql.GetDatabase(databaseIDInt)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusNotFound, gin.H{"Error": "DatabaseID not found"})
		return
	}
	_, err = api.psql.DeleteDatabase(databaseIDInt)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusNotFound, gin.H{"Error": "DatabaseID not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Result": "Success"})
	return
}
