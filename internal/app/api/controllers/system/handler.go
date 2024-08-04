package system

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type controller struct {
	buildDate time.Time
}

func New(buildDate time.Time) controller {
	return controller{
		buildDate: buildDate,
	}
}


// HandleAbout godoc
//
//	@Summary		Gives system info about api
//	@Description	api info
//	@Tags			system
//	@Produce		json
//	@Success		200	{object}	aboutDTO
//	@Router			/api/system/about [get]
func (c controller) HandleAbout(w http.ResponseWriter, r *http.Request) {
	dto := aboutDTO{
		Product:       "KBswitch website backend Api",
		Author:        "VPavliashvili",
		Version:       "1.0",
		BuildDatetime: c.buildDate.Format(time.DateTime),
	}

	json, _ := json.Marshal(dto)

	fmt.Fprint(w, string(json[:]))
}
