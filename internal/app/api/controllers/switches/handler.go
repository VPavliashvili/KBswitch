package switches

import (
	"encoding/json"
	"fmt"
	"kbswitch/internal/core/switches/models"
	"kbswitch/internal/core/switches/services"
	"net/http"
	"strconv"
)

type controller struct {
	service services.SwitchesService
}

func New(service services.SwitchesService) controller {
	return controller{
		service: service,
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
	resp, err := c.service.GetAll()
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

	resp, err := c.service.GetByID(id)
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

// HandleSwitchAdd godoc
//
//	@Summary		Add new switch
//	@Description	Add a new switch and get resource address
//	@Tags			switches
//	@Produce		json
//	@Accept		    json
//	@Param			newswitch	body   models.SwitchRequestBody   true    "Switch to add"
//	@Success		200	{object}	      string
//	@Failure		500	{object}	models.APIError
//	@Failure		400	{object}	models.APIError
//	@Router			/api/switches [post]
func (c controller) HandleSwitchAdd(w http.ResponseWriter, r *http.Request) {
	var req models.SwitchRequestBody
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&req)
	if err != nil {
		writeErr("invalid request model", http.StatusBadRequest, w)
		return
	}

	id, err := c.service.AddNew(req)
	if err != nil {
		writeErr(err.Error(), http.StatusInternalServerError, w)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "%s%s/%d", r.Host, r.URL.Path, *id)

}
