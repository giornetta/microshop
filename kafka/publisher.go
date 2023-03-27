package kafka

import (
	"context"
	"encoding/json"
	"time"

	"github.com/giornetta/microshop/events"
	"github.com/twmb/franz-go/pkg/kgo"
)

type eventPublisher struct {
	client *kgo.Client
}

func NewEventPublisher(client *kgo.Client) events.Publisher {
	return &eventPublisher{
		client: client,
	}
}

func (p *eventPublisher) Publish(e events.Event, ctx context.Context) error {
	jsonPayload, _ := json.Marshal(e)

	record := &kgo.Record{
		Key:   []byte(e.Key()),
		Value: jsonPayload,
		Headers: []kgo.RecordHeader{
			{
				Key:   "EventType",
				Value: []byte(e.Type()),
			},
		},
		Timestamp: time.Now(),
		Topic:     e.Topic().String(),
	}

	if err := p.client.ProduceSync(ctx, record).FirstErr(); err != nil {
		return err
	}

	return nil
}
