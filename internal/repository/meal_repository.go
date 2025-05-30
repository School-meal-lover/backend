package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/School-meal-lover/backend/internal/models"

	"github.com/google/uuid"
)

type MealRepository struct {
    db *sql.DB
}

func NewMealRepository(db *sql.DB) *MealRepository {
    return &MealRepository{db: db}
}

// 레스토랑 조회
func (r *MealRepository) GetRestaurantByName(name string) (*models.Restaurant, error) {
    query := `SELECT id FROM restaurants WHERE name = $1`
    
    var restaurant models.Restaurant
    err := r.db.QueryRow(query, name).Scan(&restaurant.ID)
    if err != nil {
        return nil, fmt.Errorf("restaurant not found for name %s: %w", name, err)
    }
    
    restaurant.Name = name
    return &restaurant, nil
}

// 주차 정보 삽입
func (r *MealRepository) InsertWeek(startDate time.Time, restaurantID string) (string, error) {
    weekID := uuid.New().String()
    query := `
        INSERT INTO weeks (id, start_date, restaurants_id, created_at, updated_at)
        VALUES ($1, $2, $3, now(), now())
        RETURNING id`
    
    var insertedID string
    err := r.db.QueryRow(query, weekID, startDate, restaurantID).Scan(&insertedID)
    if err != nil {
        return "", fmt.Errorf("failed to insert week: %w", err)
    }
    
    return insertedID, nil
}

// 식사 정보 삽입
func (r *MealRepository) InsertMeal(meal *models.Meal) (string, error) {
    if meal.ID == "" {
        meal.ID = uuid.New().String()
    }
    
    query := `
        INSERT INTO meals (id, weeks_id, date, day_of_week, meal_type, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, now(), now())
        RETURNING id`
    
    var insertedID string
    err := r.db.QueryRow(query, 
        meal.ID, meal.WeekID, 
        meal.Date, meal.DayOfWeek, meal.MealType).Scan(&insertedID)
    
    if err != nil {
        return "", fmt.Errorf("failed to insert meal %w", err)
    }
    
    return insertedID, nil
}

// 메뉴 아이템 일괄 삽입
func (r *MealRepository) InsertMenuItems(menuItems []models.MenuItem) error {
    if len(menuItems) == 0 {
        return nil
    }
    
    tx, err := r.db.Begin()
    if err != nil {
        return fmt.Errorf("failed to begin transaction: %w", err)
    }
    defer tx.Rollback()
    
    stmt, err := tx.Prepare(`
        INSERT INTO menu_items (id, meals_id, category, name, name_en, price, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, now(), now())`)
    if err != nil {
        return fmt.Errorf("failed to prepare statement: %w", err)
    }
    defer stmt.Close()
    
    for _, item := range menuItems {
        if item.ID == "" {
            item.ID = uuid.New().String()
        }
        
        _, err := stmt.Exec(item.ID, item.MealID, item.Category, item.Name, item.NameEn, item.Price)
        if err != nil {
            return fmt.Errorf("failed to insert menu item %s: %w", item.Name, err)
        }
    }
    
    return tx.Commit()
}

// 주차별 식사 개수 조회
func (r *MealRepository) GetMealCountByWeekID(weekID string) (int, error) {
    query := `SELECT COUNT(*) FROM meals WHERE weeks_id = $1`
    
    var count int
    err := r.db.QueryRow(query, weekID).Scan(&count)
    if err != nil {
        return 0, fmt.Errorf("failed to get meal count: %w", err)
    }
    
    return count, nil
}

// 식사별 메뉴 아이템 개수 조회
func (r *MealRepository) GetMenuItemCountByWeekID(weekID string) (int, error) {
    query := `
        SELECT COUNT(mi.*) 
        FROM menu_items mi 
        INNER JOIN meals m ON mi.meals_id = m.id 
        WHERE m.weeks_id = $1`
    
    var count int
    err := r.db.QueryRow(query, weekID).Scan(&count)
    if err != nil {
        return 0, fmt.Errorf("failed to get menu item count: %w", err)
    }
    
    return count, nil
}