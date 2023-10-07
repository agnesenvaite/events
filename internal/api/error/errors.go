package apierror

import (
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type Error struct {
	Type      string  `json:"type"`
	Parameter *string `json:"parameter"`
}

type ListedError struct {
	Errors       []*Error `json:"errors"`
	responseCode int      `json:"-"`
}

func NewErrors(responseCode int, errors ...*Error) *ListedError {
	return &ListedError{
		Errors:       errors,
		responseCode: responseCode,
	}
}

func (l *ListedError) ToHTTPResponse(ctx echo.Context, logger *zap.Logger) {
	if len(l.Errors) == 0 {
		if err := ctx.NoContent(l.responseCode); err != nil {
			logger.With(zap.Error(err)).Error("setting response on context error")
		}
	}

	if err := ctx.JSON(l.responseCode, l); err != nil {
		logger.With(zap.Error(err)).Error("setting response on context error")
	}
}
func (l *ListedError) Error() string {
	result := []string{strconv.Itoa(l.responseCode)}

	for i := range l.Errors {
		result = append(result, l.Errors[i].Type)
	}

	return strings.Join(result, ",")
}
