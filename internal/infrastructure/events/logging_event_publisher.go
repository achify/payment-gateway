package events

import (
	"context"
	"log"
)

// LoggingEventPublisher logs events instead of delivering them to Kafka.
type LoggingEventPublisher struct{}

// NewLoggingEventPublisher creates the publisher.
func NewLoggingEventPublisher() *LoggingEventPublisher {
	return &LoggingEventPublisher{}
}

// Publish logs the topic and payload.
func (p *LoggingEventPublisher) Publish(_ context.Context, topic string, payload any) error {
	log.Printf("event topic=%s payload=%+v", topic, payload)
	return nil
}
