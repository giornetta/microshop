package kafka

import (
	"context"
	"log"
	"sync"

	"github.com/giornetta/microshop/events"
	"github.com/twmb/franz-go/pkg/kgo"
)

type listener struct {
	client *kgo.Client

	wg sync.WaitGroup

	lock     sync.RWMutex
	handlers map[events.Topic][]events.Handler
}

func NewListener(client *kgo.Client) (*listener, error) {
	l := &listener{
		client:   client,
		handlers: make(map[events.Topic][]events.Handler),
	}

	return l, nil
}

func (l *listener) Handle(topic events.Topic, handler events.Handler) {
	l.lock.Lock()
	defer l.lock.Unlock()

	handlers, ok := l.handlers[topic]
	if !ok {
		handlers = make([]events.Handler, 0, 1)
		l.client.AddConsumeTopics(topic.String())
	}

	// Should check if handler is already in the list but w/e
	l.handlers[topic] = append(handlers, handler)
}

func (l *listener) Listen(ctx context.Context) error {
	for {
		fetches := l.client.PollFetches(ctx)
		if errs := fetches.Errors(); errs != nil {
			// TODO Improve error handling
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

func (l *listener) handleEvent(e events.Event) {
	l.lock.RLock()
	defer l.lock.RUnlock()

	for _, h := range l.handlers[e.Topic()] {
		l.wg.Add(1)
		go func(h events.Handler) {
			if err := h.Handle(e); err != nil {
				// TODO Do some error Handling
				log.Println(err)
			}
			l.wg.Done()
		}(h)
	}
	l.wg.Done()
}
