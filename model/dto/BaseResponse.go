package dto

type BaseResponse[B any] struct {
	Body         B      `json:"body"`
	ErrorMessage string `json:"errorMessage"`
}
