package controller

import (
	"github.com/aq-simei/coin-pilot/api/models"
	"github.com/aq-simei/coin-pilot/api/service"
	responses "github.com/aq-simei/coin-pilot/internal"
	errors "github.com/aq-simei/coin-pilot/internal/config/error"
	"github.com/gin-gonic/gin"
)

type RecordController interface {
	// Add fields and methods as needed for the RecordController
	GetRecords(ctx *gin.Context)
	CreateRecord(ctx *gin.Context)
	UpdateRecord(ctx *gin.Context)
	DeleteRecord(ctx *gin.Context)
}

type RecordControllerImpl struct {
	service service.RecordService
}

func NewRecordController(service service.RecordService) RecordController {
	return &RecordControllerImpl{
		service: service,
	}
}

func RegisterRecordRoutes(router *gin.RouterGroup, controller RecordController) {
	router.GET("/list", controller.GetRecords)
	router.POST("/new", controller.CreateRecord)
}

func (rc *RecordControllerImpl) GetRecords(ctx *gin.Context) {
	// Retrieve user_id from context
	userID, exists := ctx.Get("user_id")
	if !exists {
		responses.Unauthorized(ctx, "User ID not found in token")
		return
	}

	// Convert userID to string
	userIDStr, ok := userID.(string)
	if !ok || userIDStr == "" {
		responses.Unauthorized(ctx, "Invalid user ID in token")
		return
	}

	// Fetch records using the userID
	records, err := rc.service.GetRecords(ctx, userIDStr)
	if err != nil {
		responses.InternalServerError(ctx, "Failed to retrieve records")
		return
	}

	ctx.JSON(200, records)
}

func (rc *RecordControllerImpl) CreateRecord(ctx *gin.Context) {
	var record models.CreateRecordPayload
	userID, exists := ctx.Get("user_id")
	if !exists {
		responses.Unauthorized(ctx, "Must be logged in to create a record")
		return
	}
	if err := ctx.ShouldBindJSON(&record); err != nil {
		responses.BadRequest(ctx, "Invalid input")
		return
	}

	createdRecord, err := rc.service.CreateRecord(ctx, record, userID.(string))
	if err != nil {
		responses.InternalServerError(ctx, err.Error())
		return
	}

	responses.Success(ctx, createdRecord)
}

func (rc *RecordControllerImpl) UpdateRecord(ctx *gin.Context) {
	id := ctx.Param("id")
	var record models.Record
	if err := ctx.ShouldBindJSON(&record); err != nil {
		responses.BadRequest(ctx, "Invalid input")
		return
	}

	updatedRecord, err := rc.service.UpdateRecord(ctx, id, record)
	if err != nil {
		responses.InternalServerError(ctx, "Failed to update record")
		return
	}

	responses.Success(ctx, updatedRecord)
}

func (rc *RecordControllerImpl) DeleteRecord(ctx *gin.Context) {
	id := ctx.Param("id")

	if id == "" {
		errors.NewBadRequest("Record ID is required")
		return
	}

	err := rc.service.DeleteRecord(ctx, id)
	if err != nil {
		responses.InternalServerError(ctx, "Failed to delete record")
		return
	}

	responses.Success(ctx, "Deleted") // No content
}
