package models

import (
	"time"
)

// User table
type User struct {
	UserID           uint      `gorm:"primaryKey"`
	Nickname         string    `gorm:"size:30;not null"`
	Email            string    `gorm:"size:50;not null"`
	DiscordID        string    `gorm:"size:25"`
	Password         string    `gorm:"size:30;not null"`
	RegistrationDate time.Time `gorm:"default:CURRENT_DATE"`
	EventCount       int       `gorm:"default:0"`
	GuideCount       int       `gorm:"default:0"`
	IsModerator      bool      `gorm:"default:false"`
	IsVeteran        bool      `gorm:"default:false"`

	// Relationships (Has Many)
	Events []Event `gorm:"foreignKey:UserID"`
	Guides []Guide `gorm:"foreignKey:UserID"`
}

func (User) TableName() string { return "users" }

// Character table
type Character struct {
	CharacterID uint   `gorm:"primaryKey"`
	Name        string `gorm:"size:20;not null"`
	Description string `gorm:"size:255;not null"`
	ImageURL    string `gorm:"size:255"`
	Type        string `gorm:"size:15;not null"`

	// Relationships (Has Many)
	Guides []Guide `gorm:"foreignKey:CharacterID"`
}

func (Character) TableName() string { return "characters" }

// Guide table
type Guide struct {
	GuideID      uint      `gorm:"primaryKey"`
	Title        string    `gorm:"size:20;not null"`
	BannerURL    string    `gorm:"size:255"`
	Type         string    `gorm:"size:15"`
	Description  string    `gorm:"size:50;not null"`
	Link         string    `gorm:"size:2083"`
	CreationDate time.Time `gorm:"default:CURRENT_DATE"`
	Likes        int       `gorm:"default:0"`
	Dislikes     int       `gorm:"default:0"`

	// Foreign Keys (Obrigatórias)
	UserID      uint `gorm:"not null"`
	CharacterID uint `gorm:"not null"`

	// Relationships (Belongs To)
	User      User      `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Character Character `gorm:"foreignKey:CharacterID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (Guide) TableName() string { return "guides" }

// Event table
type Event struct {
	EventID     uint      `gorm:"primaryKey"`
	Title       string    `gorm:"size:75;not null"`
	Description string    `gorm:"size:255;not null"`
	BannerURL   string    `gorm:"size:255"`
	Day         time.Time `gorm:"type:date;not null"`

	// Foreign Key
	UserID uint `gorm:"not null"`

	// Relationships
	User       User        `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	LastEvents []LastEvent `gorm:"foreignKey:EventID"`
}

func (Event) TableName() string { return "events" }

// LastEvent table
type LastEvent struct {
	LastEventID uint   `gorm:"primaryKey"`
	Title       string `gorm:"size:75;not null"`

	// Foreign Key
	EventID uint `gorm:"not null"`

	// Relationships
	Event Event `gorm:"foreignKey:EventID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (LastEvent) TableName() string { return "last_events" }
