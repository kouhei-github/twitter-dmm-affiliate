package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"kouhei-github/sample-gin/repository"
	"kouhei-github/sample-gin/service"
)

func InsertTwitterAutoFollowHandler(ctx *gin.Context) {
	var requestBody repository.AutoFolowingEntity
	err := ctx.ShouldBindJSON(&requestBody)
	if err != nil {
		myErr := service.MyError{Message: "リクエストBodyの取得に失敗しました。"}
		ctx.JSON(500, myErr)
		return
	}

	// 今日の日付の取得
	today, err := service.GetToday()
	if err != nil {
		ctx.JSON(500, err.Error())
		return
	}
	fmt.Println(today)
	// AutoFollowingEntityの生成
	entity, err := repository.NewAutoFolowingEntity(
		requestBody.TwitterUserId,
		requestBody.Status,
		today,
	)
	fmt.Println(entity)
	if err != nil {
		ctx.JSON(500, err.Error())
		return
	}
	err = entity.Create()
	if err != nil {
		ctx.JSON(500, err)
		return
	}
	ctx.JSON(201, "Insert Completed")
}

func FinfUseridTwitterAutoFollowHandler(ctx *gin.Context) {
	userId := ctx.Query("user_id")
	found, err := repository.FindByTwitterUserId(userId)
	if err != nil {
		ctx.JSON(500, err.Error())
	}
	ctx.JSON(200, found)
}
