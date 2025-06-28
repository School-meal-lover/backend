package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/School-meal-lover/backend/app/internal/services"
	"github.com/gin-gonic/gin"
)

type ExcelHandler struct {
    excelService *services.ExcelService
}

func NewExcelHandler(excelService *services.ExcelService) *ExcelHandler {
    return &ExcelHandler{excelService: excelService}
}

// @Summary 엑셀 처리 API
// @Description 파일을 업로드 해서 식단 데이터를 디비에 저장한다. 
// @Tags excel
// @Accept multipart/form-data 
// @Param excel formData file true "업로드할 엑셀 파일 (.xlsx 또는 .xls)" - file 타입 사용
// @Success 200 {object} models.ExcelProcessResult "Excel file processed successfully"
// @Failure 500 {object} models.ErrorResponse "Failed to process Excel file"
// @Router /upload/excel [post]
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
		// TODO: 버킷에 저장하는 것으로 변경 필요
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
// TO DO: 엑셀 파일 저장소가 정해지면 변경 필요
// @Summary 로컬 엑셀 파일 처리 (개발용)
// @Description 서버 내부에 하드코딩된 엑셀 파일 경로를 사용하여 식단 데이터를 파싱하고 DB에 저장합니다.
// @Tags excel
// @Router /process/excel/local [get]
// @Success 200 {object} models.ExcelProcessResult "Excel file processed successfully."
// @Failure 500 {object} models.ErrorResponse "Failed to process Excel file"
func (h *ExcelHandler) ProcessLocalExcel(c *gin.Context) {
    filePath := "uploads/2025_5_5_ko.xlsx" 
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