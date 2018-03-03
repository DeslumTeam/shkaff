package api

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (api *API) checkTaskParameters(c *gin.Context) (setStrings map[string]interface{}, err error) {
	var errStr, setString string
	var setList []string
	var taskUpdate map[string]string
	setStrings = make(map[string]interface{})
	c.BindJSON(&taskUpdate)
	for key, val := range taskUpdate {
		switch key {
		case "task_name":
			if val == "" {
				errStr = fmt.Sprintf("In %s bad value %s", key, val)
				return nil, errors.New(errStr)
			}
			setStrings[key] = val
		case "verb":
			valInt, err := strconv.Atoi(val)
			if err != nil || valInt > 6 {
				errStr = fmt.Sprintf("In %s bad value %s", key, val)
				return nil, errors.New(errStr)
			}
			setStrings[key] = valInt
		case "thread_count":
			valInt, err := strconv.Atoi(val)
			if err != nil || valInt > 10 {
				errStr = fmt.Sprintf("In %s bad value %s", key, val)
				return nil, errors.New(errStr)
			}
			setStrings[key] = val
		case "gzip", "ipv6", "is_active":
			_, err := strconv.ParseBool(val)
			if err != nil {
				errStr = fmt.Sprintf("In %s bad value %s", key, val)
				return nil, errors.New(errStr)
			}
			setStrings[key] = val
		case "months":
			// if len(val) > 0 {
			// 	strVals := strings.Split(val, ",")
			// 	for _, valStr := range strVals {
			// 		valInt, err := strconv.Atoi(valStr)
			// 		if err != nil {
			// 			errStr = fmt.Sprintf("In %s bad value %s", key, val)
			// 			return nil, errors.New(errStr)
			// 		}
			// 		if valInt < 1 || valInt > 12 {
			// 			errStr = fmt.Sprintf("In %s bad value %s", key, val)
			// 			return nil, errors.New(errStr)
			// 		}
			// 	}
			setStrings[key] = val
			// } else {
			// 	setStrings[key] = "{}"
			// }
		case "days":
			// if len(val) > 0 {
			// 	strVals := strings.Split(val, ",")
			// 	for _, valStr := range strVals {
			// 		valInt, err := strconv.Atoi(valStr)
			// 		if err != nil {
			// 			errStr = fmt.Sprintf("In %s bad value %s", key, val)
			// 			return nil, errors.New(errStr)
			// 		}
			// 		if valInt < 1 || valInt > 31 {
			// 			errStr = fmt.Sprintf("In %s bad value %s", key, val)
			// 			return nil, errors.New(errStr)
			// 		}
			// 	}
			setStrings[key] = val
			// } else {
			// 	setStrings[key] = "{}"
			// }
		case "hours":
			// if len(val) > 0 {
			// 	strVals := strings.Split(val, ",")
			// 	for _, valStr := range strVals {
			// 		valInt, err := strconv.Atoi(valStr)
			// 		if err != nil {
			// 			errStr = fmt.Sprintf("In %s bad value %s", key, val)
			// 			return nil, errors.New(errStr)
			// 		}
			// 		if valInt < 0 || valInt > 23 {
			// 			errStr = fmt.Sprintf("In %s bad value %s", key, val)
			// 			return nil, errors.New(errStr)
			// 		}
			// 	}
			setStrings[key] = val
			// } else {
			// 	setStrings[key] = "{}"
			// }
		case "minutes":
			valInt, err := strconv.Atoi(val)
			if err != nil || valInt < 0 || valInt > 60 {
				errStr = fmt.Sprintf("In %s bad value %s", key, val)
				return nil, errors.New(errStr)
			}
			setStrings[key] = valInt
		case "db_id":
			valInt, err := strconv.Atoi(val)
			if err != nil {
				errStr = fmt.Sprintf("In %s bad value %s", key, val)
				return nil, errors.New(errStr)
			}
			setStrings[key] = valInt
		case "databases":
			setStrings[key] = val
		default:
			errStr = fmt.Sprintf("Bad field %s", key)
			return nil, errors.New(errStr)
		}
		setList = append(setList, setString)
	}
	return
}

func (api *API) checkDatabaseParameters(c *gin.Context) (setStrings map[string]interface{}, err error) {
	var errStr, setString string
	var setList []string
	var databaseUpdate map[string]string
	setStrings = make(map[string]interface{})
	c.BindJSON(&databaseUpdate)
	for key, val := range databaseUpdate {
		switch key {
		case "user_id", "type_id":
			valInt, err := strconv.Atoi(val)
			if err != nil {
				errStr = fmt.Sprintf("In %s bad value %s", key, val)
				return nil, errors.New(errStr)
			}
			setStrings[key] = valInt
		case "server_name":
			if val == "" {
				errStr = fmt.Sprintf("In %s bad value %s", key, val)
				return nil, errors.New(errStr)
			}
			setStrings[key] = val
		case "host":
			if val == "" {
				errStr = fmt.Sprintf("In %s bad value %s", key, val)
				return nil, errors.New(errStr)
			}
			setStrings[key] = val
		case "port":
			valInt, err := strconv.Atoi(val)
			if err != nil || valInt < 1024 && valInt > 65565 {
				errStr = fmt.Sprintf("In %s bad value %s", key, val)
				return nil, errors.New(errStr)
			}
			setStrings[key] = valInt
		case "is_active":
			_, err := strconv.ParseBool(val)
			if err != nil {
				errStr = fmt.Sprintf("In %s bad value %s", key, val)
				return nil, errors.New(errStr)
			}
			setStrings[key] = val
		case "db_user", "custom_name", "db_password":
			setStrings[key] = val
		default:
			errStr = fmt.Sprintf("Bad field %s", key)
			return nil, errors.New(errStr)
		}
		setList = append(setList, setString)
	}
	return
}

func (api *API) checkUserParameters(c *gin.Context) (setStrings map[string]interface{}, err error) {
	var errStr string
	var userUpdate map[string]string
	setStrings = make(map[string]interface{})
	c.BindJSON(&userUpdate)
	for key, val := range userUpdate {
		switch key {
		case "login", "api_token":
			if val == "" {
				errStr = fmt.Sprintf("In %s bad value %s", key, val)
				return nil, errors.New(errStr)
			}
			setStrings[key] = val
		case "password":
			if val == "" {
				errStr = fmt.Sprintf("In %s bad value %s", key, val)
				return nil, errors.New(errStr)
			}
			hasher := md5.New()
			hasher.Write([]byte(val))
			setStrings[key] = hex.EncodeToString(hasher.Sum(nil))
		case "is_active", "is_admin":
			_, err := strconv.ParseBool(val)
			if err != nil {
				errStr = fmt.Sprintf("In %s bad value %s", key, val)
				return nil, errors.New(errStr)
			}
			setStrings[key] = val
		case "first_name", "last_name":
			setStrings[key] = val
		default:
			errStr = fmt.Sprintf("Bad field %s", key)
			return nil, errors.New(errStr)
		}
	}
	return
}
