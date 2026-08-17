package handler

import (
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func List(db *dynamodb.Client, tableName string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := db.GetItem(r.Context(), &dynamodb.GetItemInput{
			TableName: aws.String(tableName),
		})

		if err != nil {
			w.WriteHeader(http.StatusTeapot)
		}

		w.WriteHeader(http.StatusOK)
	})
}
