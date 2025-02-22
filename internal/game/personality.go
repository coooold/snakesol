package game

import "math/rand"

// ViewRange 定义AI蛇的视野范围
const (
	ViewWidth  = 40 // 视野宽度
	ViewHeight = 40 // 视野高度
)

// ViewInfo 存储AI蛇的视野信息
type ViewInfo struct {
	Center    Position   // 视野中心点（蛇头位置）
	Food      []Position // 视野范围内的食物位置
	Snakes    []*Snake   // 视野范围内的其他蛇
	Obstacles []Position // 视野范围内的障碍物
}

// PersonalityType 定义AI蛇的性格类型
type PersonalityType int

const (
	Aggressive  PersonalityType = iota // 进攻型
	Evasive                            // 躲避型
	Balanced                           // 平衡型
	Guerrilla                          // 游击型
	Trapper                            // 陷阱型
	Survival                           // 生存型
	Cooperative                        // 协作型
)

// PersonalityColor 定义每种性格对应的颜色
var PersonalityColor = map[PersonalityType]string{
	Aggressive:  "#FF4136", // 红色，表示攻击性
	Evasive:     "#39CCCC", // 青色，表示谨慎
	Balanced:    "#2ECC40", // 绿色，表示平衡
	Guerrilla:   "#FF851B", // 橙色，表示灵活
	Trapper:     "#B10DC9", // 紫色，表示神秘
	Survival:    "#FFDC00", // 黄色，表示保守
	Cooperative: "#7FDBFF", // 蓝色，表示团队协作
}

// PersonalityWeights 定义每种性格的决策权重
type PersonalityWeights struct {
	AttackWeight    float64
	FoodWeight      float64
	SurvivalWeight  float64
	SpaceWeight     float64
	MobilityWeight  float64
	TrapWeight      float64
	CooperateWeight float64 // 协作权重
}

// GetPersonalityWeights 根据性格类型返回对应的决策权重
func GetPersonalityWeights(p PersonalityType) PersonalityWeights {
	switch p {
	case Aggressive:
		return PersonalityWeights{
			AttackWeight:   0.4,
			FoodWeight:     0.8,
			SurvivalWeight: 0.2,
			SpaceWeight:    0.1,
		}
	case Evasive:
		return PersonalityWeights{
			SurvivalWeight: 0.4,
			SpaceWeight:    0.3,
			FoodWeight:     0.2,
			AttackWeight:   0.1,
		}
	case Balanced:
		return PersonalityWeights{
			SurvivalWeight: 0.3,
			FoodWeight:     0.3,
			SpaceWeight:    0.2,
			AttackWeight:   0.2,
		}
	case Guerrilla:
		return PersonalityWeights{
			MobilityWeight: 0.4,
			FoodWeight:     0.3,
			SurvivalWeight: 0.2,
			AttackWeight:   0.1,
		}
	case Trapper:
		return PersonalityWeights{
			TrapWeight:     0.4,
			SpaceWeight:    0.3,
			SurvivalWeight: 0.2,
			FoodWeight:     0.1,
		}
	case Survival:
		return PersonalityWeights{
			SurvivalWeight: 0.5,
			SpaceWeight:    0.3,
			FoodWeight:     0.15,
			AttackWeight:   0.05,
		}
	case Cooperative:
		return PersonalityWeights{
			CooperateWeight: 0.5, // 增加协作权重，强化团队配合
			AttackWeight:    0.2, // 降低单独攻击倾向
			SurvivalWeight:  0.2,
			FoodWeight:      0.1,
			MobilityWeight:  0.2, // 添加机动性权重，提升包围能力
			SpaceWeight:     0.1, // 添加空间权重，优化位置选择
		}
	default:
		return PersonalityWeights{} // 默认返回空权重
	}
}

// RandomPersonality 随机生成一个性格类型
func RandomPersonality() PersonalityType {
	return PersonalityType(rand.Intn(7))
}
