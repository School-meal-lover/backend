package services

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/School-meal-lover/backend/app/internal/models"
	"github.com/School-meal-lover/backend/app/internal/repository"
	"github.com/School-meal-lover/backend/app/pkg/excel"
)

type ExcelService struct {
    mealRepo *repository.MealRepository
    parser   *excel.Parser
}

func NewExcelService(mealRepo *repository.MealRepository) *ExcelService {
    return &ExcelService{
        mealRepo: mealRepo,
        parser:   excel.NewParser(),
    }
}

// 엑셀 파일 전체 처리 프로세스
func (s *ExcelService) ProcessExcelFile(filePath string) (*models.ExcelProcessResult, error) {
    log.Printf("Starting to process Excel file: %s", filePath)
    
    // 1. 엑셀 파일 열기
    f, err := s.parser.OpenExcelFile(filePath)
    if err != nil {
        return nil, fmt.Errorf("failed to open Excel file: %w", err)
    }
    defer f.Close()
    
    // 2. 레스토랑 정보 읽기 
    restaurantName, err := s.parser.ReadRestaurantName(f)
    if err != nil {
        return nil, fmt.Errorf("failed to read restaurant name: %w", err)
    }
    
    restaurant, err := s.mealRepo.GetRestaurantByName(restaurantName)
    if err != nil {
        return nil, fmt.Errorf("failed to get restaurant: %w", err)
    }
    
    // 3. 주차 정보 생성 
		//weekStartDate 형식: "2006-01-02"
    weekStartDate, err := s.parser.ReadWeekStartDate(f)
    if err != nil {
        return nil, fmt.Errorf("failed to read week start date: %w", err)
    }
    
    weekID, err := s.mealRepo.InsertWeek(weekStartDate, restaurant.ID)
    if err != nil {
        return nil, fmt.Errorf("failed to insert week: %w", err)
    }
    
    // 4. 날짜 정보 구성
		sheetName, err := s.parser.GetFirstNonEmptySheet(f)
		if err != nil {
		return nil, fmt.Errorf("failed to get first non-empty sheet: %w", err)
	}
    dates, err := s.parser.BuildDatesFromExcel(f, sheetName)
    if err != nil {
        return nil, fmt.Errorf("failed to build dates: %w", err)
    }
    
    // 5. 식사 및 메뉴 데이터 처리
    totalMeals, totalMenuItems, err := s.processMealsAndMenus(f, weekID, dates)
    if err != nil {
        return nil, fmt.Errorf("failed to process meals and menus: %w", err)
    }
    
    return &models.ExcelProcessResult{
        Success:        true,
        RestaurantName: restaurantName,
        WeekStartDate:  weekStartDate.Format("2006-01-02"),
				WeekID:         weekID,
        TotalMeals:     totalMeals,
        TotalMenuItems: totalMenuItems,
        Message:        "KoreanExcel file processed successfully",
    }, nil
}
// 영어 엑셀 파일 처리

func (s *ExcelService) ProcessEnglishExcelFile(filePath string, weekID string) (*models.ExcelProcessResult, error) {
	f, err := s.parser.OpenExcelFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open English Excel file: %w", err)
	}
	defer f.Close()

	sheetName,err := s.parser.GetFirstNonEmptySheet(f)
	if err != nil {
		return nil, fmt.Errorf("failed to get first non-empty sheet: %w", err)
	}
	dates, err := s.parser.BuildDatesFromExcel(f, sheetName)
	if err != nil {
			return nil, fmt.Errorf("failed to build dates: %w", err)
	}

	mealTypes := s.getMealTypeConfigs()
	updatedCount := 0

	for _, mealType := range mealTypes {
		for _, dateInfo := range dates{
			mealID, err := s.mealRepo.GetMealIDByWeekDateAndType(weekID, dateInfo.Date, mealType.MealType)
			if err != nil {
				log.Printf("Meal not found for %s %s: %v", dateInfo.Date, mealType.MealType, err)
				continue
			}
			englishNames, err := s.parser.ReadMenuItems(f, dateInfo.Col, mealType.StartRow, mealType.EndRow)
			if err != nil {
				log.Printf("Failed to read English menu items for %s %s: %v", dateInfo.Date, mealType.MealType, err)
				continue
			}
			if len(englishNames) == 0 {
				continue
			}
			koreanMenuItems, err := s.mealRepo.GetMenuItemsByMealIDOrdered(mealID)
			if err != nil {
					log.Printf("Failed to fetch Korean menu items for mealID %s: %v", mealID, err)
					continue
			}
			minLen := min(len(englishNames), len(koreanMenuItems))
			for i := 0; i < minLen; i++ {
					err := s.mealRepo.UpdateMenuItemNameEn(koreanMenuItems[i].ID, englishNames[i])
					if err != nil {
							log.Printf("Failed to update NameEn for menu item %s: %v", koreanMenuItems[i].ID, err)
							continue
					}
					updatedCount++
			}
		}
	}
	return &models.ExcelProcessResult{
        Success:        true,
				WeekID:         weekID,
        Message:        "English Excel file processed successfully",
    }, nil
}
// 식사 및 메뉴 처리 (비즈니스 로직)
func (s *ExcelService) processMealsAndMenus(f *excel.ExcelFile, weekID string, dates []models.DateInfo) (int, int, error) {
    mealTypes := s.getMealTypeConfigs()
    totalMeals := 0
    totalMenuItems := 0
    
    for _, mealType := range mealTypes {
        for _, dateInfo := range dates {
            // 식사 정보 생성
            meal := s.buildMealFromDateInfo(weekID, dateInfo, mealType.MealType)
            
            mealID, err := s.mealRepo.InsertMeal(meal)
            if err != nil {
                log.Printf("Failed to insert meal for %s, %s: %v", mealType.MealType, dateInfo.Date, err)
                continue
            }
            totalMeals++
            
            // 메뉴 아이템 읽기 및 처리
            menuItemNames, err := s.parser.ReadMenuItems(f, dateInfo.Col, mealType.StartRow, mealType.EndRow)
            if err != nil {
                log.Printf("Failed to read menu items: %v", err)
                continue
            }
            
            if len(menuItemNames) > 0 {
                menuItems := s.buildMenuItems(mealID, menuItemNames, mealType.MealType)
                
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

// 식사 객체 생성 아침/점심/점심/저녁
func (s *ExcelService) buildMealFromDateInfo(weekID string, dateInfo models.DateInfo, mealType string) *models.Meal {
    date, err := time.Parse("2006-01-02", dateInfo.Date)
    if err != nil {
        log.Printf("Failed to parse date %s: %v", dateInfo.Date, err)
        date = time.Now() // 기본값
    }
     
    return &models.Meal{
        WeekID:    weekID,
        Date:      date,
        DayOfWeek: dateInfo.DayOfWeek,
        MealType:  mealType,
    }
}

// 메뉴 아이템 생성 
func (s *ExcelService) buildMenuItems(mealID string, itemNames []string, mealType string) []models.MenuItem {
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

// 식사 타입별 설정 반환
func (s *ExcelService) getMealTypeConfigs() []models.MealTypeConfig {
    return []models.MealTypeConfig{
        {"Breakfast", 7, 16},
        {"Lunch_1", 18, 18},
        {"Lunch_2", 21, 26},
        {"Dinner", 27, 32},
    }
}

// 식사 타입별 카테고리 반환
func (s *ExcelService) getCategoriesForMealType(mealType string) []string {
    if mealType == "Lunch_1" {
        return []string{"메인메뉴"} // 일품 메뉴
    }
		if mealType == "Breakfast" {
				return []string{"밥", "국", "반찬", "메인메뉴", "반찬", "반찬", "반찬", "반찬", "반찬"}
		}
    return []string{"밥", "국", "메인메뉴", "메인메뉴", "반찬", "반찬"}
}