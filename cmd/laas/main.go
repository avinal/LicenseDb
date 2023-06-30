// SPDX-FileCopyrightText: 2023 Kavya Shukla <kavyuushukla@gmail.com>
// SPDX-License-Identifier: GPL-2.0-only

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/fossology/LicenseDb/pkg/api"
	"github.com/fossology/LicenseDb/pkg/authenticate"
	"github.com/fossology/LicenseDb/pkg/models"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// declare flags to input the basic requirement of database connection and the path of the data file
var (
	// argument to enter the name of database host
	dbhost = flag.String("host", "localhost", "host name")
	// port number of the host
	port = flag.String("port", "5432", "port number")
	// argument to enter the database user
	user = flag.String("user", "fossy", "user name")
	// name of database to be connected
	dbname = flag.String("dbname", "fossology", "database name")
	// password of the database
	password = flag.String("password", "fossy", "password")
	// path of data file
	datafile = flag.String("datafile", "licenseRef.json", "datafile path")
	// auto-update the database
	populatedb = flag.Bool("populatedb", false, "boolean variable to update database")
)

func main() {
	flag.Parse()

	dburi := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s", *dbhost, *port, *user, *dbname, *password)
	gormConfig := &gorm.Config{}
	database, err := gorm.Open(postgres.Open(dburi), gormConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err := database.AutoMigrate(&models.License{}); err != nil {
		log.Fatalf("Failed to automigrate database: %v", err)
	}

	if err := database.AutoMigrate(&models.User{}); err != nil {
		log.Fatalf("Failed to automigrate database: %v", err)
	}
	if *populatedb {
		var licenses []models.License
		// read the file of data
		byteResult, _ := ioutil.ReadFile(*datafile)
		// unmarshal the json file and it into the struct format
		if err := json.Unmarshal(byteResult, &licenses); err != nil {
			log.Fatalf("error reading from json file: %v", err)
		}
		for _, license := range licenses {
			// populate the data in the database table
			database.Create(&license)
		}
	}
	api.DB = database

	r := gin.Default()
	r.NoRoute(api.HandleInvalidUrl)
	authorized := r.Group("/")
	authorized.Use(authenticate.AuthenticationMiddleware())
	authorized.GET("/api/license/:shortname", api.GetLicense)
	authorized.POST("/api/license", api.CreateLicense)
	authorized.PATCH("/api/license/update/:shortname", api.UpdateLicense)
	authorized.GET("/api/licenses", api.SearchInLicense)
	r.POST("/api/user", authenticate.CreateUser)
	authorized.GET("/api/users", authenticate.GetAllUser)
	authorized.GET("/api/user/:id", authenticate.GetUser)
	r.Run()
}
