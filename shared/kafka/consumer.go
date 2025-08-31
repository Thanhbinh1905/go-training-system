package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

// Consumer interface for Kafka operations
type Consumer interface {
	ConsumeTeamEvents(ctx context.Context, handler TeamEventHandler) error
	ConsumeAssetEvents(ctx context.Context, handler AssetEventHandler) error
	Close() error
}

// TeamEventHandler handles team events
type TeamEventHandler func(event *TeamEvent) error

// AssetEventHandler handles asset events
type AssetEventHandler func(event *AssetEvent) error

type AssetShareEventHandler func(event *AssetShareEvent) error

// KafkaConsumer implements Consumer interface
type KafkaConsumer struct {
	teamReader  *kafka.Reader
	assetReader *kafka.Reader
}

// NewConsumer creates a new Kafka consumer
func NewConsumer(brokers []string, groupID string) (*KafkaConsumer, error) {
	teamReader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		Topic:    TeamActivityTopic,
		GroupID:  groupID + "-team",
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
		MaxWait:  1 * time.Second,
	})

	assetReader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		Topic:    AssetChangesTopic,
		GroupID:  groupID + "-asset",
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
		MaxWait:  1 * time.Second,
	})

	return &KafkaConsumer{
		teamReader:  teamReader,
		assetReader: assetReader,
	}, nil
}

// ConsumeTeamEvents consumes team events from team.activity topic
func (c *KafkaConsumer) ConsumeTeamEvents(ctx context.Context, handler TeamEventHandler) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			m, err := c.teamReader.ReadMessage(ctx)
			if err != nil {
				log.Printf("Error reading team event: %v", err)
				continue
			}

			var event TeamEvent
			if err := json.Unmarshal(m.Value, &event); err != nil {
				log.Printf("Error unmarshaling team event: %v", err)
				continue
			}

			if err := handler(&event); err != nil {
				log.Printf("Error handling team event: %v", err)
				// Consider implementing retry logic here
				continue
			}

			log.Printf("Processed team event: %s for team %s", event.EventType, event.TeamID)
		}
	}
}

// ConsumeAssetEvents consumes asset events from asset.changes topic
func (c *KafkaConsumer) ConsumeAssetEvents(ctx context.Context, handler AssetEventHandler, shareHandler AssetShareEventHandler) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			m, err := c.assetReader.ReadMessage(ctx)
			if err != nil {
				log.Printf("Error reading asset event: %v", err)
				continue
			}

			// Unmarshal to read event type
			var base BaseEvent
			if err := json.Unmarshal(m.Value, &base); err != nil {
				log.Printf("Error unmarshaling base event: %v", err)
				continue
			}

			switch base.EventType {
			case AssetEventFolderShared, AssetEventFolderUnshared, AssetEventNoteShared, AssetEventNoteUnshared:
				var shareEvent AssetShareEvent
				if err := json.Unmarshal(m.Value, &shareEvent); err != nil {
					log.Printf("Error unmarshaling asset share event: %v", err)
					continue
				}
				if err := shareHandler(&shareEvent); err != nil {
					log.Printf("Error handling asset share event: %v", err)
					continue
				}
				log.Printf("Processed asset share event: %+v", shareEvent)

			default:
				var event AssetEvent
				if err := json.Unmarshal(m.Value, &event); err != nil {
					log.Printf("Error unmarshaling asset event: %v", err)
					continue
				}
				if err := handler(&event); err != nil {
					log.Printf("Error handling asset event: %v", err)
					continue
				}
				log.Printf("Processed asset event: %+v", event)
			}
		}
	}
}

// Close closes the Kafka consumer
func (c *KafkaConsumer) Close() error {
	var errs []error

	if err := c.teamReader.Close(); err != nil {
		errs = append(errs, fmt.Errorf("failed to close team reader: %w", err))
	}

	if err := c.assetReader.Close(); err != nil {
		errs = append(errs, fmt.Errorf("failed to close asset reader: %w", err))
	}

	if len(errs) > 0 {
		return fmt.Errorf("errors closing consumer: %v", errs)
	}

	return nil
}
