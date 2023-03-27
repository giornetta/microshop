package events

import (
	"context"
	"encoding/json"
	"errors"
)

type Topic string

func (t Topic) String() string {
	return string(t)
}

type Type string

func (t Type) String() string {
	return string(t)
}

type Key string

type Event interface {
	Type() Type
	Key() Key
	Topic() Topic
}

type Publisher interface {
	Publish(e Event, ctx context.Context) error
}

type Handler interface {
	Handle(e Event) error
}

type Decoder func(payload []byte) (Event, error)

func fromJSON[T Event](payload []byte) (Event, error) {
	var evt T
	if err := json.Unmarshal(payload, &evt); err != nil {
		return nil, err
	}

	return evt, nil
}

var decoders = make(map[Type]Decoder)

func registerEvent[T Event](eventType Type) {
	decoders[eventType] = fromJSON[T]
}

func Decode(t Type, payload []byte) (Event, error) {
	f, ok := decoders[t]
	if !ok {
		return nil, errors.New("the provided event type is not recognized")
	}

	return f(payload)
}
