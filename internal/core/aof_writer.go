package core

import "github.com/kenelite/gedis/internal/response"

type AOFWriter interface {
	Write(value response.Value) error
}
