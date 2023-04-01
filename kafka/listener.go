package kafka

import (
	"context"
	"log"
	"sync"

	"github.com/giornetta/microshop/events"
	"github.com/twmb/franz-go/pkg/kgo"
)

// Listener fetches incoming Kafka messages, offering a simple API to specify handlers for them.
type Listener struct {
	client *kgo.Client

	wg sync.WaitGroup

	lock     sync.RWMutex
	handlers map[events.Topic]events.Handler
}

// NewListener returns a new Kafka Listener.
func NewListener(client *kgo.Client) (*Listener, error) {
	l := &Listener{
		client:   client,
		handlers: make(map[events.Topic]events.Handler),
	}

	return l, nil
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
			l.wg.Wait()
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

			l.wg.Add(1)
			go l.handleEvent(event)
		}
	}

}

func (l *Listener) handleEvent(e events.Event) {
	l.lock.RLock()
	defer l.lock.RUnlock()

	handler, ok := l.handlers[e.Topic()]
	if !ok {
		return
	}

	l.wg.Add(1)
	go func(h events.Handler) {
		if err := h.Handle(e); err != nil {
			// TODO Do some error Handling
			log.Println(err)
		}
		l.wg.Done()
	}(handler)

	l.wg.Done()
}
