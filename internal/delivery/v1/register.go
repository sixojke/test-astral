package v1

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sixojke/test-astral/domain"
)

type registerUserInp struct {
	Token    string `json:"token"`
	Login    string `json:"login"`
	Password string `json:"pswd"`
}

func (r *registerUserInp) validate() error {
	if r.Token == "" {
		return domain.ErrInvalidToken
	}

	if !validateLogin(r.Login) {
		return domain.ErrInvalidLogin
	}

	if !validatePassword(r.Password) {
		return domain.ErrInvalidPassword
	}

	return nil
}

type registerUserResponse struct {
	Login string `json:"login"`
}

// @Summary Register user
// @Tags auth
// @Description Create user account
// @ModuleID registerUser
// @Accept json
// @Produce json
// @Param input body registerUserInp true "Register info"
// @Success 200 {object} swagResponse{response=registerUserResponse} "Successful registration"
// @Failure 400 {object} swagError "Bad Request"
// @Failure 500 {object} swagError "Internal Server Error"
// @Router /register [post]
func (h *Handler) registerUser(c *gin.Context) {
	var inp registerUserInp
	if err := c.BindJSON(&inp); err != nil {
		errResponse(c, http.StatusBadRequest, err.Error(), domain.ErrCantParseJSON.Error())

		return
	}

	if err := inp.validate(); err != nil {
		errResponse(c, http.StatusBadRequest, err.Error(), err.Error())

		return
	}

	if err := h.service.User.SignUp(inp.Token, inp.Login, inp.Password); err != nil {
		if errors.Is(err, domain.ErrInvalidToken) || errors.Is(err, domain.ErrLoginIsBusy) {
			errResponse(c, http.StatusBadRequest, err.Error(), err.Error())
		} else {
			errResponse(c, http.StatusInternalServerError, err.Error(), domain.ErrInternalServerError.Error())
		}

		return
	}

	newResponse(c, http.StatusOK, nil, registerUserResponse{Login: inp.Login})
}
