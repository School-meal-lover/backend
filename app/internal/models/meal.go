package models

import "time"

type Restaurant struct {
    ID   string `json:"id" db:"id"`
    Name string `json:"name" db:"name"`
}

type Week struct {
    ID            string    `json:"id" db:"id"`
    StartDate     time.Time `json:"start_date" db:"start_date"`
    RestaurantID  string    `json:"restaurant_id" db:"restaurants_id"`
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

// 응답용 구조체
type ExcelProcessResult struct {
    Success        bool   `json:"success"`
    RestaurantName string `json:"restaurant_name,omitempty"`
		WeekID 			string `json:"week_id,omitempty"`
    WeekStartDate  string `json:"week_start_date,omitempty"`
    TotalMeals     int    `json:"total_meals,omitempty"`
    TotalMenuItems int    `json:"total_menu_items,omitempty"`
    Message        string `json:"message"`
}
