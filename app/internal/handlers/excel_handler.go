package handlers

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/School-meal-lover/backend/app/internal/models"
	"github.com/School-meal-lover/backend/app/internal/services"
	"github.com/gin-gonic/gin"
)

type ExcelHandler struct {
    excelService *services.ExcelService
}
type DualExcelProcessResponse struct {
    Success   bool                        `json:"success"`
    ResultKo  models.ExcelProcessResult   `json:"result_ko"`
    ResultEn  models.ExcelProcessResult   `json:"result_en"`
}

func NewExcelHandler(excelService *services.ExcelService) *ExcelHandler {
    return &ExcelHandler{excelService: excelService}
}

// @Summary 엑셀 처리 API
// @Description 파일을 업로드 해서 식단 데이터를 디비에 저장한다. 
// @Tags excel
// @Accept multipart/form-data 
// @Param excel_ko formData file true "한국어 엑셀 파일"
// @Param excel_en formData file true "영어 엑셀 파일"
// @Success 200 {object} DualExcelProcessResponse "Excel file processed successfully"
// @Failure 500 {object} models.ErrorResponse "Failed to process Excel file"
// @Router /upload/excel [post]
func (h *ExcelHandler) UploadAndProcessExcel(c *gin.Context) {
		var files [2]*multipart.FileHeader
    var err error

		files[0], err = c.FormFile("excel_ko")
		files[1], err = c.FormFile("excel_en")

		var fileKoPath, fileEnPath string
    for i,file:= range files[:2]{
			if file == nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"success": false,
					"error":   fmt.Sprintf("file %d is missing",i),
				})
				return
			}

			ext := filepath.Ext(file.Filename)
			if ext != ".xlsx" && ext != ".xls" {
				c.JSON(http.StatusBadRequest, gin.H{
					"success": false,
					"error":  fmt.Sprintf("file %d has invalid extension", i),
				})
				return
			}
			savePath := "uploads/" + file.Filename
			if err := c.SaveUploadedFile(file, savePath); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"success": false,
					"error":   fmt.Sprintf("failed to save file %d", i),
				})
				return
			}
			if i == 0 {
            fileKoPath = savePath
        } else {
            fileEnPath = savePath
        }
		}

    resultKo, err := h.excelService.ProcessExcelFile(fileKoPath)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "success": false,
            "error":   "failed to process Korean Excel: " + err.Error(),
        })
        return
    }

    resultEn, err := h.excelService.ProcessEnglishExcelFile(fileEnPath, resultKo.WeekID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "success": false,
            "error":   "failed to process English Excel: " + err.Error(),
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "success":     true,
        "result_ko":   resultKo,
        "result_en":   resultEn,
    })
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