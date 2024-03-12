package response

const (
	StatusOK  = "OK"
	StatusErr = "Error"
)

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

func Ok() Response {
	return Response{
		Status: StatusOK,
	}
}

func Err(errMsg string) Response {
	return Response{
		Status: StatusErr,
		Error:  errMsg,
	}
}
