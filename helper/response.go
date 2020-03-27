package helper

type Response struct {
	Status int         `json:"status"`
	Data   interface{} `json:"data"`
}

type ErrorReponse struct {
	Status int         `json:"status"`
	Errs   interface{} `json:"errors"`
}

func NewResponse(status int, data interface{}) Response {
	return Response{Data: data, Status: status}
}

func NewErrorResponse(status int, err []error) ErrorReponse {
	listErrors := formatErrors(err)
	return ErrorReponse{Status: status, Errs: listErrors}
}

func formatErrors(errs []error) []string {
	var listErrors []string

	for _, err := range errs {
		listErrors = append(listErrors, err.Error())
	}
	return listErrors
}
