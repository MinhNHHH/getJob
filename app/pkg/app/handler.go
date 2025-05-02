package app

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/MinhNHHH/get-job/pkg/database/data"
)

func (app *Application) Notification(w http.ResponseWriter, r *http.Request) {
	var body NotificationRequest
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var taskInfo TaskInfo
	val, err := app.DB.RedisGet(body.TaskId)
	if err != nil {
		log.Fatal(err)
		return
	}

	err = json.Unmarshal([]byte(val), &taskInfo)
	if err != nil {
		log.Fatal(err)
		return
	}
	jobChan := make(chan TaskInfo)
	go func() {
		jobChan <- taskInfo
	}()
	go app.InsertJob(<-jobChan)
}

func (app *Application) InsertJob(taskInfo TaskInfo) {
	for _, job := range taskInfo.Results {
		log.Printf("Processing job: %+v", job)
		existed, companyId := app.DB.IsExisted(job.CompanyName)
		if !existed {
			companyId, err := app.DB.InsertCompany(&data.Companies{
				Name: job.CompanyName,
				Url:  job.ComapnyUrl,
			})
			if err != nil {
				log.Printf("Error inserting company %s: %v", job.CompanyName, err)
				log.Fatal(err)
				return
			}
			log.Printf("Inserted company %s with ID %d", job.CompanyName, companyId)
		}
		if app.DB.IsJobExisted(job.Title, job.Location, companyId) {
			continue
		}
		_, err := app.DB.InsertJob(&data.Job{
			Title:       job.Title,
			Location:    job.Location,
			Description: job.Description,
			CompanyId:   companyId,
		})
		if err != nil {
			log.Printf("Error inserting job %s: %v", job.Title, err)
			log.Fatal(err)
			return
		}
		log.Printf("Successfully inserted job %s for company %s", job.Title, job.CompanyName)
	}
}
