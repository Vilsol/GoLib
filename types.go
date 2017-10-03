package GoLib

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Status  int    `json:"-"`
}

type GenericResponse struct {
	Success bool           `json:"success"`
	Error   *ErrorResponse `json:"error,omitempty"`
	Data    interface{}    `json:"data,omitempty"`
}

var (
	ErrorInternalError    = ErrorResponse{-1, "Internal error: ", 500}
	ErrorInvalidParameter = ErrorResponse{-2, "", 400}
)
