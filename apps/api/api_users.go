package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (api *API) createUser(c *gin.Context) {
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
	setStrings, err := api.checkUserParameters(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Error": err.Error()})
		return
	}
	_, err = api.psql.CreateUser(setStrings)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, "OK")
	return
}

func (api *API) getUser(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusNotFound, gin.H{"Error": "Token is empty"})
		return
	}
	isExist, err := api.psql.GetUserByToken(token)
	if err != nil || !isExist {
		api.log.Error(err)
		c.JSON(http.StatusNotFound, gin.H{"Error": "User with this token not found"})
		return
	}
	UserID := c.Param("UserID")
	UserIDInt, err := strconv.Atoi(UserID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Error": "Bad UserID"})
		return
	}
	user, err := api.psql.GetUser(UserIDInt)
	if err != nil {
		api.log.Error(err)
		c.JSON(http.StatusNotFound, gin.H{"Error": "UserID not found"})
		return
	}
	c.JSON(http.StatusOK, user)
	return
}

func (api *API) updateUser(c *gin.Context) {
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
	userID := c.Param("UserID")
	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Error": err.Error()})
		return
	}
	_, err = api.psql.GetUser(userIDInt)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Error": "DatabaseID not found"})
		return
	}
	setStrings, err := api.checkUserParameters(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Error": err.Error()})
		return
	}
	api.log.Info(setStrings)
	_, err = api.psql.UpdateUser(userIDInt, setStrings)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Status": "OK"})
	return
}

func (api *API) deleteUser(c *gin.Context) {
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
	userID := c.Param("UserID")
	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Error": "Bad DatabaseID"})
		return
	}
	_, err = api.psql.GetUser(userIDInt)
	if err != nil {
		api.log.Error(err)
		c.JSON(http.StatusNotFound, gin.H{"Error": "DatabaseID not found"})
		return
	}
	_, err = api.psql.DeleteUser(userIDInt)
	if err != nil {
		api.log.Error(err)
		c.JSON(http.StatusNotFound, gin.H{"Error": "DatabaseID not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Result": "Success"})
	return
}
