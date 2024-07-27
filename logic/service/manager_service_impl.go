package service

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"student_manage/logic/data/request"
	"student_manage/logic/data/response"
	"student_manage/logic/model"
	"student_manage/logic/repository"
	myerrors "student_manage/utils/errors"
)

type ManagerServiceImpl struct {
	ManagerRepository repository.ManagerRepository
	Validate          *validator.Validate
}

func NewManagerServiceImpl(repository repository.ManagerRepository, validate *validator.Validate) *ManagerServiceImpl {
	return &ManagerServiceImpl{
		ManagerRepository: repository,
		Validate:          validate,
	}
}

func (r *ManagerServiceImpl) Create(manager request.CreateManagerRequest) error {
	err := r.Validate.Struct(manager)
	if err != nil {
		return myerrors.ParamErr{Err: fmt.Errorf("service: 创建管理员校验参数失败 -> %w", err)}
	}
	managerModel := model.Manager{
		Level:    manager.Level,
		Name:     manager.Name,
		Account:  manager.Account,
		Password: manager.Password,
	}
	err = r.ManagerRepository.Create(managerModel)
	if err != nil {
		return myerrors.DbErr{Err: fmt.Errorf("service: 创建管理员失败 -> %w", err)}
	}
	return nil
}

func (r *ManagerServiceImpl) Update(manager request.UpdateManagerRequest) error {
	managerData, err := r.ManagerRepository.GetById(manager.Id)
	if err != nil {
		return myerrors.DbErr{Err: fmt.Errorf("service: 更新管理员时查找管理员失败 -> %w", err)}
	}
	managerData.Level = manager.Level
	managerData.Name = manager.Name
	managerData.Account = manager.Account
	managerData.Password = manager.Password
	err = r.ManagerRepository.Update(managerData)
	if err != nil {
		return myerrors.DbErr{Err: fmt.Errorf("service: 更新管理员失败 -> %w", err)}
	}
	return nil
}

func (r *ManagerServiceImpl) Delete(id int) error {
	err := r.ManagerRepository.Delete(id)
	if err != nil {
		return myerrors.DbErr{Err: fmt.Errorf("service: 删除管理员失败 -> %w", err)}
	}
	return nil
}

func (r *ManagerServiceImpl) GetById(id int) (response.ManagerResponse, error) {
	managerData, err := r.ManagerRepository.GetById(id)
	if err != nil {
		return response.ManagerResponse{}, myerrors.DbErr{Err: fmt.Errorf("service: 查找管理员id:%d失败 -> %w", id, err)}
	}
	managerResponse := response.ManagerResponse{
		Level:    managerData.Level,
		Name:     managerData.Name,
		Account:  managerData.Account,
		Password: managerData.Password,
	}
	return managerResponse, nil
}

func (r *ManagerServiceImpl) GetAll() ([]response.ManagerResponse, error) {
	result, err := r.ManagerRepository.GetAll()
	if err != nil {
		return []response.ManagerResponse{}, myerrors.DbErr{Err: fmt.Errorf("service: 查找全部管理员失败 -> %w", err)}
	}
	var managers []response.ManagerResponse
	for _, manager := range result {
		manager := response.ManagerResponse{
			Id:       manager.Id,
			Level:    manager.Level,
			Name:     manager.Name,
			Account:  manager.Account,
			Password: manager.Password,
		}
		managers = append(managers, manager)
	}
	return managers, nil
}