package errs

type AppError struct {
	Status  int
	Message string
}

func (a AppError) Error() string {
	return a.Message
}
