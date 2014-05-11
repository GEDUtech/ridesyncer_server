package api

import (
	"encoding/json"
	"github.com/martini-contrib/render"
	"net/http"
	"time"
)

const (
	GOOGLE_API_KEY = "AIzaSyAaXv0wQpVcQDo_p_KQfd8T7dndbUlx4TY"
)

func decode(req *http.Request, render render.Render, data interface{}) error {
	decoder := json.NewDecoder(req.Body)

	err := decoder.Decode(data)
	if err != nil {
		switch err.(type) {
		case *time.ParseError:
			// ignore, will be caught by validation
			return nil
		default:
			render.JSON(http.StatusNotAcceptable, map[string]interface{}{"message": "Invalid JSON", "reason": err})
		}

	}
	return err
}
