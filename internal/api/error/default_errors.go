package apierror

const (
	TypeInvalidParameter      = "INVALID_PARAMETER_VALUE"
	TypeInvalidRequestBody    = "INVALID_REQUEST_BODY"
	TypeInternalEOF           = "INTERNAL_EOF"
	TypeInvalidDatetimeFormat = "INVALID_DATETIME_FORMAT"
)

var (
	ErrInvalidRequestBody = &Error{
		Type: TypeInvalidRequestBody,
	}

	ErrInternalEOF = &Error{
		Type: TypeInternalEOF,
	}
	ErrInvalidDatetimeFormat = &Error{
		Type: TypeInvalidDatetimeFormat,
	}
)

func NewInvalidParameterError(parameter string) *Error {
	return &Error{
		Type:      TypeInvalidParameter,
		Parameter: &parameter,
	}
}
