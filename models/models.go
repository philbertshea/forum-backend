package models

type User struct {
	ID       uint   `gorm:"primaryKey"`
	Username string `gorm:"unique;not null"`
}

type Thread struct {
	ID       uint   `gorm:"primaryKey"`
	UserID   uint   `gorm:"not null"`
	Title    string `gorm:"not null"`
	Content  string `gorm:"type:text;not null"`
	Comments []Comment
}

type Comment struct {
	ID       uint   `gorm:"primaryKey"`
	Content  string `gorm:"type:text;not null"`
	ThreadID uint   `gorm:"not null"`
	UserID   uint   `gorm:"not null"`
}
