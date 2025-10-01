package ports

import "context"

// EventPublisher publishes domain events to the configured transport.
type EventPublisher interface {
	Publish(ctx context.Context, topic string, payload any) error
}
