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

package cdn

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/dhawton/log4g"
	"github.com/gin-gonic/gin"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/vzau/common/utils"
	"github.com/vzau/thoth/internal/cachedStorage"
	"github.com/vzau/thoth/internal/server/response"
	"github.com/vzau/thoth/pkg/cache"
	"github.com/vzau/thoth/pkg/database"
	"github.com/vzau/thoth/pkg/discord"
	"github.com/vzau/thoth/pkg/storage"
	"github.com/vzau/thoth/pkg/user"
	dbTypes "github.com/vzau/types/database"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var log = log4g.Category("controllers/cdn")

func GetCDN(c *gin.Context) {
	files := []dbTypes.File{}
	resultFiles := []dbTypes.File{}

	if err := database.DB.Preload(clause.Associations).Find(&files).Error; err != nil {
		log.Error("Error fetching files: %s", err.Error())
		response.RespondMessage(c, http.StatusInternalServerError, "Error fetching files")
		return
	}

	filterArray, hasArray := c.GetQueryArray("filter")
	if hasArray {
		resultFiles = filterFiles(files, filterArray)
	} else {
		filterString, hasString := c.GetQuery("filter")
		if hasString {
			resultFiles = filterFiles(files, []string{filterString})
		} else {
			resultFiles = files
		}
	}

	response.Respond(c, http.StatusOK, struct {
		Files []FileDTO `json:"files"`
	}{Files: setDownloadUrl(resultFiles)})
}

func GetCDNFile(c *gin.Context) {
	id := c.Param("id")
	file := dbTypes.File{}
	if err := database.DB.Find(&file, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			response.RespondMessage(c, http.StatusNotFound, "File not found")
			return
		}
		log.Error("Error fetching file: %s", err.Error())
		response.RespondMessage(c, http.StatusInternalServerError, "Error fetching file")
		return
	}

	if file.Category.Name == "Staff" {
		u, exists := c.Get("x-user")
		if !exists || u == nil || !user.HasRolesWithUser(u.(*dbTypes.User), user.StaffRolesTraining) {
			response.RespondMessage(c, http.StatusForbidden, "Forbidden")
			return
		}
	}

	reader, err := cachedStorage.GetCachedFile(file)
	if err != nil {
		log.Error("Error fetching file: %s", err.Error())
		response.RespondMessage(c, http.StatusInternalServerError, "Error fetching file")
		return
	}

	if strings.Contains(file.ContentType, "image/") || strings.EqualFold(file.ContentType, "application/pdf") {
		c.Writer.Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=%s", file.Filename))
	} else {
		c.Writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", file.Filename))
	}
	c.Writer.Header().Set("Content-Type", file.ContentType)

	http.ServeContent(c.Writer, c.Request, file.Filename, file.CreatedAt, reader)
}

func setDownloadUrl(files []dbTypes.File) []FileDTO {
	var ret []FileDTO
	for _, file := range files {
		ret = append(ret, FileDTO{
			Name:        file.Name,
			Description: file.Description,
			Category:    file.Category.Name,
			URL:         fmt.Sprintf("/v1/cdn/%d", file.ID),
			CreatedAt:   file.CreatedAt,
			UpdatedAt:   file.UpdatedAt,
		})
	}

	return ret
}

func filterFiles(files []dbTypes.File, filterArray []string) []dbTypes.File {
	resultFiles := []dbTypes.File{}
	for _, filter := range filterArray {
		for _, file := range files {
			if strings.EqualFold(file.Category.Name, filter) {
				resultFiles = append(resultFiles, file)
				break
			}
		}
	}
	return resultFiles
}

func categoryExists(name string) (dbTypes.Category, bool) {
	category := dbTypes.Category{}
	if err := database.DB.Where("name = ?", name).First(&category).Error; err != nil {
		log.Error("Error looking up category %s: %s", name, err.Error())
		return dbTypes.Category{}, false
	}
	return category, true
}

func PostCDN(c *gin.Context) {
	var fileDetails CDN
	if err := c.ShouldBind(&fileDetails); err != nil {
		log.Error("Error binding: %s", err.Error())
		response.RespondMessage(c, http.StatusBadRequest, "Bad Request")
		return
	}

	category, exists := categoryExists(fileDetails.Category)
	if !exists {
		log.Error("Error fetching category %s", fileDetails.Category)
		response.RespondMessage(c, http.StatusBadRequest, "Error fetching category")
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		log.Error("Error fetching file: %s", err.Error())
		response.RespondMessage(c, http.StatusBadRequest, "Error fetching file")
		return
	}

	filename := filepath.Base(file.Filename)
	ext := filepath.Ext(filename)
	filename, _ = gonanoid.New()
	if utils.Getenv("APP_ENV", "prod") != "prod" {
		filename = "dev_" + filename
	}
	filename = filename + ext

	dbFile := dbTypes.File{
		Name:        fileDetails.Name,
		Description: fileDetails.Description,
		Category:    category,
		CategoryID:  category.ID,
		Bucket:      utils.Getenv("AWS_BUCKET", "vzau"),
		Key:         "uploads/" + filename,
		ContentType: storage.GetContentType("/tmp/" + filename),
		Size:        file.Size,
		Filename:    filepath.Base(file.Filename),
	}

	if err := c.SaveUploadedFile(file, filepath.Join("/tmp", filename)); err != nil {
		log.Error("Error saving file: %s", err.Error())
		response.RespondMessage(c, http.StatusInternalServerError, "Error saving file")
		return
	}

	err = storage.UploadFile(dbFile.Bucket, dbFile.Key, "/tmp/"+filename, dbFile.ContentType)
	if err != nil {
		log.Error("Error uploading file: %s", err.Error())
		response.RespondMessage(c, http.StatusInternalServerError, "Error uploading file")
		return
	}

	if err := database.DB.Create(&dbFile).Error; err != nil {
		log.Error("Error creating file: %s", err.Error())
		response.RespondMessage(c, http.StatusInternalServerError, "Error creating file")
		return
	}

	response.Respond(c, http.StatusCreated, struct {
		File dbTypes.File `json:"file"`
	}{dbFile})
}

func handleDelete(bucket string, key string) {
	err := storage.DeleteFile(bucket, key)
	if err != nil {
		log.Error("Error deleting file: %s", err.Error())
		discord.SendToDiscordf(utils.Getenv("DISCORD_WEBHOOK_AUDIT", ""), "Failed to delete %s/%s, delete manually!", bucket, key)
		return
	}
}

func DeleteCDN(c *gin.Context) {
	id := c.Param("id")
	file := dbTypes.File{}
	if err := database.DB.Where("id = ?", id).First(&file).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			response.RespondMessage(c, http.StatusNotFound, "File not found")
			return
		}
		log.Error("Error fetching file: %s", err.Error())
		response.RespondMessage(c, http.StatusInternalServerError, "Error fetching file")
		return
	}

	go handleDelete(file.Bucket, file.Key)
	cache.Cache.Delete(fmt.Sprintf("cdn/%d", file.ID))

	if err := database.DB.Delete(&file).Error; err != nil {
		log.Error("Error deleting file: %s", err.Error())
		response.RespondMessage(c, http.StatusInternalServerError, "Error deleting file")
		return
	}

	response.RespondBlank(c, http.StatusNoContent)
}

func PatchCDNUpdate(c *gin.Context) {
	id := c.Param("id")
	data := CDN{}
	if err := c.ShouldBind(&data); err != nil {
		log.Error("Error binding: %s", err.Error())
		response.RespondMessage(c, http.StatusBadRequest, "Bad Request")
		return
	}

	file := dbTypes.File{}
	if err := database.DB.Where("id = ?", id).Preload(clause.Associations).First(&file).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			response.RespondMessage(c, http.StatusNotFound, "File not found")
			return
		}
		log.Error("Error fetching file: %s", err.Error())
		response.RespondMessage(c, http.StatusInternalServerError, "Error fetching file")
		return
	}

	if data.Name != "" {
		file.Name = data.Name
	}

	if data.Description != "" {
		file.Description = data.Description
	}

	if data.Category != "" {
		if data.Category != file.Category.Name {
			category, exists := categoryExists(data.Category)
			if !exists {
				log.Error("Error fetching category %s", data.Category)
				response.RespondMessage(c, http.StatusNotFound, "Category not found")
				return
			}
			file.Category = category
			file.CategoryID = category.ID
		}
	}

	if err := database.DB.Save(&file).Error; err != nil {
		log.Error("Error saving file: %s", err.Error())
		response.RespondMessage(c, http.StatusInternalServerError, "Error saving file")
		return
	}

	response.Respond(c, http.StatusOK, struct {
		File dbTypes.File `json:"file"`
	}{File: file})
}

func PostCDNUpdate(c *gin.Context) {
	id := c.Param("id")
	uploadedFile, err := c.FormFile("file")
	if err != nil {
		log.Error("Error fetching file: %s", err.Error())
		response.RespondMessage(c, http.StatusBadRequest, "Error fetching file")
		return
	}

	file := dbTypes.File{}
	if err := database.DB.Where("id = ?", id).Preload(clause.Associations).Find(&file).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			response.RespondMessage(c, http.StatusNotFound, "File not found")
			return
		}
		log.Error("Error fetching file: %s", err.Error())
		response.RespondMessage(c, http.StatusInternalServerError, "Error fetching file")
		return
	}

	tmpfile, _ := gonanoid.New()
	if err := database.DB.Save(&file).Error; err != nil {
		log.Error("Error saving file: %s", err.Error())
		response.RespondMessage(c, http.StatusInternalServerError, "Error saving file")
		return
	}

	if err := c.SaveUploadedFile(uploadedFile, filepath.Join("/tmp", tmpfile)); err != nil {
		log.Error("Error saving file: %s", err.Error())
		response.RespondMessage(c, http.StatusInternalServerError, "Error saving file")
		return
	}

	file.ContentType = storage.GetContentType("/tmp/" + tmpfile)
	file.Size = uploadedFile.Size
	file.Filename = filepath.Base(uploadedFile.Filename)

	if err := database.DB.Save(&file).Error; err != nil {
		log.Error("Error saving file: %s", err.Error())
		response.RespondMessage(c, http.StatusInternalServerError, "Error saving file")
		return
	}

	err = storage.UploadFile(utils.Getenv("AWS_BUCKET", "vzau"), file.Key, "/tmp/"+tmpfile, storage.GetContentType("/tmp/"+tmpfile))
	if err != nil {
		log.Error("Error uploading file: %s", err.Error())
		response.RespondMessage(c, http.StatusInternalServerError, "Error uploading file")
		return
	}

	cache.Cache.Delete(fmt.Sprintf("cdn/%d", file.ID))

	response.RespondBlank(c, http.StatusNoContent)
}
