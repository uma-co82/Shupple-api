package service

func (e *Error) Error() string {
	return e.Message
}

func RaiseError(code int, msg string, vmsg []string) error {
	return &Error{Code: code, Message: msg, ValidationMessage: vmsg}
}
