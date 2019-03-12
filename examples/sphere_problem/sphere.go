package main

import (
	"fmt"
	"math"

	"../../pso"
)

func Sphere(x []float64) float64 {
	var sum float64 = 0
	for i := 0; i < len(x); i++ {
		sum += math.Pow(x[i], 2)
	}
	return -sum
}

func main() {
	// 设置
	var setting = pso.GetDefaultSetting()
	setting.NDim = 3
	setting.Evaluater = Sphere
	setting.Steps = 600
	setting.XLow = -100
	setting.XHigh = 100
	setting.VMax = 30
	setting.Size = pso.CalculateSwarmSize(setting.NDim, pso.PSO_MAX_SIZE)
	setting.Check()

	// 创建粒子群
	swarm := pso.NewSwarm(setting)
	result := swarm.Run()
	fmt.Printf("%v", result)
}
