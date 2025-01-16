package utilities

type Response struct {
	ID      uint   `json:"id,omitempty"`
	Token   string `json:"token,omitempty"`
	Err     string `json:"err,omitempty"`
	Message string `json:"message,omitempty"`
}

type OptFunc func(*Response)

func WithID(id uint) OptFunc {
	return func(r *Response) {
		r.ID = id
	}
}

func WithErr(err string) OptFunc {
	return func(r *Response) {
		r.Err = err
	}
}

func WithToken(token string) OptFunc {
	return func(r *Response) {
		r.Token = token
	}
}

func NewResponse(message string, opts ...OptFunc) Response {
	res := Response{
		Message: message,
	}

	for _, opt := range opts {
		opt(&res)
	}
	return res
}