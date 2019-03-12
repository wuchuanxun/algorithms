package pso

import (
	"math"
	"math/rand"
	"time"
)

const PSO_MAX_SIZE int = 100
const W_DEFAULT float64 = 0.7298
const (
	GlobalOnly  = iota
	GlobalLocal = iota
)

// PSO 算法的设置入口
type PSOSetting struct {
	// 维度
	NDim int
	// 目标函数，计算适应度
	Evaluater func(x []float64) float64
	// 群体大小
	Size int
	// 训练过程展示的周期
	ShowInterval int
	// 训练过程是否可见
	ShowDetails bool
	// 训练迭代次数
	Steps int
	// 认知系数
	C1 float64
	C2 float64
	// 初始惯性权重
	WMax float64
	// 最终惯性权重
	WMin float64
	// 轮次更新权重
	Alpha float64
	// 坐标下限
	XLow float64
	// 坐标上限
	XHigh float64
	// 是否对位置进行裁剪
	ClipPos bool
	// 更新方式选择
	UpdateMethod int
	// 速度上限
	VMax float64
	// 随机数生成器
	RandGenerater *rand.Rand
}

// 检查参数设定合法性
func (setting *PSOSetting) Check() {
	// 如果为0表示每一次都展示
	if setting.ShowDetails && setting.ShowInterval <= 0 {
		setting.ShowInterval = 1
	}
	switch {
	case setting.NDim <= 0:
		panic("The Dim should be positive number!")
	case setting.Evaluater == nil:
		panic("No evaluate func!")
	case setting.Size <= 0:
		panic("The Size should be positive number!")
	case setting.Steps <= 0:
		panic("The Steps should be positive number!")
	case setting.ClipPos && setting.XLow >= setting.XHigh:
		panic("Wrong range of axis limit")
	}
}

/**
 * PSO 算法的设置入口
 * @param {dim} 编码的维度
 * @param {max_size} 最大总群数目
 */
func CalculateSwarmSize(dim, max_size int) int {
	s := 10. + 2.*math.Sqrt(float64(dim))
	size := int(math.Floor(s + 0.5))
	if size > max_size {
		return max_size
	} else {
		return size
	}
}

// 得到一个默认的参数设置
func GetDefaultSetting() *PSOSetting {
	return &PSOSetting{
		Size:          20,
		ShowInterval:  10,
		ShowDetails:   true,
		Steps:         10000,
		C1:            1.5,
		C2:            1.5,
		Alpha:         0.36,
		WMax:          W_DEFAULT,
		WMin:          W_DEFAULT,
		XLow:          0,
		XHigh:         1,
		ClipPos:       true,
		UpdateMethod:  GlobalOnly,
		VMax:          1,
		RandGenerater: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}
