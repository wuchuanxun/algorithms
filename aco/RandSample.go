package aco

import "math/rand"

func PdfSample(pdf []float64) int {
	cdf := make([]float64, len(pdf))
	cdf[0] = pdf[0]
	for index := 1; index < len(pdf); index++ {
		cdf[index] = cdf[index-1] + pdf[index]
	}
	if cdf[len(pdf)-1] != 1 {
		scale := cdf[len(pdf)-1]
		for i, _ := range cdf {
			cdf[i] /= scale
		}
	}
	return CdfSample(cdf)
}

func CdfSample(cdf []float64) int {
	r := rand.Float64()

	bucket := 0
	for r > cdf[bucket] {
		bucket++
	}
	return bucket
}
