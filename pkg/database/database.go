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

package database

import (
	"fmt"
	golog "log"
	"os"
	"time"

	"github.com/dhawton/log4g"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB
var MaxAttempts = 10
var DelayBetweenAttempts = time.Minute * 1
var attempt = 1
var log = log4g.Category("db")

func Connect(user string, pass string, hostname string, port string, database string) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True", user, pass, hostname, port, database)
	newLogger := logger.New(
		golog.New(os.Stdout, "\r\n", golog.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Silent,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: newLogger})
	sqlDB, _ := db.DB()
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetMaxIdleConns(10)
	if err != nil {
		log.Error("Error connecting to database: " + err.Error())
		if attempt < MaxAttempts {
			log.Info("Attempt %d/%d failed. Waiting %s before trying again", attempt, MaxAttempts, DelayBetweenAttempts.String())
			time.Sleep(DelayBetweenAttempts)
			attempt += 1
			Connect(user, pass, hostname, port, database)
			return
		}
		log.Fatal("Max attempts occurred. Aborting startup.")
	}

	DB = db
}
