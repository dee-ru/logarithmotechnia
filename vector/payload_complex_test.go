package vector

import (
	"fmt"
	"math"
	"math/cmplx"
	"reflect"
	"strconv"
	"testing"
)

func TestComplex(t *testing.T) {
	emptyNA := []bool{false, false, false, false, false}

	testData := []struct {
		name          string
		data          []complex128
		na            []bool
		names         map[string]int
		expectedNames map[string]int
		isEmpty       bool
	}{
		{
			name:    "normal + na",
			data:    []complex128{1.1 + 0i, 2.2 + 2.2i, 3.3 + 3.3i, 4.4 + 4.4i, 5.5 - 5.5i},
			na:      []bool{false, false, false, false, false},
			names:   nil,
			isEmpty: false,
		},
		{
			name:    "normal + empty na",
			data:    []complex128{1.1 + 0i, 2.2 + 2.2i, 3.3 + 3.3i, 4.4 + 4.4i, 5.5 - 5.5i},
			na:      []bool{},
			names:   nil,
			isEmpty: false,
		},
		{
			name:    "normal + nil na",
			data:    []complex128{1.1 + 0i, 2.2 + 2.2i, 3.3 + 3.3i, 4.4 + 4.4i, 5.5 - 5.5i},
			na:      nil,
			names:   nil,
			isEmpty: false,
		},
		{
			name:    "normal + na",
			data:    []complex128{1.1 + 0i, 2.2 + 2.2i, 3.3 + 3.3i, 4.4 + 4.4i, 5.5 - 5.5i},
			na:      []bool{false, true, true, true, false},
			names:   nil,
			isEmpty: false,
		},
		{
			name:    "normal + incorrect sized na",
			data:    []complex128{1.1 + 0i, 2.2 + 2.2i, 3.3 + 3.3i, 4.4 + 4.4i, 5.5 - 5.5i},
			na:      []bool{false, false, false, false},
			names:   nil,
			isEmpty: true,
		},
		{
			name:          "normal + names",
			data:          []complex128{1.1 + 0i, 2.2 + 2.2i, 3.3 + 3.3i, 4.4 + 4.4i, 5.5 - 5.5i},
			na:            []bool{false, false, false, false, false},
			names:         map[string]int{"one": 1, "three": 3, "five": 5},
			expectedNames: map[string]int{"one": 1, "three": 3, "five": 5},
			isEmpty:       false,
		},
		{
			name:          "normal + incorrect names",
			data:          []complex128{1.1 + 0i, 2.2 + 2.2i, 3.3 + 3.3i, 4.4 + 4.4i, 5.5 - 5.5i},
			na:            []bool{false, false, false, false, false},
			names:         map[string]int{"zero": 0, "one": 1, "three": 3, "five": 5, "seven": 7},
			expectedNames: map[string]int{"one": 1, "three": 3, "five": 5},
			isEmpty:       false,
		},
	}

	for _, data := range testData {
		t.Run(data.name, func(t *testing.T) {
			var v Vector
			if data.names == nil {
				v = Complex(data.data, data.na)
			} else {
				config := Config{NamesMap: data.names}
				v = Complex(data.data, data.na, config).(*vector)
			}

			vv := v.(*vector)

			if data.isEmpty {
				_, ok := vv.payload.(*emptyPayload)
				if !ok {
					t.Error("Vector's payload is not empty")
				}
			} else {
				length := len(data.data)
				if vv.length != length {
					t.Error(fmt.Sprintf("Vector length (%d) is not equal to data length (%d)\n", vv.length, length))
				}

				payload, ok := vv.payload.(*complexPayload)
				if !ok {
					t.Error("Payload is not complexPayload")
				} else {
					if !reflect.DeepEqual(payload.data, data.data) {
						t.Error(fmt.Sprintf("Payload data (%v) is not equal to correct data (%v)\n",
							payload.data[1:], data.data))
					}

					if vv.length != vv.DefNameable.length || vv.length != payload.length {
						t.Error(fmt.Sprintf("Lengths are different: (vv.length - %d, "+
							"vv.DefNameable.length - %d, payload.length - %d, ",
							vv.length, vv.DefNameable.length, payload.length))
					}
				}

				if len(data.na) > 0 && len(data.na) == length {
					if !reflect.DeepEqual(payload.na, data.na) {
						t.Error(fmt.Sprintf("Payload na (%v) is not equal to correct na (%v)\n",
							payload.na[1:], data.na))
					}
				} else if len(data.na) == 0 {
					if !reflect.DeepEqual(payload.na, emptyNA) {
						t.Error(fmt.Sprintf("len(data.na) == 0 : incorrect payload.na (%v)", payload.na))
					}
				} else {
					t.Error("error")
				}

				if data.names != nil {
					if !reflect.DeepEqual(vv.names, data.expectedNames) {
						t.Error(fmt.Sprintf("Vector names (%v) is not equal to out names (%v)",
							vv.names, data.expectedNames))
					}
				}

			}
		})
	}
}

func TestComplex_Len(t *testing.T) {
	testData := []struct {
		in        []float64
		outLength int
	}{
		{[]float64{1, 2, 3, 4, 5}, 5},
		{[]float64{1, 2, 3}, 3},
		{[]float64{}, 0},
		{nil, 0},
	}

	for i, data := range testData {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			payload := Float(data.in, nil).(*vector).payload
			if payload.Len() != data.outLength {
				t.Error(fmt.Sprintf("Payloads's length (%d) is not equal to out (%d)",
					payload.Len(), data.outLength))
			}
		})
	}
}

func TestComplex_ByIndices(t *testing.T) {
	vec := Complex([]complex128{1, 2, 3, 4, 5}, []bool{false, false, false, false, true})
	testData := []struct {
		name    string
		indices []int
		out     []complex128
		outNA   []bool
	}{
		{
			name:    "all",
			indices: []int{1, 2, 3, 4, 5},
			out:     []complex128{1, 2, 3, 4, 5},
			outNA:   []bool{false, false, false, false, true},
		},
		{
			name:    "all reverse",
			indices: []int{5, 4, 3, 2, 1},
			out:     []complex128{5, 4, 3, 2, 1},
			outNA:   []bool{true, false, false, false, false},
		},
		{
			name:    "some",
			indices: []int{5, 1, 3},
			out:     []complex128{5, 1, 3},
			outNA:   []bool{true, false, false},
		},
	}

	for _, data := range testData {
		t.Run(data.name, func(t *testing.T) {
			payload := vec.ByIndices(data.indices).(*vector).payload.(*complexPayload)
			if !reflect.DeepEqual(payload.data, data.out) {
				t.Error(fmt.Sprintf("payload.data (%v) is not equal to data.out (%v)", payload.data, data.out))
			}
			if !reflect.DeepEqual(payload.na, data.outNA) {
				t.Error(fmt.Sprintf("payload.data (%v) is not equal to data.out (%v)", payload.data, data.out))
			}
		})
	}
}

func TestComplex_SupportsSelector(t *testing.T) {
	testData := []struct {
		name        string
		filter      interface{}
		isSupported bool
	}{
		{
			name:        "func(int, complex128, bool) bool",
			filter:      func(int, complex128, bool) bool { return true },
			isSupported: true,
		},
		{
			name:        "func(int, int, bool) bool",
			filter:      func(int, int, bool) bool { return true },
			isSupported: false,
		},
	}

	payload := Complex([]complex128{1}, nil).(*vector).payload
	for _, data := range testData {
		t.Run(data.name, func(t *testing.T) {
			if payload.SupportsSelector(data.filter) != data.isSupported {
				t.Error("Selector's support is incorrect.")
			}
		})
	}
}

func TestComplex_Filter(t *testing.T) {
	testData := []struct {
		name   string
		filter interface{}
		out    []bool
	}{
		{
			name:   "func(int, complex128, bool) bool",
			filter: func(idx int, val complex128, na bool) bool { return idx == 1 || val == 3 || na == true },
			out:    []bool{true, false, true, false, true},
		},
	}

	payload := Complex([]complex128{5, 4, 3, 2, 1}, []bool{false, false, false, false, true}).(*vector).payload

	for _, data := range testData {
		t.Run(data.name, func(t *testing.T) {
			filtered := payload.Select(data.filter)
			if !reflect.DeepEqual(payload.Select(data.filter), data.out) {
				t.Error(fmt.Sprintf("payload.Select() (%v) is not equal to out value (%v)",
					filtered, data.out))
			}
		})
	}
}

func TestComplex_Select(t *testing.T) {
	testData := []struct {
		name string
		fn   interface{}
		out  []bool
	}{
		{
			name: "Odd",
			fn:   func(idx int, _ complex128, _ bool) bool { return idx%2 == 1 },
			out:  []bool{true, false, true, false, true, false, true, false, true, false},
		},
		{
			name: "Even",
			fn:   func(idx int, _ complex128, _ bool) bool { return idx%2 == 0 },
			out:  []bool{false, true, false, true, false, true, false, true, false, true},
		},
		{
			name: "Nth(3)",
			fn:   func(idx int, _ complex128, _ bool) bool { return idx%3 == 0 },
			out:  []bool{false, false, true, false, false, true, false, false, true, false},
		},
		{
			name: "func() bool {return true}",
			fn:   func() bool { return true },
			out:  []bool{false, false, false, false, false, false, false, false, false, false},
		},
	}

	payload := Complex([]complex128{1, 2, 39, 4, 56, 2, 45, 90, 4, 3}, nil).(*vector).payload

	for _, data := range testData {
		t.Run(data.name, func(t *testing.T) {
			result := payload.Select(data.fn)
			if !reflect.DeepEqual(result, data.out) {
				t.Error(fmt.Sprintf("Result (%v) is not equal to out (%v)", result, data.out))
			}
		})
	}
}

func TestComplex_Booleans(t *testing.T) {
	testData := []struct {
		in    []complex128
		inNA  []bool
		out   []bool
		outNA []bool
	}{
		{
			in:    []complex128{1, 3, 0, 100, 0},
			inNA:  []bool{false, false, false, false, false},
			out:   []bool{true, true, false, true, false},
			outNA: []bool{false, false, false, false, false},
		},
		{
			in:    []complex128{10, 0, 12, 14, 1110},
			inNA:  []bool{false, false, false, true, true},
			out:   []bool{true, false, true, false, false},
			outNA: []bool{false, false, false, true, true},
		},
		{
			in:    []complex128{1, 3, 0, 100, 0, -11, -10},
			inNA:  []bool{false, false, false, false, false, false, true},
			out:   []bool{true, true, false, true, false, true, false},
			outNA: []bool{false, false, false, false, false, false, true},
		},
	}

	for i, data := range testData {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			vec := Complex(data.in, data.inNA)
			payload := vec.(*vector).payload.(*complexPayload)

			booleans, na := payload.Booleans()
			if !reflect.DeepEqual(booleans, data.out) {
				t.Error(fmt.Sprintf("Booleans (%v) are not equal to data.out (%v)\n", booleans, data.out))
			}
			if !reflect.DeepEqual(na, data.outNA) {
				t.Error(fmt.Sprintf("IsNA (%v) are not equal to data.outNA (%v)\n", na, data.outNA))
			}
		})
	}
}

func TestComplex_Integers(t *testing.T) {
	testData := []struct {
		in    []complex128
		inNA  []bool
		out   []int
		outNA []bool
	}{
		{
			in:    []complex128{1, 3, 0, 100, 0},
			inNA:  []bool{false, false, false, false, false},
			out:   []int{1, 3, 0, 100, 0},
			outNA: []bool{false, false, false, false, false},
		},
		{
			in:    []complex128{10, 0, 12, 14, 1110},
			inNA:  []bool{false, false, false, true, true},
			out:   []int{10, 0, 12, 0, 0},
			outNA: []bool{false, false, false, true, true},
		},
		{
			in:    []complex128{1, 3, 0, 100, 0, -11, -10},
			inNA:  []bool{false, false, false, false, false, false, true},
			out:   []int{1, 3, 0, 100, 0, -11, 0},
			outNA: []bool{false, false, false, false, false, false, true},
		},
	}

	for i, data := range testData {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			vec := Complex(data.in, data.inNA)
			payload := vec.(*vector).payload.(*complexPayload)

			integers, na := payload.Integers()
			if !reflect.DeepEqual(integers, data.out) {
				t.Error(fmt.Sprintf("Integers (%v) are not equal to data.out (%v)\n", integers, data.out))
			}
			if !reflect.DeepEqual(na, data.outNA) {
				t.Error(fmt.Sprintf("IsNA (%v) are not equal to data.outNA (%v)\n", na, data.outNA))
			}
		})
	}
}

func TestComplex_Floats(t *testing.T) {
	testData := []struct {
		in    []complex128
		inNA  []bool
		out   []float64
		outNA []bool
	}{
		{
			in:    []complex128{1, 3, 0, 100, 0},
			inNA:  []bool{false, false, false, false, false},
			out:   []float64{1, 3, 0, 100, 0},
			outNA: []bool{false, false, false, false, false},
		},
		{
			in:    []complex128{10, 0, 12, 14, 1110},
			inNA:  []bool{false, false, false, true, true},
			out:   []float64{10, 0, 12, math.NaN(), math.NaN()},
			outNA: []bool{false, false, false, true, true},
		},
		{
			in:    []complex128{1, 3, 0, 100, 0, -11, -10},
			inNA:  []bool{false, false, false, false, false, false, true},
			out:   []float64{1, 3, 0, 100, 0, -11, math.NaN()},
			outNA: []bool{false, false, false, false, false, false, true},
		},
	}

	for i, data := range testData {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			vec := Complex(data.in, data.inNA)
			payload := vec.(*vector).payload.(*complexPayload)

			floats, na := payload.Floats()
			correct := true
			for i := 0; i < len(floats); i++ {
				if math.IsNaN(data.out[i]) {
					if !math.IsNaN(floats[i]) {
						correct = false
					}
				} else if floats[i] != data.out[i] {
					correct = false
				}
			}
			if !correct {
				t.Error(fmt.Sprintf("Floats (%v) are not equal to data.out (%v)\n", floats, data.out))
			}
			if !reflect.DeepEqual(na, data.outNA) {
				t.Error(fmt.Sprintf("IsNA (%v) are not equal to data.outNA (%v)\n", na, data.outNA))
			}
		})
	}
}

func TestComplex_Complexes(t *testing.T) {
	testData := []struct {
		in    []complex128
		inNA  []bool
		out   []complex128
		outNA []bool
	}{
		{
			in:    []complex128{1, 3, 0, 100, 0, cmplx.NaN()},
			inNA:  []bool{false, false, false, false, false, false},
			out:   []complex128{1 + 0i, 3 + 0i, 0 + 0i, 100 + 0i, 0 + 0i, cmplx.NaN()},
			outNA: []bool{false, false, false, false, false, false},
		},
		{
			in:    []complex128{10, 0, 12, 14, 1110},
			inNA:  []bool{false, false, false, true, true},
			out:   []complex128{10 + 0i, 0 + 0i, 12 + 0i, cmplx.NaN(), cmplx.NaN()},
			outNA: []bool{false, false, false, true, true},
		},
		{
			in:    []complex128{1, 3, 0, 100, 0, -11, -10},
			inNA:  []bool{false, false, false, false, false, false, true},
			out:   []complex128{1 + 0i, 3 + 0i, 0 + 0i, 100 + 0i, 0 + 0i, -11 + 0i, cmplx.NaN()},
			outNA: []bool{false, false, false, false, false, false, true},
		},
	}

	for i, data := range testData {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			vec := Complex(data.in, data.inNA)
			payload := vec.(*vector).payload.(*complexPayload)

			complexes, na := payload.Complexes()
			correct := true
			for i := 0; i < len(complexes); i++ {
				if cmplx.IsNaN(data.out[i]) {
					if !cmplx.IsNaN(complexes[i]) {
						correct = false
					}
				} else if complexes[i] != data.out[i] {
					correct = false
				}
			}
			if !correct {
				t.Error(fmt.Sprintf("Complexes (%v) are not equal to data.out (%v)\n", complexes, data.out))
			}
			if !reflect.DeepEqual(na, data.outNA) {
				t.Error(fmt.Sprintf("IsNA (%v) are not equal to data.outNA (%v)\n", na, data.outNA))
			}
		})
	}
}

func TestComplex_Strings(t *testing.T) {
	testData := []struct {
		in    []complex128
		inNA  []bool
		out   []string
		outNA []bool
	}{
		{
			in:    []complex128{1, 3, cmplx.NaN(), 100, 0, cmplx.Inf()},
			inNA:  []bool{false, false, false, false, false, false},
			out:   []string{"(1.000+0.000i)", "(3.000+0.000i)", "NaN", "(100.000+0.000i)", "(0.000+0.000i)", "Inf"},
			outNA: []bool{false, false, false, false, false, false},
		},
		{
			in:    []complex128{10, 0, 12, 14, 1110},
			inNA:  []bool{false, false, false, true, true},
			out:   []string{"(10.000+0.000i)", "(0.000+0.000i)", "(12.000+0.000i)", "NA", "NA"},
			outNA: []bool{false, false, false, true, true},
		},
		{
			in:   []complex128{1, 3, cmplx.NaN(), 100, 0, -11, -10},
			inNA: []bool{false, false, false, false, false, false, true},
			out: []string{"(1.000+0.000i)", "(3.000+0.000i)", "NaN", "(100.000+0.000i)", "(0.000+0.000i)",
				"(-11.000+0.000i)", "NA"},
			outNA: []bool{false, false, false, false, false, false, true},
		},
	}

	for i, data := range testData {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			vec := Complex(data.in, data.inNA)
			payload := vec.(*vector).payload.(*complexPayload)

			strings, na := payload.Strings()
			if !reflect.DeepEqual(strings, data.out) {
				t.Error(fmt.Sprintf("Strings (%v) are not equal to data.out (%v)\n", strings, data.out))
			}
			if !reflect.DeepEqual(na, data.outNA) {
				t.Error(fmt.Sprintf("IsNA (%v) are not equal to data.outNA (%v)\n", na, data.outNA))
			}
		})
	}
}
