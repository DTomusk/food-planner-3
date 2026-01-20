package directive

import (
	"context"
	"foodplanner/internal/auth"

	"github.com/99designs/gqlgen/graphql"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

func AuthDirective(
	ctx context.Context,
	obj interface{},
	next graphql.Resolver,
) (interface{}, error) {
	_, err := auth.ClaimsFromContext(ctx)
	if err != nil {
		return nil, gqlerror.Errorf("unauthorized")
	}
	return next(ctx)
}
