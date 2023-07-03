// SPDX-FileCopyrightText: 2023 Kavya Shukla <kavyuushukla@gmail.com>
// SPDX-License-Identifier: GPL-2.0-only

package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/fossology/LicenseDb/pkg/authenticate"
	"github.com/fossology/LicenseDb/pkg/db"
	"github.com/fossology/LicenseDb/pkg/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Router() *gin.Engine {
	r := gin.Default()
	r.NoRoute(HandleInvalidUrl)
	authorized := r.Group("/")
	authorized.Use(authenticate.AuthenticationMiddleware())
	authorized.GET("/api/license/:shortname", GetLicense)
	authorized.POST("/api/license", CreateLicense)
	authorized.PATCH("/api/license/update/:shortname", UpdateLicense)
	authorized.GET("/api/licenses", SearchInLicense)
	r.POST("/api/user", authenticate.CreateUser)
	authorized.GET("/api/users", authenticate.GetAllUser)
	authorized.GET("/api/user/:id", authenticate.GetUser)
	return r
}

func HandleInvalidUrl(c *gin.Context) {

	er := models.LicenseError{
		Status:    http.StatusNotFound,
		Message:   "No such path exists please check URL",
		Error:     "invalid path",
		Path:      c.Request.URL.Path,
		Timestamp: time.Now().Format(time.RFC3339),
	}
	c.JSON(http.StatusNotFound, er)
}
func GetAllLicense(c *gin.Context) {
	var licenses []models.License

	err := db.DB.Find(&licenses).Error
	if err != nil {
		er := models.LicenseError{
			Status:    http.StatusBadRequest,
			Message:   "Licenses not found",
			Error:     err.Error(),
			Path:      c.Request.URL.Path,
			Timestamp: time.Now().Format(time.RFC3339),
		}
		c.JSON(http.StatusBadRequest, er)
		return
	}
	res := models.LicenseResponse{
		Data:   licenses,
		Status: http.StatusOK,
		Meta: models.Meta{
			ResourceCount: len(licenses),
		},
	}

	c.JSON(http.StatusOK, res)
}

func GetLicense(c *gin.Context) {
	var license models.License

	queryParam := c.Param("shortname")
	if queryParam == "" {
		return
	}

	err := db.DB.Where("shortname = ?", queryParam).First(&license).Error

	if err != nil {
		er := models.LicenseError{
			Status:    http.StatusBadRequest,
			Message:   fmt.Sprintf("no license with shortname '%s' exists", queryParam),
			Error:     err.Error(),
			Path:      c.Request.URL.Path,
			Timestamp: time.Now().Format(time.RFC3339),
		}
		c.JSON(http.StatusBadRequest, er)
		return
	}

	res := models.LicenseResponse{
		Data:   []models.License{license},
		Status: http.StatusOK,
		Meta: models.Meta{
			ResourceCount: 1,
		},
	}

	c.JSON(http.StatusOK, res)
}

func CreateLicense(c *gin.Context) {
	var input models.LicenseInput

	if err := c.ShouldBindJSON(&input); err != nil {
		er := models.LicenseError{
			Status:    http.StatusBadRequest,
			Message:   fmt.Sprintf("invalid request"),
			Error:     err.Error(),
			Path:      c.Request.URL.Path,
			Timestamp: time.Now().Format(time.RFC3339),
		}
		c.JSON(http.StatusBadRequest, er)
		return
	}

	if input.Active == "" {
		input.Active = "t"
	}

	license := models.License(input)

	db.DB.Create(&license)

	c.JSON(http.StatusOK, gin.H{"data": license})
}

func UpdateLicense(c *gin.Context) {
	var update models.License
	var license models.License
	shortname := c.Param("shortname")
	if err := db.DB.Where("shortname = ?", shortname).First(&license).Error; err != nil {
		er := models.LicenseError{
			Status:    http.StatusBadRequest,
			Message:   fmt.Sprintf("license not found"),
			Error:     err.Error(),
			Path:      c.Request.URL.Path,
			Timestamp: time.Now().Format(time.RFC3339),
		}
		c.JSON(http.StatusBadRequest, er)
		return
	}
	if err := c.ShouldBindJSON(&update); err != nil {
		er := models.LicenseError{
			Status:    http.StatusBadRequest,
			Message:   fmt.Sprintf("invalid request"),
			Error:     err.Error(),
			Path:      c.Request.URL.Path,
			Timestamp: time.Now().Format(time.RFC3339),
		}
		c.JSON(http.StatusBadRequest, er)
		return
	}
	if err := db.DB.Model(&license).Updates(update).Error; err != nil {
		er := models.LicenseError{
			Status:    http.StatusBadRequest,
			Message:   fmt.Sprintf("Failed to update license"),
			Error:     err.Error(),
			Path:      c.Request.URL.Path,
			Timestamp: time.Now().Format(time.RFC3339),
		}
		c.JSON(http.StatusInternalServerError, er)
		return
	}
	res := models.LicenseResponse{
		Data:   []models.License{license},
		Status: http.StatusOK,
		Meta: models.Meta{
			ResourceCount: 1,
		},
	}

	c.JSON(http.StatusOK, res)

}

func SearchInLicense(c *gin.Context) {
	feild := c.Query("feild")
	search_term := c.Query("search_term")
	search := c.Query("search")
	if feild == "" && search_term == "" {
		GetAllLicense(c)
		return
	}
	var query *gorm.DB
	var license []models.License
	if search == "fuzzy" {
		query = db.DB.Where(fmt.Sprintf("%s ILIKE ?", feild), fmt.Sprintf("%%%s%%", search_term)).Find(&license)
	} else {
		query = db.DB.Where(feild+" @@ plainto_tsquery(?)", search_term).Find(&license)
	}

	if err := query.Error; err != nil {
		er := models.LicenseError{
			Status:    http.StatusBadRequest,
			Message:   fmt.Sprintf("incorrect query to search in the database"),
			Error:     err.Error(),
			Path:      c.Request.URL.Path,
			Timestamp: time.Now().Format(time.RFC3339),
		}
		c.JSON(http.StatusBadRequest, er)
		return
	}
	res := models.LicenseResponse{
		Data:   license,
		Status: http.StatusOK,
		Meta: models.Meta{
			ResourceCount: len(license),
		},
	}
	c.JSON(http.StatusOK, res)

}
