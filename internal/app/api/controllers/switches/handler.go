package switches

import (
	"encoding/json"
	"fmt"
	"kbswitch/internal/core/models"
	"kbswitch/internal/core/repositories"
	"net/http"
	"strconv"
)

type controller struct {
	repo repositories.SwitchesRepo
}

func New(repo repositories.SwitchesRepo) controller {
	return controller{
		repo: repo,
	}
}

func writeErr(err string, status int, w http.ResponseWriter) {
	e := models.APIError{
		Status:  status,
		Message: err,
	}

	w.WriteHeader(status)
	fmt.Fprint(w, e)
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
		writeErr(err.Error(), http.StatusInternalServerError, w)
		return
	}
	if resp == nil {
		writeErr(
			"collection got nil from a repo",
			http.StatusInternalServerError,
			w,
		)
		return
	}

	dtos := make([]SwitchDTO, len(resp))
	for i, item := range resp {
		dtos[i] = AsDTO(item)
	}

	json, _ := json.Marshal(dtos)

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", string(json[:]))
}

// HandleSwitchByID godoc
//
//	@Summary		Get switch by ID
//	@Description	Gives a single switch by database ID
//	@Param			id	path		int	true	"Switch ID"
//	@Tags			switches
//	@Produce		json
//	@Success		200	{object}	      SwitchDTO
//	@Failure		500	{object}	models.APIError
//	@Failure		400	{object}	models.APIError
//	@Failure		404	{object}	models.APIError
//	@Router			/api/switches/{id} [get]
func (c controller) HandleSwitchByID(w http.ResponseWriter, r *http.Request) {
	p := r.PathValue("id")
	if p == "" {
		writeErr("request parameter is missing", http.StatusBadRequest, w)
		return
	}
	id, err := strconv.Atoi(p)
	if err != nil {
		writeErr("request parameter should be int", http.StatusBadRequest, w)
		return
	}

	resp, err := c.repo.GetByID(id)
	if err != nil {
		writeErr(err.Error(), http.StatusInternalServerError, w)
		return
	}
	if resp == nil {
		writeErr("no resource found for a given id", http.StatusNotFound, w)
		return
	}

	dto := AsDTO(*resp)
	json, _ := json.Marshal(dto)

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", string(json[:]))
}
