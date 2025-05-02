package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	getJob "github.com/MinhNHHH/get-job/pkg/app"
	"github.com/MinhNHHH/get-job/pkg/cfgs"
	"github.com/MinhNHHH/get-job/pkg/database/repository/dbrepo"
)

func main() {
	app := getJob.Application{}
	cfgs := cfgs.LoadConfigs()
	sqlConn, err := app.ConnectDB(cfgs.DB_CONNECTION_URI)
	if err != nil {
		log.Fatal(err)
		return
	}
	redisConn, err := app.ConnectReids(cfgs.REDIS_CONNECTION_URI)
	if err != nil {
		log.Fatal(err)
		return
	}
	app.DB = &dbrepo.DBRepo{SqlConn: sqlConn, RedisConn: redisConn}
	app.Cfg = &cfgs

	if len(os.Args) < 2 {
		fmt.Println("Usage: go run cmd/main.go [command]")
		fmt.Println("Commands:")
		fmt.Println("  migrate [steps] - Run migrations (optional number of steps)")
		fmt.Println("  create [name]  - Create new migration files")
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
			app.Migrate(step)
		} else {
			app.Migrate(0)
		}
	case "create":
		if len(os.Args) < 3 {
			fmt.Println("Please provide a name for the migration")
			os.Exit(1)
		}
		app.GenerateMigration(os.Args[2])
	default:
		// start the server
		err := http.ListenAndServe(":8080", app.Routes())
		if err != nil {
			log.Fatal(err)
		}
	}
}
