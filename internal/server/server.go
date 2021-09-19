/*
ZAU Thoth API
Copyright (C) 2021 Daniel A. Hawton (daniel@hawton.org)

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package server

import (
	"fmt"
	"os"
	"strconv"

	"github.com/common-nighthawk/go-figure"
	"github.com/dhawton/log4g"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/vzau/common/utils"
	"github.com/vzau/thoth/internal/controllers/auth"
	"github.com/vzau/thoth/internal/server/middleware"
	"github.com/vzau/thoth/pkg/cache"
	"github.com/vzau/thoth/pkg/database"
	"github.com/vzau/thoth/pkg/version"
	dbTypes "github.com/vzau/types/database"
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

	if utils.Getenv("SESSION_SECRET", "a") == "a" {
		log.Fatal("SESSION_SECRET is not set and is required.")
	}

	appenv := utils.Getenv("APP_ENV", "prod")
	log.Info("APP_ENV=%s", appenv)
	if appenv == "dev" {
		log4g.SetLogLevel(log4g.DEBUG)
	}
	log.Debug("Done setting log level")

	log.Info("Configuring OAuth2 Client")
	auth.Init()

	log.Info("Creating Database Connection")
	database.Connect(utils.Getenv("DB_USERNAME", "root"), utils.Getenv("DB_PASSWORD", "secret12345"), utils.Getenv("DB_HOSTNAME", "localhost"), utils.Getenv("DB_PORT", "3306"), utils.Getenv("DB_NAME", "zau"))

	log.Info("Running auto migrate")
	err := database.DB.AutoMigrate(&dbTypes.User{}, &dbTypes.Role{}, &dbTypes.Category{}, &dbTypes.File{})
	if err != nil {
		log.Fatal("Error running auto migrate: %s", err)
	}

	log.Info("Building Cache")
	defaultTTL, _ := strconv.ParseInt(utils.Getenv("CACHE_DEFAULT_TTL", "300"), 10, 32)
	if defaultTTL == 0 {
		defaultTTL = 48 * 60 * 60
	}
	err = cache.BuildCache(int(defaultTTL))
	if err != nil {
		log.Fatal("Error building cache: %s", err)
	}

	log.Info("Configuring gin webserver")
	server := NewServer(appenv)
	log.Info("Configuring routes")
	SetupRoutes(server.engine)

	log.Info("Starting, listening on :%d", port)
	server.engine.Run(fmt.Sprintf(":%d", port))
}

func NewServer(appenv string) *Server {
	log.Debug("Setting gin mode to release mode")
	gin.SetMode(gin.ReleaseMode)

	server := Server{}
	engine := gin.New()

	log.Debug("Loading Recovery middleware")
	engine.Use(gin.Recovery())

	log.Debug("Loading CORS middleware")
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowCredentials = true
	engine.Use(cors.New(corsConfig))

	log.Debug("Loading Logger middleware")
	engine.Use(middleware.Logger)

	log.Debug("Loading Session middleware")
	store := cookie.NewStore([]byte(utils.Getenv("SESSION_SECRET", "")))
	store.Options(sessions.Options{
		Path:     "/",
		HttpOnly: true,
		MaxAge:   86400 * 7,
		Domain:   utils.Getenv("SESSION_DOMAIN", ".chicagoartcc.org"),
	})
	engine.Use(sessions.Sessions(utils.Getenv("SESSION_COOKIE", "thoth"), store))

	log.Debug("Loading UpdateCookie middleware")
	engine.Use(middleware.UpdateCookie)

	// This checks for session and loads user info into Context
	log.Debug("Loading auth middleware")
	engine.Use(middleware.Auth)

	log.Debug("Loading HTML globs")
	server.engine = engine
	engine.LoadHTMLGlob("static/*")

	return &server
}
