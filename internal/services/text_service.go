package services

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/School-meal-lover/backend/internal/models"
	"github.com/School-meal-lover/backend/internal/repository"
)

type TextService struct {
	mealRepo *repository.MealRepository
}

func NewTextService(mealRepo *repository.MealRepository) *TextService {
	return &TextService{
		mealRepo: mealRepo,
	}
}

// Plain text로 한 주치 식단 데이터 처리
func (s *TextService) ProcessTextData(text string) (*models.ExcelProcessResult, error) {
	log.Printf("Starting to process text data")

	// 1. 식당 타입 파싱
	restaurantType, err := s.parseRestaurantType(text)
	if err != nil {
		return nil, fmt.Errorf("failed to parse restaurant type: %w", err)
	}

	// 2. 주 시작일 파싱
	weekStartDate, err := s.parseWeekStartDate(text)
	if err != nil {
		return nil, fmt.Errorf("failed to parse week start date: %w", err)
	}

	// 3. 주차 정보 생성
	weekID, err := s.mealRepo.FindOrCreateWeek(weekStartDate, restaurantType)
	if err != nil {
		return nil, fmt.Errorf("failed to insert week: %w", err)
	}

	// 4. 날짜 정보 구성 (식당 타입에 따라 요일 수가 다름)
	dates, err := s.buildDatesFromText(weekStartDate, restaurantType)
	if err != nil {
		return nil, fmt.Errorf("failed to build dates: %w", err)
	}

	// 5. 식사 및 메뉴 데이터 처리
	totalMeals, totalMenuItems, err := s.processMealsAndMenusFromText(text, weekID, dates)
	if err != nil {
		return nil, fmt.Errorf("failed to process meals and menus: %w", err)
	}

	return &models.ExcelProcessResult{
		Success:        true,
		RestaurantType: string(restaurantType),
		WeekStartDate:  weekStartDate.Format("2006-01-02"),
		WeekID:         weekID,
		TotalMeals:     totalMeals,
		TotalMenuItems: totalMenuItems,
		Message:        "Text data processed successfully",
	}, nil
}

// 텍스트에서 식당 타입 파싱
func (s *TextService) parseRestaurantType(text string) (models.RestaurantType, error) {
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "식당:") {
			restaurantStr := strings.TrimSpace(strings.TrimPrefix(line, "식당:"))
			restaurantStr = strings.ToUpper(restaurantStr)
			if strings.Contains(restaurantStr, "1") || restaurantStr == "RESTAURANT_1" {
				return models.Restaurant1, nil
			}
			if strings.Contains(restaurantStr, "2") || restaurantStr == "RESTAURANT_2" {
				return models.Restaurant2, nil
			}
		}
	}
	return "", fmt.Errorf("restaurant type not found in text")
}

// 텍스트에서 주 시작일 파싱
func (s *TextService) parseWeekStartDate(text string) (time.Time, error) {
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "주 시작일:") {
			dateStr := strings.TrimSpace(strings.TrimPrefix(line, "주 시작일:"))
			date, err := time.Parse("2006-01-02", dateStr)
			if err != nil {
				return time.Time{}, fmt.Errorf("invalid date format: %s, expected YYYY-MM-DD", dateStr)
			}
			return date, nil
		}
	}
	return time.Time{}, fmt.Errorf("week start date not found in text")
}

// 주 시작일로부터 날짜 정보 구성 (식당 타입에 따라 요일 수가 다름)
func (s *TextService) buildDatesFromText(weekStartDate time.Time, restaurantType models.RestaurantType) ([]models.DateInfo, error) {
	var dates []models.DateInfo
	var dayNames []string
	var dayCount int

	// 식당 타입에 따라 요일 수 결정
	if restaurantType == models.Restaurant1 {
		// 식당 1: 월화수목금 (5일)
		dayNames = []string{"월요일", "화요일", "수요일", "목요일", "금요일"}
		dayCount = 5
	} else if restaurantType == models.Restaurant2 {
		// 식당 2: 월화수목금토일 (7일)
		dayNames = []string{"월요일", "화요일", "수요일", "목요일", "금요일", "토요일", "일요일"}
		dayCount = 7
	} else {
		return nil, fmt.Errorf("unknown restaurant type: %s", restaurantType)
	}

	for i := 0; i < dayCount; i++ {
		currentDate := weekStartDate.AddDate(0, 0, i)
		dates = append(dates, models.DateInfo{
			Date:      currentDate.Format("2006-01-02"),
			DayOfWeek: dayNames[i],
			Col:       "", // 텍스트 파싱에서는 컬럼 정보가 필요 없음
		})
	}

	return dates, nil
}

// 텍스트에서 식사 및 메뉴 데이터 처리
func (s *TextService) processMealsAndMenusFromText(text string, weekID string, dates []models.DateInfo) (int, int, error) {
	mealTypes := s.getMealTypeConfigs()
	totalMeals := 0
	totalMenuItems := 0

	// 텍스트를 날짜별로 분리
	daySections := s.splitTextByDays(text)

	for _, dateInfo := range dates {
		dayText, exists := daySections[dateInfo.DayOfWeek]
		if !exists {
			log.Printf("No data found for %s", dateInfo.DayOfWeek)
			continue
		}

		for _, mealType := range mealTypes {
			// 식사 정보 생성
			meal := s.buildMealFromDateInfo(weekID, dateInfo, mealType.MealType)

			mealID, err := s.mealRepo.FindOrCreateMeal(meal)
			if err != nil {
				log.Printf("Failed to insert meal for %s, %s: %v", mealType.MealType, dateInfo.Date, err)
				return totalMeals, totalMenuItems, err
			}
			totalMeals++

			// 메뉴 아이템 파싱 (카테고리 포함)
			menuItemsWithCategory, err := s.parseMenuItemsFromDayText(dayText, mealType.MealType)
			if err != nil {
				log.Printf("Failed to parse menu items for %s %s: %v", dateInfo.DayOfWeek, mealType.MealType, err)
				continue
			}

			if len(menuItemsWithCategory) > 0 {
				menuItems := s.buildMenuItemsFromCategory(mealID, menuItemsWithCategory)

				err := s.mealRepo.InsertMenuItems(menuItems)
				if err != nil {
					log.Printf("Failed to insert menu items for meal %s: %v", mealID, err)
					continue
				}
				totalMenuItems += len(menuItems)
			}
		}
	}

	return totalMeals, totalMenuItems, nil
}

// 텍스트를 요일별로 분리
func (s *TextService) splitTextByDays(text string) map[string]string {
	daySections := make(map[string]string)
	dayNames := []string{"월요일", "화요일", "수요일", "목요일", "금요일", "토요일", "일요일"}

	lines := strings.Split(text, "\n")
	currentDay := ""
	var currentSection []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// 요일로 시작하는 라인인지 확인
		isDayLine := false
		for _, dayName := range dayNames {
			if strings.HasPrefix(line, dayName) {
				// 이전 요일 섹션 저장
				if currentDay != "" {
					daySections[currentDay] = strings.Join(currentSection, "\n")
				}
				currentDay = dayName
				currentSection = []string{}
				isDayLine = true
				break
			}
		}

		if !isDayLine && currentDay != "" {
			currentSection = append(currentSection, line)
		}
	}

	// 마지막 요일 섹션 저장
	if currentDay != "" {
		daySections[currentDay] = strings.Join(currentSection, "\n")
	}

	return daySections
}

// MenuItemWithCategory 카테고리와 메뉴명을 함께 저장하는 구조체
type MenuItemWithCategory struct {
	Category string
	Name     string
}

// 요일 텍스트에서 특정 식사 타입의 메뉴 아이템 파싱 (카테고리 포함)
func (s *TextService) parseMenuItemsFromDayText(dayText string, mealType string) ([]MenuItemWithCategory, error) {
	var items []MenuItemWithCategory
	lines := strings.Split(dayText, "\n")

	mealTypeKorean := s.getMealTypeKorean(mealType)
	inMealSection := false

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// 식사 타입 섹션 시작 확인
		if strings.HasPrefix(line, mealTypeKorean+":") {
			inMealSection = true
			continue
		}

		// 다른 식사 타입 섹션 시작 시 종료
		if inMealSection {
			if s.isMealTypeLine(line) && !strings.HasPrefix(line, mealTypeKorean) {
				break
			}

			// 메뉴 아이템 파싱 (카테고리:메뉴명 형식)
			if strings.Contains(line, ":") {
				parts := strings.SplitN(line, ":", 2)
				if len(parts) == 2 {
					category := strings.TrimSpace(parts[0])
					menuName := strings.TrimSpace(parts[1])
					if menuName != "" {
						items = append(items, MenuItemWithCategory{
							Category: category,
							Name:     menuName,
						})
					}
				}
			}
		}
	}

	return items, nil
}

// 식사 타입의 한국어 이름 반환
func (s *TextService) getMealTypeKorean(mealType string) string {
	switch mealType {
	case "Breakfast":
		return "아침"
	case "Lunch_1":
		return "점심1"
	case "Lunch_2":
		return "점심2"
	case "Dinner":
		return "저녁"
	default:
		return mealType
	}
}

// 라인이 식사 타입 라인인지 확인
func (s *TextService) isMealTypeLine(line string) bool {
	mealTypes := []string{"아침", "점심1", "점심2", "저녁"}
	for _, mealType := range mealTypes {
		if strings.HasPrefix(line, mealType+":") {
			return true
		}
	}
	return false
}

// 식사 타입별 설정 반환
func (s *TextService) getMealTypeConfigs() []models.MealTypeConfig {
	return []models.MealTypeConfig{
		{MealType: "Breakfast", StartRow: 0, EndRow: 0},
		{MealType: "Lunch_1", StartRow: 0, EndRow: 0},
		{MealType: "Lunch_2", StartRow: 0, EndRow: 0},
		{MealType: "Dinner", StartRow: 0, EndRow: 0},
	}
}

// 식사 객체 생성
func (s *TextService) buildMealFromDateInfo(weekID string, dateInfo models.DateInfo, mealType string) *models.Meal {
	date, err := time.Parse("2006-01-02", dateInfo.Date)
	if err != nil {
		log.Printf("Failed to parse date %s: %v", dateInfo.Date, err)
		date = time.Now()
	}

	return &models.Meal{
		WeekID:    weekID,
		Date:      date,
		DayOfWeek: dateInfo.DayOfWeek,
		MealType:  mealType,
	}
}

// 메뉴 아이템 생성 (카테고리 정보 포함)
func (s *TextService) buildMenuItemsFromCategory(mealID string, itemsWithCategory []MenuItemWithCategory) []models.MenuItem {
	var menuItems []models.MenuItem

	if len(itemsWithCategory) == 0 {
		return menuItems
	}

	for _, item := range itemsWithCategory {
		category := strings.TrimSpace(item.Category)
		if category == "" {
			category = "기타" // 카테고리가 없으면 기본값
		}

		menuItems = append(menuItems, models.MenuItem{
			MealID:   mealID,
			Category: category,
			Name:     strings.TrimSpace(item.Name),
			NameEn:   "",
			Price:    0.00,
		})
	}
	return menuItems
}

// 메뉴 아이템 생성 (기존 방식 - 호환성을 위해 유지)
func (s *TextService) buildMenuItems(mealID string, itemNames []string, mealType string) []models.MenuItem {
	categories := s.getCategoriesForMealType(mealType)
	var menuItems []models.MenuItem

	if len(itemNames) == 0 {
		log.Printf("No menu items found for meal type %s", mealType)
		return menuItems
	}

	for idx, name := range itemNames {
		var category string

		if idx < len(categories) {
			category = categories[idx]
		} else {
			category = "기타" // 기본 카테고리
		}
		menuItems = append(menuItems, models.MenuItem{
			MealID:   mealID,
			Category: category,
			Name:     strings.TrimSpace(name),
			NameEn:   "",
			Price:    0.00,
		})
	}
	return menuItems
}

// 식사 타입별 카테고리 반환
func (s *TextService) getCategoriesForMealType(mealType string) []string {
	if mealType == "Lunch_1" {
		return []string{"메인메뉴"}
	}
	if mealType == "Breakfast" {
		return []string{"밥", "국", "반찬", "메인메뉴", "반찬", "반찬", "반찬", "반찬", "반찬"}
	}
	return []string{"밥", "국", "메인메뉴", "메인메뉴", "반찬", "반찬"}
}

