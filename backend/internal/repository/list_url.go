package repository

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"url-shortener/internal/domain"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func (r *URLRepository) List(ctx context.Context, limit int32, cursor string) (urls []domain.URL, nextCursor string, err error) {
	input := &dynamodb.ScanInput{
		TableName: aws.String(r.tableName),
		Limit:     aws.Int32(limit),
	}

	if cursor != "" {
		startKey, err := decodeCursor(cursor)
		if err != nil {
			return nil, "", err
		}
		input.ExclusiveStartKey = startKey
	}

	out, err := r.client.Scan(ctx, input)
	if err != nil {
		return nil, "", fmt.Errorf("dynamodb scan: %w", err)
	}

	if err := attributevalue.UnmarshalListOfMaps(out.Items, &urls); err != nil {
		return nil, "", fmt.Errorf("unmarshal records: %w", err)
	}

	if len(out.LastEvaluatedKey) > 0 {
		nextCursor, err = encodeCursor(out.LastEvaluatedKey)
		if err != nil {
			return nil, "", fmt.Errorf("encode cursor: %w", err)
		}
	}

	return urls, nextCursor, nil
}

func encodeCursor(key map[string]types.AttributeValue) (string, error) {
	var rawMap map[string]interface{}
	if err := attributevalue.UnmarshalMap(key, rawMap); err != nil {
		return "", err
	}

	data, err := json.Marshal(rawMap)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(data), err
}

func decodeCursor(cursor string) (map[string]types.AttributeValue, error) {
	data, err := base64.URLEncoding.DecodeString(cursor)
	if err != nil {
		return nil, err
	}

	var rawMap map[string]interface{}
	if err := json.Unmarshal(data, &rawMap); err != nil {
		return nil, err
	}

	return attributevalue.MarshalMap(rawMap)
}
