package apierror

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type Handler struct {
	logger *zap.Logger
}

func NewErrorHandler(logger *zap.Logger) *Handler {
	return &Handler{
		logger: logger,
	}
}

func (h *Handler) Handler(err error, ctx echo.Context) {
	if handled := h.handleAsErrors(err, ctx); handled {
		return
	}

	switch typedErr := errors.Cause(err).(type) {
	case validation.Errors:
		typedErr, ok := typedErr.Filter().(validation.Errors)
		if !ok {
			return
		}

		NewErrors(http.StatusBadRequest, mapValidationErrors(typedErr, "")...).ToHTTPResponse(ctx, h.logger)
	case validation.Error:
		NewErrors(http.StatusBadRequest, &Error{
			Type: typeFromCode(typedErr.Code()),
		}).ToHTTPResponse(ctx, h.logger)
	case *echo.HTTPError:
		h.handleEchoError(typedErr, err, ctx)
	default:
		h.logger.With(zap.Error(err)).Error(err.Error())

		if responseErr := ctx.NoContent(http.StatusInternalServerError); responseErr != nil {
			h.logger.With(zap.Error(err)).Error("setting response on context error")
		}
	}
}

func (h *Handler) handleAsErrors(err error, ctx echo.Context) bool {
	var handled bool

	syntaxErr := &json.SyntaxError{}

	if errors.As(err, &syntaxErr) {
		NewErrors(http.StatusBadRequest, ErrInvalidRequestBody).ToHTTPResponse(ctx, h.logger)

		handled = true
	} else if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) {
		h.logger.With(zap.Error(err)).Error("EOF error")

		NewErrors(http.StatusServiceUnavailable, ErrInternalEOF).ToHTTPResponse(ctx, h.logger)

		handled = true
	}

	unmarshalErr := &json.UnmarshalTypeError{}

	if errors.As(err, &unmarshalErr) {
		NewErrors(
			http.StatusBadRequest,
			NewInvalidParameterError(unmarshalErr.Field),
		).ToHTTPResponse(ctx, h.logger)

		handled = true
	}

	return handled
}

func (h *Handler) handleEchoError(typedErr *echo.HTTPError, err error, ctx echo.Context) {
	if typedErr.Internal != nil {
		if _, ok := typedErr.Internal.(*time.ParseError); ok {
			NewErrors(http.StatusBadRequest, ErrInvalidDatetimeFormat).ToHTTPResponse(ctx, h.logger)
		}

		h.logger.With(zap.Error(typedErr.Internal)).Warn(typedErr.Error())
	}

	if responseErr := ctx.NoContent(typedErr.Code); responseErr != nil {
		h.logger.With(zap.Error(err)).Error("setting response on context error")
	}
}

func mapValidationErrors(errs validation.Errors, prefix string) []*Error {
	mappedErrors := make([]*Error, 0)

	for field, err := range errs {
		fieldName := prefix + field

		if typedErr, ok := err.(validation.Error); ok {
			params := typedErr.Params()
			if params == nil {
				params = make(map[string]any)
			}

			params["parameter"] = fieldName

			mappedErrors = append(mappedErrors, &Error{
				Type:      typeFromCode(typedErr.Code()),
				Parameter: &fieldName,
			})

			continue
		}

		if typedErr, ok := err.(validation.Errors); ok {
			mappedErrors = append(mappedErrors, mapValidationErrors(typedErr, fieldName+".")...)
		}
	}

	return mappedErrors
}

func typeFromCode(code string) string {
	namespaceAsBytes := []byte(code)
	separatorIndex := strings.LastIndex(code, ".")

	return strings.ToUpper(string(namespaceAsBytes[separatorIndex+1:]))
}
