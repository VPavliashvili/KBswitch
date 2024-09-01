package switches_test

import (
	"encoding/json"
	"fmt"
	"kbswitch/internal/app/api/controllers/switches"
	"kbswitch/internal/core/common"
	"kbswitch/internal/core/switches/models"
	"net/http"
	"strings"
	"testing"
)

func intptr(x int) *int {
	return &x
}

type fakeWriter struct {
	input        string
	headerStatus int
}

// Header implements http.ResponseWriter.
func (w fakeWriter) Header() http.Header {
	panic("unimplemented")
}

func (w *fakeWriter) Write(p []byte) (int, error) {
	w.input = string(p[:])
	return 0, nil
}

func (w *fakeWriter) WriteHeader(statusCode int) {
	w.headerStatus = statusCode
}

type fakeService struct {
	pluralReturner     func() ([]models.Switch, *common.AppError)
	singleReturner     func(string, string) (*models.Switch, *common.AppError)
	addSwitchAction    func(reqbody models.SwitchRequestBody) (*int, *common.AppError)
	deleteSwitchAction func(string, string) *common.AppError
	updateSwitchAction func(string, string, models.SwitchRequestBody) (*models.Switch, error)
}

func (f fakeService) Update(brand, name string, m models.SwitchRequestBody) (*models.Switch, error) {
	return f.updateSwitchAction(brand, name, m)
}

func (f fakeService) Remove(brand, name string) *common.AppError {
	return f.deleteSwitchAction(brand, name)
}

func (f fakeService) AddNew(s models.SwitchRequestBody) (*int, *common.AppError) {
	return f.addSwitchAction(s)
}

func (f fakeService) GetAll() ([]models.Switch, *common.AppError) {
	return f.pluralReturner()
}

func (f fakeService) GetSingle(brand, name string) (*models.Switch, *common.AppError) {
	return f.singleReturner(brand, name)
}

func TestHandleSwitchUpdate(t *testing.T) {
	tcases := []struct {
		service  fakeService
		w        *fakeWriter
		req      *http.Request
		expected struct {
			data         string
			headerStatus int
		}
	}{
		{
			service: fakeService{},
			w:       &fakeWriter{},
			req: func() *http.Request {
				msg := `{"tst":"wrong"}`
				rq, _ := http.NewRequest("PATCH", "", strings.NewReader(msg))

				rq.SetPathValue("brand", "tst")
				rq.SetPathValue("name", "tstname")

				return rq
			}(),
			expected: struct {
				data         string
				headerStatus int
			}{
				data: common.APIError{
					Status:  http.StatusBadRequest,
					Message: "invalid request model",
				}.Error(),
				headerStatus: http.StatusBadRequest,
			},
		},
		{
			w: &fakeWriter{},
			req: func() *http.Request {
				rq := &http.Request{}

				rq.SetPathValue("brand", "tst")
				rq.SetPathValue("name", "tstname")

				return rq
			}(),
			service: fakeService{},
			expected: struct {
				data         string
				headerStatus int
			}{
				data: common.APIError{
					Status:  http.StatusBadRequest,
					Message: "request body is entirely missing/nil",
				}.Error(),
				headerStatus: http.StatusBadRequest,
			},
		},
		{
			service: fakeService{},
			w:       &fakeWriter{},
			req: func() *http.Request {
				rq := &http.Request{}
				rq.SetPathValue("brandnotset", "tst")
				rq.SetPathValue("namenotset", "tstname")

				return rq
			}(),
			expected: struct {
				data         string
				headerStatus int
			}{
				data: common.APIError{
					Message: "request parameters 'name' and 'brand' are missing",
					Status:  http.StatusBadRequest,
				}.Error(),
				headerStatus: http.StatusBadRequest,
			},
		},
		{
			service: fakeService{},
			w:       &fakeWriter{},
			req: func() *http.Request {
				rq := &http.Request{}
				rq.SetPathValue("brand", "tst")
				rq.SetPathValue("namenotset", "tstname")

				return rq
			}(),
			expected: struct {
				data         string
				headerStatus int
			}{
				data: common.APIError{
					Status:  http.StatusBadRequest,
					Message: "request parameter 'name' is missing",
				}.Error(),
				headerStatus: http.StatusBadRequest,
			},
		},
		{
			service: fakeService{},
			w:       &fakeWriter{},
			req: func() *http.Request {
				rq := &http.Request{}
				rq.SetPathValue("brandnotset", "tst")
				rq.SetPathValue("name", "tstname")

				return rq
			}(),
			expected: struct {
				data         string
				headerStatus int
			}{
				data: common.APIError{
					Status:  http.StatusBadRequest,
					Message: "request parameter 'brand' is missing",
				}.Error(),
				headerStatus: http.StatusBadRequest,
			},
		},
		{
			service: fakeService{},
			w:       &fakeWriter{},
			req: func() *http.Request {
				rq := &http.Request{}
				rq.SetPathValue("brandnotset", "tst")
				rq.SetPathValue("name", "tstname")

				return rq
			}(),
			expected: struct {
				data         string
				headerStatus int
			}{
				data: common.APIError{
					Status:  http.StatusBadRequest,
					Message: "request parameter 'brand' is missing",
				}.Error(),
				headerStatus: http.StatusBadRequest,
			},
		},
		{
			service: fakeService{updateSwitchAction: func(string, string, models.SwitchRequestBody) (*models.Switch, error) {
				return nil, fmt.Errorf("tst")
			}},
			w: &fakeWriter{},
			req: func() *http.Request {
				s := models.SwitchRequestBody{}
				j, _ := json.Marshal(s)
				rq, _ := http.NewRequest("PATCH", "", strings.NewReader(string(j[:])))

				rq.SetPathValue("brand", "tst")
				rq.SetPathValue("name", "tstname")

				return rq
			}(),
			expected: struct {
				data         string
				headerStatus int
			}{
				data: common.APIError{
					Status:  http.StatusInternalServerError,
					Message: "tst",
				}.Error(),
				headerStatus: http.StatusInternalServerError,
			},
		},
		{
			service: fakeService{updateSwitchAction: func(string, string, models.SwitchRequestBody) (*models.Switch, error) {
				return &models.Switch{Name: "test"}, nil
			}},
			w: &fakeWriter{},
			req: func() *http.Request {
				s := models.SwitchRequestBody{Name: "test"}
				j, _ := json.Marshal(s)
				rq, _ := http.NewRequest("PATCH", "", strings.NewReader(string(j[:])))

				rq.SetPathValue("brand", "tst")
				rq.SetPathValue("name", "oldname")

				return rq
			}(),
			expected: struct {
				data         string
				headerStatus int
			}{
				data: func() string {
					resp := models.Switch{
						Name: "test",
					}
					dto := switches.AsDTO(resp)
					j, _ := json.Marshal(dto)

					return string(j[:])
				}(),
				headerStatus: http.StatusOK,
			},
		},
	}

	for _, tc := range tcases {
		handler := switches.New(tc.service)
		handler.HandleSwitchUpdate(tc.w, tc.req)
		if tc.expected.data != tc.w.input {
			t.Errorf("HandleSwitchUpdate failed\nexpected %v\ngot %s", tc.expected.data, tc.w.input)
		}
		if tc.expected.headerStatus != tc.w.headerStatus {
			t.Errorf("HandleSwitchUpdate response header failed\nexpected %v\ngot  %v",
				tc.expected.headerStatus, tc.w.headerStatus)
		}
	}
}

func TestHandleSwitches(t *testing.T) {
	tcases := []struct {
		w        *fakeWriter
		service  fakeService
		expected struct {
			data         string
			headerStatus int
		}
	}{
		{
			w: &fakeWriter{},
			service: fakeService{
				pluralReturner: func() ([]models.Switch, *common.AppError) {
					return []models.Switch{
						{
							Lifespan:         100,
							OperatingForce:   50,
							ActivationTravel: 1.9,
							TotalTravel:      4.5,
						},
					}, nil
				},
			},
			expected: struct {
				data         string
				headerStatus int
			}{
				data: func() string {
					entities := []switches.SwitchDTO{
						{
							Lifespan:         "100M",
							OperatingForce:   "50gf",
							ActivationTravel: "1.9mm",
							TotalTravel:      "4.5mm",
						},
					}

					json, _ := json.Marshal(entities)
					return string(json[:])
				}(),
				headerStatus: http.StatusOK,
			},
		},
		{
			w: &fakeWriter{},
			service: fakeService{
				pluralReturner: func() ([]models.Switch, *common.AppError) {
					e := common.NewError(common.ErrInternalServer, "tst")
					return nil, &e
				},
			},
			expected: struct {
				data         string
				headerStatus int
			}{
				data: common.APIError{
					Status:  http.StatusInternalServerError,
					Message: "tst",
				}.Error(),

				headerStatus: http.StatusInternalServerError,
			},
		},
		{
			w: &fakeWriter{},
			service: fakeService{
				pluralReturner: func() ([]models.Switch, *common.AppError) {
					sws := make([]models.Switch, 0)
					return sws, nil
				},
			},
			expected: struct {
				data         string
				headerStatus int
			}{
				data:         "[]",
				headerStatus: http.StatusOK,
			},
		},
		{
			w: &fakeWriter{},
			service: fakeService{
				pluralReturner: func() ([]models.Switch, *common.AppError) {
					return nil, nil
				},
			},
			expected: struct {
				data         string
				headerStatus int
			}{
				data: common.APIError{
					Message: "collection got nil from a service",
					Status:  http.StatusInternalServerError,
				}.Error(),
				headerStatus: http.StatusInternalServerError,
			},
		},
	}

	for _, tc := range tcases {
		handler := switches.New(tc.service)
		handler.HandleSwitches(tc.w, nil)
		if tc.expected.data != tc.w.input {
			t.Errorf("HandleSwitches failed\nexpected %v\ngot %s", tc.expected.data, tc.w.input)
		}
		if tc.expected.headerStatus != tc.w.headerStatus {
			t.Errorf("HandleSwitches response header failed\nexpected %v\ngot  %v",
				tc.expected.headerStatus, tc.w.headerStatus)
		}
	}
}

func TestHandleSingleSwitch(t *testing.T) {
	tcases := []struct {
		service  fakeService
		w        *fakeWriter
		req      *http.Request
		expected struct {
			data         string
			headerStatus int
		}
	}{
		{
			service: fakeService{},
			w:       &fakeWriter{},
			req:     &http.Request{},
			expected: struct {
				data         string
				headerStatus int
			}{
				data: common.APIError{
					Message: "request parameters 'name' and 'brand' are missing",
					Status:  http.StatusBadRequest,
				}.Error(),
				headerStatus: http.StatusBadRequest,
			},
		},
		{
			service: fakeService{},
			w:       &fakeWriter{},
			req: func() *http.Request {
				rq := &http.Request{}
				rq.SetPathValue("name", "tst")

				return rq
			}(),
			expected: struct {
				data         string
				headerStatus int
			}{
				data: common.APIError{
					Message: "request parameter 'brand' is missing",
					Status:  http.StatusBadRequest,
				}.Error(),
				headerStatus: http.StatusBadRequest,
			},
		},
		{
			service: fakeService{},
			w:       &fakeWriter{},
			req: func() *http.Request {
				rq := &http.Request{}
				rq.SetPathValue("brand", "tst")

				return rq
			}(),
			expected: struct {
				data         string
				headerStatus int
			}{
				data: common.APIError{
					Status:  http.StatusBadRequest,
					Message: "request parameter 'name' is missing",
				}.Error(),
				headerStatus: http.StatusBadRequest,
			},
		},
		{
			service: fakeService{singleReturner: func(brand, name string) (*models.Switch, *common.AppError) {
				return nil, nil
			}},
			w: &fakeWriter{},
			req: func() *http.Request {
				rq := &http.Request{}
				rq.SetPathValue("brand", "tst")
				rq.SetPathValue("name", "tst")

				return rq
			}(),
			expected: struct {
				data         string
				headerStatus int
			}{
				data: common.APIError{
					Status:  http.StatusNotFound,
					Message: "no resource found for a given name and brand",
				}.Error(),
				headerStatus: http.StatusNotFound,
			},
		},
		{
			service: fakeService{singleReturner: func(brand, name string) (*models.Switch, *common.AppError) {
				e := common.NewError(common.ErrInternalServer, "tst")
				return nil, &e
			}},
			w: &fakeWriter{},
			req: func() *http.Request {
				rq := &http.Request{}
				rq.SetPathValue("brand", "tst")
				rq.SetPathValue("name", "tst")

				return rq
			}(),
			expected: struct {
				data         string
				headerStatus int
			}{
				data: common.APIError{
					Status:  http.StatusInternalServerError,
					Message: "tst",
				}.Error(),
				headerStatus: http.StatusInternalServerError,
			},
		},
		{
			service: fakeService{singleReturner: func(brand, name string) (*models.Switch, *common.AppError) {
				return &models.Switch{
					Lifespan:         100,
					OperatingForce:   50,
					ActivationTravel: 1.9,
					TotalTravel:      4.5,
				}, nil
			}},
			w: &fakeWriter{},
			req: func() *http.Request {
				rq := &http.Request{}
				rq.SetPathValue("brand", "tst")
				rq.SetPathValue("name", "tst")

				return rq
			}(),
			expected: struct {
				data         string
				headerStatus int
			}{
				data: func() string {
					dto := switches.SwitchDTO{
						Lifespan:         "100M",
						OperatingForce:   "50gf",
						ActivationTravel: "1.9mm",
						TotalTravel:      "4.5mm",
					}
					j, _ := json.Marshal(dto)

					return string(j[:])
				}(),
				headerStatus: http.StatusOK,
			},
		},
	}

	for _, tc := range tcases {
		handler := switches.New(tc.service)
		handler.HandleSingleSwitch(tc.w, tc.req)
		if tc.expected.data != tc.w.input {
			t.Errorf("HandleSingleSwitch failed\nexpected %v\ngot %s", tc.expected.data, tc.w.input)
		}
		if tc.expected.headerStatus != tc.w.headerStatus {
			t.Errorf("HandleSingleSwitch response header failed\nexpected %v\ngot  %v",
				tc.expected.headerStatus, tc.w.headerStatus)
		}
	}
}

func TestHandleSwitchRemove(t *testing.T) {
	tcases := []struct {
		service  fakeService
		w        *fakeWriter
		req      *http.Request
		expected struct {
			data         string
			headerStatus int
		}
	}{
		{
			service: fakeService{},
			w:       &fakeWriter{},
			req:     &http.Request{},
			expected: struct {
				data         string
				headerStatus int
			}{
				data: common.APIError{
					Message: "request parameters 'name' and 'brand' are missing",
					Status:  http.StatusBadRequest,
				}.Error(),
				headerStatus: http.StatusBadRequest,
			},
		},
		{
			service: fakeService{},
			w:       &fakeWriter{},
			req: func() *http.Request {
				rq := &http.Request{}
				rq.SetPathValue("brand", "tst")
				rq.SetPathValue("namenotset", "tstname")

				return rq
			}(),
			expected: struct {
				data         string
				headerStatus int
			}{
				data: common.APIError{
					Status:  http.StatusBadRequest,
					Message: "request parameter 'name' is missing",
				}.Error(),
				headerStatus: http.StatusBadRequest,
			},
		},
		{
			service: fakeService{},
			w:       &fakeWriter{},
			req: func() *http.Request {
				rq := &http.Request{}
				rq.SetPathValue("brandnotset", "tst")
				rq.SetPathValue("name", "tstname")

				return rq
			}(),
			expected: struct {
				data         string
				headerStatus int
			}{
				data: common.APIError{
					Status:  http.StatusBadRequest,
					Message: "request parameter 'brand' is missing",
				}.Error(),
				headerStatus: http.StatusBadRequest,
			},
		},
		{
			service: fakeService{deleteSwitchAction: func(brand, name string) *common.AppError {
				e := common.NewError(common.ErrInternalServer, "tst")
				return &e
			}},
			w: &fakeWriter{},
			req: func() *http.Request {
				rq := &http.Request{}
				rq.SetPathValue("brand", "tstbrand")
				rq.SetPathValue("name", "tstname")

				return rq
			}(),
			expected: struct {
				data         string
				headerStatus int
			}{
				data: common.APIError{
					Status:  http.StatusInternalServerError,
					Message: "tst",
				}.Error(),
				headerStatus: http.StatusInternalServerError,
			},
		},
		{
			service: fakeService{deleteSwitchAction: func(s1, s2 string) *common.AppError {
				return nil
			}},
			w: &fakeWriter{},
			req: func() *http.Request {
				rq := &http.Request{}
				rq.SetPathValue("brand", "tstbrand")
				rq.SetPathValue("name", "tstname")

				return rq
			}(),
			expected: struct {
				data         string
				headerStatus int
			}{
				data:         "",
				headerStatus: http.StatusNoContent,
			},
		},
	}

	for _, tc := range tcases {
		handler := switches.New(tc.service)
		handler.HandleSwitchRemove(tc.w, tc.req)
		if tc.expected.data != tc.w.input {
			t.Errorf("HandleSwitchRemove failed\nexpected %v\ngot %s", tc.expected.data, tc.w.input)
		}
		if tc.expected.headerStatus != tc.w.headerStatus {
			t.Errorf("HandleSwitchRemove response header failed\nexpected %v\ngot  %v",
				tc.expected.headerStatus, tc.w.headerStatus)
		}
	}
}

func TestHandleSwitchAdd(t *testing.T) {
	tcases := []struct {
		service  fakeService
		w        *fakeWriter
		req      *http.Request
		expected struct {
			data         string
			headerStatus int
		}
	}{
		{
			service: fakeService{},
			w:       &fakeWriter{},
			req: func() *http.Request {
				rq, _ := http.NewRequest("POST", "", strings.NewReader(""))

				return rq
			}(),
			expected: struct {
				data         string
				headerStatus int
			}{
				data: common.APIError{
					Status:  http.StatusBadRequest,
					Message: "invalid request model", // entirelly missing body
				}.Error(),
				headerStatus: http.StatusBadRequest,
			},
		},
		{
			service: fakeService{},
			w:       &fakeWriter{},
			req: func() *http.Request {
				msg := `{"tst":"value"}`
				rq, _ := http.NewRequest("POST", "", strings.NewReader(msg))

				return rq
			}(),
			expected: struct {
				data         string
				headerStatus int
			}{
				data: common.APIError{
					Status:  http.StatusBadRequest,
					Message: "invalid request model", // wrong structure request body
				}.Error(),
				headerStatus: http.StatusBadRequest,
			},
		},
		{
			service: fakeService{
				addSwitchAction: func(reqbody models.SwitchRequestBody) (*int, *common.AppError) {
					e := common.NewError(common.ErrBadRequest, "tst")
					return nil, &e
				},
			},
			w: &fakeWriter{},
			req: func() *http.Request {
				s := models.SwitchRequestBody{}
				j, _ := json.Marshal(s)
				msg := string(j[:])
				rq, _ := http.NewRequest("POST", "", strings.NewReader(msg))

				return rq
			}(),
			expected: struct {
				data         string
				headerStatus int
			}{
				data: common.APIError{
					Status:  http.StatusBadRequest,
					Message: "tst",
				}.Error(),
				headerStatus: http.StatusBadRequest,
			},
		},
		{
			service: fakeService{
				addSwitchAction: func(reqbody models.SwitchRequestBody) (*int, *common.AppError) {
					return intptr(123), nil
				},
			},
			w: &fakeWriter{},
			req: func() *http.Request {
				s := models.SwitchRequestBody{}
				j, _ := json.Marshal(s)
				msg := string(j[:])
				rq, _ := http.NewRequest("POST", "", strings.NewReader(msg))
				rq.Host = "tsthost:tstport/"
				rq.URL.Path = "api/switches"

				return rq
			}(),
			expected: struct {
				data         string
				headerStatus int
			}{
				data: func() string {
					return "tsthost:tstport/api/switches/123"
				}(),
				headerStatus: http.StatusCreated,
			},
		},
	}

	for _, tc := range tcases {
		handler := switches.New(tc.service)
		handler.HandleSwitchAdd(tc.w, tc.req)
		if tc.expected.data != tc.w.input {
			t.Errorf("HandleSwitchAdd failed\nexpected %v\ngot %s", tc.expected.data, tc.w.input)
		}
		if tc.expected.headerStatus != tc.w.headerStatus {
			t.Errorf("HandleSwitchAdd response header failed\nexpected %v\ngot  %v",
				tc.expected.headerStatus, tc.w.headerStatus)
		}
	}
}
