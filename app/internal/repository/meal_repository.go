package repository

import (
	"database/sql"
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/School-meal-lover/backend/app/internal/models"

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

func (r *MealRepository) HandleRepositoryError(err error, notFoundCode, notFoundMessage string) (*models.RestaurantMealsResponse, error) {
	if err != nil {
		return &models.RestaurantMealsResponse{
			Success: false,
			Error:   notFoundMessage,
			Code:    notFoundCode,
		}, nil
	}
	return nil, nil
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
	defer func() {
		if err := tx.Rollback(); err != nil  && err != sql.ErrTxDone {
			log.Printf("failed to rollback transaction: %v", err)
		}	
		}()

	stmt, err := tx.Prepare(`
				INSERT INTO menu_items (id, meals_id, category, name, name_en, price, created_at, updated_at)
				VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
				ON CONFLICT (meals_id, category, name) DO UPDATE SET
					name_en = EXCLUDED.name_en,
					price = EXCLUDED.price,
					updated_at = NOW();
			`)
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

// 레스토랑 정보 조회
func (r *MealRepository) GetRestaurantInfo(restaurantID string) (*models.RestaurantInfo, error) {
	restaurant := &models.RestaurantInfo{}
	query := `SELECT id, name, COALESCE(name_en, '') FROM restaurants WHERE id = $1`

	err := r.db.QueryRow(query, restaurantID).Scan(
		&restaurant.ID, &restaurant.Name, &restaurant.NameEn)

	return restaurant, err
}
func (r *MealRepository) GetWeekInfo(restaurantId, date string) (*models.WeekInfo, error) {
	week := &models.WeekInfo{}
	var startDate time.Time
	query := `
		SELECT id, start_date
        FROM weeks 
        WHERE restaurants_id = $1 
        AND $2 >= start_date 
        AND $2 <= start_date + INTERVAL '6 days'
        ORDER BY start_date DESC
        LIMIT 1`

	err := r.db.QueryRow(query, restaurantId, date).Scan(&week.ID, &startDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get week by date: %w", err)
	}
	week.StartDate = startDate.Format("2006-01-02")
	week.EndDate = startDate.AddDate(0, 0, 6).Format("2006-01-02")

	return week, nil
}

func (r *MealRepository) GetMealsData(weekID string) ([]*models.DayMeals, *models.MealsSummary, error) {
	query := `
			SELECT
								m.id, m.date, m.day_of_week, m.meal_type, mi.category,
								COALESCE(mi.id, '00000000-0000-0000-0000-000000000000'::UUID) as menu_id,
								COALESCE(mi.name, '') as menu_name,
								COALESCE(mi.name_en, '') as menu_name_en,
								COALESCE(mi.price, 0) as price
						FROM meals m
						LEFT JOIN menu_items mi ON m.id = mi.meals_id
						WHERE m.weeks_id = $1
						ORDER BY m.date,
										CASE m.meal_type
												WHEN 'Breakfast' THEN 1 WHEN 'Lunch_1' THEN 2
												WHEN 'Lunch_2' THEN 3 WHEN 'Dinner' THEN 4 END`

	rows, err := r.db.Query(query, weekID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get meals data: %w", err)
	}
	defer rows.Close()

	mealsByDay := make(map[string]*models.DayMeals)
	mealMap := make(map[string]*models.MealInfo)
	totalMeals := 0
	totalMenuItems := 0

	for rows.Next() {
		var mealID, dayOfWeek, mealType, menuID, category, menuName, menuNameEn string
		var date time.Time
		var price float64

		err := rows.Scan(&mealID, &date, &dayOfWeek, &mealType,&category,&menuID, &menuName, &menuNameEn, &price)
		if err != nil {
			return nil, nil, err
		}
		dateStr := date.Format("2006-01-02")
		//요일별 데이터
		if mealsByDay[dayOfWeek] == nil {
			mealsByDay[dayOfWeek] = &models.DayMeals{
				Date:      dateStr,
				DayOfWeek: dayOfWeek,
				Meals:     make(map[string]*models.MealInfo),
			}
		}
		// 식사별 데이터 초기화
		mealKey := mealID
		if mealMap[mealKey] == nil {
			mealMap[mealKey] = &models.MealInfo{
				MealID:    mealID,
				MealType:  mealType,
				MenuItems: []*models.MenuItemResponse{},
			}
			mealsByDay[dayOfWeek].Meals[mealType] = mealMap[mealKey]
			totalMeals++
		}
		// 메뉴 아이템 넣기
		if menuID != "" {
			menuItem := &models.MenuItemResponse{
				ID:     menuID,
				Category: category,
				Name:   menuName,
				NameEn: menuNameEn,
				Price:  price,
			}
			mealMap[mealKey].MenuItems = append(mealMap[mealKey].MenuItems, menuItem)
			totalMenuItems++
		}
	}
	var orderedDays []*models.DayMeals
	for _, dayMeal := range mealsByDay {
		orderedDays = append(orderedDays, dayMeal)
  }
	sort.Slice(orderedDays, func(i,j int) bool {
    dateI, errI := time.Parse(("2006-01-02"), orderedDays[i].Date)
    dateJ, errJ := time.Parse(("2006-01-02"), orderedDays[j].Date)
    if errI != nil || errJ != nil {
      log.Printf("failed to parse date: %v, %v", errI, errJ)
      return orderedDays[i].Date < orderedDays[j].Date
    }
    return dateI.Before(dateJ)
  })

	summary := &models.MealsSummary{
		TotalDays:      len(orderedDays),
		TotalMeals:     totalMeals,
		TotalMenuItems: totalMenuItems,
	}

	return orderedDays, summary, nil
}
