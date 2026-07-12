package service

import (
	"bluespot/internal/model"
	"bluespot/internal/pkg/postal"
)

type MailService struct{}

func NewMailService() *MailService {
	return &MailService{}
}

func (s *MailService) SendMail(param *model.SendMailRequest) error {
	return postal.SendMail(param.To, param.Subject, param.Body)
}
