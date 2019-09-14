package models

import (
	"strings"
	"time"
)

// User 用户登录信息
type User struct {
	ID         int64     `gorm:"primary_key" json:"id,omitempty"`
	Password   string    `gorm:"column:password;type:varchar(20)" json:"-"`
	Name       string    `gorm:"column:name;type:varchar(100)" json:"name,omitempty"`
	Token      string    `gorm:"-" json:"token,omitempty"`
	IconURL    string    `gorm:"column:icon_url" json:"iconUrl"`
	TimeCreate time.Time `gorm:"column:time_create" json:"timeCreate"`
	Level      int       `gorm:"column:level" json:"-"`
}

//  `json:"-"` 把struct编码成json字符串时，会忽略这个字段
//	`json:"id,omitempty"` //如果这个字段是空值，则不编码到JSON里面，否则用id为名字编码
//	`json:",omitempty"`   //如果这个字段是空值，则不编码到JSON里面，否则用属性名为名字编码

// TableName ...
func (u User) TableName() string {
	return "user"
}

func initSystemUser() {
	tx := dbOrmDefault.Model(&User{}).Begin()
	tx.Create(User{
		ID:         15625045984,
		Password:   "123456",
		Name:       "乂末",
		Level:      9,
		TimeCreate: time.Now(),
	})
	tx.Create(User{
		ID:         15688888888,
		Password:   "123456",
		Name:       "Sghen",
		Level:      9,
		TimeCreate: time.Now(),
	})
	tx.Create(User{
		ID:         15622222222,
		Password:   "123456",
		Name:       "Morge",
		Level:      9,
		TimeCreate: time.Now(),
	})
	tx.Create(User{
		ID:         15666666666,
		Password:   "123456",
		Name:       "SghenMorge",
		Level:      9,
		TimeCreate: time.Now(),
	})
	tx.Commit()
}

// CreateUser ...
func CreateUser(ID int64, password string, name string) (*User, error) {
	user := &User{
		ID:         ID,
		Password:   password,
		Name:       name,
		TimeCreate: time.Now(),
		Level:      1,
	}

	err := dbOrmDefault.Model(&User{}).Create(user).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

// QueryUser ...
func QueryUser(ID int64) (*User, error) {
	user := &User{
		ID: ID,
	}
	err := dbOrmDefault.Model(&User{}).Find(user).Error
	if err == nil {
		return user, nil
	} else {
		return nil, err
	}
}

// QueryUsers ...
func QueryUsers(IDs []int64) ([]*User, error) {
	list := make([]*User, 0)
	err := dbOrmDefault.Model(&User{}).Select("id, name, icon_url").Where("id in (?)", IDs).Find(&list).Error
	if err == nil {
		return list, nil
	}
	return nil, err
}

// UpdateUser ...
func UpdateUser(ID int64, password string, name string, iconURL string) (*User, error) {
	user := &User{
		ID: ID,
	}

	err := dbOrmDefault.Model(&User{}).Find(user).Error
	if err == nil {
		if len(strings.TrimSpace(password)) > 0 {
			user.Password = password
		}
		if len(strings.TrimSpace(name)) > 0 {
			user.Name = name
		}
		if len(strings.TrimSpace(iconURL)) > 0 {
			user.IconURL = iconURL
		}

		err = dbOrmDefault.Model(&User{}).Save(user).Error
		if err != nil {
			return nil, err
		} else {
			return user, nil
		}
	} else {
		return nil, err
	}
}

// DeleteUser ...
func DeleteUser(ID int64) error {
	user := &User{
		ID: ID,
	}

	err := dbOrmDefault.Model(&User{}).Delete(user).Error
	return err
}
