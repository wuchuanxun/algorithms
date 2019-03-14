package aco

import (
	"math"
	"math/rand"
)

// 全局变量，信息素矩阵
var PheromoneMatrix = make(map[int]float64)

type ACOSetting struct {
	// 源维度
	SourceDim int
	// 目标 维度
	TargetDim int
	// 信息启发式因子
	Alpha float64
	// 期望启发式因子
	Beta float64
	// 蚁群数量
	Size int
	// 迭代次数
	Steps int
	// 信息挥发因子
	Rho float64
	// 信息素上限
	TaoMax float64
	// 信息素下限
	TaoMin float64
	// 信息素更新方法
	PheromoneUpdateMethod func(i, j, step int, dtao float64)
	// 目标函数，计算函数值
	Evaluater func(x MatchVector) float64
	// 获取信息素,从状态i选择状态j
	GetPheromone func(i, j, step int) float64
	// 获取默认期望
	GetExpect func(i, j int) float64
	// 根据函数值给出信息素更新量
	GetPheromoneDelta func(v float64) float64
	// 蚂蚁的游走方法
	AntWalkMethod func(ant *Ant, step int, pheromone func(i, j, step int) float64, expect func(i, j int) float64)
}

func DefaultAntWalkEven(setting *ACOSetting) func(ant *Ant, step int, pheromone func(i, j, step int) float64, expect func(i, j int) float64) {
	SDim := setting.SourceDim
	TDim := setting.TargetDim
	return func(ant *Ant, step int, pheromone func(i, j, step int) float64, expect func(i, j int) float64) {
		for i := 0; i < SDim; i++ {
			pdf := make([]float64, TDim)
			for j := 0; j < TDim; j++ {
				pdf[j] = math.Pow(pheromone(i, j, step), setting.Alpha) * math.Pow(expect(i, j), setting.Beta)
			}
			t := PdfSample(pdf)
			ant.Matches[i] = Connection{
				source: i,
				target: t,
			}
		}
	}
}

func DefaultAntWalkPerm(setting *ACOSetting) func(ant *Ant, step int, pheromone func(i, j, step int) float64, expect func(i, j int) float64) {
	if setting.SourceDim != setting.TargetDim {
		panic("The Dim of source and target not match!")
	}
	Dim := setting.SourceDim
	return func(ant *Ant, step int, pheromone func(i, j, step int) float64, expect func(i, j int) float64) {
		setting := ant.Setting
		// 禁忌表
		visited := make(map[int]bool)
		s := rand.Intn(Dim)
		initials := s
		var t int
		visited[s] = true
		for len(visited) < Dim {
			pdf := make([]float64, Dim)
			for j := 0; j < Dim; j++ {
				pdf[j] = math.Pow(pheromone(s, j, step), setting.Alpha) * math.Pow(expect(s, j), setting.Beta)
			}
			for k, _ := range visited {
				pdf[k] = 0
			}
			t = PdfSample(pdf)
			ant.Matches[len(visited)-1] = Connection{
				source: s,
				target: t,
			}
			s = t
		}
		ant.Matches[Dim-1] = Connection{
			source: t,
			target: initials,
		}
	}
}

func GetDefaultSetting(sdim, tdim int) *ACOSetting {
	setting := &ACOSetting{
		SourceDim: sdim,
		TargetDim: tdim,
		Alpha:     2,
		Beta:      1,
		Size:      30,
		Steps:     1000,
		Rho:       0.06,
		TaoMax:    2,
		TaoMin:    0.1,
	}
	setting.PheromoneUpdateMethod = DefaultPheromoneUpdater(setting)
	setting.GetPheromone = DefaultPheromone(setting)
	setting.GetExpect = DefaultExpectValue(setting)
	setting.GetPheromoneDelta = DefaultPheromoneDeltaCalculate(setting)
	setting.AntWalkMethod = DefaultAntWalkEven(setting)
	return setting
}

func (setting *ACOSetting) Check() {
	switch {
	case setting.SourceDim <= 0:
		panic("Invalid Source Dimension!")
	case setting.TargetDim <= 0:
		panic("Invalid Target DImension!")
	case setting.Alpha < 0:
		panic("Invalid Alpha!")
	case setting.Beta < 0:
		panic("Invalid Beta!")
	case setting.Rho <= 0 || setting.Rho >= 1:
		panic("Invalid Rho!")
	case setting.PheromoneUpdateMethod == nil:
		panic("No PheromoneUpdateMethod!")
	case setting.Evaluater == nil:
		panic("No Evaluater!")
	case setting.GetPheromone == nil:
		panic("No GetPheromone func!")
	case setting.GetExpect == nil:
		panic("No GetExpect func!")
	case setting.GetPheromoneDelta == nil:
		panic("No GetPheromoneDelta func!")
	}
}

// 裁剪函数
func Clip(v, l, u float64) float64 {
	if v < l {
		return l
	}
	if v > u {
		return u
	}
	return v
}

func DefaultPheromoneUpdater(setting *ACOSetting) func(i, j, step int, dtao float64) {
	YDim := setting.TargetDim
	TaoMax := setting.TaoMax
	TaoMin := setting.TaoMin
	Rho := setting.Rho
	return func(i, j, step int, dtao float64) {
		p, found := PheromoneMatrix[i*YDim+j]
		if found {
			p = (1-Rho)*p + Rho*dtao
			PheromoneMatrix[i*YDim+j] = Clip(p, TaoMin, TaoMax)
			return
		}
		if dtao != 0 {
			p = TaoMax * math.Pow(1-Rho, float64(step))
			p = (1-Rho)*p + Rho*dtao
			PheromoneMatrix[i*YDim+j] = Clip(p, TaoMin, TaoMax)
		}
	}
}

func DefaultPheromone(setting *ACOSetting) func(i, j, step int) float64 {
	YDim := setting.TargetDim
	TaoMax := setting.TaoMax
	TaoMin := setting.TaoMin
	Rho := setting.Rho
	return func(i, j, step int) float64 {
		p, found := PheromoneMatrix[i*YDim+j]
		if found {
			return p
		}
		p = TaoMax * math.Pow(1-Rho, float64(step))
		PheromoneMatrix[i*YDim+j] = Clip(p, TaoMin, TaoMax)
		return p
	}
}

func DefaultPheromoneDeltaCalculate(setting *ACOSetting) func(v float64) float64 {
	var Q float64 = -1.0
	var Initialized = false
	taomin := math.Max(setting.TaoMin, 0.5)
	return func(v float64) float64 {
		if !Initialized {
			Initialized = true
			Q = v * taomin
		}
		return Q / v
	}
}

func DefaultExpectValue(setting *ACOSetting) func(i, j int) float64 {
	var YDim float64 = float64(setting.TargetDim)
	return func(i, j int) float64 {
		return 1.0 / YDim
	}
}
