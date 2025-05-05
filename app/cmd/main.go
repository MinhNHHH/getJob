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
	"github.com/MinhNHHH/get-job/pkg/llm"
)

func initializeApp(cfgs cfgs.Configs) getJob.Application {
	app := getJob.Application{}

	sqlConn, err := app.ConnectDB(cfgs.DB_CONNECTION_URI)
	if err != nil {
		log.Fatal(err)
	}
	redisConn, err := app.ConnectReids(cfgs.REDIS_CONNECTION_URI)
	if err != nil {
		log.Fatal(err)
	}
	app.DB = &dbrepo.DBRepo{SqlConn: sqlConn, RedisConn: redisConn}
	app.LLM = llm.NewLLM(cfgs)
	app.Cfgs = cfgs
	return app
}

func startServer(app getJob.Application) {
	log.Println("Starting server on :8080")
	err := http.ListenAndServe(":8080", app.Routes())
	if err != nil {
		log.Fatal(err)
	}
}

func handleMigrate(app getJob.Application, args []string) {
	step := 0
	if len(args) == 1 {
		var err error
		step, err = strconv.Atoi(args[0])
		if err != nil {
			log.Fatal(err)
		}
	}
	app.Migrate(step, app.Cfgs.DB_CONNECTION_URI)
}

func handleCreate(app getJob.Application, args []string) {
	if len(args) < 1 {
		fmt.Println("Please provide a name for the migration")
		os.Exit(1)
	}
	app.GenerateMigration(args[0])
}

func handleCommand(app getJob.Application, command string, args []string) {
	switch command {
	case "migrate":
		handleMigrate(app, args)
	case "create":
		handleCreate(app, args)
	default:
		startServer(app)
	}
}

func main() {
	cfgs := cfgs.LoadConfigs()
	app := initializeApp(cfgs)

	if len(os.Args) < 2 {
		fmt.Println("Usage: go run cmd/main.go [command]")
		fmt.Println("Commands:")
		fmt.Println("  migrate [steps] - Run migrations (optional number of steps)")
		fmt.Println("  create [name]   - Create new migration files")
		fmt.Println("  start           - Run server")
		os.Exit(1)
	}
	command := os.Args[1]
	args := os.Args[2:]

	handleCommand(app, command, args)
}
