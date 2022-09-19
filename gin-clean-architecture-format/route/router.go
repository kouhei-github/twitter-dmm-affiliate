package route

import (
	"github.com/gin-gonic/gin"
	"kouhei-github/sample-gin/controller"
)

func GetRouter() *gin.Engine {
	r := gin.Default()

	// デレバリー項目の追加
	r.POST("/api/v1/delivery", controller.InsertDeliveryHandler)

	// 商品情報の追加
	r.POST("/api/v1/merchandise-batch", controller.BulkInsertMerchandiseHandler)
	r.PUT("/api/v1/merchandise-batch", controller.BulkUpdateMerchandiseHandler)

	// TwitterのフォローしたユーザーIDと期間の保存
	r.POST("/api/v1/twitter", controller.InsertTwitterAutoFollowHandler)
	r.GET("/api/v1/twitter", controller.FinfUseridTwitterAutoFollowHandler)
	return r
}
