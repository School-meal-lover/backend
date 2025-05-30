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
        SELECT COUNT(*) 
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

func (r *MealRepository) GetRestaurantWeekMeals(restaurantID string, date string) (*models.RestaurantMealsData, error) {
	// 레스토랑 정보 조회
	restaurant, err := r.getRestaurantInfo(restaurantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get restaurant: %w", err)
	}
	// 주차 정보 조회
	week, err := r.getWeekInfo(restaurantID, date)
	if err != nil {
		return nil, fmt.Errorf("failed to get week: %w", err)
	}
	// 식단 조회
	mealsByDay, summary, err := r.getMealsData(week.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get meals data: %w", err)
	}
	return &models.RestaurantMealsData{
			Restaurant: restaurant,
			Week:       week,
			MealsByDay: mealsByDay,
			Summary:    summary,
	}, nil
}

// 레스토랑 정보 조회
func (r *MealRepository) getRestaurantInfo(restaurantID string) (*models.RestaurantInfo, error) {
    restaurant := &models.RestaurantInfo{}
    query := `SELECT id, name, COALESCE(name_en, '') FROM restaurants WHERE id = $1`
    
    err := r.db.QueryRow(query, restaurantID).Scan(
        &restaurant.ID, &restaurant.Name, &restaurant.NameEn)
    
    return restaurant, err
}

func (r *MealRepository) getWeekInfo(restaurantId, date string) (*models.WeekInfo, error) {
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

func (r *MealRepository) getMealsData(weekID string) (map[string]*models.DayMeals, *models.MealsSummary, error) {
	query := `
		SELECT 
            m.id, m.date, m.day_of_week, m.meal_type,
            COALESCE(mi.id, '') as menu_id,
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

	for rows.Next(){
		var mealID, dayOfWeek, mealType, menuID, menuName, menuNameEn string
		var date time.Time
		var price float64

		err := rows.Scan(&mealID, &date, &dayOfWeek, &mealType, 
                        &menuID, &menuName, &menuNameEn, &price)
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
					MenuItems: []*models.MenuItem{},
			}
			mealsByDay[dayOfWeek].Meals[mealType] = mealMap[mealKey]
			totalMeals++
		}
		// 메뉴 아이템 넣기
		if menuID != "" {
			menuItem := &models.MenuItem{
					ID:     menuID,
					Name:   menuName,
					NameEn: menuNameEn,
					Price:  price,
			}
			mealMap[mealKey].MenuItems = append(mealMap[mealKey].MenuItems, menuItem)
			totalMenuItems++
		}
	}

	summary := &models.MealsSummary{
        TotalDays:      len(mealsByDay),
        TotalMeals:     totalMeals,
        TotalMenuItems: totalMenuItems,
    }

    return mealsByDay, summary, nil
}


