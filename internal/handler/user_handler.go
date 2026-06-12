package handler

import (
	"strconv"
	"user-api/internal/logger"
	"user-api/internal/models"
	"user-api/internal/service"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type UserHandler struct {
	service  *service.UserService
	validate *validator.Validate
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{
		service:  svc,
		validate: validator.New(),
	}
}

func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	var req models.CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "invalid request body"})
	}

	if err := h.validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: err.Error()})
	}

	logger.Log.Info("creating user", zap.String("name", req.Name), zap.String("dob", req.Dob))

	user, err := h.service.CreateUser(c.Context(), req)
	if err != nil {
		logger.Log.Error("failed to create user", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: err.Error()})
	}

	logger.Log.Info("user created", zap.Int32("id", user.ID))
	return c.Status(fiber.StatusCreated).JSON(user)
}

func (h *UserHandler) GetUserByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id < 1 {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "invalid user id"})
	}

	logger.Log.Info("fetching user", zap.Int("id", id))

	user, err := h.service.GetUserByID(c.Context(), int32(id))
	if err != nil {
		logger.Log.Warn("user not found", zap.Int("id", id))
		return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{Error: err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(user)
}

func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id < 1 {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "invalid user id"})
	}

	var req models.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "invalid request body"})
	}

	if err := h.validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: err.Error()})
	}

	logger.Log.Info("updating user", zap.Int("id", id))

	user, err := h.service.UpdateUser(c.Context(), int32(id), req)
	if err != nil {
		logger.Log.Error("failed to update user", zap.Int("id", id), zap.Error(err))
		return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{Error: err.Error()})
	}

	logger.Log.Info("user updated", zap.Int("id", id))
	return c.Status(fiber.StatusOK).JSON(user)
}

func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id < 1 {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "invalid user id"})
	}

	logger.Log.Info("deleting user", zap.Int("id", id))

	if err := h.service.DeleteUser(c.Context(), int32(id)); err != nil {
		logger.Log.Warn("user not found for delete", zap.Int("id", id))
		return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{Error: err.Error()})
	}

	logger.Log.Info("user deleted", zap.Int("id", id))
	return c.SendStatus(fiber.StatusNoContent)
}

func (h *UserHandler) ListUsers(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	logger.Log.Info("listing users", zap.Int("page", page), zap.Int("limit", limit))

	result, err := h.service.ListUsers(c.Context(), page, limit)
	if err != nil {
		logger.Log.Error("failed to list users", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Error: err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(result)
}
