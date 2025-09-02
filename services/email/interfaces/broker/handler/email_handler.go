package handler

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/sayyidinside/monorepo-gofiber-clean/services/email/domain/service"
	"github.com/sayyidinside/monorepo-gofiber-clean/shared/interfaces/model"
)

type EmailHandler interface {
	SendEmail(body []byte) error
}

type emailHandler struct {
	service service.EmailService
}

func NewEmailHandler(service service.EmailService) EmailHandler {
	return &emailHandler{
		service: service,
	}
}

func (h *emailHandler) SendEmail(body []byte) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	email := model.Email{}
	if err := json.Unmarshal(body, &email); err != nil {
		log.Fatal(err)
		return err
	}

	h.service.SendEmail(ctx, email)

	return nil
}
