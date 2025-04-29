package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	getJob "github.com/MinhNHHH/get-job/pkg/app"
	"github.com/MinhNHHH/get-job/pkg/cfgs"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	app := getJob.Application{}
	cfgs := cfgs.LoadConfigs()
	db := getJob.NewDB(cfgs)
	app.DB = db
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run cmd/main.go [command]")
		fmt.Println("Commands:")
		fmt.Println("  up        - Run all up migrations")
		fmt.Println("  down      - Run all down migrations")
		fmt.Println("  create    - Create new migration files")
		os.Exit(1)
	}

	command := os.Args[1]
	switch command {
	case "migrate":
		if len(os.Args) == 3 {
			step, err := strconv.Atoi(os.Args[2])
			if err != nil {
				log.Fatal(err)
			}
			db.Migrate(step)
		} else {
			db.Migrate(0)
		}
	case "create":
		if len(os.Args) < 3 {
			fmt.Println("Please provide a name for the migration")
			os.Exit(1)
		}
		db.GenerateMigration(os.Args[2])
	default:
		// start the server
		err := http.ListenAndServe(":8080", app.Routes())
		if err != nil {
			log.Fatal(err)
		}
	}
}
