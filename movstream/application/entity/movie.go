package entity

import "time"

type Movie struct {
	ID          int    `gorm:"column:id;type:int;primaryKey;autoIncrement:true;unique" json:"id"`
	Title       string `gorm:"column:title;type:string;size:255" json:"title"`
	Description string `gorm:"column:description;type:string;size:255" json:"description"`
	Duration    int64  `gorm:"column:duration;type:int;" json:"duration"`
	Artist      string `gorm:"column:artist;type:string;size:100" json:"artist"`
	Genre       string `gorm:"column:genre;type:string;size:100" json:"genre"`
	WatchUrl    string `gorm:"column:watch_url;type:string;size:255" json:"watch_url"`

	CreatedAt time.Time `gorm:"column:created_at;" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;" json:"updated_at"`
	CreatedId int       `gorm:"column:created_id;type:int" json:"created_id"`

	// custom
	IsEmpty bool `gorm:"-" json:"-"`
}

func (t Movie) TableName() string {
	return "movies"
}

type AllMovie struct {
	ID          int    `gorm:"column:id;type:int;primaryKey;autoIncrement:true;unique" json:"id"`
	Title       string `gorm:"column:title;type:string;size:255" json:"title"`
	Description string `gorm:"column:description;type:string;size:255" json:"description"`
	Duration    int64  `gorm:"column:duration;type:int;" json:"duration"`
	Artist      string `gorm:"column:artist;type:string;size:100" json:"artist"`
	Genre       string `gorm:"column:genre;type:string;size:100" json:"genre"`
	WatchUrl    string `gorm:"column:watch_url;type:string;size:255" json:"watch_url"`

	CreatedAt time.Time `gorm:"column:created_at;" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;" json:"updated_at"`

	IsViewed *int `gorm:"column:is_viewed;type:int;" json:"is_viewed"`
	IsVoted  *int `gorm:"column:is_voted;type:int;" json:"is_voted"`
}
