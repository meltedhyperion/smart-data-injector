package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/go-chi/chi/v5"
	"github.com/meltedhyperion/smart-data-injector/util"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func initHandler(app *App, r *chi.Mux) {
	r.Get("/health", func(rw http.ResponseWriter, r *http.Request) {
		sendResponse(rw, 200, nil, "Server is running")
	})

	r.Post("/upload", func(rw http.ResponseWriter, r *http.Request) {
		const maxFileSize int64 = 1024 * 1024 * 10

		r.Body = http.MaxBytesReader(rw, r.Body, maxFileSize)

		err := r.ParseMultipartForm(maxFileSize)
		if err != nil {
			sendErrorResponse(rw, http.StatusBadRequest, err.Error(), "File not found in the request")
			return
		}

		file, handler, err := r.FormFile("file")
		if err != nil {
			sendErrorResponse(rw, http.StatusBadRequest, err.Error(), "Error retrieving the file")
			return
		}
		defer file.Close()

		if !util.IsValidFileType(handler.Filename) {
			sendErrorResponse(rw, http.StatusBadRequest, "Invalid file type", "Only .json or .csv files are allowed")
			return
		}

		filePath := filepath.Join("./uploads", handler.Filename)
		out, err := os.Create(filePath)
		if err != nil {
			sendErrorResponse(rw, http.StatusInternalServerError, err.Error(), "Error creating the file")
			return
		}
		defer out.Close()

		_, err = io.Copy(out, file)
		if err != nil {
			sendErrorResponse(rw, http.StatusInternalServerError, err.Error(), "Error saving the file")
			return
		}
		svc := s3.New(Sess)
		srcFile, err := handler.Open()
		if err != nil {
			sendErrorResponse(rw, http.StatusInternalServerError, err.Error(), "Error opening file")
		}
		defer srcFile.Close()
		uploadFileName := util.GenerateUniqueFileName(handler.Filename)
		key := "data/" + uploadFileName

		params := &s3.PutObjectInput{
			Bucket: aws.String(os.Getenv("AWS_S3_BUCKET")),
			Key:    aws.String(key),
			Body:   srcFile,
		}

		_, err = svc.PutObject(params)
		log.Print(err)
		if err != nil {
			sendErrorResponse(rw, http.StatusInternalServerError, err.Error(), "Error uploading file")
		}
		logsCollection, err := GetCollection(os.Getenv("DB_NAME"), "logs")
		if err != nil {
			sendErrorResponse(rw, http.StatusInternalServerError, err.Error(), "Error getting logs collection")
		}
		insertID, err := logsCollection.InsertOne(r.Context(), map[string]interface{}{
			"file_name": uploadFileName,
			"file_upload": map[string]interface{}{
				"status":    true,
				"timestamp": time.Now(),
			},
			"s3_trigger_completed": map[string]interface{}{
				"status":    false,
				"timestamp": 0,
			},
			"parsed_to_json": map[string]interface{}{
				"status":    false,
				"timestamp": 0,
			},
			"source_schema_metadata": map[string]interface{}{
				"status":    false,
				"timestamp": 0,
			},
			"target_schema_metadata": map[string]interface{}{
				"status":    false,
				"timestamp": 0,
			},
			"association_mapping_generated": map[string]interface{}{
				"status":    false,
				"timestamp": 0,
			},
			"data_injection_completed": map[string]interface{}{
				"status":    false,
				"timestamp": 0,
			},
		})

		if err != nil {
			sendErrorResponse(rw, http.StatusInternalServerError, err.Error(), "Error saving log")
		}

		err = os.Remove(filePath)
		if err != nil {
			log.Printf("Warning: Failed to delete file %s: %v", filePath, err)
		} else {
			log.Printf("File %s successfully deleted from server", filePath)
		}

		url := "https://smart-data-injector.s3.ap-south-1.amazonaws.com/" + key
		logId := insertID.InsertedID.(primitive.ObjectID).Hex()
		data := map[string]interface{}{
			"url":   url,
			"logId": logId,
		}
		sendResponse(rw, 200, data, "File uploaded successfully with public read access.")
	})
}
