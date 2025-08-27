package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/segmentio/kafka-go"
)

// Producer interface for Kafka operations
type Producer interface {
	PublishTeamEvent(ctx context.Context, event *TeamEvent) error
	PublishAssetEvent(ctx context.Context, event *AssetEvent) error
	PublishAssetShareEvent(ctx context.Context, event *AssetShareEvent) error
	Close() error
}

// KafkaProducer implements Producer interface
type KafkaProducer struct {
	writer *kafka.Writer
}

// NewProducer creates a new Kafka producer
func NewProducer(brokers []string) (*KafkaProducer, error) {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Balancer: &kafka.LeastBytes{},
	}

	return &KafkaProducer{
		writer: writer,
	}, nil
}

// PublishTeamEvent publishes team events to team.activity topic
func (p *KafkaProducer) PublishTeamEvent(ctx context.Context, event *TeamEvent) error {
	eventBytes, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal team event: %w", err)
	}

	err = p.writer.WriteMessages(ctx, kafka.Message{
		Topic: TeamActivityTopic,
		Key:   []byte(event.TeamID),
		Value: eventBytes,
	})

	if err != nil {
		return fmt.Errorf("failed to publish team event: %w", err)
	}

	log.Printf("Published team event: %s for team %s", event.EventType, event.TeamID)
	return nil
}

// PublishAssetEvent publishes asset events to asset.changes topic
func (p *KafkaProducer) PublishAssetEvent(ctx context.Context, event *AssetEvent) error {
	eventBytes, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal asset event: %w", err)
	}

	err = p.writer.WriteMessages(ctx, kafka.Message{
		Topic: AssetChangesTopic,
		Key:   []byte(event.AssetID),
		Value: eventBytes,
	})

	if err != nil {
		return fmt.Errorf("failed to publish asset event: %w", err)
	}

	log.Printf("Published asset event: %s for %s %s", event.EventType, event.AssetType, event.AssetID)
	return nil
}

// PublishAssetEvent publishes asset events to asset.changes topic
func (p *KafkaProducer) PublishAssetShareEvent(ctx context.Context, event *AssetShareEvent) error {
	eventBytes, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal asset event: %w", err)
	}

	err = p.writer.WriteMessages(ctx, kafka.Message{
		Topic: AssetChangesTopic,
		Key:   []byte(event.AssetID),
		Value: eventBytes,
	})

	if err != nil {
		return fmt.Errorf("failed to publish asset event: %w", err)
	}

	log.Printf("Published asset event: %s for %s %s", event.EventType, event.AssetType, event.AssetID)
	return nil
}

// Close closes the Kafka producer
func (p *KafkaProducer) Close() error {
	return p.writer.Close()
}
