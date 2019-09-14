package controllers

import (
	"strings"

	"github.com/astaxie/beego/validation"
)

//
// 有token验证的，均采用query传递该参数；
// 跳转路由前，在路由InsertFilter中，手动设置解析token中信息，
// 并写入对应的输入结构体中。
//

// user create和login共用输入结构体
type getCreateUserParams struct {
	ID   int64  `form:"uId" json:"uId" valid:"Required"`
	Pw   string `form:"pw" json:"pw" valid:"Required"`
	Name string `form:"name"`
	Type string `form:"type"`
}

func (params *getCreateUserParams) Valid(v *validation.Validation) {
	if params.ID <= 0 {
		v.SetError("uId", "不能为空")
	} else if len(strings.TrimSpace(params.Pw)) == 0 {
		v.SetError("pw", "不能为空")
	}
}
