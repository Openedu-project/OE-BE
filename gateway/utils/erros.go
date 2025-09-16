package utils

type AppError struct {
	Status  int
	Message string
}

func (e *AppError) Error() string {
	return e.Message
}
