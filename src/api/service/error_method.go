package service

func (e *Error) Error() string {
	return e.Message
}

func RaiseError(code int, msg string, validationMessages []string) error {
	return &Error{Code: code, Message: msg, ValidationMessage: validationMessages}
}

func RaiseDBError() error {
	return &Error{Code: 500, Message: "DB Error Occur", ValidationMessage: nil}
}
