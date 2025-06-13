package models

import (
	"time"

	"github.com/lib/pq"
)

type RecordType string

const (
	// TypeExpense represents an expense record
	TypeExpense RecordType = "expense"
	// TypeIncome represents an income record
	TypeIncome RecordType = "income"
)

type Record struct {
	ID          string         `json:"id" gorm:"type:string;default:gen_random_uuid();primaryKey"`
	Name        string         `json:"name" gorm:"not null"`
	Description string         `json:"description" gorm:"type:text;default:''"`
	Date        time.Time      `json:"date"`
	Tags        pq.StringArray `json:"tags" gorm:"type:text[]"`
	Type        RecordType     `json:"type" gorm:"type:record_type;not null;index"`
	Amount      int64          `json:"amount" gorm:"not null"`
	CreatedAt   time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt   *time.Time     `json:"deleted_at,omitempty" gorm:"index"`
	UserID      string         `json:"user_id" gorm:"not null;index;constraint:OnDelete:CASCADE"`
	User        User           `json:"user" gorm:"foreignKey:UserID"` // Foreign key relationship
}

type CreateRecordPayload struct {
	Name        string         `json:"name" binding:"required"`
	Description string         `json:"description"`
	Date        time.Time      `json:"date" binding:"required"`
	Tags        pq.StringArray `json:"tags" gorm:"type:text[]"`
	Type        RecordType     `json:"type" binding:"required"`
	Amount      int64          `json:"amount" binding:"required"`
}
