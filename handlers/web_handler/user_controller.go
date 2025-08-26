package web_handler

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"laotop_final/environment"
	"laotop_final/handlers"
	"laotop_final/middlewares"
	"laotop_final/usecase/user_service"
	"laotop_final/utils"
	"laotop_final/validation"
	"time"
)

type UserController interface {
	GetPermission(ctx *fiber.Ctx) error

	GetRole(ctx *fiber.Ctx) error
	CreatedRole(ctx *fiber.Ctx) error
	UpdateRole(ctx *fiber.Ctx) error
	DeleteRole(ctx *fiber.Ctx) error
	DeleteUser(ctx *fiber.Ctx) error

	SignIn(ctx *fiber.Ctx) error
	GetUserInfos(ctx *fiber.Ctx) error
	ChangePassword(ctx *fiber.Ctx) error
	SignOut(ctx *fiber.Ctx) error

	GetUser(ctx *fiber.Ctx) error
	CreateUser(ctx *fiber.Ctx) error
	UpdateUser(ctx *fiber.Ctx) error
	ResetPasswordUser(ctx *fiber.Ctx) error
}

type userController struct {
	serviceUser user_service.UserService
}

func (c *userController) GetPermission(ctx *fiber.Ctx) error {
	response, err := c.serviceUser.GetPermission()
	if err != nil {
		return handlers.NewErrorResponses(ctx, err)
	}
	return handlers.NewSuccessResponse(ctx, response)
}

func (c *userController) GetRole(ctx *fiber.Ctx) error {
	response, err := c.serviceUser.GetRole()
	if err != nil {
		return handlers.NewErrorResponses(ctx, err)
	}
	return handlers.NewSuccessResponse(ctx, response)
}

func (c *userController) CreatedRole(ctx *fiber.Ctx) error {
	req := new(user_service.RoleCreate)
	_ = ctx.BodyParser(req)
	errValidate := validation.Validate(req)
	if errValidate != nil {
		return handlers.NewErrorValidate(ctx, errValidate[0].Error)
	}
	err := c.serviceUser.CreateRole(*req)
	if err != nil {
		return handlers.NewErrorResponses(ctx, err)
	}
	return handlers.NewSuccessResponse(ctx, "success")
}

func (c *userController) UpdateRole(ctx *fiber.Ctx) error {
	req := new(user_service.RoleUpdate)
	_ = ctx.BodyParser(req)
	errValidate := validation.Validate(req)
	if errValidate != nil {
		return handlers.NewErrorValidate(ctx, errValidate[0].Error)
	}
	err := c.serviceUser.UpdateRole(*req)
	if err != nil {
		return handlers.NewErrorResponses(ctx, err)
	}
	return handlers.NewSuccessResponse(ctx, "success")
}

func (c *userController) DeleteRole(ctx *fiber.Ctx) error {
	req := new(user_service.RoleDelete)
	_ = ctx.BodyParser(req)
	errValidate := validation.Validate(req)
	if errValidate != nil {
		return handlers.NewErrorValidate(ctx, errValidate[0].Error)
	}
	err := c.serviceUser.DeleteRole(*req)
	if err != nil {
		return handlers.NewErrorResponses(ctx, err)
	}
	return handlers.NewSuccessMsg(ctx)
}

func (c *userController) DeleteUser(ctx *fiber.Ctx) error {
	req := new(user_service.UserDelect)
	if err := ctx.BodyParser(req); err != nil {
		return handlers.NewErrorResponses(ctx, err)
	}

	errValidate := validation.Validate(req)
	if errValidate != nil {
		return handlers.NewErrorValidate(ctx, errValidate[0].Error)
	}

	err := c.serviceUser.DeleteUser(*req)
	if err != nil {
		return handlers.NewErrorResponses(ctx, err)
	}
	return handlers.NewSuccessMsg(ctx)
}

func (c *userController) SignIn(ctx *fiber.Ctx) error {
	req := new(user_service.SignIn)
	if err := ctx.BodyParser(req); err != nil {
		return handlers.NewErrorResponses(ctx, err)
	}
	errValidate := validation.Validate(req)
	if errValidate != nil {
		return handlers.NewErrorValidate(ctx, errValidate[0].Error)
	}
	response, token, err := c.serviceUser.SignIn(*req)
	if err != nil {
		return handlers.NewErrorResponses(ctx, err)
	}

	ctx.Cookie(&fiber.Cookie{
		Name:     environment.VITE_WEB_HEADER_NAME,
		Value:    token,
		Path:     "/",
		HTTPOnly: false,
		Secure:   false,
		SameSite: "Lax",
		Expires:  time.Now().Add(24 * time.Hour),
	})

	return handlers.NewSuccessResponseSignIn(ctx, response, token)
}

func (c *userController) GetUserInfos(ctx *fiber.Ctx) error {
	id := middlewares.GetUserID(ctx)
	fmt.Printf("err %v", id)
	response, err := c.serviceUser.GetUserInfos(user_service.RoleDelete{
		ID: utils.AtoI(id),
	})
	if err != nil {
		return handlers.NewErrorResponses(ctx, err)
	}
	return handlers.NewSuccessResponse(ctx, response)
}

func (c *userController) ChangePassword(ctx *fiber.Ctx) error {
	req := new(user_service.ChangePassword)
	_ = ctx.BodyParser(req)
	errValidate := validation.Validate(req)
	if errValidate != nil {
		return handlers.NewErrorValidate(ctx, errValidate[0].Error)
	}
	id := middlewares.GetUserID(ctx)
	req.ID = utils.AtoI(id)
	response, err := c.serviceUser.ChangePassword(*req)
	if err != nil {
		return handlers.NewErrorResponses(ctx, err)
	}
	return handlers.NewSuccessResponse(ctx, response)
}
func (c *userController) SignOut(ctx *fiber.Ctx) error {
	strKey := middlewares.GetUserRoleID(ctx)
	key := middlewares.GetUserID(ctx)
	response, err := c.serviceUser.SignOut(strKey, key)
	if err != nil {
		return handlers.NewErrorResponses(ctx, err)
	}
	return handlers.NewSuccessResponse(ctx, response)
}

func (c *userController) GetUser(ctx *fiber.Ctx) error {
	response, err := c.serviceUser.GetUser()
	if err != nil {
		return handlers.NewErrorResponses(ctx, err)
	}
	return handlers.NewSuccessResponse(ctx, response)
}

func (c *userController) CreateUser(ctx *fiber.Ctx) error {
	req := new(user_service.CreateUserRequest)
	_ = ctx.BodyParser(req)
	errValidate := validation.Validate(req)
	if errValidate != nil {
		return handlers.NewErrorValidate(ctx, errValidate[0].Error)
	}
	response, err := c.serviceUser.CreateUser(*req)
	if err != nil {
		return handlers.NewErrorResponses(ctx, err)
	}
	return handlers.NewSuccessResponse(ctx, response)
}

func (c *userController) UpdateUser(ctx *fiber.Ctx) error {
	req := new(user_service.Update)
	_ = ctx.BodyParser(req)
	errValidate := validation.Validate(req)
	if errValidate != nil {
		return handlers.NewErrorValidate(ctx, errValidate[0].Error)
	}
	response, err := c.serviceUser.UpdateUser(*req)
	if err != nil {
		return handlers.NewErrorResponses(ctx, err)
	}
	return handlers.NewSuccessResponse(ctx, response)
}

func (c *userController) ResetPasswordUser(ctx *fiber.Ctx) error {
	req := new(user_service.ResetPassword)
	_ = ctx.BodyParser(req)
	errValidate := validation.Validate(req)
	if errValidate != nil {
		return handlers.NewErrorValidate(ctx, errValidate[0].Error)
	}
	id := middlewares.GetUserID(ctx)
	if utils.AtoI(id) == req.ID {
		return handlers.NewErrorResponses(ctx, errors.New("can not reset password to the same user_service"))
	}
	response, err := c.serviceUser.ResetPassword(*req)
	if err != nil {
		return handlers.NewErrorResponses(ctx, err)
	}
	return handlers.NewSuccessResponse(ctx, response)
}

func NewUserController(
	serviceUser *user_service.UserService,
) UserController {
	return &userController{
		serviceUser: *serviceUser,
	}
}
