package main

import (
	"math"

	v1 "github.com/klinakuf/caus/pkg/apis/caus.rss.uni-stuttgart.de/v1"
)

// AdjustBufferAmount adjusts the buffer size depedning on the current publishing rate
// publishingRate - the measuered publishing rate
// currentCapacity allocated - this is the total number of allocated instances
// currentBuffer - this is the current buffer size
// currentPerf - this is the  performance metric
// bufferThreshold - the threshold to increase the buffer size
func AdjustBufferAmount(publishingRate float64, currentReplicas float64, currentBuffer float64, currentPerf float64, bufferThreshold float64) (newBufferSize int32) {
	usage := publishingRate / ((currentReplicas - currentBuffer) * currentPerf)

	bufferThresh := bufferThreshold / 100.0

	//if the usage is touching the buffer check how much.
	if usage > 1 {
		difference := publishingRate - ((currentReplicas - currentBuffer) * currentPerf)
		bufferUsage := difference / (currentBuffer * currentPerf)
		if bufferUsage > bufferThresh {
			return int32(currentBuffer + 1)
		}
	} else {
		// if usage is less then we need to scale down the buffer
		// TODO: check this to add initialbuffer
		return int32(math.Max(1, currentBuffer-1))
	}
	return int32(currentBuffer)
}

// BaseWorkload calculates the base workload neeedd to cope with current publishing rate
// calculation methods
func BaseWorkload(publishingRate float64, currentPerf float64) (numRep int32) {
	return int32(math.Ceil(publishingRate / currentPerf))
}

// CalcReplicas for the given publishing Rate it will either return:
// - the minimum capacity if the publishingRate is less then the capacity
// - current number of replicas allocated
// - the maximum capacity if the workload+buffer exceeds the limit
// The logic behind the controller
func CalcReplicas(elasticity *v1.Elasticity, publishingRate float64, currentReplicas float64) (newNumRep int32, bufferSize int32) {

	// minimum capacity -------
	if publishingRate < float64(elasticity.Spec.Deployment.Capacity) {
		minReplicas := int32(1)
		if elasticity.Spec.Deployment.MinReplicas != nil {
			minReplicas = *elasticity.Spec.Deployment.MinReplicas
		}
		return minReplicas + elasticity.Spec.Buffer.Initial, elasticity.Spec.Buffer.Initial
	}

	var bufferForCalc float64
	if elasticity.Status.BufferedReplicas == 0 {
		bufferForCalc = float64(elasticity.Spec.Buffer.Initial)
	} else {
		bufferForCalc = float64(elasticity.Status.BufferedReplicas)
	}

	// current capacity ---
	baseWorkload := BaseWorkload(publishingRate, float64(elasticity.Spec.Deployment.Capacity))
	bufferSize = AdjustBufferAmount(publishingRate, float64(currentReplicas), bufferForCalc, float64(elasticity.Spec.Deployment.Capacity), float64(elasticity.Spec.Buffer.Threshold))
	totalReplicas := baseWorkload + bufferSize

	// maximum capacity ---
	if totalReplicas > elasticity.Spec.Deployment.MaxReplicas {
		totalReplicas = elasticity.Spec.Deployment.MaxReplicas
		bufferSize = elasticity.Spec.Buffer.Initial
	}

	return totalReplicas, bufferSize

}
