package handle

import (
	"github.com/kenelite/gedis/internal/response"
	"github.com/kenelite/gedis/internal/storage"
)

func Save(args []response.Value) response.Value {
	if err := storage.SaveRDB("dump.rdb"); err != nil {
		return response.Value{Typ: "error", Str: "ERR SAVE failed: " + err.Error()}
	}
	return response.Value{Typ: "string", Str: "OK"}
}

func BgSave(args []response.Value) response.Value {
	go func() {
		_ = storage.SaveRDB("dump.rdb") // Optionally log the error
	}()
	return response.Value{Typ: "string", Str: "Background saving started"}
}
