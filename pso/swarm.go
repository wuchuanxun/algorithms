package pso

import (
	"fmt"
	"math"
)

type Swarm struct {
	// 全局最佳位置
	GBest FVector
	// 全局最佳适应度
	GBFitness float64
	// 第i次迭代最佳位置
	IBest FVector
	// 第i次迭代最佳适应度
	IBFitness float64
	// 迭代次数
	Step int
	// 粒子群
	ParticleGroup []*Particle
	// 全局设置
	setting *PSOSetting
}

// 每一代的初始化
func (swarm *Swarm) ResetIthFitness() {
	swarm.IBFitness = -math.MaxFloat32
}

// 更新最佳的位置
func (swarm *Swarm) UpdateBest(particle *Particle) {
	if particle.Fitness > swarm.IBFitness {
		swarm.IBFitness = particle.Fitness
		swarm.IBest = particle.Location.Copy()
		if swarm.IBFitness > swarm.GBFitness {
			swarm.GBFitness = swarm.IBFitness
			swarm.GBest = swarm.IBest.Copy()
		}
	}
}

// 生成一个粒子群
func NewSwarm(setting *PSOSetting) *Swarm {
	var swarm = &Swarm{
		Step:          0,
		IBFitness:     -math.MaxFloat32,
		GBFitness:     -math.MaxFloat32,
		setting:       setting,
		ParticleGroup: make([]*Particle, setting.Size),
	}
	for i := 0; i < setting.Size; i++ {
		swarm.ParticleGroup[i] = NewParticle(setting)
		swarm.UpdateBest(swarm.ParticleGroup[i])
	}
	return swarm
}

func (swarm *Swarm) Run() FVector {
	for swarm.Step = 0; swarm.Step < swarm.setting.Steps; swarm.Step += 1 {
		swarm.ResetIthFitness()
		for i, particle := range swarm.ParticleGroup {
			particle.Update(i, swarm.IBest, swarm.GBest)
			swarm.UpdateBest(particle)
		}
		if swarm.setting.ShowDetails && (swarm.Step%swarm.setting.ShowInterval == 0) {
			fmt.Printf("Step %d :: min err=%.5e\n", swarm.Step, swarm.GBFitness)
		}
	}
	return swarm.GBest
}
