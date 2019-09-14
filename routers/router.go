// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"go-game/controllers"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

func init() {
	beego.Router("/", &controllers.BaseController{}, "GET:BaseGetTest")
	beego.Router("/ws", &controllers.GameController{})

	beego.InsertFilter("/ws", beego.BeforeRouter, func(ctx *context.Context) {
		controllers.GatewayAccessUser(ctx, false)
	})
	beego.InsertFilter("/v1/peotry/create", beego.BeforeRouter, func(ctx *context.Context) {
		controllers.GatewayAccessUser(ctx, false)
	})

	//详见　https://beego.me/docs/mvc/controller/router.md
	nsv := beego.NewNamespace("/v1",
		beego.NSNamespace("/user",
			beego.NSRouter("/create", &controllers.UserController{}, "post:CreateUser"),
			beego.NSRouter("/login", &controllers.UserController{}, "post:LoginUser"),
		),
	)

	beego.AddNamespace(nsv)
}
