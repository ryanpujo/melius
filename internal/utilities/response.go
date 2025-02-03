package utilities

import "github.com/ryanpujo/melius/internal/models"

type Response struct {
	ID        uint              `json:"id,omitempty"`
	Token     string            `json:"token,omitempty"`
	Country   *models.Country   `json:"country,omitempty"`
	Countries []*models.Country `json:"countries,omitempty"`
	State     *models.State     `json:"state,omitempty"`
	City      *models.City      `json:"city,omitempty"`
	Address   *models.Address   `json:"address,omitempty"`
	Err       string            `json:"err,omitempty"`
	Message   string            `json:"message,omitempty"`
}

type OptFunc func(*Response)

func WithCountry(country *models.Country) OptFunc {
	return func(r *Response) {
		r.Country = country
	}
}

func WithCountries(countries []*models.Country) OptFunc {
	return func(r *Response) {
		r.Countries = countries
	}
}

func WithState(state *models.State) OptFunc {
	return func(r *Response) {
		r.State = state
	}
}

func WithCity(city *models.City) OptFunc {
	return func(r *Response) {
		r.City = city
	}
}

func WithAddress(address *models.Address) OptFunc {
	return func(r *Response) {
		r.Address = address
	}
}

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
