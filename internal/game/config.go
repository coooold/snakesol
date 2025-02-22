package game

import (
	"math/rand"
	"time"
)

// 游戏配置
var (
	colors = []string{"#FF0000", "#00FF00", "#0000FF", "#FFFF00", "#FF00FF", "#00FFFF"}
	names  = []string{
		"蛇王", "小青", "大白", "闪电", "飞毛腿", "追风",
		"青龙", "白虎", "玄武", "朱雀", "麒麟", "凤凰",
		"游侠", "猎手", "勇士", "斗士", "战神", "幻影",
		"飞龙", "神蛇",
	}
	// 蛇的称号前缀
	name_prefixs = []string{
		"无敌", "神速", "致命", "狂暴", "霸王",
		"至尊", "傲世", "绝世", "天下", "不败",
		"无双", "王者", "霸主", "至强", "无上",
		"绝命", "无情", "狂野", "战神", "传说",
	}
)

// getRandName 随机生成一个蛇的称号
func getRandName() string {
	rand.Seed(time.Now().UnixNano())
	name := name_prefixs[rand.Intn(len(name_prefixs))] + names[rand.Intn(len(names))]
	return name
}

// DefaultConfig 返回默认游戏配置
func DefaultConfig() *GameConfig {
	return &GameConfig{
		// 游戏区域的列数
		Cols: 100,
		// 游戏区域的行数
		Rows: 100,
		// 蛇的初始长度
		InitialSnakeLength: 10,
		// AI蛇生成的时间间隔(秒)
		AISpawnInterval: 10,
		// 游戏更新的时间间隔(毫秒)
		UpdateInterval: 150,
		// 初始AI蛇的数量
		InitialAICount: 50,
		// AI蛇的最大数量
		MaxAICount: 100,
		// 苹果生成的时间间隔(秒)
		AppleSpawnInterval: 1,
		// 苹果的存活时间(秒)
		AppleLifetime: 10,
	}
}
