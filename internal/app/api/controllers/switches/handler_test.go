package switches_test

import (
	"encoding/json"
	"fmt"
	"kbswitch/internal/app/api/controllers/switches"
	"kbswitch/internal/core/models"
	"net/http"
	"testing"
)

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

type fakeRepo struct {
	pluralReturner func() ([]models.SwitchEntity, error)
}

func (r fakeRepo) GetAll() ([]models.SwitchEntity, error) {
	return r.pluralReturner()
}

func TestHandleSwitches(t *testing.T) {
	tcases := []struct {
		w        *fakeWriter
		repo     *fakeRepo
		expected struct {
			data         string
			headerStatus int
		}
	}{
		{
			w: &fakeWriter{},
			repo: &fakeRepo{
				pluralReturner: func() ([]models.SwitchEntity, error) {
					return []models.SwitchEntity{
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
			repo: &fakeRepo{
				pluralReturner: func() ([]models.SwitchEntity, error) {
					return nil, fmt.Errorf("tst")
				},
			},
			expected: struct {
				data         string
				headerStatus int
			}{
				data:         "tst",
				headerStatus: http.StatusInternalServerError,
			},
		},
		{
			w: &fakeWriter{},
			repo: &fakeRepo{
				pluralReturner: func() ([]models.SwitchEntity, error) {
					entities := make([]models.SwitchEntity, 0)
					return entities, nil
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
			repo: &fakeRepo{
				pluralReturner: func() ([]models.SwitchEntity, error) {
					return nil, nil
				},
			},
			expected: struct {
				data         string
				headerStatus int
			}{
				data:         "collection got nil from a repo",
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
