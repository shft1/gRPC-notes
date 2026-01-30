package eventbus

import (
	"context"

	"github.com/shft1/grpc-notes/internal/domain/notes"
	"github.com/shft1/grpc-notes/observability/logger"
)

type eventBus struct {
	log logger.Logger
	bus chan *notes.NoteEvent
}

func NewEventBus(log logger.Logger, cap int) *eventBus {
	return &eventBus{
		log: log,
		bus: make(chan *notes.NoteEvent, cap),
	}
}

func (b *eventBus) Produce(event *notes.NoteEvent) {
	select {
	case b.bus <- event:
		b.log.Info(
			"event has been sent to the bus",
			logger.NewField("id", event.ID),
			logger.NewField("title", event.Title),
		)
	default:
		b.log.Warn(
			"the bus is full",
			logger.NewField("id", event.ID),
			logger.NewField("title", event.Title),
		)
	}
}

func (b *eventBus) Consume(ctx context.Context) (*notes.NoteEvent, error) {
	select {
	case event := <-b.bus:
		return event, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}
