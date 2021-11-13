package app

import (
	"context"
	"time"

	"github.com/leksss/hw-test/hw12_13_14_15_calendar/internal/domain/entities"
	"github.com/leksss/hw-test/hw12_13_14_15_calendar/internal/domain/errors"
	"github.com/leksss/hw-test/hw12_13_14_15_calendar/internal/domain/interfaces"
	"github.com/leksss/hw-test/hw12_13_14_15_calendar/internal/infrastructure/logger"
	pb "github.com/leksss/hw-test/hw12_13_14_15_calendar/pb/event"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type App struct {
	logger  logger.Log
	storage interfaces.Storage
}

func New(logger logger.Log, storage interfaces.Storage) *App {
	return &App{
		logger:  logger,
		storage: storage,
	}
}

func (a *App) CreateEvent(ctx context.Context, in *pb.CreateEventRequest) (*pb.CreateEventResponse, error) {
	if err := validateEvent(in.Event); err != nil {
		return &pb.CreateEventResponse{
			Success: false,
			Errors:  toProtoError([]*pb.Error{}, err),
		}, nil
	}

	uuid, err := a.storage.CreateEvent(ctx, toEntityEvent(in.Event))
	if err != nil {
		return &pb.CreateEventResponse{
			Success: false,
			Errors:  toProtoError([]*pb.Error{}, err),
		}, nil
	}

	return &pb.CreateEventResponse{
		Success: true,
		UUID:    uuid,
	}, nil
}

func (a *App) UpdateEvent(ctx context.Context, in *pb.UpdateEventRequest) (*pb.UpdateEventResponse, error) {
	var pbErrors []*pb.Error
	if in.UUID == "" {
		pbErrors = toProtoError(pbErrors, errors.ErrEventUUIDIsRequired)
	}
	if err := validateEvent(in.Event); err != nil {
		pbErrors = toProtoError(pbErrors, err)
	}
	if len(pbErrors) > 0 {
		return &pb.UpdateEventResponse{
			Success: false,
			Errors:  pbErrors,
		}, nil
	}

	rowsCnt, err := a.storage.UpdateEvent(ctx, in.UUID, toEntityEvent(in.Event))
	if err != nil {
		return &pb.UpdateEventResponse{
			Success: false,
			Errors:  toProtoError([]*pb.Error{}, err),
		}, nil
	}
	if rowsCnt == 0 {
		return &pb.UpdateEventResponse{
			Success: false,
			Errors:  toProtoError([]*pb.Error{}, errors.ErrEventNotFound),
		}, nil
	}

	return &pb.UpdateEventResponse{
		Success: true,
	}, nil
}

func (a *App) DeleteEvent(ctx context.Context, in *pb.DeleteEventRequest) (*pb.DeleteEventResponse, error) {
	if in.UUID == "" {
		return &pb.DeleteEventResponse{
			Success: false,
			Errors:  toProtoError([]*pb.Error{}, errors.ErrEventUUIDIsRequired),
		}, nil
	}

	err := a.storage.DeleteEvent(ctx, in.UUID)
	if err != nil {
		return &pb.DeleteEventResponse{
			Success: false,
			Errors:  toProtoError([]*pb.Error{}, err),
		}, nil
	}

	return &pb.DeleteEventResponse{
		Success: true,
	}, nil
}

func (a *App) GetEventList(ctx context.Context, in *pb.GetEventListRequest) (*pb.GetEventListEventResponse, error) {
	filter := entities.Filter{}
	if in.Limit > 0 {
		filter.Limit = in.Limit
	}
	if in.Offset > 0 {
		filter.Offset = in.Offset
	}
	if in.UUID != "" {
		filter.UUID = in.UUID
	}

	events, err := a.storage.GetEventList(ctx, filter)
	if err != nil {
		return &pb.GetEventListEventResponse{
			Success: false,
			Errors:  toProtoError([]*pb.Error{}, err),
		}, nil
	}

	pbEvents := make([]*pb.Event, 0)
	for _, e := range events {
		pbEvents = append(pbEvents, toProtoEvent(e))
	}

	return &pb.GetEventListEventResponse{
		Success: true,
		Events:  pbEvents,
	}, nil
}

func toProtoEvent(event *entities.Event) *pb.Event {
	return &pb.Event{
		UUID:      event.UUID,
		OwnerID:   event.OwnerID,
		Title:     event.Title,
		StartedAt: timestamppb.New(*event.StartedAt),
		EndedAt:   timestamppb.New(*event.EndedAt),
		Text:      event.Text,
		NotifyFor: event.NotifyFor,
	}
}

func toProtoError(errs []*pb.Error, err error) []*pb.Error {
	return append(errs, &pb.Error{
		Code: "event",
		Msg:  err.Error(),
	})
}

func toEntityEvent(event *pb.Event) *entities.Event {
	startedAt := time.Unix(event.StartedAt.Seconds, 0)
	endedAt := time.Unix(event.EndedAt.Seconds, 0)
	return &entities.Event{
		OwnerID:   event.OwnerID,
		Title:     event.Title,
		StartedAt: &startedAt,
		EndedAt:   &endedAt,
		Text:      event.Text,
		NotifyFor: event.NotifyFor,
	}
}

func validateEvent(event *pb.Event) error {
	if event.OwnerID == "" {
		return errors.ErrEventOwnerIDIsRequired
	}
	if event.Title == "" {
		return errors.ErrEventTitleIsRequired
	}
	if event.StartedAt == nil {
		return errors.ErrEventStartedAtIsRequired
	}
	if event.EndedAt == nil {
		return errors.ErrEventEndedAtIsRequired
	}
	return nil
}
