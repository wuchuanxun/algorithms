package aco

import (
	"math"
)

type Swarm struct {
	// 本次最佳蚂蚁
	IBestAnt *Ant
	// 迭代次数
	Step int
	// 粒子群
	AntGroup []*Ant
	// 全局设置
	Setting *ACOSetting
}

// 每一代的初始化
func (swarm *Swarm) ResetIth() {
	swarm.IBestAnt.EvaluatedValue = math.MaxFloat64
}

// 更新最佳的位置
func (swarm *Swarm) UpdateBest(ant *Ant) {
	if ant.EvaluatedValue < swarm.IBestAnt.EvaluatedValue {
		swarm.IBestAnt.EvaluatedValue = ant.EvaluatedValue
		swarm.IBestAnt.Matches = ant.Matches.Copy()
	}
}

// 生成一个粒子群
func NewSwarm(setting *ACOSetting) *Swarm {
	var swarm = &Swarm{
		Step:     0,
		IBestAnt: NewAnt(setting),
		Setting:  setting,
		AntGroup: make([]*Ant, setting.Size),
	}
	return swarm
}

func (swarm *Swarm) Run() MatchVector {
	for swarm.Step = 0; swarm.Step < swarm.Setting.Steps; swarm.Step += 1 {
		swarm.ResetIth()
		for _, ant := range swarm.AntGroup {
			// 实现一只蚂蚁的游走
			swarm.Setting.AntWalkMethod(ant, swarm.Step, swarm.Setting.GetPheromone, swarm.Setting.GetExpect)
			ant.EvaluatedValue = swarm.Setting.Evaluater(ant.Matches)
			swarm.UpdateBest(ant)
		}
		// 更新信息素矩阵
		dtaoMap := make(map[int]float64)
		dtao := swarm.Setting.GetPheromoneDelta(swarm.IBestAnt.EvaluatedValue)
		for _, Connc := range swarm.IBestAnt.Matches {
			dtaoMap[Connc.source*swarm.Setting.TargetDim+Connc.target] = dtao
		}
		for i := 0; i < swarm.Setting.SourceDim; i++ {
			for j := 0; j < swarm.Setting.TargetDim; j++ {
				swarm.Setting.PheromoneUpdateMethod(i, j, swarm.Step, dtaoMap[i*swarm.Setting.TargetDim+j])
			}
		}
	}
	return swarm.IBestAnt.Matches
}
