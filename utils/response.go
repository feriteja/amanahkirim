package utils

type AppResponse struct {
	Message string
	Data    interface{}
}

type MutationResponse struct {
	Error_message []string
	Data          interface{}
}

func (e *AppResponse) Error() string {
	return e.Message
}

func (e *AppResponse) Response() AppResponse {
	return *e
}
