package model

import (
	"github.com/google/uuid"
)

type Supplier struct {
	ID   uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Name string    `gorm:"type:varchar(255);not null" json:"name"`
}

func (m *Supplier) TableName() string {
	return "suppliers"
}
