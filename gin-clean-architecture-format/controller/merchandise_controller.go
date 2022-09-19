package controller

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"kouhei-github/sample-gin/repository"
)

func BulkInsertMerchandiseHandler(c *gin.Context) {
	buf := make([]byte, 2048)
	// ここでRequest.Bodyを読み切る
	n, _ := c.Request.Body.Read(buf)

	// リクエストBodyの内容を保存する構造体
	var requestBody []repository.MerchandiseEntity
	err := json.Unmarshal(buf[0:n], &requestBody)
	if err != nil {
		c.JSON(500, err)
		return
	}
	err = repository.CreateMerchandiseList(requestBody)
	if err != nil {
		c.JSON(500, err)
		return
	}
	c.JSON(200, "Batch Insert Completed")
}

func BulkUpdateMerchandiseHandler(c *gin.Context) {
	buf := make([]byte, 2048)
	// ここでRequest.Bodyを読み切る
	n, _ := c.Request.Body.Read(buf)

	// リクエストBodyの内容を保存する構造体
	var requestBody []repository.MerchandiseEntity
	err := json.Unmarshal(buf[0:n], &requestBody)
	fmt.Println(requestBody[0].IsUpload)
	if err != nil {
		c.JSON(500, err)
		return
	}
	err = repository.UpdateMerchandiseList(requestBody)
	if err != nil {
		c.JSON(500, err)
		return
	}
	c.JSON(200, "Batch update Completed")
}
