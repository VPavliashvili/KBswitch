package switches

import (
	"context"
	"fmt"
	"kbswitch/internal/core/common"
	"kbswitch/internal/core/common/logging"
	"kbswitch/internal/core/switches"
	"kbswitch/internal/core/switches/models"
)

var (
	ErrNoSwitch      = common.NewError(common.ErrNotFound, "resource with given brand and name not found")
	ErrAlreadyExists = common.NewError(common.ErrBadRequest, "switch with given brand and name already exist")
	ErrErrorMissing  = common.NewError(common.ErrInternalServer, "no error returned when response was missing")
)

func New(logger logging.Logger, repo switches.Repo) switches.Service {
	return service{
		repo:   repo,
		logger: logger,
	}
}

type service struct {
	repo   switches.Repo
	logger logging.Logger
}

func (s service) AddNew(ctx context.Context, reqbody models.SwitchRequestBody) (*int, *common.AppError) {
	switchID, err := s.repo.GetID(ctx, reqbody.Brand, reqbody.Name)
	if err != nil {
		s.logger.LogError(err.Error())
		return nil, common.Wrap(err)
	}
	if switchID != nil {
		s.logger.LogError(fmt.Sprintf("switch on given ID %v, already exist", switchID))
		return nil, &ErrAlreadyExists
	}

	entity := models.SwitchEntity{
		Manufacturer:     reqbody.Brand,
		ActuationType:    reqbody.ActuationType,
		Lifespan:         reqbody.Lifespan,
		Model:            reqbody.Name,
		Image:            []byte(reqbody.Image),
		OperatingForce:   reqbody.OperatingForce,
		ActivationTravel: reqbody.ActivationTravel,
		TotalTravel:      reqbody.TotalTravel,
		SoundProfile:     reqbody.SoundProfile,
		TriggerMethod:    reqbody.TriggerMethod,
		Profile:          reqbody.Profile,
	}

	resp, err := s.repo.AddNew(ctx, entity)
	if err != nil {
		s.logger.LogError(err.Error())
		return nil, common.Wrap(err)
	}

	s.logger.LogTrace(fmt.Sprintf("result is %v", resp))
	return resp, nil
}

func (s service) Update(ctx context.Context, brand, name string, body models.SwitchRequestBody) (*models.Switch, *common.AppError) {
	switchID, err := s.repo.GetID(ctx, brand, name)
	if err != nil {
		s.logger.LogError(err.Error())
		return nil, common.Wrap(err)
	}
	if switchID == nil {
		s.logger.LogError("switchID from repo was nil")
		return nil, &ErrNoSwitch
	}

	entity := models.SwitchEntity{
		Manufacturer:     body.Brand,
		ActuationType:    body.ActuationType,
		Lifespan:         body.Lifespan,
		Model:            body.Name,
		Image:            []byte(body.Image),
		OperatingForce:   body.OperatingForce,
		ActivationTravel: body.ActivationTravel,
		TotalTravel:      body.TotalTravel,
		SoundProfile:     body.SoundProfile,
		TriggerMethod:    body.TriggerMethod,
		Profile:          body.Profile,
	}
	resp, err := s.repo.Update(ctx, *switchID, entity)
	if err != nil {
		s.logger.LogError(err.Error())
		return nil, common.Wrap(err)
	}
	if resp == nil {
		return nil, nil
	}

	res := models.Switch{
		Brand:            resp.Manufacturer,
		ActuationType:    resp.ActuationType,
		Lifespan:         resp.Lifespan,
		Name:             resp.Model,
		Image:            string(resp.Image[:]),
		OperatingForce:   resp.OperatingForce,
		ActivationTravel: resp.ActivationTravel,
		TotalTravel:      resp.TotalTravel,
		SoundProfile:     resp.SoundProfile,
		TriggerMethod:    resp.TriggerMethod,
		Profile:          resp.Profile,
	}
	s.logger.LogTrace(fmt.Sprintf("result is %v", res))

	return &res, nil
}

func (s service) Remove(ctx context.Context, brand, name string) *common.AppError {
	switchID, err := s.repo.GetID(ctx, brand, name)
	if err != nil {
		return common.Wrap(err)
	}
	if switchID == nil {
		return &ErrNoSwitch
	}

	err = s.repo.Remove(ctx, *switchID)
	if err != nil {
		return common.Wrap(err)
	}

	return nil
}

func (s service) GetAll(ctx context.Context) ([]models.Switch, *common.AppError) {
	res := []models.Switch{}
	resp, err := s.repo.GetAll(ctx)
	if err != nil {
		s.logger.LogError(err.Error())
		return res, common.Wrap(err)
	}

	for _, item := range resp {
		s := models.Switch{
			Brand:            item.Manufacturer,
			ActuationType:    item.ActuationType,
			Lifespan:         item.Lifespan,
			Name:             item.Model,
			Image:            string(item.Image[:]),
			OperatingForce:   item.OperatingForce,
			ActivationTravel: item.ActivationTravel,
			TotalTravel:      item.TotalTravel,
			SoundProfile:     item.SoundProfile,
			TriggerMethod:    item.TriggerMethod,
			Profile:          item.Profile,
		}

		res = append(res, s)
	}
	s.logger.LogTrace(fmt.Sprintf("result is %v", res))

	return res, nil
}

func (s service) GetSingle(ctx context.Context, brand, name string) (*models.Switch, *common.AppError) {
	switchID, err := s.repo.GetID(ctx, brand, name)

	if err != nil {
		s.logger.LogError(err.Error())
		return nil, common.Wrap(err)
	}
	if switchID == nil {
		s.logger.LogError("switchID from repo was nil")
		return nil, &ErrNoSwitch
	}

	resp, err := s.repo.GetSingle(ctx, *switchID)
	if err != nil {
		s.logger.LogError(err.Error())
		return nil, common.Wrap(err)
	}
	if resp == nil {
		s.logger.LogError("response from repo was nil")
		return nil, &ErrErrorMissing
	}
	res := models.Switch{
		Brand:            resp.Manufacturer,
		ActuationType:    resp.ActuationType,
		Lifespan:         resp.Lifespan,
		Name:             resp.Model,
		Image:            string(resp.Image[:]),
		OperatingForce:   resp.OperatingForce,
		ActivationTravel: resp.ActivationTravel,
		TotalTravel:      resp.TotalTravel,
		SoundProfile:     resp.SoundProfile,
		TriggerMethod:    resp.TriggerMethod,
		Profile:          resp.Profile,
	}
	s.logger.LogTrace(fmt.Sprintf("result is %v", res))

	return &res, nil
}
