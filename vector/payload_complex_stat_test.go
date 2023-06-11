package vector

import (
	"fmt"
	"logarithmotechnia/internal/util"
	"math/cmplx"
	"reflect"
	"testing"
)

func TestComplexPayload_Sum(t *testing.T) {
	testData := []struct {
		name    string
		payload *complexPayload
		sumData []complex128
		sumNA   []bool
	}{
		{
			name:    "without na",
			payload: ComplexPayload([]complex128{-20 + 10i, 10 - 5i, 4 + 2i, -20 + 20i, 27 - 26i}, nil).(*complexPayload),
			sumData: []complex128{1 + 1i},
			sumNA:   []bool{false},
		},
		{
			name: "with na",
			payload: ComplexPayload([]complex128{-20 + 10i, 10 - 5i, 4 + 2i, -20 + 20i, 27 - 26i},
				[]bool{false, false, true, false, false}).(*complexPayload),
			sumData: []complex128{cmplx.NaN()},
			sumNA:   []bool{true},
		},
	}

	for _, data := range testData {
		t.Run(data.name, func(t *testing.T) {
			sumPayload := data.payload.Sum().(*complexPayload)

			if !util.EqualComplexArrays(sumPayload.data, data.sumData) {
				t.Error(fmt.Sprintf("Sum data (%v) is not equal to expected (%v)",
					sumPayload.data, data.sumData))
			}

			if !reflect.DeepEqual(sumPayload.NA, data.sumNA) {
				t.Error(fmt.Sprintf("Sum data (%v) is not equal to expected (%v)",
					sumPayload.NA, data.sumNA))
			}
		})
	}
}

func TestComplexPayload_Mean(t *testing.T) {
	testData := []struct {
		name    string
		payload *complexPayload
		data    []complex128
		NA      []bool
	}{
		{
			name:    "without na",
			payload: ComplexPayload([]complex128{-10 + 10i, 10 - 5i, 4 - 5i, -20 + 20i, 26 - 10i}, nil).(*complexPayload),
			data:    []complex128{2 + 2i},
			NA:      []bool{false},
		},
		{
			name: "with na",
			payload: ComplexPayload([]complex128{-20 + 10i, 10 - 5i, 4 + 2i, -20 + 20i, 26 - 26i},
				[]bool{false, false, true, false, false}).(*complexPayload),
			data: []complex128{cmplx.NaN()},
			NA:   []bool{true},
		},
	}

	for _, data := range testData {
		t.Run(data.name, func(t *testing.T) {
			payload := data.payload.Mean().(*complexPayload)

			if !util.EqualComplexArrays(payload.data, data.data) {
				t.Error(fmt.Sprintf("Sum data (%v) is not equal to expected (%v)",
					payload.data, data.data))
			}

			if !reflect.DeepEqual(payload.NA, data.NA) {
				t.Error(fmt.Sprintf("Sum data (%v) is not equal to expected (%v)",
					payload.NA, data.NA))
			}
		})
	}
}

func TestComplexPayload_Prod(t *testing.T) {
	testData := []struct {
		name    string
		payload *complexPayload
		data    []complex128
		NA      []bool
	}{
		{
			name:    "without na",
			payload: ComplexPayload([]complex128{-10 + 10i, 10 - 5i, 4 - 5i, -20 + 20i, 26 - 10i}, nil).(*complexPayload),
			data:    []complex128{-788000 + 124000i},
			NA:      []bool{false},
		},
		{
			name: "with na",
			payload: ComplexPayload([]complex128{-20 + 10i, 10 - 5i, 4 + 2i, -20 + 20i, 26 - 26i},
				[]bool{false, false, true, false, false}).(*complexPayload),
			data: []complex128{cmplx.NaN()},
			NA:   []bool{true},
		},
		{
			name:    "one element",
			payload: ComplexPayload([]complex128{-10 + 10i}, nil).(*complexPayload),
			data:    []complex128{-10 + 10i},
			NA:      []bool{false},
		},
		{
			name:    "zero element",
			payload: ComplexPayload([]complex128{}, nil).(*complexPayload),
			data:    []complex128{0},
			NA:      []bool{false},
		},
	}

	for _, data := range testData {
		t.Run(data.name, func(t *testing.T) {
			payload := data.payload.Prod().(*complexPayload)

			if !util.EqualComplexArrays(payload.data, data.data) {
				t.Error(fmt.Sprintf("Prod data (%v) is not equal to expected (%v)",
					payload.data, data.data))
			}

			if !reflect.DeepEqual(payload.NA, data.NA) {
				t.Error(fmt.Sprintf("Prod data (%v) is not equal to expected (%v)",
					payload.NA, data.NA))
			}
		})
	}
}

func TestComplexPayload_CumSum(t *testing.T) {
	testData := []struct {
		name    string
		payload *complexPayload
		data    []complex128
		na      []bool
	}{
		{
			name:    "without na",
			payload: ComplexPayload([]complex128{-10 + 10i, 10 - 5i, 4 - 5i, -20 + 20i, 26 - 10i}, nil).(*complexPayload),
			data:    []complex128{-10 + 10i, 0 + 5i, 4 + 0i, -16 + 20i, 10 + 10i},
			na:      []bool{false, false, false, false, false},
		},
		{
			name: "with na",
			payload: ComplexPayload([]complex128{-20 + 10i, 10 - 5i, 4 + 2i, -20 + 20i, 26 - 26i},
				[]bool{false, false, true, false, false}).(*complexPayload),
			data: []complex128{-20 + 10i, -10 + 5i, cmplx.NaN(), cmplx.NaN(), cmplx.NaN()},
			na:   []bool{false, false, true, true, true},
		},
		{
			name:    "one element",
			payload: ComplexPayload([]complex128{-10 + 10i}, nil).(*complexPayload),
			data:    []complex128{-10 + 10i},
			na:      []bool{false},
		},
		{
			name:    "zero element",
			payload: ComplexPayload([]complex128{}, nil).(*complexPayload),
			data:    []complex128{},
			na:      []bool{},
		},
	}

	for _, data := range testData {
		t.Run(data.name, func(t *testing.T) {
			payload := data.payload.CumSum().(*complexPayload)

			if !util.EqualComplexArrays(payload.data, data.data) {
				t.Error(fmt.Sprintf("Prod data (%v) is not equal to expected (%v)",
					payload.data, data.data))
			}

			if !reflect.DeepEqual(payload.NA, data.na) {
				t.Error(fmt.Sprintf("Prod data (%v) is not equal to expected (%v)",
					payload.NA, data.na))
			}
		})
	}
}

func TestComplexPayload_CumProd(t *testing.T) {
	testData := []struct {
		name    string
		payload *complexPayload
		data    []complex128
		na      []bool
	}{
		{
			name:    "without na",
			payload: ComplexPayload([]complex128{-10 + 10i, 10 - 5i, 4 - 5i, -20 + 20i, 26 - 10i}, nil).(*complexPayload),
			data:    []complex128{-10 + 10i, -50 + 150i, 550 + 850i, -28000 - 6000i, -788000 + 124000i},
			na:      []bool{false, false, false, false, false},
		},
		{
			name: "with na",
			payload: ComplexPayload([]complex128{-20 + 10i, 10 - 5i, 4 + 2i, -20 + 20i, 26 - 26i},
				[]bool{false, false, true, false, false}).(*complexPayload),
			data: []complex128{-20 + 10i, -150 + 200i, cmplx.NaN(), cmplx.NaN(), cmplx.NaN()},
			na:   []bool{false, false, true, true, true},
		},
		{
			name:    "one element",
			payload: ComplexPayload([]complex128{-10 + 10i}, nil).(*complexPayload),
			data:    []complex128{-10 + 10i},
			na:      []bool{false},
		},
		{
			name:    "zero element",
			payload: ComplexPayload([]complex128{}, nil).(*complexPayload),
			data:    []complex128{},
			na:      []bool{},
		},
	}

	for _, data := range testData {
		t.Run(data.name, func(t *testing.T) {
			payload := data.payload.CumProd().(*complexPayload)

			if !util.EqualComplexArrays(payload.data, data.data) {
				t.Error(fmt.Sprintf("Prod data (%v) is not equal to expected (%v)",
					payload.data, data.data))
			}

			if !reflect.DeepEqual(payload.NA, data.na) {
				t.Error(fmt.Sprintf("Prod data (%v) is not equal to expected (%v)",
					payload.NA, data.na))
			}
		})
	}
}
