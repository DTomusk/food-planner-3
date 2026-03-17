package errors

import (
	"github.com/vektah/gqlparser/v2/gqlerror"
)

func NewUnauthenticatedError(message string) *gqlerror.Error {
	return &gqlerror.Error{
		Message: message,
		Extensions: map[string]interface{}{
			"code": "UNAUTHENTICATED",
		},
	}
}
