package handler

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type Server struct {
	db        *dynamodb.Client
	tableName string
}

func NewServer(db *dynamodb.Client, tableName string) *Server {
	return &Server{db: db, tableName: tableName}
}
