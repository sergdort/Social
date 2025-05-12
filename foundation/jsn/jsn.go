package jsn

import (
	"encoding/json"
	"net/http"
)

func ReadJSON(r *http.Request, data any) error {

	var decoder = json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	defer r.Body.Close()

	return decoder.Decode(data)
}
