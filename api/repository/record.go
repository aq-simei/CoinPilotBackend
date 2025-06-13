package repository

import (
	"github.com/aq-simei/coin-pilot/api/models"
	"gorm.io/gorm"
)

type RecordRepository interface {
	GetRecords(userID string) ([]models.Record, error)
	CreateRecord(record models.CreateRecordPayload, userID string) (*models.CreateRecordPayload, error)
	UpdateRecord(id string, record models.Record) (*models.Record, error)
	DeleteRecord(id string) error
}

type RecordRepositoryImpl struct {
	db *gorm.DB
}

func NewRecordRepository(db *gorm.DB) RecordRepository {
	return &RecordRepositoryImpl{db: db}
}

func (r *RecordRepositoryImpl) GetRecords(userID string) ([]models.Record, error) {
	var records []models.Record
	result := r.db.Where("user_id = ?", userID).Find(&records)
	if result.Error != nil {
		return nil, result.Error
	}
	return records, nil
}

func (r *RecordRepositoryImpl) CreateRecord(record models.CreateRecordPayload, userID string) (*models.CreateRecordPayload, error) {
	// Map the CreateRecordPayload to a Record
	newRecord := &models.Record{
		Name:        record.Name,
		Date:        record.Date,
		Description: record.Description,
		Tags:        record.Tags,
		Type:        record.Type,
		Amount:      record.Amount,
		UserID:      userID,
	}
	result := r.db.Create(newRecord)
	if result.Error != nil {
		return nil, result.Error
	}
	return &record, nil
}

func (r *RecordRepositoryImpl) UpdateRecord(id string, record models.Record) (*models.Record, error) {
	existingRecord := &models.Record{}
	result := r.db.First(existingRecord, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}

	// Update the existing record with new values
	existingRecord.Type = record.Type
	existingRecord.Amount = record.Amount
	existingRecord.Description = record.Description

	if err := r.db.Save(existingRecord).Error; err != nil {
		return nil, err
	}

	return existingRecord, nil
}

func (r *RecordRepositoryImpl) DeleteRecord(id string) error {
	result := r.db.Delete(&models.Record{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
