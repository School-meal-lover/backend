package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/School-meal-lover/backend/internal/services"
	"github.com/gin-gonic/gin"
)

type ExcelHandler struct {
    excelService *services.ExcelService
}

func NewExcelHandler(excelService *services.ExcelService) *ExcelHandler {
    return &ExcelHandler{excelService: excelService}
}

// POST /api/v1/upload/excel - 엑셀 파일 업로드 및 처리
//메소드 리시버: struct 타입의 인스턴스 받아서 이 메소드 실행 
func (h *ExcelHandler) UploadAndProcessExcel(c *gin.Context) {
    // 파일 업로드 처리
    file, err := c.FormFile("excel")
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "success": false,
            "error":   "excel file is required",
        })
        return
    }
    
    // 파일 확장자 검증
    ext := filepath.Ext(file.Filename)
    if ext != ".xlsx" && ext != ".xls" {
        c.JSON(http.StatusBadRequest, gin.H{
            "success": false,
            "error":   "only excel files (.xlsx, .xls) are allowed",
        })
        return
    }
    
    // 파일 저장
    filename := "uploads/" + file.Filename
    if err := c.SaveUploadedFile(file, filename); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "success": false,
            "error":   "failed to save file",
        })
        return
    }
    
    // 엑셀 처리 서비스 호출
    result, err := h.excelService.ProcessExcelFile(filename)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "success": false,
            "error":   err.Error(),
        })
        return
    }
    
    c.JSON(http.StatusOK, result)
}

// GET /api/v1/process/excel/local - 테스트용 api
func (h *ExcelHandler) ProcessLocalExcel(c *gin.Context) {
    filePath := "uploads/2025_6_1_ko.xlsx" 
    cwd, err := os.Getwd()
    fmt.Printf("DEBUG: Current Working Directory (CWD): %s\n", cwd)
 		result, err := h.excelService.ProcessExcelFile(filePath)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "success": false,
            "error":   fmt.Sprintf("Failed to process Excel file: %v", err),
        })
        return
    }
	
    c.JSON(http.StatusOK, gin.H{
        "success": true,
        "message": "Excel file processed successfully.",
        "data":    result,
    })
}