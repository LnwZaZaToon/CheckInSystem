package teacher

import (
	"myapp/core/teacher"

	"github.com/gofiber/fiber/v2"
)

type HttpTeacherHandler struct {
	service teacher.TeacherService
}

func NewHttpTeacherHandler(service teacher.TeacherService) *HttpTeacherHandler {
	return &HttpTeacherHandler{service: service}
}

func (h *HttpTeacherHandler) ImportStudentsFromCSV(c *fiber.Ctx) error {
	// Implementation goes here
	return nil
}

func (h *HttpTeacherHandler) GetAllStudents(c *fiber.Ctx) error {
	ctx := c.UserContext()
	teacherID, _ := c.ParamsInt("teacher_id")
	students, err := h.service.GetAllStudents(ctx, teacherID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(students)
}

func (h *HttpTeacherHandler) GetAllClasses(c *fiber.Ctx) error {
	// Implementation goes here
	return nil
}

func (h *HttpTeacherHandler) Login(c *fiber.Ctx) error {
	ctx := c.UserContext()
	type LoginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	err := h.service.Login(ctx, req.Username, req.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid username or password",
		})
	}
	return c.JSON(fiber.Map{
		"message": "Login successful",
	})
}
func (h *HttpTeacherHandler) Register(c *fiber.Ctx) error {
	ctx := c.UserContext()
	type RegisterRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	var req RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	err := h.service.Register(ctx, req.Username, req.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Registration successful",
	})
}
