package repository

import (
	"context"
	"errors"
	"fmt"
	"url-shortener/internal/domain"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type URLRepository struct {
	client    *dynamodb.Client
	tableName string
}

func NewURLRepository(client *dynamodb.Client, tableName string) *URLRepository {
	return &URLRepository{client: client, tableName: tableName}
}

func (r *URLRepository) Get(code string) (*domain.URL, error) {
	return nil, nil
}

func (r *URLRepository) Find(ctx context.Context, code string) (*domain.URL, error) {
	id := "Code#" + code

	out, err := r.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
		ConsistentRead: aws.Bool(true),
	})

	if err != nil {
		return nil, fmt.Errorf("dynamodb get item: %w", err)
	}

	if len(out.Item) == 0 {
		return nil, domain.ErrNotFound
	}

	var record domain.URL
	if err := attributevalue.UnmarshalMap(out.Item, &record); err != nil {
		return nil, fmt.Errorf("unmarshal record: %w", err)
	}

	return &record, nil
}

func (r *URLRepository) Save(ctx context.Context, url *domain.URL) error {
	item, err := attributevalue.MarshalMap(url)
	if err != nil {
		return fmt.Errorf("marshal record: %w", err)
	}

	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName:           aws.String(r.tableName),
		Item:                item,
		ConditionExpression: aws.String("attribute_not_exists(pk)"),
	})

	if err != nil {
		var cfe *types.ConditionalCheckFailedException
		if errors.As(err, &cfe) {
			return domain.ErrAlreadyExists
		}
		return fmt.Errorf("dynamodb put item: %w", err)
	}

	return nil
}
