package handlers

import (
	"restapi/database"
	"restapi/types"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	validate = validator.New()
)

func ValidateTask(task types.Task) []*types.ValidateErrorResponse {
	var errors []*types.ValidateErrorResponse
	err := validate.Struct(task)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element types.ValidateErrorResponse
			element.FailedField = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}
	return errors
}

func CreateTask(ctx *fiber.Ctx) error {
	task := types.Task{}

	if err := ctx.BodyParser(&task); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	errors := ValidateTask(task)
	if errors != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(errors)
	}

	task, err := database.InsertTask(task)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return ctx.Status(fiber.StatusConflict).JSON(types.ErrorResponse{
				Success: false,
				Message: "Duplicate",
			})
		}

		return ctx.Status(fiber.StatusInternalServerError).JSON(err)
	}

	return ctx.JSON(task)
}


func GetTasks(ctx *fiber.Ctx) error {
	tasks, err := database.AllTasks()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(err)
	}

	return ctx.JSON(tasks)
}

func GetTask(ctx *fiber.Ctx) error {
	id := ctx.Params("id", "")
	uuid, err := uuid.Parse(id)
	if err != nil {
		return ctx.JSON(types.ErrorResponse{
			Success: false,
			Message: "UUID is incorrect",
		})
	}

	task, err := database.FindTask(uuid)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(err)
	}

	return ctx.JSON(task)
}

func UpdateTask(ctx *fiber.Ctx) error {
	id := ctx.Params("id", "")
	uuid, err := uuid.Parse(id)
	if err != nil {
		return ctx.JSON(types.ErrorResponse{
			Success: false,
			Message: "UUID is incorrect",
		})
	}

	task := types.Task{}

	if err := ctx.BodyParser(&task); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	errors := ValidateTask(task)
	if errors != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(errors)
	}

	task, err = database.UpdateTask(uuid, task)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(err)
	}

	return ctx.JSON(task)
}

func DeleteTask(ctx *fiber.Ctx) error {
	id := ctx.Params("id", "")
	uuid, err := uuid.Parse(id)
	if err != nil {
		return ctx.JSON(types.ErrorResponse{
			Success: false,
			Message: "UUID is incorrect",
		})
	}

	task, err := database.DeleteTask(uuid)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(err)
	}

	return ctx.JSON(task)
}
