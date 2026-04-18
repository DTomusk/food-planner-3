package graph

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strings"

	"foodplanner/internal/auth"
	"foodplanner/internal/events"
	"foodplanner/internal/logging"
	"foodplanner/internal/middleware"

	"github.com/99designs/gqlgen/graphql"
	"github.com/google/uuid"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

const complexityLimitExceededReason = "complexity_limit_exceeded"

func NewComplexityLimitErrorPresenter(publisher events.Publisher, maxAcceptedComplexity int) graphql.ErrorPresenterFunc {
	return func(ctx context.Context, err error) *gqlerror.Error {
		presented := graphql.DefaultErrorPresenter(ctx, err)
		publishComplexityRejectedEvent(ctx, publisher, maxAcceptedComplexity, presented)
		return presented
	}
}

func publishComplexityRejectedEvent(ctx context.Context, publisher events.Publisher, maxAcceptedComplexity int, presented *gqlerror.Error) {
	if publisher == nil || !isComplexityLimitError(presented) {
		return
	}

	logger := logging.FromContext(ctx)
	var actorUserID *uuid.UUID
	if claims, ok := auth.ClaimsFromContext(ctx); ok {
		if parsedID, err := uuid.Parse(claims.UserID); err == nil {
			actorUserID = &parsedID
		}
	}

	opCtx := graphql.GetOperationContext(ctx)
	operationName := ""
	operationType := ""
	queryHash := ""
	if opCtx != nil {
		operationName = opCtx.OperationName
		queryHash = hashQuery(opCtx.RawQuery)
		if opCtx.Operation != nil {
			if operationName == "" {
				operationName = opCtx.Operation.Name
			}
			operationType = string(opCtx.Operation.Operation)
		}
	}

	ipAddress, _ := ctx.Value(middleware.IPKey).(string)
	userAgent, _ := ctx.Value(middleware.UserAgentKey).(string)
	path := ""
	if req, ok := ctx.Value(middleware.RequestKey).(*http.Request); ok && req != nil {
		path = req.URL.Path
		if userAgent == "" {
			userAgent = req.UserAgent()
		}
	}

	event := events.NewGraphQLRequestRejectedEvent(
		uuid.New(),
		actorUserID,
		operationName,
		operationType,
		ipAddress,
		userAgent,
		path,
		complexityLimitExceededReason,
		queryHash,
		maxAcceptedComplexity,
		presented.Message,
	)
	if err := publisher.Publish(ctx, event); err != nil {
		logger.Warn("Failed to publish GraphQL request rejected event", "eventType", event.Meta.Type, "reason", event.Reason, "err", err)
	}
}

func isComplexityLimitError(err *gqlerror.Error) bool {
	if err == nil {
		return false
	}

	message := strings.ToLower(err.Message)
	return strings.Contains(message, "complexity") && strings.Contains(message, "limit")
}

func hashQuery(rawQuery string) string {
	if rawQuery == "" {
		return ""
	}

	hash := sha256.Sum256([]byte(rawQuery))
	return hex.EncodeToString(hash[:])
}
