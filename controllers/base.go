package controllers

import (
	"errors"
	"fmt"
	"go-game/models"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/plugins/cors"

	jwt "github.com/dgrijalva/jwt-go"
)

// BaseController ...
type BaseController struct {
	beego.Controller
}

func init() {
	fmt.Println("basecontroller init")
	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
		AllowCredentials: true,
	}))
}

// respToJSON 缁熶竴鎺ュ彛杩斿洖
func (c *BaseController) respToJSON(data ResponseData) {
	respMsg, ok := data[models.STR_MSG]
	if !ok || (ok && len(respMsg.(string)) <= 0) {
		data[models.STR_MSG] = models.MConfig.CodeMsgMap[data[models.STR_CODE].(int)]
	}
	// c.Ctx.Output.SetStatus(201)
	c.Data["json"] = data
	c.ServeJSON()
}

// BaseGetTest 鍩虹娴嬭瘯璋冪敤
func (c *BaseController) BaseGetTest() {
	data := c.GetResponseData()

	ip := c.Ctx.Request.Header.Get("X-Forwarded-For")
	if strings.Contains(ip, "127.0.0.1") || ip == "" {
		ip = c.Ctx.Request.Header.Get("X-real-ip")
	}

	if ip == "" {
		ip = "127.0.0.1"
	}
	data["ip"] = ip

	c.respToJSON(data)
}

// GatewayAccessUser ...
func GatewayAccessUser(ctx *context.Context) {
	datas := ResponseData{}
	token := ctx.Input.Query("token")

	if len(token) <= 0 {
		datas[models.STR_CODE] = models.CODE_ERR
		datas[models.STR_MSG] = "token不能为空"
		ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")
		ctx.Output.JSON(datas, false, true)
		return
	}

	claims, err := CheckUserToken(token)
	if err != nil {
		datas[models.STR_CODE] = models.CODE_ERR_TOKEN
		errStr := err.Error()

		if strings.Contains(errStr, "expired") {
			datas[models.STR_MSG] = "token失效，请重新登录"
		} else {
			datas[models.STR_MSG] = "token参数错误"
		}

		ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")
		ctx.Output.JSON(datas, false, false)
		return
	}

	userId, _ := strconv.ParseInt(claims["userId"].(string), 10, 64)
	ctx.Input.SetData("userId", userId)
	ctx.Input.SetData("userName", claims["userName"])
	ctx.Input.SetData("level", claims["level"])
	// ctx.Input.Context.Request.Form.Add("userId", claims["userId"].(string))
	return
}

// CheckUserToken 检测解析token
func CheckUserToken(tokenString string) (map[string]interface{}, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(models.MConfig.JwtSecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, errors.New("token get mapcliams err")
	}
}

// ResponseData ...
type ResponseData map[string]interface{}

// GetResponseData ...
func (c *BaseController) GetResponseData() ResponseData {
	return ResponseData{models.STR_CODE: models.CODE_OK}
}
