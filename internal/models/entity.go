package models

import "time"

type RestaurantType string
const (
	Restaurant1 RestaurantType = "RESTAURANT_1"
	Restaurant2 RestaurantType = "RESTAURANT_2"
)

type Week struct {
	ID         string         `json:"id" db:"id"`
	StartDate  time.Time      `json:"start_date" db:"start_date"`
	RestaurantType RestaurantType `json:"restaurant" db:"restaurants"`
}

type Meal struct {
	ID        string    `json:"id" db:"id"`
	WeekID    string    `json:"week_id" db:"weeks_id"`
	Date      time.Time `json:"date" db:"date"`
	DayOfWeek string    `json:"day_of_week" db:"day_of_week"`
	MealType  string    `json:"meal_type" db:"meal_type"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type MenuItem struct {
	ID       string  `json:"id" db:"id"`
	MealID   string  `json:"meal_id" db:"meals_id"`
	Category string  `json:"category" db:"category"`
	Name     string  `json:"name" db:"name"`
	NameEn   string  `json:"name_en" db:"name_en"`
	Price    float64 `json:"price" db:"price"`
}

// 엑셀 파싱용 구조체
type DateInfo struct {
	Date      string `json:"date"`
	DayOfWeek string `json:"day_of_week"`
	Col       string `json:"col"`
}

type MealTypeConfig struct {
	MealType string `json:"meal_type"`
	StartRow int    `json:"start_row"`
	EndRow   int    `json:"end_row"`
}
