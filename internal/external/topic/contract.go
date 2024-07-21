package topic

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/sns"
)

type Service struct {
	client *sns.Client

	TopicName string
	TopicArn  string
}

type TopicService interface {
	UpdateTopicArn(ctx context.Context) error
	Publish(ctx context.Context, data *Message) (*string, error)
}

type EventType string

type Message struct {
	Data interface{} `json:"data"`
}

func NewMessage(data interface{}) *Message {
	return &Message{
		Data: data,
	}
}
