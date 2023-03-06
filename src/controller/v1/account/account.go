package account

import (
	"log"
	"strings"
	"net/http"

	echo "github.com/labstack/echo/v4"
	entity "go-rest-api/src/http"
	"github.com/forkyid/go-utils/v1/validation"
	"github.com/pkg/errors"
	"go-rest-api/src/constant"
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

// @Summary Get User Data
// @Description Get User Data
// @Tags Accounts
// @Produce application/json
// @Param Authorization header string true "Bearer Token"
// @Success 200 {object} http.GetUser
// @Failure 400 {string} string "Bad Request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal Server Error"
// @Router /v1/accounts [get]
func (ctrl *Controller) Get(ctx echo.Context) {
	accountID, err := jwt.ExtractID(ctx.Request().Header.Get("Authorization"))
	if err != nil {
		rest.ResponseMessage(ctx, http.StatusUnauthorized)
		return
	}

	response, err := ctrl.svc.TakeAccountByID(accountID)
	if err != nil {
		rest.ResponseMessage(ctx, http.StatusInternalServerError)
		log.Println("get account by id:", err)
		return
	}

	rest.ResponseData(ctx, http.StatusOK, response)
	return
}

// Register godoc
// @Summary Register Account
// @Description Register Account
// @Tags Accounts
// @Param Payload body http.RegisterUser true "Payload"
// @Success 201 {object} string "Created"
// @Failure 400 {string} string "Bad Request"
// @Failure 409 {string} string "Resource Conflict"
// @Failure 500 {string} string "Internal Server Error"
// @Router /v1/accounts/register [post]
func (ctrl *Controller) Register(ctx echo.Context) {
	req := new(entity.RegisterUser)
	if err := ctx.Bind(req); err != nil {
		rest.ResponseError(ctx, http.StatusBadRequest, map[string]string{
			"body": constant.ErrInvalidFormat.Error()})
		return
	}

	// required tapi tidak diisi akan return bad request
	if err := validation.Validator.Struct(req); err != nil {
		log.Println("validate struct:", err, "request:", req)
		rest.ResponseError(ctx, http.StatusBadRequest, err)
		return
	}

	request := *req
	if request.Username == "" {
		rest.ResponseError(ctx, http.StatusBadRequest, map[string]string{
			"username": constant.ErrInvalidFormat.Error()})
		return
	}
	request.Username = strings.ToLower(request.Username)
	err := ctrl.svc.Create(request)
	if errors.Is(err, constant.ErrAccountExist) {
		rest.ResponseError(ctx, http.StatusConflict, map[string]string{
			"account": constant.ErrAccountExist.Error()})
	} else if err != nil {
		log.Println("register:", err.Error())
		rest.ResponseMessage(ctx, http.StatusInternalServerError)
	} else {
		rest.ResponseMessage(ctx, http.StatusCreated)
	}
}

// Update godoc
// @Summary Update Account
// @Description Update Account
// @Tags Accounts
// @Param Authorization header string true "Bearer Token"
// @Param Payload body http.UpdateUser true "Payload"
// @Success 200 {string} string "Success"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /v1/accounts [patch]
func (ctrl *Controller) Update(ctx echo.Context) {
	req := new(entity.UpdateUser)
	err := ctx.Bind(req)
	if err != nil {
		log.Println("bind json:", err, "request:", req)
		rest.ResponseError(ctx, http.StatusBadRequest, map[string]string{
			"body": constant.ErrInvalidFormat.Error()})
		return
	}

	// required tapi tidak diisi akan return bad request
	if err := validation.Validator.Struct(req); err != nil {
		log.Println("validate struct:", err, "request:", req)
		rest.ResponseError(ctx, http.StatusBadRequest, err)
		return
	}

	accountID, err := jwt.ExtractID(ctx.Request().Header.Get("Authorization"))
	if err != nil {
		rest.ResponseMessage(ctx, http.StatusUnauthorized)
		return
	}

	request := *req
	*request.Username = strings.ToLower(*request.Username)
	err = ctrl.svc.Update(accountID, request)
	if err != nil {
		if errors.Is(err, constant.ErrAccountNotRegistered) {
			rest.ResponseError(ctx, http.StatusBadRequest, map[string]string{
				"accounts": constant.ErrAccountNotRegistered.Error()})
			return
		} else if errors.Is(err, constant.ErrUsernameCannotBeEmpty) {
			rest.ResponseError(ctx, http.StatusBadRequest, map[string]string{
				"accounts": constant.ErrUsernameCannotBeEmpty.Error()})
			return
		} else if errors.Is(err, constant.ErrPasswordCannotBeEmpty) {
			rest.ResponseError(ctx, http.StatusBadRequest, map[string]string{
				"accounts": constant.ErrPasswordCannotBeEmpty.Error()})
			return
		} else if errors.Is(err, constant.ErrUsernameAlreadyExist) {
			rest.ResponseError(ctx, http.StatusBadRequest, map[string]string{
				"accounts": constant.ErrUsernameAlreadyExist.Error()})
			return
		} else if errors.Is(err, constant.ErrEmailAlreadyExist) {
			rest.ResponseError(ctx, http.StatusBadRequest, map[string]string{
				"accounts": constant.ErrEmailAlreadyExist.Error()})
			return
		} else if errors.Is(err, constant.ErrKTPNumberAlreadyExist) {
			rest.ResponseError(ctx, http.StatusBadRequest, map[string]string{
				"accounts": constant.ErrKTPNumberAlreadyExist.Error()})
			return
		} else if errors.Is(err, constant.ErrPhoneNumberAlreadyExist) {
			rest.ResponseError(ctx, http.StatusBadRequest, map[string]string{
				"accounts": constant.ErrPhoneNumberAlreadyExist.Error()})
			return
		} else if errors.Is(err, constant.ErrInvalidDOBFormat) {
			rest.ResponseError(ctx, http.StatusBadRequest, map[string]string{
				"accounts": constant.ErrInvalidDOBFormat.Error()})
			return
		}
		rest.ResponseMessage(ctx, http.StatusInternalServerError)
		log.Println("update account: ", err.Error())
		return
	}

	rest.ResponseMessage(ctx, http.StatusOK)
}

// Delete godoc
// @Summary Delete Account
// @Description Delete Account By User Itself
// @Tags Accounts
// @Param Authorization header string true "Bearer Token"
// @Success 200 {string} string "Success"
// @Failure 400 {string} string "Bad Request"
// @Failure 409 {string} string "Resource Conflict"
// @Failure 500 {string} string "Internal Server Error"
// @Router /v1/accounts [delete]
func (ctrl *Controller) Delete(ctx echo.Context) {
	accountID, err := jwt.ExtractID(ctx.Request().Header.Get("Authorization"))
	if err != nil {
		rest.ResponseMessage(ctx, http.StatusUnauthorized)
		return
	}

	err = ctrl.svc.Delete(accountID)
	if err != nil {
		if errors.Is(err, constant.ErrAccountNotRegistered) {
			rest.ResponseError(ctx, http.StatusBadRequest, map[string]string{
				"accounts": constant.ErrAccountNotRegistered.Error()})
			return
		}
		rest.ResponseMessage(ctx, http.StatusInternalServerError)
		log.Println("delete account: ", err.Error())
		return
	}
		
	rest.ResponseMessage(ctx, http.StatusOK)
}
