package errors

import (
	"fmt"
)

type Code string

type StatusMessage struct {
	Message    string
	StatusCode int
}

func GetStatusMessage(message string, statusCode int) StatusMessage {
	return StatusMessage{
		Message:    message,
		StatusCode: statusCode,
	}
}

const (
	// ErrorCodeResourceNotFound error when the requested resource was not found
	ErrorCodeResourceNotFound = "resource_not_exist"
	// ErrorCodeInvalidArgs error when the request data is incorrect or incomplete
	ErrorCodeInvalidArgs = "invalid_args"
	// ErrorCodeBadRequest error when the needed input is not provided
	ErrorCodeBadRequest = "bad_request"
	//ErrorCodeResourceAlreadyExists when a resource already exists
	ErrorCodeResourceAlreadyExists = "resource_already_exist"
	// ErrorCodeUnknownIssue when the issue is unknown
	ErrorCodeUnknownIssue = "unknown_issue"
	// ErrorCodePermissionDenied when user dont have necessary permissions
	ErrorCodePermissionDenied = "permission_denied"
	// ErrorCodeResourceStateConflict when the resource is in another state and generates a conflict
	ErrorCodeResourceStateConflict = "resource_state_conflict"
	// ErrorCodeNotImplemented when some method is not implemented
	ErrorCodeNotImplemented = "not_implemented"
	// ErrorCodeUnauthorized when user is not authorized
	ErrorCodeUnauthorized = "unauthorized"
	// ErrorCodeNotFound when is not found
	ErrorCodeNotFound = "not_found"
	// ErrorCodeElasticSearchError when using a elastic search and it returned an error
	ErrorCodeElasticSearchError = "elasticsearch_error"
	// ErrorCodeDatabaseError when using database commands and it returned an error
	ErrorCodeDatabaseError = "database_error"
)

var (
	ErrorMessageList = map[Code]StatusMessage{
		ErrorCodeResourceNotFound:      GetStatusMessage("Resource not found", 404),
		ErrorCodeInvalidArgs:           GetStatusMessage("Invalid parameters were passed.", 400),
		ErrorCodeBadRequest:            GetStatusMessage("Required data not valid.", 400),
		ErrorCodeResourceAlreadyExists: GetStatusMessage("The posted resource already existed.", 400),
		ErrorCodeUnknownIssue:          GetStatusMessage("Unknown issue was caught and message was not specified.", 500),
		ErrorCodePermissionDenied:      GetStatusMessage("Current user has no permission to perform the action.", 403),
		ErrorCodeResourceStateConflict: GetStatusMessage("The posted resource already existed.", 409),
		ErrorCodeNotImplemented:        GetStatusMessage("Method not implemented", 501),
		ErrorCodeElasticSearchError:    GetStatusMessage("Elastic search error.", 500),
		ErrorCodeDatabaseError:         GetStatusMessage("Database error.", 500),
	}
)

func AddError(code Code, message string, statusCode int) {
	ErrorMessageList[code] = GetStatusMessage(message, statusCode)
}

type Error struct {
	Source     string                `json:"source"`
	Message    string                `json:"message"`
	Code       Code                  `json:"code"`
	Fields     []map[string][]string `json:"fields,omitempty"`
	StatusCode int                   `json:"-"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("code: %s, message: %s", e.Code, e.Message)
}

func New(source string, code Code) *Error {
	var (
		message    = "Error not described"
		statusCode = 999
	)
	if val, ok := ErrorMessageList[code]; ok {
		message = val.Message
		statusCode = val.StatusCode
	}
	return &Error{
		Source:     source,
		Message:    message,
		Code:       code,
		StatusCode: statusCode,
	}
}

func NewCommon(source string, err error) *Error {
	if err == nil {
		return nil
	}
	if e, ok := err.(*Error); ok {
		return e
	}
	return &Error{
		Source:     source,
		Message:    err.Error(),
		Code:       ErrorCodeUnknownIssue,
		StatusCode: ErrorMessageList[ErrorCodeUnknownIssue].StatusCode,
	}
}

func (e *Error) SetCode(code Code) *Error {
	e.Code = code
	return e
}

func (e *Error) SetMessage(format string, args ...interface{}) *Error {
	e.Message = fmt.Sprintf(format, args...)
	return e
}

func (e *Error) AddFieldError(field string, message ...string) *Error {
	e.initializeFields()
	if value, ok := e.Fields[0][field]; ok {
		value = append(value, message...)
		return e
	}
	e.Fields[0][field] = message
	return e
}

func (e *Error) initializeFields() {
	if e.Fields == nil {
		e.Fields = []map[string][]string{
			map[string][]string{},
		}
	}
}
