package app

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func (app *Application) Notification(w http.ResponseWriter, r *http.Request) {
	var body NotificationRequest
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var taskInfo TaskInfo
	val, err := app.DB.RedisClient.Get(ctx, body.TaskId).Result()
	if err != nil {
		log.Fatal(err)
		return
	}

	err = json.Unmarshal([]byte(val), &taskInfo)
	if err != nil {
		log.Fatal(err)
		return
	}
	go app.GetJobs(taskInfo)
	fmt.Println(taskInfo)
}

func (app *Application) GetJobs(taskInfo TaskInfo) {

}
