package routers

import (
	"tunaiku/controllers"
	"github.com/astaxie/beego"
)

func init() {
    beego.Router("/stock", &controllers.StockController{}, "get:Index")
    beego.Router("/stock/new", &controllers.StockController{}, "get:New")
    beego.Router("/stock/new", &controllers.StockController{}, "post:Create")
    beego.Router("/stock/calculate", &controllers.StockController{}, "get:Calculate")
    beego.Router("/stock/export", &controllers.StockController{}, "get:Export")
}
