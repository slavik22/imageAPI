package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/slavik22/imageAPI/db/sqlc"
	"github.com/slavik22/imageAPI/token"
	"github.com/slavik22/imageAPI/util"
	"net/http"
	"net/url"
)

func (server *Server) getImages(ctx *gin.Context) {
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	images, err := server.store.GetImages(ctx, authPayload.UserId)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, images)
}

type imageRequest struct {
	ImageUrl string `json:"image_url" binding:"required"`
}

func (server *Server) createImage(ctx *gin.Context) {
	var requestData imageRequest

	if err := ctx.ShouldBindJSON(&requestData); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if _, err := url.ParseRequestURI(requestData.ImageUrl); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	path, err := util.DownloadImage(requestData.ImageUrl, "../images")

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	arg := db.CreateImageParams{
		UserID:    authPayload.UserId,
		ImageUrl:  requestData.ImageUrl,
		ImagePath: path,
	}

	image, err := server.store.CreateImage(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, image)
}
