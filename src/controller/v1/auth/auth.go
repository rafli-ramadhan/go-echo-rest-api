package auth

import (
	"fmt"
	"log"
	"net/http"

	echo "github.com/labstack/echo/v4"
	entity "go-rest-api/src/http"
	"github.com/forkyid/go-utils/v1/aes"
	"github.com/forkyid/go-utils/v1/validation"
	"github.com/pkg/errors"
	"go-rest-api/src/constant"
	"go-rest-api/src/pkg/bcrypt"
	"go-rest-api/src/pkg/jwt"
	"go-rest-api/src/pkg/rest"
	"go-rest-api/src/service/v1/account"
)

type Controller struct {
	svc account.Servicer
}

func NewController(
	servicer account.Servicer,
) *Controller {
	return &Controller{
		svc: servicer,
	}
}

// @Summary User Login
// @Description User Login
// @Tags Auth
// @Produce application/json
// @Param Payload body http.Auth true "Payload"
// @Success 200 {object} http.Token
// @Failure 400 {string} string "Bad Request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal Server Error"
// @Router /v1/auth [post]
func (ctrl *Controller) Login(ctx echo.Context) {
	req := new(entity.Auth)
	err := ctx.Bind(req)
	if err != nil {
		log.Println("bind json:", err, "request:", req)
		rest.ResponseError(ctx, http.StatusBadRequest, map[string]string{
			"body": constant.ErrInvalidFormat.Error()})
		return
	}

	if err := validation.Validator.Struct(req); err != nil {
		log.Println("validate struct:", err, "request:", req)
		rest.ResponseError(ctx, http.StatusBadRequest, err)
		return
	}

	exist, _ := ctrl.svc.CheckAccountByUsername(req.Username)
	if !exist {
		rest.ResponseError(ctx, http.StatusBadRequest, map[string]string{
			"accounts": constant.ErrAccountNotRegistered.Error()})
		return
	}

	account, err:= ctrl.svc.TakeAccountByUsername(req.Username)
	if err != nil {
		rest.ResponseMessage(ctx, http.StatusInternalServerError)
		return
	}

	err = bcrypt.ComparePassword(account.Password, req.Password)
	if err != nil {
		rest.ResponseError(ctx, http.StatusBadRequest, map[string]string{
			"accounts": constant.ErrInvalidPassword.Error()})
		return
	}

	token, err := jwt.GenerateJWT(aes.Encrypt(int(account.ID)))
	if err != nil {
		rest.ResponseMessage(ctx, http.StatusInternalServerError)
		return
	}

	rest.ResponseData(ctx, http.StatusOK, entity.Token{
		Token: fmt.Sprintf("Bearer %v", token),
	})
}

// @Summary Update User Password
// @Description Update User Password
// @Tags Auth
// @Produce application/json
// @Param Payload body http.ForgotPassword true "Payload"
// @Success 200 {string} string "Success"
// @Failure 400 {string} string "Bad Request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal Server Error"
// @Router /v1/auth/forgot [patch]
func (ctrl *Controller) ForgotPassword(ctx echo.Context) {
	req := new(entity.ForgotPassword)
	err := ctx.Bind(req)
	if err != nil {
		log.Println("bind json:", err, "request:", req)
		rest.ResponseError(ctx, http.StatusBadRequest, map[string]string{
			"body": constant.ErrInvalidFormat.Error()})
		return
	}

	if err := validation.Validator.Struct(req); err != nil {
		log.Println("validate struct:", err, "request:", req)
		rest.ResponseError(ctx, http.StatusBadRequest, err)
		return
	}

	request := *req
	err = ctrl.svc.UpdatePassword(request)
	if err != nil {
		if errors.Is(err, constant.ErrPasswordCannotBeEmpty) {
			rest.ResponseError(ctx, http.StatusBadRequest, map[string]string{
				"accounts": constant.ErrPasswordCannotBeEmpty.Error()})
			return
		}
		rest.ResponseMessage(ctx, http.StatusInternalServerError)
		log.Println("update account password:", err)
		return
	}

	rest.ResponseMessage(ctx, http.StatusOK)
}