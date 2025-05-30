package excel

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/School-meal-lover/backend/internal/models"

	"github.com/xuri/excelize/v2"
)

type Parser struct{}

type ExcelFile struct {
    *excelize.File
}

func NewParser() *Parser {
    return &Parser{}
}

// 엑셀 파일 열기
func (p *Parser) OpenExcelFile(filePath string) (*ExcelFile, error) {
    f, err := excelize.OpenFile(filePath)
    if err != nil {
        return nil, fmt.Errorf("failed to open Excel file: %w", err)
    }
    
    return &ExcelFile{File: f}, nil
}

// 레스토랑 이름 읽기
func (p *Parser) ReadRestaurantName(f *ExcelFile) (string, error) {
    cell, err := f.GetCellValue("12", "D2")
    if err != nil {
        return "", fmt.Errorf("failed to read cell D2: %w", err)
    }
    
    cell = strings.TrimSpace(cell)
    if cell == "" {
        return "", fmt.Errorf("cell D2 (restaurant name) is empty")
    }
    
    return cell, nil
}

// 주차 시작 날짜 읽기
func (p *Parser) ReadWeekStartDate(f *ExcelFile) (time.Time, error) {
    cell, err := f.GetCellValue("12", "D6")
    if err != nil {
        return time.Time{}, fmt.Errorf("failed to read cell D6: %w", err)
    }
    
    cell = strings.TrimSpace(cell)
    if cell == "" {
        return time.Time{}, fmt.Errorf("cell D6 (week start date) is empty")
    }
    
    // 엑셀 날짜 형태: "Mon 5/26"
    parts := strings.SplitN(cell, " ", 2)
    if len(parts) < 2 {
        return time.Time{}, fmt.Errorf("invalid date format in cell D6: %s, expected 'Day MM/DD'", cell)
    }
    
    fullDate := fmt.Sprintf("2025 %s", parts[1])
    parsedDate, err := time.Parse("2006 1/2", fullDate)
    if err != nil {
        return time.Time{}, fmt.Errorf("failed to parse date %s: %w", fullDate, err)
    }
    
    return parsedDate, nil
}

// 엑셀에서 날짜 정보 구성
func (p *Parser) BuildDatesFromExcel(f *ExcelFile) ([]models.DateInfo, error) {
    cols := []string{"D", "E", "F", "G", "H"}
    var dates []models.DateInfo
    
    for _, col := range cols {
        cell, err := f.GetCellValue("12", col+"6")
        if err != nil {
            return nil, fmt.Errorf("failed to read cell %s6: %w", col, err)
        }
        
        cell = strings.TrimSpace(cell)
        if cell == "" {
            continue
        }
        
        parts := strings.SplitN(cell, " ", 2)
        if len(parts) < 2 {
            continue
        }
        
        dateParts := strings.Split(parts[1], "/")
        if len(dateParts) < 2 {
            continue
        }
        
        month, err := strconv.Atoi(dateParts[0])
        if err != nil {
            continue
        }
        
        day, err := strconv.Atoi(dateParts[1])
        if err != nil {
            continue
        }
        
        formatted := fmt.Sprintf("2025-%02d-%02d", month, day)
        
        date, err := time.Parse("2006-01-02", formatted)
        if err != nil {
            continue
        }
        
        dates = append(dates, models.DateInfo{
            Date:      date.Format("2006-01-02"),
            DayOfWeek: parts[0],
            Col:       col,
        })
    }
    
    return dates, nil
}

// 메뉴 아이템 읽기
func (p *Parser) ReadMenuItems(f *ExcelFile, col string, startRow, endRow int) ([]string, error) {
    var items []string
    
    for rowIdx := startRow; rowIdx <= endRow; rowIdx++ {
        cell, err := f.GetCellValue("12", fmt.Sprintf("%s%d", col, rowIdx))
        if err != nil {
            continue // 에러가 있는 셀은 스킵
        }
        
        cell = strings.TrimSpace(cell)
        if cell != "" {
            items = append(items, cell)
        }
    }
    
    return items, nil
}