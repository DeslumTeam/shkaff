package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (api *API) createDatabase(c *gin.Context) {
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
	setStrings, err := api.checkDatabaseParameters(c)
	if err != nil {
		api.log.Error(err)
		c.JSON(http.StatusNotFound, gin.H{"Error": err.Error()})
		return
	}
	_, err = api.psql.CreateDatabase(setStrings)
	if err != nil {
		api.log.Error(err)
		c.JSON(http.StatusNotFound, gin.H{"Error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, "OK")
	return
}

func (api *API) getDatabase(c *gin.Context) {
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
	DatabaseID := c.Param("DatabaseID")
	DatabaseIDInt, err := strconv.Atoi(DatabaseID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Error": "Bad DatabaseID"})
		return
	}
	database, err := api.psql.GetDatabase(DatabaseIDInt)
	if err != nil {
		api.log.Error(err)
		c.JSON(http.StatusNotFound, gin.H{"Error": "DatabaseID not found"})
		return
	}
	c.JSON(http.StatusOK, database)
	return
}

func (api *API) updateDatabase(c *gin.Context) {
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
	databaseID := c.Param("DatabaseID")
	databaseIDInt, err := strconv.Atoi(databaseID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Error": "Bad DatabaseID"})
		return
	}
	_, err = api.psql.GetDatabase(databaseIDInt)
	if err != nil {
		api.log.Error(err)
		c.JSON(http.StatusNotFound, gin.H{"Error": "DatabaseID not found"})
		return
	}
	_, err = api.psql.DeleteDatabase(databaseIDInt)
	if err != nil {
		api.log.Error(err)
		c.JSON(http.StatusNotFound, gin.H{"Error": "DatabaseID not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Result": "Success"})
	return
}
