package v1

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sixojke/test-astral/domain"
)

type authUserInp struct {
	Login    string `json:"login"`
	Password string `json:"pswd"`
}

func (a *authUserInp) validate() error {
	// INFO: remove validation if login or password verification in registration is changed
	if !validateLogin(a.Login) {
		return domain.ErrInvalidLogin
	}

	if !validatePassword(a.Password) {
		return domain.ErrInvalidPassword
	}

	return nil
}

type authUserResponse struct {
	Token string `json:"token"`
}

// @Summary Auth user
// @Tags auth
// @Description User login
// @ModuleID authUser
// @Accept json
// @Produce json
// @Param input body authUserInp true "Register info"
// @Success 200 {object} swagResponse{response=authUserResponse} "Successful login"
// @Failure 400 {object} swagError "Bad Request"
// @Failure 500 {object} swagError "Internal Server Error"
// @Router /auth [post]
func (h *Handler) authUser(c *gin.Context) {
	var inp authUserInp
	if err := c.BindJSON(&inp); err != nil {
		errResponse(c, http.StatusBadRequest, err.Error(), domain.ErrCantParseJSON.Error())

		return
	}

	if err := inp.validate(); err != nil {
		errResponse(c, http.StatusBadRequest, err.Error(), err.Error())

		return
	}

	token, err := h.service.SignIn(inp.Login, inp.Password)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			errResponse(c, http.StatusBadRequest, err.Error(), err.Error())
		} else {
			errResponse(c, http.StatusInternalServerError, err.Error(), domain.ErrInternalServerError.Error())
		}

		return
	}

	newResponse(c, http.StatusOK, nil, authUserResponse{
		Token: token,
	})
}

// @Summary Delete session by token
// @Tags auth
// @Description Delete session by token
// @ModuleID deleteSession
// @Accept json
// @Produce json
// @Param token path string true "Session token"
// @Success 200 {object} swagResponse{response=map[string]bool} "Success"
// @Failure 400 {object} swagError "Bad Request"
// @Failure 500 {object} swagError "Internal Server Error"
// @Router /auth/{token} [delete]
func (h *Handler) deleteSession(c *gin.Context) {
	token := c.Param("token")

	if token == "" {
		errResponse(c, http.StatusBadRequest, domain.ErrParameterIsEmpty.Error(), domain.ErrParameterIsEmpty.Error())

		return
	}

	if err := h.service.User.DeleteSession(token); err != nil {
		errResponse(c, http.StatusInternalServerError, err.Error(), domain.ErrInternalServerError.Error())

		return
	}

	newResponse(c, http.StatusOK, nil, map[string]bool{
		token: true,
	})
}
