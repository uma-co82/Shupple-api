package domain

/**
 * エラーが発生した場合にフロントへ返却するError構造体
 */
type Error struct {
	Code              int      `json:"code"`
	Message           string   `json:"message"`
	ValidationMessage []string `json:"validationMessage"`
	FrontMessage      string   `json:"frontMessage"`
}

func (e *Error) Error() string {
	return e.Message
}

func RaiseError(code int, msg string, validationMessages []string) error {
	return &Error{Code: code, Message: msg, ValidationMessage: validationMessages}
}

func RaiseDBError() error {
	return &Error{Code: 500, Message: "DB Error Occur", ValidationMessage: nil}
}
