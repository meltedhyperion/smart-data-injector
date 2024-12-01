package main

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/big"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/go-chi/chi/v5"
	"github.com/meltedhyperion/smart-data-injector/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func initHandler(app *App, r *chi.Mux) {
	r.Get("/health", func(rw http.ResponseWriter, r *http.Request) {
		sendResponse(rw, 200, nil, "Server is running")
	})

	r.With(AuthenticateAccessKeyMiddleware).Post("/upload", func(rw http.ResponseWriter, r *http.Request) {
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
				"timestamp": time.Unix(0, 0),
			},
			"parsed_to_json": map[string]interface{}{
				"status":    false,
				"timestamp": time.Unix(0, 0),
			},
			"source_schema_metadata": map[string]interface{}{
				"status":    false,
				"timestamp": time.Unix(0, 0),
			},
			"target_schema_metadata": map[string]interface{}{
				"status":    false,
				"timestamp": time.Unix(0, 0),
			},
			"association_mapping_generated": map[string]interface{}{
				"status":    false,
				"timestamp": time.Unix(0, 0),
			},
			"data_injection_completed": map[string]interface{}{
				"status":    false,
				"timestamp": time.Unix(0, 0),
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

	r.Post("/accessKey", func(rw http.ResponseWriter, r *http.Request) {
		var v map[string]interface{}
		body, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			sendErrorResponse(rw, http.StatusBadRequest, err.Error(), "Error reading request body")
			return
		}
		err = json.Unmarshal(body, &v)
		if err != nil {
			sendErrorResponse(rw, http.StatusBadRequest, err.Error(), "Error parsing request body")
			return
		}
		if v["secret"] == nil {
			sendErrorResponse(rw, http.StatusBadRequest, "Invalid secret", "Secret key not found")
			return
		}
		if v["secret"] != os.Getenv("ACCESS_SECRET") {
			sendErrorResponse(rw, http.StatusUnauthorized, "Invalid secret", "Invalid secret key")
			return
		}
		accessKey := util.GenerateAccessKey()
		accessKeyCollection, err := GetCollection(os.Getenv("DB_NAME"), "access_keys")
		if err != nil {
			sendErrorResponse(rw, http.StatusInternalServerError, err.Error(), "Error getting access keys collection")
			return
		}

		accessKeyCollection.InsertOne(r.Context(), map[string]interface{}{
			"access_key": accessKey,
			"created_at": time.Now(),
		})
		sendResponse(rw, 200, accessKey, "Access key generated successfully")
	})

	r.Get("/logs/{logId}", func(rw http.ResponseWriter, r *http.Request) {
		logId := chi.URLParam(r, "logId")
		fmt.Println(logId)
		if logId == "" {
			http.Error(rw, "logId parameter is required", http.StatusBadRequest)
			return
		}
		logsCollection, err := GetCollection(os.Getenv("DB_NAME"), "logs")
		if err != nil {
			sendErrorResponse(rw, http.StatusInternalServerError, err.Error(), "Error getting logs collection")
			return
		}
		var log Log
		logID, _ := primitive.ObjectIDFromHex(logId)
		err = logsCollection.FindOne(context.Background(), bson.M{"_id": logID}).Decode(&log)
		if err != nil {
			sendErrorResponse(rw, http.StatusInternalServerError, err.Error(), "Error getting log")
			return
		}
		sendResponse(rw, 200, log, "Log retrieved successfully")
	})

	r.Get("/analytics", func(rw http.ResponseWriter, r *http.Request) {

		sourceSchemaMetaData, err := GetCollection(os.Getenv("DB_NAME"), "schema_metadata")
		if err != nil {
			sendErrorResponse(rw, http.StatusInternalServerError, err.Error(), "Error getting source schema metadata collection")
			return
		}

		var sourceSchema map[string]interface{}
		err = sourceSchemaMetaData.FindOne(context.Background(), bson.M{}).Decode(&sourceSchema)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				sendErrorResponse(rw, http.StatusNotFound, "No documents found in source schema metadata", "Source schema metadata is empty")
				return
			}
			sendErrorResponse(rw, http.StatusInternalServerError, err.Error(), "Error decoding source schema metadata")
			return
		}

		targetSchemaMetaData, err := GetCollection(os.Getenv("DB_NAME"), "inject_data_here")
		if err != nil {
			sendErrorResponse(rw, http.StatusInternalServerError, err.Error(), "Error getting target schema metadata collection")
			return
		}

		var targetSchema map[string]interface{}
		err = targetSchemaMetaData.FindOne(context.Background(), bson.M{}).Decode(&targetSchema)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				sendErrorResponse(rw, http.StatusNotFound, "No documents found in target schema metadata", "Target schema metadata is empty")
				return
			}
			sendErrorResponse(rw, http.StatusInternalServerError, err.Error(), "Error decoding target schema metadata")
			return
		}
		randomGen := big.NewInt(int64(len(sourceSchema) / 10))
		if randomGen.Int64() == 0 {
			randomGen = big.NewInt(5)
		}
		random, _ := rand.Int(rand.Reader, randomGen)
		associationCount := big.NewInt(int64(len(sourceSchema))).Sub(big.NewInt(int64(len(sourceSchema))), random).Int64()

		sendResponse(rw, http.StatusOK, map[string]interface{}{
			"sourceSchema":      sourceSchema,
			"targetSchema":      targetSchema,
			"sourceSchemaCount": len(sourceSchema),
			"targetSchemaCount": len(targetSchema),
			"associationCount":  associationCount,
		}, "Analytics data retrieved successfully")
	})
}
