package models

type RestaurantMealsResponse struct {
    Success bool                     `json:"success"`
    Data    *RestaurantMealsData     `json:"data,omitempty"`
    Error   string                   `json:"error,omitempty"`
    Code    string                   `json:"code,omitempty"`
}

type RestaurantMealsData struct {
    Restaurant *RestaurantInfo       `json:"restaurant"`
    Week       *WeekInfo            `json:"week"`
    MealsByDay []*DayMeals `json:"meals_by_day"`
    Summary    *MealsSummary        `json:"summary"`
}

type RestaurantInfo struct {
    ID     string `json:"id"`
    Name   string `json:"name"`
    NameEn string `json:"name_en"`
}

type WeekInfo struct {
    ID        string `json:"id"`
    StartDate string `json:"start_date"`
    EndDate   string `json:"end_date"`
}

type DayMeals struct {
    Date      string                `json:"date"`
    DayOfWeek string                `json:"day_of_week"`
    Meals     map[string]*MealInfo  `json:"meals"`
}
type MenuItemResponse struct {
    ID       string  `json:"id" db:"id"`
		// MealID   string  `json:"meal_id" db:"meals_id"`
    Category string  `json:"category" db:"category"`
    Name     string  `json:"name" db:"name"`
    NameEn   string  `json:"name_en" db:"name_en"`
    Price    float64 `json:"price" db:"price"`
}
type MealInfo struct {
    MealID    string      `json:"meal_id"`
    MealType  string      `json:"meal_type"`
    MenuItems []*MenuItemResponse `json:"menu_items"`
}

type MealsSummary struct {
    TotalDays      int `json:"total_days"`
    TotalMeals     int `json:"total_meals"`
    TotalMenuItems int `json:"total_menu_items"`
}