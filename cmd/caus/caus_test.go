package main

import "testing"
import v1 "github.com/klinakuf/caus/pkg/apis/caus.rss.uni-stuttgart.de/v1"
import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

func TestAdjustBufferAmount(t *testing.T) {
	currentPerf := 8.0
	cases := []struct {
		publishingRate, totalReplicas, currentBuffer float64
		threshold                                    int32
		newBufferSize                                int32
	}{
		{float64(20), float64(2), float64(1), 50, 2}, // it increases to one
		{float64(20), float64(4), float64(2), 50, 2}, // it stays the same because buffer is touched less then 0.5
		{float64(22), float64(5), float64(2), 50, 1}, // it should decrease the buffer
	}
	for _, c := range cases {
		result := AdjustBufferAmount(c.publishingRate, c.totalReplicas, c.currentBuffer, currentPerf, float64(c.threshold), 1)
		if result != c.newBufferSize {
			t.Errorf("adjustBufferAmount(%v) = %v, want %v", c.publishingRate, result, c.newBufferSize)
		}
	}
}

func TestCalcReplicas(t *testing.T) {
	minReplicas := int32(1)

	example := &v1.Elasticity{
		ObjectMeta: metav1.ObjectMeta{
			Name: "choreographer-elasticity",
		},
		Spec: v1.ElasticitySpec{
			Deployment: v1.DeploymentSpec{
				Name:        "com-rss-choreography-deployment",
				Capacity:    8,
				MinReplicas: &minReplicas,
				MaxReplicas: 20,
			},
			Workload: v1.WorkloadSpec{
				Queue: "NGP.Choreography.SimplePayrollRun",
			},
			Buffer: v1.BufferSpec{
				Initial:   1,
				Threshold: 50,
			},
		},
		Status: v1.ElasticityStatus{
			Message:          "Created, not processed yet",
			BufferedReplicas: 1,
		},
	}

	example2 := &v1.Elasticity{
		ObjectMeta: metav1.ObjectMeta{
			Name: "choreographer-elasticity",
		},
		Spec: v1.ElasticitySpec{
			Deployment: v1.DeploymentSpec{
				Name:        "com-sap-ngp-xx-choreography-deployment",
				Capacity:    8,
				MinReplicas: &minReplicas,
				MaxReplicas: 20,
			},
			Workload: v1.WorkloadSpec{
				Queue: "NGP.Choreography.SimplePayrollRun",
			},
			Buffer: v1.BufferSpec{
				Initial:   1,
				Threshold: 50,
			},
		},
		Status: v1.ElasticityStatus{
			Message: "Created, not processed yet",
		},
	}

	cases := []struct {
		testname                        string
		publishingRate, currentReplicas float64
		newNumRep                       int32
		newBufferSize                   int32
	}{
		{"c1", float64(20), float64(3), 4, 1},  //  the rate 20 is exactly at 0.5 so buffer should not increase
		{"c11", float64(22), float64(3), 5, 2}, // the rate 22 exceeds the buffer above 0.5 as defined in example Elasticity
		{"c2", float64(20), float64(4), 4, 1},  // the rate is 20 the current number of replicas is 3+1 it should remain 3+1
		{"c3", float64(22), float64(5), 4, 1},  // rate 22 4+1 replicas scale down to 3+1
	}
	for _, c := range cases {
		total, buffer := CalcReplicas(example, c.publishingRate, c.currentReplicas)
		if total != c.newNumRep {
			t.Errorf("Test with name %v  failed for total not expected CalcReplicas(%v) = %v, want %v", c.testname, c.publishingRate, total, c.newNumRep)
		}
		if buffer != c.newBufferSize {
			t.Errorf("Test with name %v  failed buffer not expected CalcReplicas(%v) = %v, want %v", c.testname, c.publishingRate, buffer, c.newBufferSize)
		}
	}

	// these tests, test the logic of scaling out and in of the behavior in the start when the rate is less then capacity, more  & below the buffer and more & above the buffer with threshold .5
	cases2 := []struct {
		testname                        string
		publishingRate, currentReplicas float64
		newNumRep                       int32
		newBufferSize                   int32
	}{
		{"c1start", float64(5), float64(1), 2, 1},  // if the controller starts and the first measured rate is 5 and 1 is the current replica we want to immediatly scale out
		{"c2start", float64(15), float64(1), 4, 2}, // the controller should start two and two as buffer because as the first value it is above the buffer and the threshold
		{"c3start", float64(19), float64(1), 5, 2}, // the controller should start 3 and two as a buffer
		{"c4start", float64(23), float64(1), 5, 2}, // the controller should start 3 and two as a buffer

	}
	for _, c := range cases2 {
		total, buffer := CalcReplicas(example2, c.publishingRate, c.currentReplicas)
		if total != c.newNumRep {
			t.Errorf("Test with name %v  failed for total not expected CalcReplicas(%v) = %v, want %v", c.testname, c.publishingRate, total, c.newNumRep)
		}
		if buffer != c.newBufferSize {
			t.Errorf("Test with name %v  failed buffer not expected CalcReplicas(%v) = %v, want %v", c.testname, c.publishingRate, buffer, c.newBufferSize)
		}
	}

}
