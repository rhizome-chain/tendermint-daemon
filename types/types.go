package types

type NoDataError struct {
	error
}

func (err NoDataError) Error() string{
	return "no data"
}

func NewNoDataError() NoDataError {
	return NoDataError{}
}

func IsNoDataError(err error) bool {
	_,ok := err.(NoDataError)
	return ok
}