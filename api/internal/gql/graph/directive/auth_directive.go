package directive

import (
	"context"
	"foodplanner/internal/auth"
	"foodplanner/internal/gql/graph/errors"

	"github.com/99designs/gqlgen/graphql"
)

func AuthDirective(
	ctx context.Context,
	obj interface{},
	next graphql.Resolver,
) (interface{}, error) {
	_, ok := auth.ClaimsFromContext(ctx)
	if !ok {
		return nil, errors.NewUnauthenticatedError("user is not authenticated")
	}
	return next(ctx)
}
