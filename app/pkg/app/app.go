package app

import (
	"context"
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

var ctx = context.Background()

func (app *Application) Routes() http.Handler {
	mux := chi.NewRouter()
	// register middleware
	mux.Use(middleware.Recoverer)
	// register routes
	mux.Post("/notification", app.Notification)
	return mux
}
