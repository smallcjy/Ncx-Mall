package shop

import "github.com/gin-gonic/gin"

type ShopRouter struct{}

func (s *ShopRouter) InitShopRouter(Router *gin.RouterGroup) {
	shopRouter := Router.Group("shop")

	{
		shopRouter.POST("createShopOrder", shopApi.CreateShopOrder) // 创建ShopOrder
		shopRouter.DELETE("deleteShopOrder", shopApi.DeleteShopOrder) // 删除ShopOrder
		shopRouter.GET("getShopOrders", shopApi.GetUserOrders) // 根据ID获取ShopOrder
	}
}