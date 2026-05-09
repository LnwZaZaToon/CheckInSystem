package class

import (
	"errors"
	"myapp/core/class"

	"github.com/gofiber/fiber/v2"
)

type classResponse struct {
	ID          int    `json:"id"`
	SubjectName string `json:"subject_name"`
}

type HttpClassHandler struct {
	service class.ClassService
}

func NewHttpClassHandler(service class.ClassService) *HttpClassHandler {
	return &HttpClassHandler{service: service}
}

func (h *HttpClassHandler) GetAllClasses(c *fiber.Ctx) error {
	ctx := c.UserContext()
	teacherID, err := c.ParamsInt("teacher_id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid teacher id"})
	}
	classes, err := h.service.GetAllClassesByTeacherID(ctx, teacherID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error"})
	}
	resp := make([]classResponse, len(classes))
	for i, cl := range classes {
		resp[i] = classResponse{ID: cl.ID, SubjectName: cl.SubjectName}
	}
	return c.JSON(resp)
}

func (h *HttpClassHandler) CreateClass(c *fiber.Ctx) error {
	ctx := c.UserContext()
	teacherID, err := c.ParamsInt("teacher_id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid teacher id"})
	}
	type createClassRequest struct {
		SubjectName string `json:"subject_name"`
	}
	var req createClassRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}
	id, err := h.service.CreateClasses(ctx, teacherID, req.SubjectName)
	if err != nil {
		if errors.Is(err, class.ErrSubjectNameRequired) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error"})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"id": id})
}
