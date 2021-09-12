package server

import (
	"fmt"
	"os"

	"github.com/common-nighthawk/go-figure"
	"github.com/dhawton/log4g"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/vzau/common/utils"
	"github.com/vzau/thoth/internal/server/middleware"
	"github.com/vzau/thoth/pkg/database"
	"github.com/vzau/thoth/pkg/version"
)

type Server struct {
	engine *gin.Engine
}

var log = log4g.Category("server")

func Run(port int) {
	intro := figure.NewFigure("Thoth", "", false).Slicify()
	for i := 0; i < len(intro); i++ {
		log.Info(intro[i])
	}
	log.Info("Thoth %s", version.FriendlyVersion())
	log.Info("Checking for .env, loading if exists")
	if _, err := os.Stat(".env"); err == nil {
		log.Info("Found .env, loading")
		err := godotenv.Load()
		if err != nil {
			log.Error("Error loading .env file")
		}
	}

	appenv := utils.Getenv("APP_ENV", "prod")
	log.Info("APP_ENV=%s", appenv)
	if appenv == "dev" {
		log4g.SetLogLevel(log4g.DEBUG)
	}
	log.Debug("Done setting log level")

	log.Info("Creating Database Connection")
	database.Connect(utils.Getenv("DB_USERNAME", "root"), utils.Getenv("DB_PASSWORD", "secret12345"), utils.Getenv("DB_HOSTNAME", "localhost"), utils.Getenv("DB_PORT", "3306"), utils.Getenv("DB_NAME", "zau"))

	log.Info("Configuring gin webserver")
	server := NewServer(appenv)
	log.Info("Configuring routes")
	SetupRoutes(server.engine)

	log.Info("Starting, listening on :%d", port)
	server.engine.Run(fmt.Sprintf(":%d", port))
}

func NewServer(appenv string) *Server {
	server := Server{}
	engine := gin.New()
	engine.Use(gin.Recovery())

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowCredentials = true
	engine.Use(cors.New(corsConfig))

	engine.Use(middleware.Logger)

	server.engine = engine
	engine.LoadHTMLGlob("static/*")

	return &server
}
