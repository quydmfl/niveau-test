package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/quydmfl/niveau-test/internal/helper"
	"gorm.io/gorm"
)

type Product struct {
	ID         uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Reference  string    `gorm:"type:varchar(50);not null;unique" json:"reference"`
	Name       string    `gorm:"type:varchar(255);not null" json:"name"`
	DateAdded  time.Time `gorm:"column:added_date;type:date;default:CURRENT_DATE" json:"added_date"`
	Status     string    `gorm:"type:varchar(50);check:status IN ('Available', 'Out of Stock', 'On Order')" json:"status"`
	CategoryID uuid.UUID `gorm:"type:uuid;not null" json:"category_id"`
	Price      float64   `gorm:"type:numeric(10,2);default:0" json:"price"`
	StockCity  string    `gorm:"type:varchar(100);default:null" json:"stock_city"`
	SupplierID uuid.UUID `gorm:"type:uuid;default:null" json:"supplier_id"`
	Quantity   int       `gorm:"type:int;default:0" json:"quantity"`

	// Relationship
	Supplier  Supplier    `gorm:"foreignKey:SupplierID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"supplier,omitempty"`
	Documents []Documents `gorm:"foreignKey:ProductID" json:"documents,omitempty"`
	Category  Category    `gorm:"foreignKey:CategoryID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"category"`
}

func (p *Product) TableName() string {
	return "products"
}

func (p *Product) BeforeCreate(tx *gorm.DB) error {
	if p.Reference == "" {
		for {
			p.Reference = helper.GenerateProductReference()

			var count int64
			tx.Model(&Product{}).Where("reference = ?", p.Reference).Count(&count)
			if count == 0 {
				break // Exit loop if unique
			}
		}
	}

	return nil
}
