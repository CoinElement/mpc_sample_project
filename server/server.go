package server

import (
	"fmt"
	"mpc_sample_project/controllers"
	"mpc_sample_project/models"
	"mpc_sample_project/services"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func CreateServer() *http.Server {

	/*
	 configure Logger
	*/
	log := logrus.New()
	log.Out = os.Stdout
	log.Level = logrus.DebugLevel

	if err := godotenv.Load(); err != nil {
		log.Fatal(err.Error())
	}
	LOG_FILE_LOCATION, exists := os.LookupEnv("LOG_FILE_LOCATION")
	if !exists {
		log.Fatal("missing LOG_FILE_LOCATION")
	}
	logfile, err := os.OpenFile(LOG_FILE_LOCATION, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("failed to open file")
	} else {
		log.Out = logfile
		log.Formatter = &logrus.JSONFormatter{}
	}

	/*
	 configure Database
	*/
	db_type, exists := os.LookupEnv("DB_TYPE")
	if !exists {
		log.Fatal("missing DB_TYPE environment variable")
	}
	db_path, exists := os.LookupEnv("DB_PATH")
	if !exists {
		log.Fatal("missing DB_PATH environment variable")
	}
	db := models.NewDB(db_type, db_path)
	if err := db.Connect(); err != nil {
		log.Fatal("db connection failed")
	}

	ADDR_PORT, exists := os.LookupEnv("ADDR_PORT")
	if !exists {
		log.Fatal("missing ADDR_PORT")
		ADDR_PORT = "8080"
	}

	/*
		Initialize Services
	*/
	ms := services.NewMpcService(log)

	/*
		Initialize Controllers
	*/
	mpcController := controllers.NewMpcController(log, ms)

	/*
		Initialize gin
	*/
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(CORSMiddleware())

	r.GET("/start", mpcController.HandleStart)
	r.GET("/ping", mpcController.HandlePing)
	r.POST("/commit", mpcController.HandleCommitment)         // 接收commitment
	r.POST("/result", mpcController.HandleCommitment)         // 接收上一家的 result 或 party_k 的 noise
	r.POST("/notification", mpcController.HandleNotification) // 接收参与计算的请求

	/*
		Start HTTP Server
	*/
	// initialize server
	addr := fmt.Sprintf("%s:%s", "0.0.0.0", ADDR_PORT)
	server := makeServer(addr, r)

	go handleGracefulShutdown(server)

	return server
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, X-Auth-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
