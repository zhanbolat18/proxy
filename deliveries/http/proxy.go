package http

import (
	"github.com/gin-gonic/gin"
	"github.com/zhanbolat18/proxy/internal/proxy"
	"github.com/zhanbolat18/proxy/internal/proxy/entities"
	"net/http"
)

type ProxyController struct {
	uc proxy.ProxyUsecase
}

func NewProxyController(uc proxy.ProxyUsecase) *ProxyController {
	return &ProxyController{uc: uc}
}

func (p *ProxyController) Proxy(ctx *gin.Context) {
	req := entities.Request{}
	err := ctx.BindJSON(&req)
	if err != nil {
		_ = ctx.Error(err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	res, err := p.uc.Proxy(ctx, req)
	if err != nil {
		_ = ctx.Error(err)
		// @TODO need handle error and set correct status by error
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	ctx.JSON(http.StatusOK, res)
}
