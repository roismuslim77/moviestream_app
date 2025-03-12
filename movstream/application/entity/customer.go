package entity

import (
	"github.com/golang-jwt/jwt/v4"
	"time"
)

type Customer struct {
	ID                int       `gorm:"column:id;type:int;primaryKey;autoIncrement:true;unique" json:"id"`
	Email             string    `gorm:"column:email;type:string;size:255;unique" json:"email"`
	FullName          string    `gorm:"column:name;type:string;size:255" json:"name"`
	BirthPlace        string    `gorm:"column:birth_place;type:string;size:255" json:"birth_place"`
	BirthDate         time.Time `gorm:"column:birth_date;" json:"birth_date"`
	IdentityPhotoLink string    `gorm:"column:identity_photo_link;type:string;size:255" json:"identity_photo_link"`
	Type              *string   `gorm:"column:type;type:string;size:100;" json:"type"`

	CreatedAt time.Time `gorm:"column:created_at;" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;" json:"updated_at"`

	// custom
	IsEmpty bool `gorm:"-" json:"-"`
}

func (t Customer) TableName() string {
	return "customers"
}

type Claims struct {
	CustomerId int    `json:"customer_id"`
	Type       string `json:"type"`
	jwt.RegisteredClaims
}
