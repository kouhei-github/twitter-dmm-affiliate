package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"kouhei-github/sample-gin/repository"
	"kouhei-github/sample-gin/service"
	"net/url"
	"os"
	"strconv"
	"strings"
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

func PostMediaIdHandoler(ctx *gin.Context) {
	oauth := service.NewTOAuth1()
	dir, err := os.Getwd()
	if err != nil {
		ctx.JSON(500, "Imageのパスの取得に失敗しました.")
		return
	}
	files, err := os.ReadDir(dir + "/public")
	var images []*service.ImageIdRequest
	for _, file := range files {
		image, err := oauth.GetMediaId(dir + "/public/" + file.Name())
		if err != nil {
			ctx.JSON(500, err.Error())
			return
		}
		images = append(images, image)
	}
	imageRequests := service.SeparateArray(2, images)
	err = oauth.PostToTwitterWithAttachment(imageRequests)
	if err != nil {
		ctx.JSON(500, err.Error())
	}
	ctx.JSON(200, "成功しました")
}

func SearchTagAndAutoFollowHandler(c *gin.Context) {
	search := c.Query("search")
	search = strings.ReplaceAll(search, ",", " #")
	search = "#" + search
	search = url.QueryEscape(search)
	oauth := service.NewTOAuth1()
	tweet, err := oauth.SearchHashTagOnTwitter(search, 10)
	if err != nil {
		c.JSON(500, "検索に失敗しました")
		return
	}
	// 今日の日付の取得
	today, err := service.GetToday()
	if err != nil {
		c.JSON(500, err.Error())
		return
	}

	var entities []repository.AutoFolowingEntity
	for _, tweetIncludeUser := range tweet.Data {
		err := oauth.FollowTwitterUser(tweetIncludeUser.AuthorId)
		if err != nil {
			c.JSON(500, "フォローできませんでした")
			return
		}
		entity, err := repository.NewAutoFolowingEntity(
			tweetIncludeUser.AuthorId,
			1,
			today,
		)
		if err != nil {
			c.JSON(500, "AutoFolowingEntityを製造できませんでした")
			return
		}
		entities = append(entities, *entity)
	}

	err = repository.BulkInsertAutoFollowing(entities)
	if err != nil {
		c.JSON(500, "AutoFolowingEntityを保存できませんでした")
		return
	}
	c.JSON(200, tweet)
}

func CreateReportOnIncreasingRateFromBefore(c *gin.Context) {
	entity, err := repository.FindByLastRecord()
	if err != nil {
		c.JSON(500, "最後のレコードを取得できませんでした")
		return
	}

	if err != nil {
		c.JSON(500, err)
		return
	}

	date, err := service.GetCarbonDate(entity.ExpireDate)
	if err != nil {
		c.JSON(500, err)
		return
	}
	today := service.CarbonToString(date)
	todayEntities, err := repository.FindByExpireDate(today) //本日

	subDay := date.SubDays(1)
	subDate := service.CarbonToString(subDay)
	subEntities, err := repository.FindByExpireDate(subDate) //前日
	if err != nil {
		c.JSON(500, err)
		return
	}

	increaseStruct := service.NewIncrease(len(todayEntities), len(subEntities))
	rate := increaseStruct.IncreaseRate()
	increase := increaseStruct.Today - increaseStruct.Yesterday
	fmt.Println(increaseStruct.Today)
	fmt.Println(increaseStruct.Yesterday)
	fmt.Println(increase)
	fmt.Println(rate)

	records := [][]string{
		{"昨日", "今日", "増加数", "増加率"},
		{strconv.Itoa(increaseStruct.Yesterday), strconv.Itoa(increaseStruct.Today), strconv.Itoa(increase), strconv.Itoa(rate)},
	}
	fileName, err := service.WriteCsv(records)
	if err != nil {
		c.JSON(500, err.Error())
		return
	}
	url := "http://localhost:8000/" + fileName

	c.JSON(200, url)
}
