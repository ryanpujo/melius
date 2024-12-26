package utilities

type Response struct {
	ID      uint   `json:"id,omitempty"`
	Token   string `json:"token,omitempty"`
	Err     string `json:"err,omitempty"`
	Message string `json:"message,omitempty"`
}
