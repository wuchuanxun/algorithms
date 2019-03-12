package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"

	"../../pso"
)

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

func main() {
	var setting = pso.GetDefaultSetting()
	fmt.Printf("%v", setting)
	ImportData("../data/tsp/ch71009.tsp.txt")
}
