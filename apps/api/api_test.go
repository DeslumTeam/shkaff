package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
)

const (
	HOST = "0.0.0.0"
	PORT = 8080
)

func TestGetDatabaseRequests(t *testing.T) {
	reqString := []string{"GetUser", "GetDatabase", "GetTask"}
	for _, req := range reqString {
		url := fmt.Sprintf("http://%s:%d/api/v1/%s/1?token=12345", HOST, PORT, req)
		resp, err := http.Get(url)
		if err != nil {
			log.Panicln(err)
		}
		getResponse(resp)
	}
}

func TestUpdateUser(t *testing.T) {
	form := map[string]string{
		"login":      "forTest",
		"password":   "hopley",
		"api_token":  "12345",
		"first_name": "Mike",
		"last_name":  "Polonsky",
		"is_active":  "true",
		"is_admin":   "false",
	}
	jsonStr, err := json.Marshal(form)
	if err != nil {
		log.Panicln(err)
	}
	url := fmt.Sprintf("http://%s:%d/api/v1/UpdateUser/1?token=12345", HOST, PORT)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	getResponse(resp)
	defer resp.Body.Close()

}

func TestUpdateDatabase(t *testing.T) {
	form := map[string]string{

		"server_name": "ModifedByAPI",
		"custom_name": "CustomByAPI",
		"host":        "127.0.0.2",
		"port":        "27019",
		"is_active":   "false",
		"db_user":     "test",
		"db_password": "nopassword",
	}
	jsonStr, err := json.Marshal(form)
	if err != nil {
		log.Panicln(err)
	}
	url := fmt.Sprintf("http://%s:%d/api/v1/UpdateDatabase/1?token=12345", HOST, PORT)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	getResponse(resp)
	defer resp.Body.Close()

}

func TestUpdateTask(t *testing.T) {
	form := map[string]string{

		"verb":         "3",
		"task_name":    "TaskModifedByAPI",
		"is_active":    "false",
		"thread_count": "5",
		"ipv6":         "false",
		"databases":    "{}",
		"gzip":         "false",
		"months":       "{1,2,3,4,5,6,7,8,9}",
		"days":         "{12,14,17}",
		"hours":        "{1}",
		"minutes":      "50",
	}
	jsonStr, err := json.Marshal(form)
	if err != nil {
		log.Panicln(err)
	}
	url := fmt.Sprintf("http://%s:%d/api/v1/UpdateTask/1?token=12345", HOST, PORT)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	getResponse(resp)
	defer resp.Body.Close()

}

func TestCreateUser(t *testing.T) {
	form := map[string]string{
		"login":      "apitest",
		"password":   "ifull",
		"api_token":  "123451",
		"first_name": "Orly",
		"last_name":  "Nike",
		"is_active":  "true",
		"is_admin":   "false",
	}
	jsonStr, err := json.Marshal(form)
	if err != nil {
		log.Panicln(err)
	}
	url := fmt.Sprintf("http://%s:%d/api/v1/CreateUser?token=12345", HOST, PORT)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	getResponse(resp)
	defer resp.Body.Close()

}

func TestCreateDatabase(t *testing.T) {
	form := map[string]string{
		"user_id":     "1",
		"type_id":     "1",
		"server_name": "ModifedByAPI1",
		"custom_name": "CustomByAPI1",
		"host":        "127.0.0.4",
		"port":        "27019",
		"is_active":   "true",
		"db_user":     "test",
		"db_password": "nopassword",
	}
	jsonStr, err := json.Marshal(form)
	if err != nil {
		log.Panicln(err)
	}
	url := fmt.Sprintf("http://%s:%d/api/v1/CreateDatabase?token=12345", HOST, PORT)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != 200 {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		log.Fatal(string(bodyBytes))
	}
	defer resp.Body.Close()

}

func TestCreateTask(t *testing.T) {
	form := map[string]string{
		"db_id":        "1",
		"verb":         "3",
		"task_name":    "TaskModifedByAPI1",
		"is_active":    "false",
		"thread_count": "5",
		"ipv6":         "false",
		"databases":    "{}",
		"gzip":         "false",
		"months":       "{1,2,3,4,5,6,7,8,9}",
		"days":         "{12,14,17}",
		"hours":        "{1}",
		"minutes":      "50",
	}
	jsonStr, err := json.Marshal(form)
	if err != nil {
		log.Panicln(err)
	}
	url := fmt.Sprintf("http://%s:%d/api/v1/CreateTask?token=12345", HOST, PORT)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	getResponse(resp)
	defer resp.Body.Close()

}

func TestDeleteRequests(t *testing.T) {
	client := &http.Client{}
	reqString := []string{"DeleteDatabase", "DeleteTask", "DeleteUser"}
	for _, req := range reqString {
		url := fmt.Sprintf("http://%s:%d/api/v1/%s/1?token=12345", HOST, PORT, req)
		req, err := http.NewRequest("DELETE", url, nil)
		if err != nil {
			log.Fatalln(err)
		}
		resp, err := client.Do(req)
		if err != nil {
			log.Fatalln(err)
		}
		getResponse(resp)
		defer resp.Body.Close()
	}
}

func getResponse(resp *http.Response) {
	if resp.StatusCode != 200 {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		log.Fatal(string(bodyBytes))
	}
}
