package entity

import "time"

type MovieVote struct {
	ID         int `gorm:"column:id;type:int;primaryKey;autoIncrement:true;unique" json:"id"`
	MovieId    int `gorm:"column:movie_id;type:int" json:"movie_id"`
	CustomerId int `gorm:"column:customer_id;type:int" json:"customer_id"`

	CreatedAt time.Time `gorm:"column:created_at;" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;" json:"updated_at"`

	Movie Movie `json:"movies" gorm:"foreignKey:movie_id;references:ID"`

	// custom
	IsEmpty bool `gorm:"-" json:"-"`
}

func (t MovieVote) TableName() string {
	return "movie_votes"
}
