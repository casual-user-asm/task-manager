package models

type Task struct {
	ID          uint   `json:"id" gorm:"primaryKey"`
	Title       string `json:"title"`
	Description string `'json:"description"`
	Date        string `json:"date"`
	Status      bool   `json:"status"`
	UserID      uint   `json:"userid"`
	User        User   `json:"user"`
}
