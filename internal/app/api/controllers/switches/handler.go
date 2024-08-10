package switches

import (
	"encoding/json"
	"errors"
	"fmt"
	"kbswitch/internal/core/models"
	"kbswitch/internal/core/repositories"
	"net/http"
)

type controller struct {
	repo repositories.SwitchesRepo
}

func New(repo repositories.SwitchesRepo) controller {
	return controller{
		repo: repo,
	}
}

func writeErr(err error, status int, w http.ResponseWriter) {
	e := models.APIError{
		Status:  status,
		Message: err.Error(),
	}

	w.WriteHeader(e.Status)
	fmt.Fprint(w, err.Error())
}

// HandleSwitches godoc
//
//	@Summary		Get all switches
//	@Description	Gives array of all keyboard switches
//	@Tags			switches
//	@Produce		json
//	@Success		200	{array}	    SwitchDTO
//	@Failure		500	{object}	models.APIError
//	@Router			/api/switches [get]
func (c controller) HandleSwitches(w http.ResponseWriter, r *http.Request) {
	resp, err := c.repo.GetAll()
	if err != nil {
		writeErr(err, http.StatusInternalServerError, w)
		return
	}
	if resp == nil {
		writeErr(
			errors.New("collection got nil from a repo"),
			http.StatusInternalServerError,
			w,
		)
		return
	}

	dtos := make([]SwitchDTO, len(resp))
	for i, item := range resp {
		dtos[i] = asDTO(item)
	}

	json, _ := json.Marshal(dtos)

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", string(json[:]))
}

// // ImageTest godoc
// //
// // @Summary idk
// // @Tags switches
// // @Produce image/png
// func (c controller) ImageTest(w http.ResponseWriter, r *http.Request) {
// 	var imageTemplate string = `<!DOCTYPE html>
//     <html lang="en"><head></head>
//     <body><img src="data:image/jpg;base64,{{.Image}}"></body>`
//
//     buffer := new(bytes.Buffer )
//     str := base64.StdEncoding.EncodeToString(buffer.Bytes())
//
// 	tmpl, _ := template.New("image").Parse(imageTemplate)
//     data := map[string]any{"Image": str}
//
//     tmpl.Execute(w, data)
// }
