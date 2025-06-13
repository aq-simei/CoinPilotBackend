package service

import (
	"github.com/aq-simei/coin-pilot/api/models"
	"github.com/aq-simei/coin-pilot/api/repository"
	"github.com/gin-gonic/gin"
)

type RecordService interface {
	GetRecords(ctx *gin.Context, userID string) ([]models.Record, error)
	CreateRecord(ctx *gin.Context, record models.CreateRecordPayload, userID string) (*models.CreateRecordPayload, error)
	UpdateRecord(ctx *gin.Context, id string, record models.Record) (*models.Record, error)
	DeleteRecord(ctx *gin.Context, id string) error
}

type RecordServiceImpl struct {
	repository repository.RecordRepository
}

func NewRecordService(repository repository.RecordRepository) RecordService {
	return &RecordServiceImpl{
		repository: repository,
	}
}

func (s *RecordServiceImpl) GetRecords(ctx *gin.Context, userID string) ([]models.Record, error) {
	records, err := s.repository.GetRecords(userID)
	if err != nil {
		return nil, err
	}
	return records, nil
}

func (s *RecordServiceImpl) CreateRecord(ctx *gin.Context, record models.CreateRecordPayload, userID string) (*models.CreateRecordPayload, error) {
	createdRecord, err := s.repository.CreateRecord(record, userID)
	if err != nil {
		return nil, err
	}
	return createdRecord, nil
}

func (s *RecordServiceImpl) UpdateRecord(ctx *gin.Context, id string, record models.Record) (*models.Record, error) {
	updatedRecord, err := s.repository.UpdateRecord(id, record)
	if err != nil {
		return nil, err
	}
	return updatedRecord, nil
}

func (s *RecordServiceImpl) DeleteRecord(ctx *gin.Context, id string) error {
	err := s.repository.DeleteRecord(id)
	if err != nil {
		return err
	}
	return nil
}
