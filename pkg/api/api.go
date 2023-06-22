// SPDX-FileCopyrightText: 2023 Kavya Shukla <kavyuushukla@gmail.com>
// SPDX-License-Identifier: GPL-2.0-only

package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/fossology/LicenseDb/pkg/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var DB *gorm.DB

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

	err := DB.Find(&licenses).Error
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

	err := DB.Where("shortname = ?", queryParam).First(&license).Error

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

	DB.Create(&license)

	c.JSON(http.StatusOK, gin.H{"data": license})
}

func UpdateLicense(c *gin.Context) {
	var update models.License
	var license models.License
	shortname := c.Param("shortname")
	if err := DB.First(&license, shortname).Error; err != nil {
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
	if err := DB.Model(&license).Updates(update).Error; err != nil {
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
