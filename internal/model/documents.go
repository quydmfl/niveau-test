package model

import (
	"time"

	"github.com/google/uuid"
)

type Documents struct {
	ID       uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Filename string    `gorm:"type:varchar(255);not null" json:"filename"`
	Path     string    `gorm:"type:text;not null" json:"path"`

	ProductID *uuid.UUID `gorm:"type:uuid;default:null" json:"product_id"`
	Product   *Product   `gorm:"foreignKey:ProductID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"product,omitempty"`

	UploadedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"uploaded_at"`
}

func (f *Documents) TableName() string {
	return "documents"
}
