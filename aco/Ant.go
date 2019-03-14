package aco

// 一个连接
type Connection struct {
	// 源节点
	source int
	// 目标节点
	target int
}

func (c Connection) Copy() Connection {
	return Connection{
		source: c.source,
		target: c.target,
	}
}

// 所有的连接
type MatchVector []Connection

func (s MatchVector) Copy() MatchVector {
	t := make(MatchVector, len(s))
	for i, v := range s {
		t[i] = v.Copy()
	}
	return t
}

type Ant struct {
	// 所有匹配
	Matches MatchVector
	// 目标值
	EvaluatedValue float64
	// 设置参数索引
	Setting *ACOSetting
}

func NewAnt(setting *ACOSetting) *Ant {
	return &Ant{
		Setting: setting,
		Matches: make(MatchVector, setting.SourceDim),
	}
}
