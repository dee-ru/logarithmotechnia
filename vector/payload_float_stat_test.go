package vector

import (
	"fmt"
	"logarithmotechnia/util"
	"math"
	"reflect"
	"testing"
)

func TestFloatPayload_Sum(t *testing.T) {
	testData := []struct {
		name    string
		payload *floatPayload
		data    []float64
		sumNA   []bool
	}{
		{
			name:    "without na",
			payload: FloatPayload([]float64{-20, 10, 4, -20, 27}, nil).(*floatPayload),
			data:    []float64{1},
			sumNA:   []bool{false},
		},
		{
			name:    "with na",
			payload: FloatPayload([]float64{-20, 10, 4, -20, 27}, []bool{false, false, true, false, false}).(*floatPayload),
			data:    []float64{math.NaN()},
			sumNA:   []bool{true},
		},
	}

	for _, data := range testData {
		t.Run(data.name, func(t *testing.T) {
			payload := data.payload.Sum().(*floatPayload)

			if !util.EqualFloatArrays(payload.data, data.data) {
				t.Error(fmt.Sprintf("Sum (%v) is not equal to expected (%v)",
					payload.data, data.data))
			}

			if !reflect.DeepEqual(payload.na, data.sumNA) {
				t.Error(fmt.Sprintf("Sum (%v) is not equal to expected (%v)",
					payload.na, data.sumNA))
			}
		})
	}
}

func TestFloatPayload_Max(t *testing.T) {
	testData := []struct {
		name    string
		payload *floatPayload
		data    []float64
		sumNA   []bool
	}{
		{
			name:    "without na",
			payload: FloatPayload([]float64{-20, 10, 4, -20, 27}, nil).(*floatPayload),
			data:    []float64{27},
			sumNA:   []bool{false},
		},
		{
			name:    "with na",
			payload: FloatPayload([]float64{-20, 10, 4, -20, 27}, []bool{false, false, true, false, false}).(*floatPayload),
			data:    []float64{math.NaN()},
			sumNA:   []bool{true},
		},
	}

	for _, data := range testData {
		t.Run(data.name, func(t *testing.T) {
			payload := data.payload.Max().(*floatPayload)

			if !util.EqualFloatArrays(payload.data, data.data) {
				t.Error(fmt.Sprintf("Max (%v) is not equal to expected (%v)",
					payload.data, data.data))
			}

			if !reflect.DeepEqual(payload.na, data.sumNA) {
				t.Error(fmt.Sprintf("Max (%v) is not equal to expected (%v)",
					payload.na, data.sumNA))
			}
		})
	}
}

func TestFloatPayload_Min(t *testing.T) {
	testData := []struct {
		name    string
		payload *floatPayload
		data    []float64
		sumNA   []bool
	}{
		{
			name:    "without na",
			payload: FloatPayload([]float64{-20, 10, 4, -20, 27}, nil).(*floatPayload),
			data:    []float64{-20},
			sumNA:   []bool{false},
		},
		{
			name:    "with na",
			payload: FloatPayload([]float64{-20, 10, 4, -20, 27}, []bool{false, false, true, false, false}).(*floatPayload),
			data:    []float64{math.NaN()},
			sumNA:   []bool{true},
		},
	}

	for _, data := range testData {
		t.Run(data.name, func(t *testing.T) {
			payload := data.payload.Min().(*floatPayload)

			if !util.EqualFloatArrays(payload.data, data.data) {
				t.Error(fmt.Sprintf("Min data (%v) is not equal to expected (%v)",
					payload.data, data.data))
			}

			if !reflect.DeepEqual(payload.na, data.sumNA) {
				t.Error(fmt.Sprintf("Min na (%v) is not equal to expected (%v)",
					payload.na, data.sumNA))
			}
		})
	}
}

func TestFloatPayload_Mean(t *testing.T) {
	testData := []struct {
		name    string
		payload *floatPayload
		data    []float64
		na      []bool
	}{
		{
			name:    "without na",
			payload: FloatPayload([]float64{-10, 10, 4, -20, 26}, nil).(*floatPayload),
			data:    []float64{2},
			na:      []bool{false},
		},
		{
			name:    "with na",
			payload: FloatPayload([]float64{-20, 10, 4, -20, 26}, []bool{false, false, true, false, false}).(*floatPayload),
			data:    []float64{math.NaN()},
			na:      []bool{true},
		},
	}

	for _, data := range testData {
		t.Run(data.name, func(t *testing.T) {
			payload := data.payload.Mean().(*floatPayload)

			if !util.EqualFloatArrays(payload.data, data.data) {
				t.Error(fmt.Sprintf("Meann data (%v) is not equal to expected (%v)",
					payload.data, data.data))
			}

			if !reflect.DeepEqual(payload.na, data.na) {
				t.Error(fmt.Sprintf("Mean na (%v) is not equal to expected (%v)",
					payload.na, data.na))
			}
		})
	}
}

func TestFloatPayload_Median(t *testing.T) {
	testData := []struct {
		name    string
		payload *floatPayload
		data    []float64
		na      []bool
	}{
		{
			name:    "without na odd",
			payload: FloatPayload([]float64{10, 26, 4, -20, 26}, nil).(*floatPayload),
			data:    []float64{10},
			na:      []bool{false},
		},
		{
			name:    "without na even",
			payload: FloatPayload([]float64{10, 26, 4, -20, 26, 16}, nil).(*floatPayload),
			data:    []float64{13},
			na:      []bool{false},
		},
		{
			name:    "with na",
			payload: FloatPayload([]float64{10, 26, 4, -20, 26}, []bool{false, false, true, false, false}).(*floatPayload),
			data:    []float64{math.NaN()},
			na:      []bool{true},
		},
		{
			name:    "one element",
			payload: FloatPayload([]float64{10}, nil).(*floatPayload),
			data:    []float64{10},
			na:      []bool{false},
		},
		{
			name:    "zero elements",
			payload: FloatPayload([]float64{}, nil).(*floatPayload),
			data:    []float64{math.NaN()},
			na:      []bool{true},
		},
	}

	for _, data := range testData {
		t.Run(data.name, func(t *testing.T) {
			payload := data.payload.Median().(*floatPayload)

			if !util.EqualFloatArrays(payload.data, data.data) {
				t.Error(fmt.Sprintf("Median data (%v) is not equal to expected (%v)",
					payload.data, data.data))
			}

			if !reflect.DeepEqual(payload.na, data.na) {
				t.Error(fmt.Sprintf("Mediann na (%v) is not equal to expected (%v)",
					payload.na, data.na))
			}
		})
	}
}
