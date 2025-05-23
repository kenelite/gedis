package handle

import (
	"github.com/kenelite/gedis/internal/response"
)

func Ping(args []response.Value) response.Value {
	if len(args) == 0 {
		return response.Value{Typ: "string", Str: "PONG"}
	}

	return response.Value{Typ: "string", Str: args[0].Bulk}
}
