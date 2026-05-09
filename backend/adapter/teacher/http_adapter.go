package teacher

import (
	"errors"
	"myapp/core/teacher"

	"github.com/gofiber/fiber/v2"
)

type studentResponse struct {
	ID        int    `json:"id"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	School    string `json:"school"`
}

type HttpTeacherHandler struct {
	service teacher.TeacherService
}

func NewHttpTeacherHandler(service teacher.TeacherService) *HttpTeacherHandler {
	return &HttpTeacherHandler{service: service}
}

func (h *HttpTeacherHandler) ImportStudentsFromCSV(c *fiber.Ctx) error {
	return nil
}

func (h *HttpTeacherHandler) GetAllStudents(c *fiber.Ctx) error {
	ctx := c.UserContext()
	teacherID, err := c.ParamsInt("teacher_id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid teacher id",
		})
	}
	students, err := h.service.GetAllStudents(ctx, teacherID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal Server Error",
		})
	}
	resp := make([]studentResponse, len(students))
	for i, s := range students {
		resp[i] = studentResponse{
			ID:        s.ID,
			Firstname: s.Firstname,
			Lastname:  s.Lastname,
			School:    s.School,
		}
	}
	return c.JSON(resp)
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
	token, err := h.service.Login(ctx, req.Username, req.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid username or password",
		})
	}
	return c.JSON(fiber.Map{
		"message": "Login successful",
		"token":   token,
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
		if errors.Is(err, teacher.ErrUsernameExists) {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "username already exists",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "internal server error",
		})
	}
	return c.JSON(fiber.Map{
		"message": "Registration successful",
	})
}
