# Kafka Integration Guide

## Overview

This document describes the Kafka integration implemented in the Go Training System v2. The system now uses Kafka for event-driven communication between microservices, enabling real-time updates and asynchronous processing.

## Architecture

### Event Topics

- **`team.activity`**: Team-related events (creation, member changes, manager changes)
- **`asset.changes`**: Asset-related events (folders, notes, sharing)

### Event Types

#### Team Events
- `TEAM_CREATED`: New team created
- `MEMBER_ADDED`: Member added to team
- `MEMBER_REMOVED`: Member removed from team
- `MANAGER_ADDED`: Manager added to team
- `MANAGER_REMOVED`: Manager removed from team

#### Asset Events
- `FOLDER_CREATED`, `FOLDER_UPDATED`, `FOLDER_DELETED`
- `NOTE_CREATED`, `NOTE_UPDATED`, `NOTE_DELETED`
- `FOLDER_SHARED`, `FOLDER_UNSHARED`
- `NOTE_SHARED`, `NOTE_UNSHARED`

## Implementation

### 1. Shared Kafka Package

The Kafka functionality is implemented in the `shared/kafka` package:

- **`events.go`**: Event definitions and constants
- **`producer.go`**: Kafka producer for publishing events
- **`consumer.go`**: Kafka consumer for receiving events
- **`consumer_service.go`**: Business logic handlers for events

### 2. Service Integration

#### Team Service
- Emits events when teams are created, members/managers are added/removed
- Events are published to `team.activity` topic

#### Asset Service
- Emits events when folders/notes are created, updated, deleted, or shared
- Events are published to `asset.changes` topic

### 3. Event Payload Structure

#### Team Event Example
```json
{
  "eventType": "MEMBER_ADDED",
  "teamId": "uuid",
  "performedBy": "userId",
  "targetUserId": "userId",
  "timestamp": "2024-01-01T00:00:00Z"
}
```

#### Asset Event Example
```json
{
  "eventType": "NOTE_UPDATED",
  "assetType": "note",
  "assetId": "uuid",
  "ownerId": "userId",
  "actionBy": "userId",
  "timestamp": "2024-01-01T00:00:00Z"
}
```

## Usage

### Starting the System

1. **Start infrastructure services**:
   ```bash
   docker-compose up -d redis zookeeper kafka
   ```

2. **Create Kafka topics** (if not auto-created):
   ```bash
   task kafka-reset
   ```

3. **Start microservices**:
   ```bash
   docker-compose up -d user-service team-service asset-service
   ```

### Testing Kafka

#### List Topics
```bash
task kafka-topics
```

#### Monitor Team Events
```bash
task kafka-console-consumer -- team.activity
```

#### Monitor Asset Events
```bash
task kafka-console-consumer -- asset.changes
```

#### View Kafka Logs
```bash
task kafka-logs
```

### Manual Event Production

#### Produce Team Event
```bash
task kafka-console-producer -- team.activity
# Then type: {"eventType":"TEAM_CREATED","teamId":"test-123","performedBy":"user-456","timestamp":"2024-01-01T00:00:00Z"}
```

#### Produce Asset Event
```bash
task kafka-console-producer -- asset.changes
# Then type: {"eventType":"NOTE_CREATED","assetType":"note","assetId":"note-123","ownerId":"user-456","actionBy":"user-456","timestamp":"2024-01-01T00:00:00Z"}
```

## Development

### Adding New Event Types

1. **Define event constants** in `shared/kafka/events.go`
2. **Update event structs** if needed
3. **Add event emission** in relevant service methods
4. **Add event handling** in `consumer_service.go`

### Example: Adding User Events

```go
// In events.go
const (
    UserEventCreated = "USER_CREATED"
    UserEventUpdated = "USER_UPDATED"
)

type UserEvent struct {
    BaseEvent
    UserID string `json:"userId"`
}

// In user service
userEvent := kafka.NewUserEvent(kafka.UserEventCreated, userID.String(), userID.String())
s.kafkaProducer.PublishUserEvent(ctx, userEvent)
```

### Error Handling

- **Producer errors**: Logged but don't fail the main operation
- **Consumer errors**: Logged and continue processing other events
- **Retry logic**: Can be implemented in consumer handlers

### Monitoring

- **Prometheus metrics**: Available for Kafka operations
- **Logs**: All events are logged with structured information
- **Health checks**: Kafka connectivity is monitored

## Best Practices

### 1. Event Design
- Keep events small and focused
- Include all necessary context
- Use consistent naming conventions
- Version events when schema changes

### 2. Error Handling
- Don't fail main operations due to Kafka errors
- Implement retry mechanisms for failed events
- Log errors for debugging
- Consider dead letter queues for failed events

### 3. Performance
- Use appropriate batch sizes
- Monitor Kafka lag
- Implement backpressure handling
- Use async processing when possible

### 4. Security
- Validate event payloads
- Implement authentication for producers
- Use secure Kafka configuration in production
- Monitor for suspicious event patterns

## Troubleshooting

### Common Issues

1. **Kafka connection refused**:
   - Check if Kafka container is running
   - Verify network configuration
   - Check Kafka logs: `task kafka-logs`

2. **Events not being consumed**:
   - Verify consumer group configuration
   - Check if topics exist: `task kafka-topics`
   - Monitor consumer logs

3. **High latency**:
   - Check Kafka broker configuration
   - Monitor system resources
   - Adjust batch sizes and timeouts

### Debug Commands

```bash
# Check Kafka status
docker exec -it kafka kafka-topics --bootstrap-server kafka:29092 --describe --topic team.activity

# Check consumer groups
docker exec -it kafka kafka-consumer-groups --bootstrap-server kafka:29092 --list

# Check consumer lag
docker exec -it kafka kafka-consumer-groups --bootstrap-server kafka:29092 --describe --group team-service-team
```

## Future Enhancements

1. **Event Schema Registry**: Implement Avro schema validation
2. **Dead Letter Queues**: Handle failed events gracefully
3. **Event Replay**: Ability to replay events for debugging
4. **Event Sourcing**: Use events as source of truth
5. **CQRS**: Separate read and write models
6. **Event Streaming Analytics**: Real-time analytics on events

## References

- [Kafka Documentation](https://kafka.apache.org/documentation/)
- [segmentio/kafka-go](https://github.com/segmentio/kafka-go)
- [Event-Driven Architecture Patterns](https://martinfowler.com/articles/201701-event-driven.html)
