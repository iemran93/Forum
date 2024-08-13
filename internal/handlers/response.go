package handlers

type Response struct {
    Message string `json:"message"`
    Error   string `json:"error,omitempty"`
}
