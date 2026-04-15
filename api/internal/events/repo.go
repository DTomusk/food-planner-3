package events

import (
	"context"
	"foodplanner/internal/db"
)

type EventsRepo struct {
}

const checkProcessedEventQuery = `
SELECT EXISTS (
	SELECT 1
	FROM processed_events
	WHERE event_id = $1
		AND consumer_group = $2
		AND handler_name = $3
)
`

const markProcessedEventQuery = `
INSERT INTO processed_events (
	event_id,
	consumer_group,
	handler_name
)
VALUES ($1, $2, $3)
ON CONFLICT (event_id, consumer_group, handler_name)
DO NOTHING
`

func NewEventsRepo() *EventsRepo {
	return &EventsRepo{}
}

func (r *EventsRepo) checkEventProcessed(ctx context.Context, db db.DBTX, eventID string, consumerGroup string, handlerName string) (bool, error) {
	var processed bool
	err := db.QueryRowContext(ctx, checkProcessedEventQuery, eventID, consumerGroup, handlerName).Scan(&processed)
	if err != nil {
		return false, err
	}

	return processed, nil
}

func (r *EventsRepo) markEventProcessed(ctx context.Context, db db.DBTX, eventID string, consumerGroup string, handlerName string) error {
	_, err := db.ExecContext(ctx, markProcessedEventQuery, eventID, consumerGroup, handlerName)
	return err
}
