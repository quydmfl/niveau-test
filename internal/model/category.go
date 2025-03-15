package model

import (
	"time"

	"github.com/google/uuid"
)

type Category struct {
	ID     uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Name   string    `gorm:"type:varchar(255);not null" json:"name"`
	Status string    `gorm:"type:varchar(25);not null" json:"status"`

	CreatedAt time.Time `gorm:"type:timestamp;not null;default:now()" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamp;not null;default:now()" json:"updated_at"`

	// Relationship
	Products []Product `gorm:"foreignKey:CategoryID" json:"products,omitempty"`
}

func (m *Category) TableName() string {
	return "product_categories"
}
