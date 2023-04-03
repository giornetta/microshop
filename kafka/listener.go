package kafka

import (
	"context"
	"errors"
	"sync"

	"github.com/giornetta/microshop/events"
	"github.com/twmb/franz-go/pkg/kgo"
)

// Listener fetches incoming Kafka messages, offering a simple API to specify handlers for them.
type Listener struct {
	client *kgo.Client

	lock     sync.RWMutex
	handlers map[events.Topic]events.Handler
}

// NewListener returns a new Kafka Listener.
func NewListener(client *kgo.Client) *Listener {
	l := &Listener{
		client:   client,
		handlers: make(map[events.Topic]events.Handler),
	}

	return l
}

// Handle registers the given handler for the provided topic.
// Calling this method a second time on the same topic will replace the handler.
func (l *Listener) Handle(topic events.Topic, handler events.Handler) {
	l.lock.Lock()
	defer l.lock.Unlock()

	_, ok := l.handlers[topic]
	if !ok {
		l.client.AddConsumeTopics(topic.String())
	}

	l.handlers[topic] = handler
}

// Listener is a blocking method that will fetch incoming kafka messages
// until either an error occurs or the given context is canceled.
func (l *Listener) Listen(ctx context.Context) error {
	for {
		fetches := l.client.PollFetches(ctx)
		if errs := fetches.Errors(); errs != nil {
			if errs[0].Err == context.Canceled {
				return nil
			}

			return errs[0].Err
		}

		iter := fetches.RecordIter()
		for !iter.Done() {
			record := iter.Next()

			t := events.Type(record.Headers[0].Value)
			event, _ := events.Decode(t, record.Value)

			// Events are handled sequentially, in a blocking manner, to ensure ordering.
			if err := l.handleEvent(event, ctx); err != nil {
				return err
			}

			l.client.CommitRecords(ctx, record)
		}
	}
}

func (l *Listener) handleEvent(evt events.Event, ctx context.Context) error {
	l.lock.RLock()
	defer l.lock.RUnlock()

	handler, ok := l.handlers[evt.Topic()]
	if !ok {
		return errors.New("couldn't find handler")
	}

	if err := handler.Handle(evt, ctx); err != nil {
		return err
	}

	return nil
}
