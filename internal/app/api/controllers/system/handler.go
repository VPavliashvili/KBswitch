package system

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type controller struct {
}

func New() controller {
	return controller{}
}

func (c controller) HandleAbout(w http.ResponseWriter, r *http.Request) {
	dto := aboutDTO{
		Product:       "Books Api",
		Author:        "VPavliashvili",
		Version:       "1.0",
		BuildDatetime: "tbd",
	}

	json, _ := json.Marshal(dto)

	fmt.Fprint(w, string(json[:]))
}
