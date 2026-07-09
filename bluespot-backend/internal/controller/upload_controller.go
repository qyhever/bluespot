package controller

import (
	"errors"
	"strconv"
	"time"

	"bluespot/internal/middleware"
	"bluespot/internal/model"
	"bluespot/internal/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UploadController struct {
	uploadService *service.UploadService
}

func NewUploadController(uploadService *service.UploadService) *UploadController {
	return &UploadController{uploadService: uploadService}
}

func (uc *UploadController) Verify(c *gin.Context) {
	userID, ok := getCurrentUserID(c)
	if !ok {
		ResponseFailed(c, CodeNeedLogin)
		return
	}

	var req model.UploadVerifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ResponseFailed(c, CodeInvalidParam)
		return
	}

	result, err := uc.uploadService.Verify(userID, req)
	if err != nil {
		handleUploadError(c, "verify upload failed", err)
		return
	}
	ResponseSuccess(c, result)
}

func (uc *UploadController) UploadChunk(c *gin.Context) {
	userID, ok := getCurrentUserID(c)
	if !ok {
		ResponseFailed(c, CodeNeedLogin)
		return
	}

	chunk, err := c.FormFile("chunk")
	if err != nil {
		ResponseFailed(c, CodeInvalidParam)
		return
	}
	chunkIndex, err := strconv.Atoi(c.PostForm("chunkIndex"))
	if err != nil {
		ResponseFailed(c, CodeInvalidParam)
		return
	}

	err = uc.uploadService.SaveChunk(
		userID,
		c.PostForm("uploadId"),
		c.PostForm("fileMd5"),
		c.PostForm("fileName"),
		chunkIndex,
		chunk,
	)
	if err != nil {
		handleUploadError(c, "upload chunk failed", err)
		return
	}
	// 等待2s后返回
	time.Sleep(2 * time.Second)
	ResponseSuccess(c, nil)
}

func (uc *UploadController) Merge(c *gin.Context) {
	userID, ok := getCurrentUserID(c)
	if !ok {
		ResponseFailed(c, CodeNeedLogin)
		return
	}

	var req model.UploadMergeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ResponseFailed(c, CodeInvalidParam)
		return
	}

	result, err := uc.uploadService.Merge(userID, req)
	if err != nil {
		handleUploadError(c, "merge upload failed", err)
		return
	}
	ResponseSuccess(c, result)
}

func getCurrentUserID(c *gin.Context) (uint64, bool) {
	value, exists := c.Get(middleware.ContextUserIDKey)
	if !exists {
		return 0, false
	}
	switch userID := value.(type) {
	case uint64:
		return userID, true
	case int64:
		if userID < 0 {
			return 0, false
		}
		return uint64(userID), true
	case int:
		if userID < 0 {
			return 0, false
		}
		return uint64(userID), true
	default:
		return 0, false
	}
}

func handleUploadError(c *gin.Context, logMsg string, err error) {
	if errors.Is(err, service.ErrInvalidUploadParam) ||
		errors.Is(err, service.ErrUploadIDMismatch) ||
		errors.Is(err, service.ErrChunkMissing) {
		ResponseFailed(c, CodeInvalidParam)
		return
	}
	zap.L().Error(logMsg, zap.Error(err))
	ResponseFailed(c, CodeServerBusy)
}
