package handler

import (
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func (s *Server) List() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := s.db.GetItem(r.Context(), &dynamodb.GetItemInput{
			TableName: aws.String(s.tableName),
		})

		if err != nil {
			w.WriteHeader(http.StatusTeapot)
		}

		w.WriteHeader(http.StatusOK)
	})
}
