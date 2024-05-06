package dingtalk

import (
	"Server/common"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"net"
	"strings"
)

// handleAddAsset 指令内容格式分割
func splint(line string) (string, string, error) {
	parts := strings.SplitN(line, "-", 2)
	if len(parts) == 2 {
		return parts[0], parts[1], nil
	}
	return line, "", errors.New("格式化字符串错误！")
}

// handleWhitePort 处理添加资产指令
func handleAddAsset(user common.User, content string) (string, string) {
	var accMsg, errMsg string
	department, ip, err := splint(content)
	if err != nil {
		errMsg = fmt.Sprintf("\n> | %s |  | 字符串格式化错误！ | ", content)
		return accMsg, errMsg
	}

	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		errMsg = fmt.Sprintf("\n> | %s | %s | 无效的 IP 地址！ | ", department, ip)
		return accMsg, errMsg
	}

	var asset common.Asset
	result := common.DB.Where("ip = ?", ip).First(&asset) // 使用 First 而不是 Find

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// 没有找到记录，可以添加新资产
			asset = common.Asset{
				Department: department,
				IP:         ip,
				UserID:     user.ID,
			}
			result = common.DB.Create(&asset)
			if result.Error != nil {
				errMsg = fmt.Sprintf("\n> | %s | %s | %s！| ", department, ip, result.Error)
				return accMsg, errMsg
			} else {
				accMsg = fmt.Sprintf("\n> | %s | %s | 添加成功 | ", department, ip)
				return accMsg, errMsg
			}
		} else {
			// 其他数据库错误
			errMsg = fmt.Sprintf("\n> | %s | %s | 数据库错误：%s！| ", department, ip, result.Error)
			return accMsg, errMsg
		}
	} else {
		// 已找到记录
		errMsg = fmt.Sprintf("\n> | %s | %s | IP重复入库！| ", department, ip)
		return accMsg, errMsg
	}
}

// handleDelAsset 处理删除资产指令
func handleDelAsset(user common.User, ip string) (string, string) {
	var accMsg, errMsg string
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		errMsg = fmt.Sprintf("\n> |  | %s | 无效的 IP 地址！ | ", ip)
		return accMsg, errMsg
	}

	var asset common.Asset
	result := common.DB.Where("user_id = ? AND ip = ?", user.ID, ip).First(&asset)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			errMsg = fmt.Sprintf("\n> |  | %s | 资产不存在或权限不足！| ", ip)
			return accMsg, errMsg
		} else {
			errMsg = fmt.Sprintf("\n> |  | %s | 数据库错误：%s！| ", ip, result.Error)
			return accMsg, errMsg
		}
	}

	// 删除所有与asset相关联的Port记录
	result = common.DB.Where("asset_id = ?", asset.ID).Delete(&common.Port{})
	if result.Error != nil {
		errMsg = fmt.Sprintf("\n> |  | %s | 删除Port时发生数据库错误：%s！| ", ip, result.Error)
		return accMsg, errMsg
	}

	// 删除Asset记录
	result = common.DB.Where("user_id = ? AND ip = ?", user.ID, ip).Delete(&asset)
	if result.Error != nil {
		errMsg = fmt.Sprintf("\n> |  | %s | 删除Asset时发生数据库错误：%s！| ", ip, result.Error)
		return accMsg, errMsg
	}

	department := asset.Department
	accMsg = fmt.Sprintf("\n> | %s | %s | 删除资产及相应端口成功 | ", department, ip)
	return accMsg, errMsg
}

// handleUpAsset 处理更新资产指令
func handleUpAsset(user common.User, content string) (string, string) {
	var accMsg, errMsg string
	department, ip, err := splint(content)
	if err != nil {
		errMsg = fmt.Sprintf("\n> | %s |  | 字符串格式化错误！ | ", content)
		return accMsg, errMsg
	}

	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		errMsg = fmt.Sprintf("\n> | %s | %s | 无效的 IP 地址！ | ", department, ip)
		return accMsg, errMsg
	}

	var asset common.Asset
	result := common.DB.Where("user_id = ? AND ip = ?", user.ID, ip).First(&asset)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			errMsg = fmt.Sprintf("\n> |  | %s | 资产不存在或权限不足！| ", ip)
			return accMsg, errMsg
		} else {
			errMsg = fmt.Sprintf("\n> |  | %s | 数据库错误：%s！| ", ip, result.Error)
			return accMsg, errMsg
		}
	} else {
		result = common.DB.Model(&common.Asset{}).Where("user_id = ? AND ip = ?", user.ID, ip).Update("department", department)
		if result.Error != nil {
			errMsg = fmt.Sprintf("\n> |  | %s | 数据库错误：%s！| ", ip, result.Error)
			return accMsg, errMsg
		} else {
			accMsg = fmt.Sprintf("\n> | %s | %s | 更新成功 | ", department, ip)
			return accMsg, errMsg
		}
	}
}

// handleSelectAsset 处理查询资产指令
func handleSelectAsset(user common.User) (string, string) {
	var accMsg, errMsg string
	var assets []common.Asset
	result := common.DB.Find(&assets, "user_id = ?", user.ID)
	if result.Error != nil {
		errMsg = fmt.Sprintf("\n> |  |  | 数据库错误：%s！| ", result.Error)
		return accMsg, errMsg
	} else {
		for _, asset := range assets {
			accMsg += fmt.Sprintf("\n> | %s | %s |  | ", asset.Department, asset.IP)
		}
		return accMsg, errMsg
	}
}

// handleDelPort 处理删除端口指令
func handleDelPort(user common.User, content string) (string, string) {
	var accMsg, errMsg string
	ip, strPort, err := splint(content)
	if err != nil {
		errMsg = fmt.Sprintf("\n> | %s |  |  |  |  |  |  | 字符串格式化错误！| ", content)
		return accMsg, errMsg
	}

	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		errMsg = fmt.Sprintf("\n> |  | %s | %s |  |  |  |  | 无效的 IP 地址！| ", ip, strPort)
		return accMsg, errMsg
	}

	var asset common.Asset
	result := common.DB.Where("user_id = ? AND ip = ?", user.ID, ip).First(&asset)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			errMsg = fmt.Sprintf("\n> |  | %s | %s |  |  |  |  | 资产不存在或权限不足！| ", ip, strPort)
			return accMsg, errMsg
		} else {
			errMsg = fmt.Sprintf("\n> |  | %s | %s |  |  |  |  | 数据库错误：%s！| ", ip, strPort, result.Error)
			return accMsg, errMsg
		}
	} else {
		var port common.Port
		result := common.DB.Where("port = ? AND asset_id = ?", strPort, asset.ID).First(&port)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				errMsg = fmt.Sprintf("\n> |  | %s | %s |  |  |  |  | 端口不存在！| ", ip, strPort)
				return accMsg, errMsg
			} else {
				errMsg = fmt.Sprintf("\n> |  | %s | %s |  |  |  |  | 数据库错误：%s！| ", ip, strPort, result.Error)
				return accMsg, errMsg
			}
		} else {
			result := common.DB.Where("port = ? AND asset_id = ?", strPort, asset.ID).Delete(&port)
			if result.Error != nil {
				errMsg = fmt.Sprintf("\n> |  | %s | %s |  |  |  |  | 数据库错误：%s！| ", ip, strPort, result.Error)
				return accMsg, errMsg
			}
			errMsg = fmt.Sprintf("\n> | %s | %s | %s | %s | %s | %s | %v | 删除成功 | ", asset.Department, ip, strPort, port.Protocol, port.State, port.Service, port.White)
			return accMsg, errMsg
		}
	}
}

// handleWhitePort 处理添加白名单指令
func handleWhitePort(user common.User, content string) (string, string) {
	var accMsg, errMsg string
	ip, strPort, err := splint(content)
	if err != nil {
		errMsg = fmt.Sprintf("\n> | %s |  |  |  |  |  |  | 字符串格式化错误！| ", content)
		return accMsg, errMsg
	}

	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		errMsg = fmt.Sprintf("\n> |  | %s | %s |  |  |  |  | 无效的 IP 地址！| ", ip, strPort)
		return accMsg, errMsg
	}

	var asset common.Asset
	result := common.DB.Where("user_id = ? AND ip = ?", user.ID, ip).First(&asset)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			errMsg = fmt.Sprintf("\n> |  | %s | %s |  |  |  |  | 资产不存在或权限不足！| ", ip, strPort)
			return accMsg, errMsg
		} else {
			errMsg = fmt.Sprintf("\n> |  | %s | %s |  |  |  |  | 数据库错误：%s！| ", ip, strPort, result.Error)
			return accMsg, errMsg
		}
	} else {
		var port common.Port
		result := common.DB.Where("port = ? AND asset_id = ?", strPort, asset.ID).First(&port)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				errMsg = fmt.Sprintf("\n> |  | %s | %s |  |  |  |  | 端口不存在！| ", ip, strPort)
				return accMsg, errMsg
			} else {
				errMsg = fmt.Sprintf("\n> |  | %s | %s |  |  |  |  | 数据库错误：%s！| ", ip, strPort, result.Error)
				return accMsg, errMsg
			}
		} else {
			result := common.DB.Model(&common.Port{}).Where("port = ? AND asset_id = ?", strPort, asset.ID).Update("white", true)
			if result.Error != nil {
				errMsg = fmt.Sprintf("\n> |  | %s | %s |  |  |  |  | 数据库错误：%s！| ", ip, strPort, result.Error)
				return accMsg, errMsg
			}
			errMsg = fmt.Sprintf("\n> | %s | %s | %s | %s | %s | %s | %v | 加白成功 | ", asset.Department, ip, strPort, port.Protocol, port.State, port.Service, port.White)
			return accMsg, errMsg
		}
	}
}

// handleSelectPort 处理查询端口指令
func handleSelectPort(user common.User) (string, string) {
	var accMsg, errMsg string
	var assets []common.Asset
	result := common.DB.Find(&assets, "user_id = ?", user.ID)
	if result.Error != nil {
		errMsg = fmt.Sprintf("\n> |  |  | 数据库错误：%s！| ", result.Error)
		return accMsg, errMsg
	} else {
		for _, asset := range assets {
			var ports []common.Port
			result := common.DB.Find(&ports, "asset_id = ?", asset.ID)
			if result.Error != nil {
				errMsg = fmt.Sprintf("\n> |  |  | 数据库错误：%s！| ", result.Error)
				return accMsg, errMsg
			}
			for _, port := range ports {
				accMsg += fmt.Sprintf("\n> | %s | %s | %d | %s | %s | %s | %v |  | ", asset.Department, asset.IP, port.Port, port.Protocol, port.State, port.Service, port.White)
			}
		}
		return accMsg, errMsg
	}
}
