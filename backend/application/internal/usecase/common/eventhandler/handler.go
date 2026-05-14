package eventhandler

import (
	"context"
	"log/slog"
)

type Handler struct {
	decoder Decoder

	resumeUpsertedUsecase ResumeUpsertedUsecase
}

func NewHandler(decoder Decoder, resumeUpsertedUsecase ResumeUpsertedUsecase) *Handler {
	return &Handler{decoder: decoder, resumeUpsertedUsecase: resumeUpsertedUsecase}
}

func (h *Handler) HandleEvent(ctx context.Context, event Event) {
	eventType, ok := event.Headers["event_type"]
	if !ok {
		return
	}
	var err error
	switch string(eventType) {
	case "ResumeUpsertedEvent":
		err = h.handleResumeUpsertedEvent(ctx, event.Payload)
	}

	if err != nil {
		slog.Info("oh error", "error", err)
	}
}

func (h *Handler) handleResumeUpsertedEvent(ctx context.Context, payload []byte) error {
	event, err := h.decoder.DecodeResumeUpsertedEvent(ctx, payload)
	if err != nil {
		return err
	}

	err = h.resumeUpsertedUsecase.Execute(ctx, event)
	if err != nil {
		return err
	}
	return nil
}
