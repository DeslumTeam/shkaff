package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (api *API) createTask(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusNotFound, gin.H{"Error": "Token is empty"})
		return
	}
	isExist, err := api.psql.GetUserByToken(token)
	if err != nil || !isExist {
		c.JSON(http.StatusNotFound, gin.H{"Error": "User with this token not found"})
		return
	}
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
		api.log.Error(err)
		c.JSON(http.StatusNotFound, gin.H{"Error": "TaskID not found"})
		return
	}
	c.JSON(http.StatusOK, task)
	return
}

func (api *API) updateTask(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusNotFound, gin.H{"Error": "Token is empty"})
		return
	}
	isExist, err := api.psql.GetUserByToken(token)
	if err != nil || !isExist {
		c.JSON(http.StatusNotFound, gin.H{"Error": "User with this token not found"})
		return
	}
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
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusNotFound, gin.H{"Error": "Token is empty"})
		return
	}
	isExist, err := api.psql.GetUserByToken(token)
	if err != nil || !isExist {
		c.JSON(http.StatusNotFound, gin.H{"Error": "User with this token not found"})
		return
	}
	taskID := c.Param("TaskID")
	taskIDInt, err := strconv.Atoi(taskID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Error": "Bad taskID"})
		return
	}
	task, err := api.psql.GetTask(taskIDInt, false)
	if err != nil {
		api.log.Error(err)
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

func (api *API) getAllStat(c *gin.Context) {
	// token := c.Query("token")
	// if token == "" {
	// 	c.JSON(http.StatusNotFound, gin.H{"Error": "Token is empty"})
	// 	return
	// }
	// isExist, err := api.psql.GetUserByToken(token)
	// if err != nil || !isExist {
	// 	c.JSON(http.StatusNotFound, gin.H{"Error": "User with this token not found"})
	// 	return
	// }
	// taskID := c.Param("TaskID")
	// _, err = strconv.Atoi(taskID)
	// if err != nil {
	// 	c.JSON(http.StatusNotFound, gin.H{"Error": "Bad taskID"})
	// 	return
	// }
	taskStat, err := api.report.StandartStatSelect()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, taskStat)
	return

}

func (api *API) getTasksStatus(c *gin.Context) {
	tasksCount, err := api.psql.GetTasksCount()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tasksCount)
}

func (api *API) getErrors(c *gin.Context) {
	result, err := api.report.SelectDailyErrors()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)

}
