package domain

import (
	"time"
)

type URL struct {
	PK            string    `dynamodbav:"pk"`
	Code          string    `dynamodbav:"code"`
	TargetURL     string    `dynamodbav:"target_url"`
	CreateAt      time.Time `dynamodbav:"create_at"`
	ExpiresAt     time.Time `dynamodbav:"expires_at"`
	Clicks        int64     `dynamodbav:"clicks"`
	LastClickedAt int64     `dynamodbav:"last_clicked_at,omitempty"`
}

func NewURL(code string, targetURL string, expiresInDays int) *URL {
	now := time.Now()
	return &URL{
		PK:        code,
		Code:      code,
		TargetURL: targetURL,
		CreateAt:  now,
		ExpiresAt: now.AddDate(0, 0, expiresInDays),
	}
}
