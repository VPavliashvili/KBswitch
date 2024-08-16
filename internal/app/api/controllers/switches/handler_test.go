package switches_test

import (
	"encoding/json"
	"fmt"
	"kbswitch/internal/app/api/controllers/switches"
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
	pluralReturner  func() ([]models.Switch, error)
	singleReturner  func(id int) (*models.Switch, error)
	addSwitchAction func(reqbody models.SwitchRequestBody) (*int, error)
}

func (f fakeService) AddNew(s models.SwitchRequestBody) (*int, error) {
	return f.addSwitchAction(s)
}

func (f fakeService) GetAll() ([]models.Switch, error) {
	return f.pluralReturner()
}

func (f fakeService) GetByID(id int) (*models.Switch, error) {
	return f.singleReturner(id)
}

func TestHandleSwitches(t *testing.T) {
	tcases := []struct {
		w        *fakeWriter
		repo     *fakeService
		expected struct {
			data         string
			headerStatus int
		}
	}{
		{
			w: &fakeWriter{},
			repo: &fakeService{
				pluralReturner: func() ([]models.Switch, error) {
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
			repo: &fakeService{
				pluralReturner: func() ([]models.Switch, error) {
					return nil, fmt.Errorf("tst")
				},
			},
			expected: struct {
				data         string
				headerStatus int
			}{
				data: models.APIError{
					Status:  http.StatusInternalServerError,
					Message: "tst",
				}.Error(),

				headerStatus: http.StatusInternalServerError,
			},
		},
		{
			w: &fakeWriter{},
			repo: &fakeService{
				pluralReturner: func() ([]models.Switch, error) {
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
			repo: &fakeService{
				pluralReturner: func() ([]models.Switch, error) {
					return nil, nil
				},
			},
			expected: struct {
				data         string
				headerStatus int
			}{
				data: models.APIError{
					Message: "collection got nil from a repo",
					Status:  http.StatusInternalServerError,
				}.Error(),
				headerStatus: http.StatusInternalServerError,
			},
		},
	}

	for _, tc := range tcases {
		handler := switches.New(tc.repo)
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

func TestHandleSwitchByID(t *testing.T) {
	tcases := []struct {
		repo     fakeService
		w        *fakeWriter
		req      *http.Request
		expected struct {
			data         string
			headerStatus int
		}
	}{
		{
			repo: fakeService{},
			w:    &fakeWriter{},
			req:  &http.Request{},
			expected: struct {
				data         string
				headerStatus int
			}{
				data: models.APIError{
					Message: "request parameter is missing",
					Status:  http.StatusBadRequest,
				}.Error(),
				headerStatus: http.StatusBadRequest,
			},
		},
		{
			repo: fakeService{},
			w:    &fakeWriter{},
			req: func() *http.Request {
				rq := &http.Request{}
				rq.SetPathValue("id", "wrongtype")

				return rq
			}(),
			expected: struct {
				data         string
				headerStatus int
			}{
				data: models.APIError{
					Status:  http.StatusBadRequest,
					Message: "request parameter should be int",
				}.Error(),
				headerStatus: http.StatusBadRequest,
			},
		},
		{
			repo: fakeService{singleReturner: func(id int) (*models.Switch, error) {
				return nil, nil
			}},
			w: &fakeWriter{},
			req: func() *http.Request {
				rq := &http.Request{}
				rq.SetPathValue("id", "123")

				return rq
			}(),
			expected: struct {
				data         string
				headerStatus int
			}{
				data: models.APIError{
					Status:  http.StatusNotFound,
					Message: "no resource found for a given id",
				}.Error(),
				headerStatus: http.StatusNotFound,
			},
		},
		{
			repo: fakeService{singleReturner: func(id int) (*models.Switch, error) {
				return nil, fmt.Errorf("tst")
			}},
			w: &fakeWriter{},
			req: func() *http.Request {
				rq := &http.Request{}
				rq.SetPathValue("id", "123")

				return rq
			}(),
			expected: struct {
				data         string
				headerStatus int
			}{
				data: models.APIError{
					Status:  http.StatusInternalServerError,
					Message: "tst",
				}.Error(),
				headerStatus: http.StatusInternalServerError,
			},
		},
		{
			repo: fakeService{singleReturner: func(id int) (*models.Switch, error) {
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
				rq.SetPathValue("id", "123")

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
		handler := switches.New(tc.repo)
		handler.HandleSwitchByID(tc.w, tc.req)
		if tc.expected.data != tc.w.input {
			t.Errorf("HandleSwitchByID failed\nexpected %v\ngot %s", tc.expected.data, tc.w.input)
		}
		if tc.expected.headerStatus != tc.w.headerStatus {
			t.Errorf("HandleSwitchByID response header failed\nexpected %v\ngot  %v",
				tc.expected.headerStatus, tc.w.headerStatus)
		}
	}
}

func TestHandleSwitchPatch(t *testing.T) {
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
				data: models.APIError{
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
				data: models.APIError{
					Status:  http.StatusBadRequest,
					Message: "invalid request model", // wrong structure request body
				}.Error(),
				headerStatus: http.StatusBadRequest,
			},
		},
		{
			service: fakeService{
				addSwitchAction: func(reqbody models.SwitchRequestBody) (*int, error) {
					return nil, fmt.Errorf("tst")
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
				data: models.APIError{
					Status:  http.StatusInternalServerError,
					Message: "tst",
				}.Error(),
				headerStatus: http.StatusInternalServerError,
			},
		},
		{
			service: fakeService{
				addSwitchAction: func(reqbody models.SwitchRequestBody) (*int, error) {
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
