package handlers

import (
	"os"

	"github.com/rs/zerolog"

	"github.com/JonPulfer/gofiggy/pkg/events"
)

type LoggingHandler struct {
	logger zerolog.Logger
}

func NewMockHandler() LoggingHandler {
	return LoggingHandler{logger: zerolog.New(os.Stderr).With().Timestamp().Logger()}
}

func (lh LoggingHandler) ObjectCreated(obj interface{}) {
	ev := events.New(obj, "created")
	lh.logger.Log().Fields(map[string]interface{}{"event": ev}).Msg("received created event")
}

func (lh LoggingHandler) ObjectDeleted(obj interface{}) {
	ev := events.New(obj, "deleted")
	lh.logger.Log().Fields(map[string]interface{}{"event": ev}).Msg("received deleted event")
}

func (lh LoggingHandler) ObjectUpdated(oldObj interface{}, newObj interface{}) {
	ev := events.New(newObj, "updated")
	lh.logger.Log().Fields(map[string]interface{}{"event": ev}).Msg("received updated event")
}
