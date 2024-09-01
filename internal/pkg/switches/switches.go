package switches

import (
	"kbswitch/internal/core/common"
	"kbswitch/internal/core/switches"
	"kbswitch/internal/core/switches/models"
)

var (
	ErrNoSwitch      = common.NewError(common.ErrNotFound, "resource with given brand and name not found")
	ErrAlreadyExists = common.NewError(common.ErrBadRequest, "switch with given brand and name already exist")
	ErrErrorMissing  = common.NewError(common.ErrInternalServer, "no error returned when response was missing")
)

func New(repo switches.Repo) switches.Service {
	return service{repo: repo}
}

type service struct {
	repo switches.Repo
}

func (s service) AddNew(reqbody models.SwitchRequestBody) (*int, *common.AppError) {
	switchID, err := s.repo.GetID(reqbody.Brand, reqbody.Name)
	if err != nil {
		return nil, common.Wrap(err)
	}
	if switchID != nil {
		return nil, &ErrAlreadyExists
	}

	entity := models.SwitchEntity(reqbody)

	resp, err := s.repo.AddNew(entity)
	if err != nil {
		return nil, common.Wrap(err)
	}

	return resp, nil
}

func (s service) Update(brand, name string, body models.SwitchRequestBody) (*models.Switch, *common.AppError) {
	switchID, err := s.repo.GetID(brand, name)
	if err != nil {
		return nil, common.Wrap(err)
	}
	if switchID == nil {
		return nil, &ErrNoSwitch
	}

	entity := models.SwitchEntity(body)
	resp, err := s.repo.Update(*switchID, entity)
	if err != nil {
		return nil, common.Wrap(err)
	}
	if resp == nil {
		return nil, nil
	}

	res := models.Switch(*resp)

	return &res, nil
}

func (s service) Remove(brand, name string) *common.AppError {
	switchID, err := s.repo.GetID(brand, name)
	if err != nil {
		return common.Wrap(err)
	}
	if switchID == nil {
		return &ErrNoSwitch
	}

	err = s.repo.Remove(*switchID)
	if err != nil {
		return common.Wrap(err)
	}

	return nil
}

func (s service) GetAll() ([]models.Switch, *common.AppError) {
	res := []models.Switch{}
	resp, err := s.repo.GetAll()

	for _, item := range resp {
		s := models.Switch(item)

		res = append(res, s)
	}
	if err != nil {
		return res, common.Wrap(err)
	}

	return res, nil
}

func (s service) GetSingle(brand, name string) (*models.Switch, *common.AppError) {
	switchID, err := s.repo.GetID(brand, name)

	if err != nil {
		return nil, common.Wrap(err)
	}
	if switchID == nil {
		return nil, &ErrNoSwitch
	}

	resp, err := s.repo.GetSingle(*switchID)
	if err != nil {
		return nil, common.Wrap(err)
	}
	if resp == nil {
		return nil, &ErrErrorMissing
	}
	res := models.Switch(*resp)

	return &res, nil
}
