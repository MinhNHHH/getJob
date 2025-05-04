package app

import (
	"net/http"

	"github.com/MinhNHHH/get-job/pkg/cfgs"
	"github.com/MinhNHHH/get-job/pkg/database/repository"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type Application struct {
	DB  repository.DatabaseRepo
	Cfg *cfgs.Config
}

type NotificationRequest struct {
	TaskId   string `json:"task_id"`
	TaskName string `json:"task_name"`
}

type TaskInfo struct {
	Id        string      `json:"id"`
	TaskName  string      `json:"task_name"`
	StartTime float64     `json:"start_time"`
	EndTime   float64     `json:"end_time"`
	Results   []JobDetail `json:"results"`
}

type JobDetail struct {
	Title       string `json:"title"`
	CompanyName string `json:"company_name"`
	ComapnyUrl  string `json:"company_uri"`
	Location    string `json:"location"`
	Description string `json:"description"`
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers for every response
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")

		// Handle preflight OPTIONS request
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Continue to actual handler
		next.ServeHTTP(w, r)
	})
}

// var ctx = context.Background()

func (app *Application) Routes() http.Handler {
	mux := chi.NewRouter()
	// register middleware
	mux.Use(middleware.Recoverer)
	mux.Use(withCORS)
	// register routes
	mux.Post("/notification", app.Notification)
	mux.Route("/api", func(r chi.Router) {
		r.Post("/generate-cover-letter", app.GenerateCoverLetter)
	})
	return mux
}
