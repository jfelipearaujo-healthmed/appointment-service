package topic

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

func NewService(topicName string, config aws.Config) TopicService {
	return &Service{
		client:    sns.NewFromConfig(config),
		TopicName: topicName,
	}
}

func (svc *Service) UpdateTopicArn(ctx context.Context) error {
	output, err := svc.client.ListTopics(ctx, &sns.ListTopicsInput{})
	if err != nil {
		return err
	}

	for _, topic := range output.Topics {
		if strings.Contains(*topic.TopicArn, svc.TopicName) {
			svc.TopicArn = *topic.TopicArn
			return nil
		}
	}

	return errors.New("topic not found")
}

func (svc *Service) Publish(ctx context.Context, data *Message) (*string, error) {
	body, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	req := &sns.PublishInput{
		TopicArn: aws.String(svc.TopicArn),
		Message:  aws.String(string(body)),
	}

	out, err := svc.client.Publish(ctx, req)
	if err != nil {
		return nil, err
	}

	slog.InfoContext(ctx, "message published", "topic", svc.TopicName, "message_id", *out.MessageId, "message", string(body))

	return out.MessageId, nil
}
