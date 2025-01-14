package user

import (
	"net/url"

	"github.com/gofiber/fiber/v2"
	domain "github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/server/http/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/service"
	"github.com/pkg/errors"
)

// @Summary	Request to login
// @Descriptio	Request to login to the system. Returns a Login URL
// @Tags		User
// @Produce	json
// @Param		redirect_url	query		string	true	"Redirect URL"
// @Success	200				{object}	entities.LoginResponse
// @Failure	400				{object}	HTTPError
// @Failure	500				{object}	HTTPError
// @Router		/v1/user/login [get]
func Login(svc service.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.Context()
		redirectURL, err := url.ParseRequestURI(c.Query("redirect_url"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(service.NewError(service.BadRequest, errors.Wrap(err, "failed to parse redirect url").Error()))
		}

		req := domain.LoginRequest{
			RedirectURL: redirectURL,
		}

		resp, err := svc.LoginRequest(ctx, &req)
		if err != nil {
			return err
		}

		response := entities.LoginResponse{
			LoginURL: resp.LoginURL.String(),
		}

		return c.JSON(response)
	}
}

// @Summary	Logout from the system
// @Descriptio	Logout from the system
// @Tags		User
// @Param		body	body		entities.LogoutRequest	true	"Logout information"
// @Success	200		{string}	string					"OK"
// @Failure	400		{object}	HTTPError
// @Failure	500		{object}	HTTPError
// @Router		/v1/user/logout [post]
func Logout(svc service.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.Context()
		req := entities.LogoutRequest{}
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(service.NewError(service.BadRequest, errors.Wrap(err, "failed to parse request").Error()))
		}

		domainReq := domain.Logout{
			RefreshToken: req.RefreshToken,
		}

		err := svc.LogoutRequest(ctx, &domainReq)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(service.NewError(service.InternalError, errors.Wrap(err, "failed to logout").Error()))
		}

		return c.SendStatus(fiber.StatusOK)
	}
}

// @Summary	Validate login code and request a access token
// @Descriptio	Validate login code and request a access token
// @Tags		User
// @Accept		json
// @Produce	json
// @Param		body			body		entities.LoginTokenRequest	true	"Callback information"
// @Param		redirect_url	query		string						true	"Redirect URL"
// @Success	200				{object}	entities.ClientTokenResponse
// @Failure	400				{object}	HTTPError
// @Failure	500				{object}	HTTPError
// @Router		/v1/user/login/token [post]
func RequestToken(svc service.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.Context()
		req := entities.LoginTokenRequest{}
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(service.NewError(service.BadRequest, errors.Wrap(err, "failed to parse request").Error()))
		}

		redirectURL, err := parseURL(c.Query("redirect_url"))
		if err != nil {
			return err
		}

		domainReq := domain.LoginCallback{
			Code:        req.Code,
			RedirectURL: redirectURL,
		}

		token, err := svc.ClientTokenCallback(ctx, &domainReq)
		if err != nil {
			return err
		}

		response := entities.ClientTokenResponse{
			AccessToken:  token.AccessToken,
			ExpiresIn:    token.ExpiresIn,
			RefreshToken: token.RefreshToken,
			TokenType:    token.TokenType,
		}

		return c.JSON(response)
	}
}

// @Summary		Register a new user
// @Description	Register a new user
// @Tags			User
// @Accept			json
// @Produce		json
// @Param			user	body		entities.UserRegisterRequest	true	"User information"
// @Success		201		{object}	entities.UserResponse
// @Failure		400		{object}	HTTPError
// @Failure		500		{object}	HTTPError
// @Router			/v1/user [post]
// @Param			Authorization	header	string	true	"Insert your access token"	default(Bearer <Add access token here>)
func Register(svc service.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.Context()
		req := entities.UserRegisterRequest{}
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": errors.Wrap(err, "failed to parse request").Error(),
			})
		}

		domainUser := domain.RegisterUser{
			User: domain.User{
				Email:     req.Email,
				FirstName: req.FirstName,
				LastName:  req.LastName,
				Username:  req.Username,
			},
			Password: req.Password,
			Roles:    req.Roles,
		}

		u, err := svc.Register(ctx, &domainUser)
		if err != nil {
			return err
		}

		response := entities.UserResponse{
			ID:            u.ID.String(),
			CreatedAt:     u.CreatedAt,
			Email:         u.Email,
			FirstName:     u.FirstName,
			LastName:      u.LastName,
			Username:      u.Username,
			EmployeeID:    u.EmployeeID,
			PhoneNumber:   u.PhoneNumber,
			EmailVerified: u.EmailVerified,
		}

		return c.Status(fiber.StatusCreated).JSON(response)
	}
}

func parseURL(rawURL string) (*url.URL, error) {
	return url.ParseRequestURI(rawURL)
}

// @Summary		Get all users
// @Description	Get all users
// @Tags			User
// @Produce		json
// @Success		200		{object}	entities.UserListResponse
// @Failure		400		{object}	HTTPError
// @Failure		500		{object}	HTTPError
// @Param			page	query		string	false	"Page"
// @Param			limit	query		string	false	"Limit"
// @Router			/v1/user [get]
// @Param			Authorization	header	string	true	"Insert your access token"	default(Bearer <Add access token here>)
func GetAllUsers(_ service.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusNotImplemented)
	}
}

// @Summary		Get a user by ID
// @Description	Get a user by ID
// @Tags			User
// @Produce		json
// @Success		200		{object}	entities.UserResponse
// @Failure		400		{object}	HTTPError
// @Failure		404		{object}	HTTPError
// @Failure		500		{object}	HTTPError
// @Param			user_id	path		string	true	"User ID"
// @Router			/v1/user/{user_id} [get]
// @Param			Authorization	header	string	true	"Insert your access token"	default(Bearer <Add access token here>)
func GetUserByID(_ service.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusNotImplemented)
	}
}

// @Summary		Update a user by ID
// @Description	Update a user by ID
// @Tags			User
// @Accept			json
// @Produce		json
// @Success		200		{object}	entities.UserResponse
// @Failure		400		{object}	HTTPError
// @Failure		404		{object}	HTTPError
// @Failure		500		{object}	HTTPError
// @Param			user_id	path		string						true	"User ID"
// @Param			user	body		entities.UserUpdateRequest	true	"User information"
// @Router			/v1/user/{user_id} [put]
// @Param			Authorization	header	string	true	"Insert your access token"	default(Bearer <Add access token here>)
func UpdateUserByID(_ service.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusNotImplemented)
	}
}

// @Summary		Delete a user by ID
// @Description	Delete a user by ID
// @Tags			User
// @Produce		json
// @Success		200		{string}	string
// @Failure		400		{object}	HTTPError
// @Failure		404		{object}	HTTPError
// @Failure		500		{object}	HTTPError
// @Param			user_id	path		string	true	"User ID"
// @Router			/v1/user/{user_id} [delete]
// @Param			Authorization	header	string	true	"Insert your access token"	default(Bearer <Add access token here>)
func DeleteUserByID(_ service.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusNotImplemented)
	}
}

// @Summary		Get user roles
// @Description	Get user roles
// @Tags			User
// @Produce		json
// @Success		200		{object}	entities.RoleListResponse
// @Failure		400		{object}	HTTPError
// @Failure		500		{object}	HTTPError
// @Param			user_id	path		string	true	"User ID"
// @Param			page	query		string	false	"Page"
// @Param			limit	query		string	false	"Limit"
// @Router			/v1/user/{user_id}/roles [get]
// @Param			Authorization	header	string	true	"Insert your access token"	default(Bearer <Add access token here>)
func GetUserRoles(_ service.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusNotImplemented)
	}
}
