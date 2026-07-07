package controller

import (
	"errors"

	"bluespot/internal/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AttachController struct {
	attachService *service.AttachService
}

func NewAttachController(attachService *service.AttachService) *AttachController {
	return &AttachController{
		attachService: attachService,
	}
}

// Upload godoc
// @Summary 上传附件
// @Description 上传文件到本地存储，并返回公开访问地址。
// @Tags attach
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param file formData file true "上传文件"
// @Success 200 {object} SwaggerAttachUploadResponse
// @Router /attach/upload [post]
func (ac *AttachController) Upload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		ResponseFailed(c, CodeInvalidParam)
		return
	}

	result, err := ac.attachService.Upload(file)
	if err != nil {
		if errors.Is(err, service.ErrInvalidAttachFile) {
			ResponseFailed(c, CodeInvalidParam)
			return
		}
		zap.L().Error("upload attach failed", zap.Error(err))
		ResponseFailed(c, CodeServerBusy)
		return
	}

	ResponseSuccess(c, result)
}
