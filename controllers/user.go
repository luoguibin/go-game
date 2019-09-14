package controllers

import (
	"errors"
	"fmt"
	"go-game/helper"
	"go-game/models"
	"strconv"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

// UserController operations for User
type UserController struct {
	BaseController
}

// 创建user
func (c *UserController) CreateUser() {
	data := c.GetResponseData()
	params := &getCreateUserParams{}
	if c.CheckPostParams(data, params) {
		user, err := models.CreateUser(params.ID, params.Pw, params.Name)
		if err == nil {
			if params.Type == "game" {
				models.CreateGameData(params.ID, params.Name, 100, 10000, 500, 0, 0)
				models.CreateGameSpear(params.ID, 1000, 100, 0, 0, 0, 0, 0)
				models.CreateGameShield(params.ID, 1000, 100, 0, 0, 0, 0, 0)
			}
			createUserToken(user, data)
		} else {
			data[models.STR_CODE] = models.CODE_ERR
			errStr := err.Error()
			if strings.Contains(errStr, "PRIMARY") {
				data[models.STR_MSG] = "已存在该用户"
			} else {
				data[models.STR_MSG] = "用户注册失败"
			}
		}
	}

	c.respToJSON(data)
}

// user登录
func (c *UserController) LoginUser() {
	data := c.GetResponseData()
	params := &getCreateUserParams{}
	if c.CheckPostParams(data, params) {
		user, err := models.QueryUser(params.ID)
		if err == nil {
			compare := strings.Compare(helper.HmacMd5(user.Password, models.MConfig.JwtSecretKey), params.Pw)
			if compare == 0 {
				createUserToken(user, data)
			} else {
				data[models.STR_CODE] = models.CODE_ERR
				data[models.STR_MSG] = "用户账号或密码错误"
			}
		} else {
			data[models.STR_CODE] = models.CODE_ERR
			data[models.STR_MSG] = "用户账号或密码错误"
		}
	}
	c.respToJSON(data)
}

// 创建用户token，基于json web token
func createUserToken(user *models.User, data ResponseData) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := make(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(24)).Unix()
	claims["iat"] = time.Now().Unix()
	claims["uId"] = strconv.FormatInt(user.ID, 10)
	claims["Level"] = strconv.Itoa(user.Level)

	token.Claims = claims

	tokenString, err := token.SignedString([]byte(models.MConfig.JwtSecretKey))
	if err != nil {
		data[models.STR_CODE] = models.CODE_ERR
		data[models.STR_MSG] = "用户id签名失败"
		return
	}

	user.Token = tokenString
	data[models.STR_DATA] = user
}

// 检测解析token
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
