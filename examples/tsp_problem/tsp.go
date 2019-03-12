package main

import (
	"bufio"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"

	"../../pso"
)

type Slice struct {
	sort.Interface
	idx []int
}

func (s Slice) Swap(i, j int) {
	s.Interface.Swap(i, j)
	s.idx[i], s.idx[j] = s.idx[j], s.idx[i]
}

func NewSlice(start, count, step int) []int {
	s := make([]int, count)
	for i := range s {
		s[i] = start
		start += step
	}
	return s
}

func EUDistance(p1, p2 [2]float64) float64 {
	distance := math.Pow(p1[0]-p2[0], 2) + math.Pow(p1[1]-p2[1], 2)
	distance = math.Sqrt(distance)
	return distance
}

func ImportData(filename string) [][2]float64 {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	// 读取数据点数目
	var dim int = 0
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "DIMENSION") {
			strArray := strings.Split(line, ":")
			dim, _ = strconv.Atoi(strings.Trim(strArray[1], " "))
		}
		if strings.HasPrefix(line, "NODE_COORD_SECTION") {
			break
		}
	}

	// 生成输出的数据集合
	var datas = make([][2]float64, dim)
	index := 0
	for scanner.Scan() && index < dim {
		line := scanner.Text()
		strArray := strings.Split(line, " ")
		datas[index][0], _ = strconv.ParseFloat(strings.Trim(strArray[1], " "), 64)
		datas[index][1], _ = strconv.ParseFloat(strings.Trim(strArray[2], " "), 64)
		index += 1
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return datas
}

func GetEvaluater() func(x []float64) float64 {
	var points = ImportData("../data/tsp/ch71009.tsp.txt")
	return func(x []float64) float64 {
		keyvalue := Slice{
			Interface: sort.Float64Slice(x),
			idx:       NewSlice(0, len(x), 1),
		}
		sort.Sort(keyvalue)
		order := keyvalue.idx
		var sum float64 = 0
		for i := 0; i < len(order)-1; i++ {
			sum += EUDistance(points[order[i]], points[order[i+1]])
		}
		return -sum
	}
}

func main() {
	// 设置
	var setting = pso.GetDefaultSetting()
	setting.NDim = 71009
	setting.Evaluater = GetEvaluater()
	setting.Steps = 1000
	setting.XLow = -1
	setting.XHigh = 1
	setting.VMax = 1
	setting.Size = 100
	setting.Check()

	// 创建粒子群
	swarm := pso.NewSwarm(setting)
	swarm.Run()
}
