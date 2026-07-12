package controller

import (
	"bluespot/internal/model"
	"bluespot/internal/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type MailController struct {
	mailService *service.MailService
}

func NewMailController(mailService *service.MailService) *MailController {
	return &MailController{
		mailService: mailService,
	}
}

// Send godoc
// @Summary 发送邮件
// @Description 根据收件人、标题和正文发送邮件。正文支持 HTML。
// @Tags mail
// @Accept json
// @Produce json
// @Param request body model.SendMailRequest true "邮件参数"
// @Success 200 {object} ResponseData
// @Router /mail [post]
func (mc *MailController) Send(c *gin.Context) {
	var param model.SendMailRequest
	if err := c.ShouldBindJSON(&param); err != nil {
		ResponseFailedWithMsg(c, CodeInvalidParam, "请求参数错误: "+err.Error())
		return
	}

	if err := mc.mailService.SendMail(&param); err != nil {
		zap.L().Error("send mail failed", zap.Error(err))
		ResponseFailed(c, CodeServerBusy)
		return
	}

	ResponseSuccess(c, nil)
}
