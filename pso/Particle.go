package pso

type FVector []float64

// 复制的函数
func (s FVector) Copy() FVector {
	t := make(FVector, len(s))
	// copy the slice
	copy(t, s)
	// return the copy target
	return t
}

// 单个粒子
type Particle struct {
	// 所在的位置
	Location FVector
	// 曾经最好的位置
	PreBestLoc FVector
	// 速度
	Velocity FVector
	// 适应度
	Fitness float64
	// 最佳的适应度
	BestFitness float64
	// 设置参数索引
	Setting *PSOSetting
}

// 随机生成一个切片
func GenRandomFVevtor(setting *PSOSetting) FVector {
	// 取出参数方便操作
	xlow := setting.XLow
	xhigh := setting.XHigh
	cap := setting.NDim
	// 初始化
	slice := make(FVector, cap)
	for i := 0; i < cap; i++ {
		slice[i] = xlow + (xhigh-xlow)*setting.RandGenerater.Float64()
	}
	return slice
}

// 创建一个新粒子
func NewParticle(setting *PSOSetting) *Particle {
	var particle = &Particle{
		Location: GenRandomFVevtor(setting),
	}
	particle.PreBestLoc = particle.Location.Copy()
	// 设置初始的速度
	particle.Velocity = make(FVector, setting.NDim)
	P1 := GenRandomFVevtor(setting)
	P2 := GenRandomFVevtor(setting)
	for i := 0; i < setting.NDim; i++ {
		particle.Velocity[i] = (P1[i] - P2[i]) / 2.0
	}
	// 设置参数
	particle.Setting = setting
	// 计算适应度
	particle.Fitness = particle.Setting.Evaluater(particle.Location)
	particle.BestFitness = particle.Fitness
	return particle
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

// 更新一个粒子，ibest是第i代的最佳位置，gbest是全局最佳位置
func (particle *Particle) Update(step int, ibest, gbest FVector) {
	setting := particle.Setting
	// 更新惯性系数
	var w float64 = (setting.WMax-setting.WMin)*float64(setting.Steps-step)/float64(setting.Steps) + setting.WMin
	// 计算速度
	var dv = make(FVector, setting.NDim)
	switch setting.UpdateMethod {
	case GlobalOnly:
		for i := 0; i < setting.NDim; i++ {
			dv1 := w * particle.Velocity[i]
			dv2 := setting.C1 * setting.RandGenerater.Float64() * (particle.PreBestLoc[i] - particle.Location[i])
			dv3 := setting.C2 * setting.RandGenerater.Float64() * (gbest[i] - particle.Location[i])
			dv[i] = dv1 + dv2 + dv3
		}
		break
	case GlobalLocal:
		for i := 0; i < setting.NDim; i++ {
			dv1 := w * particle.Velocity[i]
			dv2 := setting.C1 * setting.RandGenerater.Float64() * (setting.Alpha*(particle.PreBestLoc[i]-
				particle.Location[i]) + (1-setting.Alpha)*(ibest[i]-particle.Location[i]))
			dv3 := setting.C2 * setting.RandGenerater.Float64() * (gbest[i] - particle.Location[i])
			dv[i] = dv1 + dv2 + dv3
		}
		break
	default:
		panic("Undefined update method!")
	}
	// 更新速度和位置
	for i := 0; i < setting.NDim; i++ {
		particle.Velocity[i] = Clip(dv[i], -setting.VMax, setting.VMax)
		particle.Location[i] += particle.Velocity[i]
		if setting.ClipPos {
			particle.Location[i] = Clip(particle.Location[i], setting.XLow, setting.XHigh)
		}
	}
	// 更新适应度函数
	particle.Fitness = setting.Evaluater(particle.Location)
	// 更新历史最佳
	if particle.Fitness > particle.BestFitness {
		particle.BestFitness = particle.Fitness
		particle.PreBestLoc = particle.Location.Copy()
	}
}
