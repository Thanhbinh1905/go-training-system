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

	return &KafkaProducer{writer: writer}, nil
}

// publish is an internal helper to send events to Kafka
func (p *KafkaProducer) publish(ctx context.Context, topic, key string, event any) error {
	eventBytes, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	err = p.writer.WriteMessages(ctx, kafka.Message{
		Topic: topic,
		Key:   []byte(key),
		Value: eventBytes,
	})
	if err != nil {
		return fmt.Errorf("failed to publish event to topic %s: %w", topic, err)
	}

	return nil
}

// PublishTeamEvent publishes team events to team.activity topic
func (p *KafkaProducer) PublishTeamEvent(ctx context.Context, event *TeamEvent) error {
	if err := p.publish(ctx, TeamActivityTopic, event.TeamID, event); err != nil {
		return err
	}
	log.Printf("[Kafka] Published team event: %s for team %s", event.EventType, event.TeamID)
	return nil
}

// PublishAssetEvent publishes asset events to asset.changes topic
func (p *KafkaProducer) PublishAssetEvent(ctx context.Context, event *AssetEvent) error {
	if err := p.publish(ctx, AssetChangesTopic, event.AssetID, event); err != nil {
		return err
	}
	log.Printf("[Kafka] Published asset event: %s for %s %s", event.EventType, event.AssetType, event.AssetID)
	return nil
}

// PublishAssetShareEvent publishes asset share events to asset.shares topic
func (p *KafkaProducer) PublishAssetShareEvent(ctx context.Context, event *AssetShareEvent) error {
	if err := p.publish(ctx, AssetChangesTopic, event.AssetID, event); err != nil {
		return err
	}
	log.Printf("[Kafka] Published asset share event: %s for %s %s", event.EventType, event.AssetType, event.AssetID)
	return nil
}

// Close closes the Kafka producer
func (p *KafkaProducer) Close() error {
	return p.writer.Close()
}
