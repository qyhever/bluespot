package controller

import (
	"bluespot/internal/model"
	"bluespot/internal/service"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type TelegramController struct {
	telegramService *service.TelegramService
}

func NewTelegramController(telegramService *service.TelegramService) *TelegramController {
	return &TelegramController{telegramService: telegramService}
}

// Send godoc
// @Summary 发送 Telegram 消息
// @Description 使用服务端配置的 Bot Token 和 Chat ID 发送文本消息。
// @Tags telegram
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body model.SendTelegramMessageRequest true "Telegram 消息参数"
// @Success 200 {object} ResponseData
// @Router /telegram [post]
func (tc *TelegramController) Send(c *gin.Context) {
	var param model.SendTelegramMessageRequest
	if err := c.ShouldBindJSON(&param); err != nil {
		ResponseFailedWithMsg(c, CodeInvalidParam, "请求参数错误: "+err.Error())
		return
	}
	if strings.TrimSpace(param.Text) == "" {
		ResponseFailedWithMsg(c, CodeInvalidParam, "请求参数错误: text 不能为空")
		return
	}

	if err := tc.telegramService.SendMessage(c.Request.Context(), param.Text); err != nil {
		zap.L().Error("send telegram message failed", zap.Error(err))
		ResponseFailed(c, CodeServerBusy)
		return
	}

	ResponseSuccess(c, nil)
}
