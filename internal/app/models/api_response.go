package models

type MessageResponse struct {
	Message string `json:"message"`
}

type DataResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type ErrorResponse struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}
