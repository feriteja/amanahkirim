package utils

type AppResponse struct {
	Code    int
	Message string
}

func (e *AppResponse) Error() string {
	return e.Message
}

func (e *AppResponse) Response() AppResponse {
	return *e
}

type StatusMessage struct {
	Code    int
	Message string
}

// var StatusMessageResponse = map[string]StatusMessage{
// 	"User has created": {Code: http.StatusCreated, Message: "user created"},
// }
