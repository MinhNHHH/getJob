package app

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/MinhNHHH/get-job/pkg/database/data"
	"github.com/MinhNHHH/get-job/pkg/llm"
)

type GenerateCoverLetterRequest struct {
	JobTitle    string `json:"job_title"`
	CompanyName string `json:"company_name"`
	Location    string `json:"location"`
	Description string `json:"description"`
}

func (app *Application) Notification(w http.ResponseWriter, r *http.Request) {
	var body NotificationRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	val, err := app.DB.RedisGet(body.TaskId)
	if err != nil {
		log.Println("Redis error:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	var taskInfo TaskInfo
	if err := json.Unmarshal([]byte(val), &taskInfo); err != nil {
		log.Println("Unmarshal error:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	go app.InsertJob(taskInfo)
}

func (app *Application) InsertJob(taskInfo TaskInfo) {
	for _, job := range taskInfo.Results {
		log.Printf("Processing job: %+v", job)

		existed, companyId := app.DB.IsExisted(job.CompanyName)
		if !existed {
			var err error
			companyId, err = app.DB.InsertCompany(&data.Companies{
				Name: job.CompanyName,
				Url:  job.CompanyUrl,
			})
			if err != nil {
				log.Printf("Error inserting company %s: %v", job.CompanyName, err)
				return
			}
			log.Printf("Inserted company %s with ID %d", job.CompanyName, companyId)
		}

		if app.DB.IsJobExisted(job.Title, job.Location, companyId) {
			log.Printf("Job already exists: %s at %s (Company ID: %d)", job.Title, job.Location, companyId)
			continue
		}

		_, err := app.DB.InsertJob(&data.Job{
			Title:       job.Title,
			Location:    job.Location,
			Description: job.Description,
			CompanyId:   companyId,
		})
		if err != nil {
			log.Printf("Error inserting job %s: %v", job.Title, err) // fixed incorrect log.Fatal usage
			return
		}
		log.Printf("Successfully inserted job %s for company %s", job.Title, job.CompanyName)
	}
}

func (app *Application) GenerateCoverLetter(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	// Parse job details
	jobDetailsStr := r.FormValue("job_details")
	if jobDetailsStr == "" {
		http.Error(w, "Missing job details", http.StatusBadRequest)
		return
	}

	var jobDetails GenerateCoverLetterRequest
	if err := json.Unmarshal([]byte(jobDetailsStr), &jobDetails); err != nil {
		http.Error(w, "Invalid job details", http.StatusBadRequest)
		return
	}

	// Handle resume file
	var resumeContent string
	file, _, err := r.FormFile("resume")
	if err != nil && err != http.ErrMissingFile {
		log.Printf("Error reading resume file: %v", err)
		http.Error(w, "Error reading resume file", http.StatusInternalServerError)
		return
	}
	if file != nil {
		defer file.Close()
		fileBytes, err := io.ReadAll(file)
		if err != nil {
			log.Printf("Error reading resume file: %v", err)
			http.Error(w, "Error reading resume file", http.StatusInternalServerError)
			return
		}
		resumeContent, err = ExtractTextFromBytes(fileBytes)
		if err != nil {
			log.Printf("Error extracting text from resume: %v", err)
			http.Error(w, "Error processing resume", http.StatusInternalServerError)
			return
		}
	} else {
		log.Println("No resume file uploaded")
	}

	// Generate the cover letter
	coverLetter, err := app.LLM.GenerateCoverLetter(llm.JobInfo{
		JobTitle:    jobDetails.JobTitle,
		CompanyName: jobDetails.CompanyName,
		Description: jobDetails.Description,
		Location:    jobDetails.Location,
	}, resumeContent)

	if err != nil {
		log.Printf("Failed to generate cover letter: %v", err)
		http.Error(w, "Failed to generate cover letter", http.StatusInternalServerError)
		return
	}

	// Respond with a structured JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(coverLetter)
}

func GetAllJobs(w http.ResponseWriter, r *http.Request) {

}
