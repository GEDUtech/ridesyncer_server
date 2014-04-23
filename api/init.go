package api

import (
	"encoding/json"
	"github.com/martini-contrib/render"
	"net/http"
)

func decode(req *http.Request, render render.Render, data interface{}) error {
	decoder := json.NewDecoder(req.Body)

	err := decoder.Decode(data)
	if err != nil {
		render.JSON(http.StatusNotAcceptable, map[string]string{"message": "Invalid JSON"})
	}
	return err
}
