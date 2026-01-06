package controller

import (
	"ai-notetaking-be/internal/pkg/serverutils"
	"ai-notetaking-be/internal/service"

	"github.com/gofiber/fiber/v2"
)

type IChatbotController interface {
	RegisterRoutes(r fiber.Router)
	CreateSession(ctx *fiber.Ctx) error
	GetAllSessions(ctx *fiber.Ctx) error
}

type chatbotController struct {
	service service.IChatbotService
}

func NewChatbotController(service service.IChatbotService) IChatbotController {
	return &chatbotController{service: service}
}

func (c *chatbotController) RegisterRoutes(r fiber.Router) {
	h := r.Group("/chatbot/v1")
	h.Get("session", c.GetAllSessions)
	h.Post("create-session", c.CreateSession)
}

func (c *chatbotController) CreateSession(ctx *fiber.Ctx) error {
	res, err := c.service.CreateSession(ctx.Context())
	if err != nil {
		return err
	}

	return ctx.JSON(serverutils.SuccessResponse("Success create session", res))
}

func (c *chatbotController) GetAllSessions(ctx *fiber.Ctx) error {
	res, err := c.service.GetAllSessions(ctx.Context())
	if err != nil {
		return err
	}

	return ctx.JSON(serverutils.SuccessResponse("Success get all session", res))
}
