package switches

import (
	"encoding/json"
	"fmt"
	"kbswitch/internal/core/common"
	"kbswitch/internal/core/switches"
	"kbswitch/internal/core/switches/models"
	"net/http"
)

type controller struct {
	service switches.Service
}

func New(service switches.Service) controller {
	return controller{
		service: service,
	}
}

func writeErr(err string, status int, w http.ResponseWriter) {
	e := common.APIError{
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
//	@Success		200	{array}		SwitchDTO
//	@Failure		500	{object}	models.APIError
//	@Router			/api/switches [get]
func (c controller) HandleSwitches(w http.ResponseWriter, r *http.Request) {
	resp, err := c.service.GetAll()
	if err != nil {
		e := common.ToAPIErr(*err)
		w.WriteHeader(e.Status)
		fmt.Fprint(w, e.Error())
		return
	}
	if resp == nil {
		writeErr(
			"collection got nil from a service",
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

// HandleSingleSwitch godoc
//
//	@Summary		Get switch by ID
//	@Description	Gives a single switch by database ID
//	@Param			id	path	int	true	"Switch ID"
//	@Tags			switches
//	@Produce		json
//	@Success		200	{object}	SwitchDTO
//	@Failure		500	{object}	models.APIError
//	@Failure		400	{object}	models.APIError
//	@Failure		404	{object}	models.APIError
//	@Router			/api/switches/{brand}/{name} [get]
func (c controller) HandleSingleSwitch(w http.ResponseWriter, r *http.Request) {
	brand := r.PathValue("brand")
	name := r.PathValue("name")
	if brand == "" && name == "" {
		msg := "request parameters 'name' and 'brand' are missing"
		writeErr(msg, http.StatusBadRequest, w)
		return
	}
	if brand == "" {
		msg := "request parameter 'brand' is missing"
		writeErr(msg, http.StatusBadRequest, w)
		return
	}
	if name == "" {
		msg := "request parameter 'name' is missing"
		writeErr(msg, http.StatusBadRequest, w)
		return
	}

	resp, err := c.service.GetSingle(brand, name)
	if err != nil {
		writeErr(err.Error(), http.StatusInternalServerError, w)
		return
	}
	if resp == nil {
		writeErr("no resource found for a given name and brand", http.StatusNotFound, w)
		return
	}

	dto := AsDTO(*resp)
	json, _ := json.Marshal(dto)

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", string(json[:]))
}

// RemoveSwitch godoc
//
//	@Summary		Remove switch by its name and brand
//	@Description	removes switch
//	@Tags			switches
//	@Accept			json
//	@Produce		json
//	@Param			brand	path	int	true	"brand of the switch to delete"
//	@Param			name	path	int	true	"name of the switch to delete"
//	@Success		204
//	@Failure		500	{object}	models.APIError
//	@Failure		400	{object}	models.APIError
//	@Failure		404	{object}	models.APIError
//	@Router			/api/switches/{brand}/{name} [delete]
func (c controller) HandleSwitchRemove(w http.ResponseWriter, r *http.Request) {
	brand := r.PathValue("brand")
	name := r.PathValue("name")
	if brand == "" || name == "" {
		msg := "request parameters 'name' and 'brand' are missing"
		if brand != "" {
			msg = "request parameter 'name' is missing"
		} else if name != "" {
			msg = "request parameter 'brand' is missing"
		}
		writeErr(msg, http.StatusBadRequest, w)
		return
	}

	err := c.service.Remove(brand, name)
	if err != nil {
		e := common.ToAPIErr(*err)
		w.WriteHeader(e.Status)
		fmt.Fprint(w, e.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
	fmt.Fprintf(w, "")
}

// HandleSwitchUpdate godoc
//
//	@Summary		Modify/Update existing switch
//	@Description	Update a switch and get modified resource
//	@Tags			switches
//	@Produce		json
//	@Accept			json
//	@Param			brand	path		int	true	"brand of the switch to update"
//	@Param			name	path		int	true	"name of the switch to update"
//	@Success		200		{object}	SwitchDTO
//	@Failure		500		{object}	models.APIError
//	@Failure		400		{object}	models.APIError
//	@Router			/api/switches/{brand}/{name} [patch]
func (c controller) HandleSwitchUpdate(w http.ResponseWriter, r *http.Request) {
	brand := r.PathValue("brand")
	name := r.PathValue("name")
	if brand == "" && name == "" {
		msg := "request parameters 'name' and 'brand' are missing"
		writeErr(msg, http.StatusBadRequest, w)
		return
	}
	if brand == "" {
		msg := "request parameter 'brand' is missing"
		writeErr(msg, http.StatusBadRequest, w)
		return
	}
	if name == "" {
		msg := "request parameter 'name' is missing"
		writeErr(msg, http.StatusBadRequest, w)
		return
	}

	var req models.SwitchRequestBody
	if r.Body != nil {

		decoder := json.NewDecoder(r.Body)
		defer r.Body.Close()
		decoder.DisallowUnknownFields()
		err := decoder.Decode(&req)
		if err != nil {
			writeErr("invalid request model", http.StatusBadRequest, w)
			return
		}
	} else {
		writeErr("request body is entirely missing/nil", http.StatusBadRequest, w)
		return
	}

	resp, err := c.service.Update(brand, name, req)
	if err != nil {
		writeErr(err.Error(), http.StatusInternalServerError, w)
		return
	}
	j, _ := json.Marshal(AsDTO(*resp))
	result := string(j[:])

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, result)
}

// HandleSwitchAdd godoc
//
//	@Summary		Add new switch
//	@Description	Add a new switch and get resource address
//	@Tags			switches
//	@Produce		json
//	@Accept			json
//	@Param			newswitch	body		models.SwitchRequestBody	true	"Switch to add"
//	@Success		200			{object}	string
//	@Failure		500			{object}	models.APIError
//	@Failure		400			{object}	models.APIError
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
