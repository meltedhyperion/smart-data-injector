package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"go.mongodb.org/mongo-driver/bson"
)

func AuthenticateAccessKeyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		accessToken := r.Header.Get("Authorization")
		if accessToken == "" {
			sendErrorResponse(rw, http.StatusBadRequest, "No access token provided", "No access token provided")
			return
		}

		accessKeyCollection, err := GetCollection(os.Getenv("DB_NAME"), "access_keys")
		if err != nil {
			sendErrorResponse(rw, http.StatusInternalServerError, err.Error(), "Error getting access keys collection")
			return
		}

		var body map[string]interface{}
		err = accessKeyCollection.FindOne(context.Background(), bson.M{"access_key": accessToken}).Decode(&body)
		fmt.Println(err)
		if err != nil {
			sendErrorResponse(rw, http.StatusUnauthorized, "Invalid access token", "Invalid access token")
			return
		}

		next.ServeHTTP(rw, r)
	})
}
