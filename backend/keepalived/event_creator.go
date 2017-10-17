package keepalived

import (
	"time"

	"github.com/sensu/sensu-go/backend/messaging"
	"github.com/sensu/sensu-go/types"
)

// EventCreator creates alert and resolution events given an entity.
type EventCreator interface {
	Warn(entity *types.Entity) error
	Critical(entity *types.Entity) error
	Pass(entity *types.Entity) error
}

// MessageBusEventCreator publishes a message to WizardBus when Alert or Resolve
// is called, formatting the messages approriately based on the entity.
type MessageBusEventCreator struct {
	MessageBus messaging.MessageBus
}

// Warn sends a check with status of warn for a keepalive.
func (creatorPtr *MessageBusEventCreator) Warn(entity *types.Entity) error {
	event := createEvent(entity)
	event.Check.Status = 1

	return creatorPtr.MessageBus.Publish(messaging.TopicEventRaw, event)
}

// Critical sends a check with status of critical for a keepalive.
func (creatorPtr *MessageBusEventCreator) Critical(entity *types.Entity) error {
	event := createEvent(entity)
	event.Check.Status = 2

	return creatorPtr.MessageBus.Publish(messaging.TopicEventRaw, event)
}

// Pass sends a check with a status of OK for a keepalive.
func (creatorPtr *MessageBusEventCreator) Pass(entity *types.Entity) error {
	event := createEvent(entity)
	event.Check.Status = 0

	return creatorPtr.MessageBus.Publish(messaging.TopicEventRaw, event)
}

func createEvent(entity *types.Entity) *types.Event {
	keepaliveCheck := &types.Check{
		Config: &types.CheckConfig{
			Name:          "keepalive",
			Interval:      entity.KeepaliveTimeout,
			Subscriptions: []string{""},
			Command:       "",
			Handlers:      []string{"keepalive"},
			Environment:   entity.Environment,
			Organization:  entity.Organization,
		},
		Status: 1,
	}
	keepaliveEvent := &types.Event{
		Timestamp: time.Now().Unix(),
		Entity:    entity,
		Check:     keepaliveCheck,
	}

	return keepaliveEvent
}
