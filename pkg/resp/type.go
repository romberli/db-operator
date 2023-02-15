package resp

type ErrResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func NewErrResponse(code int, message string) *ErrResponse {
	return &ErrResponse{
		Code:    code,
		Message: message,
	}
}
