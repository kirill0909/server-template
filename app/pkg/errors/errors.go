package errors

type CustomError struct {
	desc       string
	statusCode int
	showToUser bool
}

func NewCustomError(desc string, status int, show bool) error {
	return &CustomError{
		desc:       desc,
		statusCode: status,
		showToUser: show,
	}
}

func (e *CustomError) Error() string {
	return e.desc
}

func (e *CustomError) GetStatusCode() int {
	return e.statusCode
}

func (e *CustomError) IsNeedToBeShown() bool {
	return e.showToUser
}
