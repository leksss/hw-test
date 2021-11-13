package internalgrpc

import (
	"context"

	"github.com/leksss/hw-test/hw12_13_14_15_calendar/internal/app"
	pb "github.com/leksss/hw-test/hw12_13_14_15_calendar/pb/event"
)

type EventService struct {
	pb.UnimplementedEventServiceServer
	app *app.App
}

func NewEventService(app *app.App) *EventService {
	return &EventService{app: app}
}

func (s *EventService) CreateEvent(ctx context.Context, in *pb.CreateEventRequest) (*pb.CreateEventResponse, error) {
	return s.app.CreateEvent(ctx, in)
}

func (s *EventService) UpdateEvent(ctx context.Context, in *pb.UpdateEventRequest) (*pb.UpdateEventResponse, error) {
	return s.app.UpdateEvent(ctx, in)
}

func (s *EventService) DeleteEvent(ctx context.Context, in *pb.DeleteEventRequest) (*pb.DeleteEventResponse, error) {
	return s.app.DeleteEvent(ctx, in)
}

func (s *EventService) GetEventList(ctx context.Context, in *pb.GetEventListRequest) (*pb.GetEventListEventResponse, error) {
	return s.app.GetEventList(ctx, in)
}
