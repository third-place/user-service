package util

type InputFieldError struct {
	error
	Input   string `json:"input"`
	Message string `json:"message"`
}

func NewInputFieldError(input string, message string) *InputFieldError {
	return &InputFieldError{
		Input:   input,
		Message: message,
	}
}
