package common

import (
	"errors"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

var DB *gorm.DB

// ConnectDB 初始化数据库连接
func ConnectDB() error {
	dsn := DBConfig.User + ":" + DBConfig.Password + "@tcp(" + DBConfig.Host + ":" + DBConfig.Port + ")/" +
		DBConfig.DBName + "?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	DB = db
	return nil
}

// MigrateModels 自动迁移模型到数据库
func MigrateModels() error {
	return DB.AutoMigrate(&User{}, &Asset{}, &Port{})
}

// InitDB 初始化数据库，包括连接和迁移
func InitDB() {
	err := ConnectDB()
	if err != nil {
		log.Fatalf("无法连接到数据库: %v", err)
	}

	err = MigrateModels()
	if err != nil {
		log.Fatalf("无法迁移数据库模型: %v", err)
	}
}

// SelectUser 查询User，如果找不到则根据staffId和name创建

func SelectUser(staffId string, name string) (User, error) {

	var user User
	result := DB.Find(&user, "staff_id = ?", staffId)
	if result.Error != nil {
		return User{}, fmt.Errorf("failed to retrieve user: %w", result.Error)
	}
	// 检查是否找到记录
	if result.RowsAffected == 0 {
		newUser := User{StaffId: staffId, Name: name}
		result = DB.Create(&newUser)

		return newUser, nil
	}

	return user, nil
}

// SelectAllUser 查询全部User
func SelectAllUser() []User {
	var user []User
	result := DB.Find(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return user
		}
		panic("failed to retrieve user: " + result.Error.Error())
	}

	return user
}

// SelectAllAsset 查询全部资产
func SelectAllAsset() []Asset {
	var assets []Asset
	result := DB.Find(&assets)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return assets
		}
		panic("failed to retrieve assets: " + result.Error.Error())
	}

	return assets
}

// SelectNotWhitePort 查询非白名单端口
func SelectNotWhitePort(user User) (bool, string, string) {
	var accMsg, errMsg string
	var assets []Asset
	var sentMsg = false
	result := DB.Find(&assets, "user_id = ?", user.ID)
	if result.Error != nil {
		errMsg = fmt.Sprintf("\n> |  |  |  |  |  |  |  | 数据库错误：%s！| ", result.Error)
		return true, accMsg, errMsg
	} else {
		for _, asset := range assets {
			var ports []Port
			result := DB.Find(&ports, "asset_id = ? AND white = ?", asset.ID, false)
			if result.Error != nil {
				errMsg = fmt.Sprintf("\n> |  |  |  |  |  |  |  | 数据库错误：%s！| ", result.Error)
				return true, accMsg, errMsg
			}
			for _, port := range ports {
				accMsg += fmt.Sprintf("\n> | %s | %s | %d | %s | %s | %s | %v |  | ", asset.Department, asset.IP, port.Port, port.Protocol, port.State, port.Service, port.White)
			}
			if len(ports) != 0 {
				sentMsg = true
			}
		}
		return sentMsg, accMsg, errMsg
	}
}
