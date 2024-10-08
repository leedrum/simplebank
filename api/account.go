package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/leedrum/simplebank/db/sqlc"
	"github.com/leedrum/simplebank/token"
)

type CreateAccountRequest struct {
	Currency string `json:"currency" binding:"required,currency"`
}

func (server *Server) createAccount(ctx *gin.Context) {
	var req CreateAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorHandler(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayload).(*token.Payload)
	arg := db.CreateAccountParams{
		Owner:    authPayload.Username,
		Currency: req.Currency,
		Balance:  0,
	}

	account, err := server.store.CreateAccount(ctx, arg)
	if err != nil {
		if db.ErrorCode(err) == db.UniqueViolation || db.ErrorCode(err) == db.ForeignKeyViolation {
			ctx.JSON(http.StatusForbidden, errorHandler(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorHandler(err))
		return
	}

	ctx.JSON(http.StatusCreated, account)
}

type getAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getAccount(ctx *gin.Context) {
	var req getAccountRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorHandler(err))
		return
	}

	account, err := server.store.GetAccount(ctx, req.ID)
	fmt.Printf("account: %v\n", account)
	fmt.Printf("err: %v\n", err)
	if err != nil {
		if errors.Is(err, db.ErrorRecordNotFound) {
			ctx.JSON(http.StatusNotFound, errorHandler(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorHandler(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayload).(*token.Payload)
	if account.Owner != authPayload.Username {
		ctx.JSON(http.StatusUnauthorized, errorHandler(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

type listAccountsRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) listAccounts(ctx *gin.Context) {
	var req listAccountsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorHandler(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayload).(*token.Payload)
	arg := db.ListAccountsParams{
		Owner:  authPayload.Username,
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}
	accounts, err := server.store.ListAccounts(ctx, arg)
	if err != nil {
		if errors.Is(err, db.ErrorRecordNotFound) {
			ctx.JSON(http.StatusNotFound, errorHandler(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorHandler(err))
		return
	}

	ctx.JSON(http.StatusOK, accounts)
}

type updateAccountRequest struct {
	Balance int64 `json:"balance" binding:"required"`
}

func (server *Server) updateAccount(ctx *gin.Context) {
	var req_get getAccountRequest
	if err := ctx.ShouldBindUri(&req_get); err != nil {
		ctx.JSON(http.StatusBadRequest, errorHandler(err))
		return
	}

	var req_update updateAccountRequest
	if err := ctx.ShouldBindJSON(&req_update); err != nil {
		ctx.JSON(http.StatusBadRequest, errorHandler(err))
		return
	}

	arg := db.UpdateAccountParams{
		ID:      req_get.ID,
		Balance: req_update.Balance,
	}

	_, err := server.store.GetAccount(ctx, arg.ID)
	if err != nil {
		if errors.Is(err, db.ErrorRecordNotFound) {
			ctx.JSON(http.StatusNotFound, errorHandler(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorHandler(err))
		return
	}

	account, err := server.store.UpdateAccount(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorHandler(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

type deleteAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) deleteAccount(ctx *gin.Context) {
	var req deleteAccountRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorHandler(err))
		return
	}

	_, err := server.store.GetAccount(ctx, req.ID)
	if err != nil {
		if errors.Is(err, db.ErrorRecordNotFound) {
			ctx.JSON(http.StatusNotFound, errorHandler(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorHandler(err))
		return
	}

	err = server.store.DeleteAccount(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorHandler(err))
		return
	}

	ctx.JSON(http.StatusOK, nil)
}
