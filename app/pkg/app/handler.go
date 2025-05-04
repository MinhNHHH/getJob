package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/MinhNHHH/get-job/pkg/database/data"
	"github.com/MinhNHHH/get-job/pkg/llm"
	"rsc.io/pdf"
)

type GenerateCoverLetterRequest struct {
	JobTitle    string `json:"job_title"`
	CompanyName string `json:"company_name"`
	Location    string `json:"location"`
	Description string `json:"description"`
}

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

func (app *Application) GenerateCoverLetter(w http.ResponseWriter, r *http.Request) {
	var resumeContent []byte
	// Parse multipart form
	err := r.ParseMultipartForm(10 << 20) // 10 MB max file size
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	// Get job details
	jobDetailsStr := r.FormValue("job_details")
	var jobDetails GenerateCoverLetterRequest
	err = json.Unmarshal([]byte(jobDetailsStr), &jobDetails)
	if err != nil {
		http.Error(w, "Invalid job details", http.StatusBadRequest)
		return
	}

	// Get file
	file, _, err := r.FormFile("resume")
	if err != nil {
		file = nil
	}
	defer file.Close()

	// Check if file is nil before attempting to close
	if file == nil {
		log.Println("No resume file uploaded")
	} else {
		// Read the file into memory
		resumeContent, err = io.ReadAll(file)
		if err != nil {
			log.Printf("Error reading resume file: %v", err)
			http.Error(w, "Error reading resume file", http.StatusInternalServerError)
			return
		}

		// Check if it's a PDF file
		fileHeader := r.MultipartForm.File["resume"][0]
		if strings.HasSuffix(strings.ToLower(fileHeader.Filename), ".pdf") {
			// Parse PDF content
			pdfReader, err := pdf.NewReader(bytes.NewReader(resumeContent), int64(len(resumeContent)))
			if err != nil {
				log.Printf("Error parsing PDF: %v", err)
				http.Error(w, "Error parsing PDF file", http.StatusInternalServerError)
				return
			}
			var textBuilder strings.Builder
			for pageNum := 1; pageNum <= pdfReader.NumPage(); pageNum++ {
				page := pdfReader.Page(pageNum)
				content := page.Content()

				// Extract text and add spaces between words
				var pageText strings.Builder
				for _, text := range content.Text {
					pageText.WriteString(text.S)
				}

				// Split into sentences using regex
				sentences := regexp.MustCompile(`[.!?]+`).Split(pageText.String(), -1)
				for _, sentence := range sentences {
					sentence = strings.TrimSpace(sentence)
					textBuilder.WriteString(sentence + "\n")
				}
			}

			// Replace binary content with extracted text
			resumeContent = []byte(textBuilder.String())
		}

		log.Printf("Resume file read successfully, size: %d bytes", len(resumeContent))
	}

	fmt.Printf("Resume content: %s", string(resumeContent))
	// Process the request...
	coverLetter, err := llm.GenerateCoverLetter(jobDetails.JobTitle, jobDetails.CompanyName, jobDetails.Location, jobDetails.Description)
	if err != nil {
		http.Error(w, "Failed to generate cover letter", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(coverLetter)
}
