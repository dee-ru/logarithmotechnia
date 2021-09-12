package vector

import (
	"fmt"
	"math"
	"math/cmplx"
	"reflect"
	"strconv"
	"testing"
)

func TestString(t *testing.T) {
	emptyNA := []bool{false, false, false, false, false}

	testData := []struct {
		name    string
		data    []string
		na      []bool
		outData []string
		isEmpty bool
	}{
		{
			name:    "normal + false na",
			data:    []string{"one", "two", "three", "four", "five"},
			na:      []bool{false, false, false, false, false},
			outData: []string{"one", "two", "three", "four", "five"},
			isEmpty: false,
		},
		{
			name:    "normal + empty na",
			data:    []string{"one", "two", "three", "four", "five"},
			na:      []bool{},
			outData: []string{"one", "two", "three", "four", "five"},
			isEmpty: false,
		},
		{
			name:    "normal + nil na",
			data:    []string{"one", "two", "three", "four", "five"},
			na:      nil,
			outData: []string{"one", "two", "three", "four", "five"},
			isEmpty: false,
		},
		{
			name:    "normal + mixed na",
			data:    []string{"one", "two", "three", "four", "five"},
			na:      []bool{false, true, true, true, false},
			outData: []string{"one", "", "", "", "five"},
			isEmpty: false,
		},
		{
			name:    "normal + incorrect sized na",
			data:    []string{"one", "two", "three", "four", "five"},
			na:      []bool{false, false, false, false},
			isEmpty: true,
		},
	}

	for _, data := range testData {
		t.Run(data.name, func(t *testing.T) {
			v := String(data.data, data.na)

			vv := v.(*vector)

			if data.isEmpty {
				naPayload, ok := vv.payload.(*naPayload)
				if !ok || naPayload.Len() > 0 {
					t.Error("Vector's payload is not empty")
				}
			} else {
				length := len(data.data)
				if vv.length != length {
					t.Error(fmt.Sprintf("Vector length (%d) is not equal to data length (%d)\n", vv.length, length))
				}

				payload, ok := vv.payload.(*stringPayload)
				if !ok {
					t.Error("Payload is not stringPayload")
				} else {
					if !reflect.DeepEqual(payload.data, data.outData) {
						t.Error(fmt.Sprintf("Payload data (%v) is not equal to correct data (%v)\n",
							payload.data, data.outData))
					}
				}

				if len(data.na) > 0 && len(data.na) == length {
					if !reflect.DeepEqual(payload.na, data.na) {
						t.Error(fmt.Sprintf("Payload na (%v) is not equal to correct na (%v)\n",
							payload.na, data.na))
					}
				} else if len(data.na) == 0 {
					if !reflect.DeepEqual(payload.na, emptyNA) {
						t.Error(fmt.Sprintf("len(data.na) == 0 : incorrect payload.na (%v)", payload.na))
					}
				} else {
					t.Error("error")
				}
			}
		})
	}
}

func TestStringPayload_Type(t *testing.T) {
	vec := String([]string{}, nil)
	if vec.Type() != "string" {
		t.Error("Type is incorrect.")
	}
}

func TestStringPayload_Len(t *testing.T) {
	testData := []struct {
		in        []string
		outLength int
	}{
		{[]string{"one", "two", "three", "four", "five"}, 5},
		{[]string{"one", "two", "three"}, 3},
		{[]string{}, 0},
		{nil, 0},
	}

	for i, data := range testData {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			payload := String(data.in, nil).(*vector).payload
			if payload.Len() != data.outLength {
				t.Error(fmt.Sprintf("Payloads's length (%d) is not equal to out (%d)",
					payload.Len(), data.outLength))
			}
		})
	}
}

func TestStringPayload_Booleans(t *testing.T) {
	testData := []struct {
		in    []string
		inNA  []bool
		out   []bool
		outNA []bool
	}{
		{
			in:    []string{"1", "3", "", "100", ""},
			inNA:  []bool{false, false, false, false, false},
			out:   []bool{true, true, false, true, false},
			outNA: []bool{false, false, false, false, false},
		},
		{
			in:    []string{"10", "", "12", "14", "1110"},
			inNA:  []bool{false, false, false, true, true},
			out:   []bool{true, false, true, false, false},
			outNA: []bool{false, false, false, true, true},
		},
		{
			in:    []string{"1", "3", "", "100", "", "-11", "-10"},
			inNA:  []bool{false, false, false, false, false, false, true},
			out:   []bool{true, true, false, true, false, true, false},
			outNA: []bool{false, false, false, false, false, false, true},
		},
	}

	for i, data := range testData {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			vec := String(data.in, data.inNA)
			payload := vec.(*vector).payload.(*stringPayload)

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

func TestStringPayload_Integers(t *testing.T) {
	testData := []struct {
		in    []string
		inNA  []bool
		out   []int
		outNA []bool
	}{
		{
			in:    []string{"1", "3", "", "100", ""},
			inNA:  []bool{false, false, false, false, false},
			out:   []int{1, 3, 0, 100, 0},
			outNA: []bool{false, false, false, false, false},
		},
		{
			in:    []string{"10", "", "12", "14", "1110"},
			inNA:  []bool{false, false, false, true, true},
			out:   []int{10, 0, 12, 0, 0},
			outNA: []bool{false, false, false, true, true},
		},
		{
			in:    []string{"1", "3", "", "100", "", "-11", "-10"},
			inNA:  []bool{false, false, false, false, false, false, true},
			out:   []int{1, 3, 0, 100, 0, -11, 0},
			outNA: []bool{false, false, false, false, false, false, true},
		},
	}

	for i, data := range testData {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			vec := String(data.in, data.inNA)
			payload := vec.(*vector).payload.(*stringPayload)

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

func TestStringPayload_Floats(t *testing.T) {
	testData := []struct {
		in    []string
		inNA  []bool
		out   []float64
		outNA []bool
	}{
		{
			in:    []string{"1", "3", "", "100", ""},
			inNA:  []bool{false, false, false, false, false},
			out:   []float64{1, 3, 0, 100, 0},
			outNA: []bool{false, false, false, false, false},
		},
		{
			in:    []string{"10", "", "12", "14", "1110"},
			inNA:  []bool{false, false, false, true, true},
			out:   []float64{10, 0, 12, math.NaN(), math.NaN()},
			outNA: []bool{false, false, false, true, true},
		},
		{
			in:    []string{"1", "3", "", "100", "", "-11", "-10"},
			inNA:  []bool{false, false, false, false, false, false, true},
			out:   []float64{1, 3, 0, 100, 0, -11, math.NaN()},
			outNA: []bool{false, false, false, false, false, false, true},
		},
	}

	for i, data := range testData {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			vec := String(data.in, data.inNA)
			payload := vec.(*vector).payload.(*stringPayload)

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

func TestStringPayload_Complexes(t *testing.T) {
	testData := []struct {
		in    []string
		inNA  []bool
		out   []complex128
		outNA []bool
	}{
		{
			in:    []string{"1+1i", "3-3i", "0", "100 + 50i", "0+0i"},
			inNA:  []bool{false, false, false, false, false},
			out:   []complex128{1 + 1i, 3 - 3i, 0 + 0i, cmplx.NaN(), 0 + 0i},
			outNA: []bool{false, false, false, false, false},
		},
		{
			in:    []string{"10+10i", "0", "12+6i", "14+7i", "1110+0i"},
			inNA:  []bool{false, false, false, true, true},
			out:   []complex128{10 + 10i, 0 + 0i, 12 + 6i, cmplx.NaN(), cmplx.NaN()},
			outNA: []bool{false, false, false, true, true},
		},
	}

	for i, data := range testData {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			vec := String(data.in, data.inNA)
			payload := vec.(*vector).payload.(*stringPayload)

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

func TestStringPayload_Strings(t *testing.T) {
	testData := []struct {
		in    []string
		inNA  []bool
		out   []string
		outNA []bool
	}{
		{
			in:    []string{"1", "3", "0", "100", ""},
			inNA:  []bool{false, false, false, false, false},
			out:   []string{"1", "3", "0", "100", ""},
			outNA: []bool{false, false, false, false, false},
		},
		{
			in:    []string{"10", "", "12", "14", "1110"},
			inNA:  []bool{false, false, false, true, true},
			out:   []string{"10", "", "12", "", ""},
			outNA: []bool{false, false, false, true, true},
		},
		{
			in:    []string{"1", "3", "0", "100", "", "-11", "-10"},
			inNA:  []bool{false, false, false, false, false, false, true},
			out:   []string{"1", "3", "0", "100", "", "-11", ""},
			outNA: []bool{false, false, false, false, false, false, true},
		},
	}

	for i, data := range testData {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			vec := String(data.in, data.inNA)
			payload := vec.(*vector).payload.(*stringPayload)

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

func TestStringPayload_Interfaces(t *testing.T) {
	testData := []struct {
		in    []string
		inNA  []bool
		out   []interface{}
		outNA []bool
	}{
		{
			in:    []string{"1", "3", "0", "100", ""},
			inNA:  []bool{false, false, false, false, false},
			out:   []interface{}{"1", "3", "0", "100", ""},
			outNA: []bool{false, false, false, false, false},
		},
		{
			in:    []string{"10", "", "12", "14", "1110"},
			inNA:  []bool{false, false, false, true, true},
			out:   []interface{}{"10", "", "12", nil, nil},
			outNA: []bool{false, false, false, true, true},
		},
		{
			in:    []string{"1", "3", "0", "100", "", "-11", "-10"},
			inNA:  []bool{false, false, false, false, false, false, true},
			out:   []interface{}{"1", "3", "0", "100", "", "-11", nil},
			outNA: []bool{false, false, false, false, false, false, true},
		},
	}

	for i, data := range testData {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			vec := String(data.in, data.inNA)
			payload := vec.(*vector).payload.(*stringPayload)

			interfaces, na := payload.Interfaces()
			if !reflect.DeepEqual(interfaces, data.out) {
				t.Error(fmt.Sprintf("Interfaces (%v) are not equal to data.out (%v)\n", interfaces, data.out))
			}
			if !reflect.DeepEqual(na, data.outNA) {
				t.Error(fmt.Sprintf("IsNA (%v) are not equal to data.outNA (%v)\n", na, data.outNA))
			}
		})
	}
}

func TestStringPayload_ByIndices(t *testing.T) {
	vec := String([]string{"1", "2", "3", "4", "5"}, []bool{false, false, false, false, true})
	testData := []struct {
		name    string
		indices []int
		out     []string
		outNA   []bool
	}{
		{
			name:    "all",
			indices: []int{1, 2, 3, 4, 5},
			out:     []string{"1", "2", "3", "4", ""},
			outNA:   []bool{false, false, false, false, true},
		},
		{
			name:    "all reverse",
			indices: []int{5, 4, 3, 2, 1},
			out:     []string{"", "4", "3", "2", "1"},
			outNA:   []bool{true, false, false, false, false},
		},
		{
			name:    "some",
			indices: []int{5, 1, 3},
			out:     []string{"", "1", "3"},
			outNA:   []bool{true, false, false},
		},
	}

	for _, data := range testData {
		t.Run(data.name, func(t *testing.T) {
			payload := vec.ByIndices(data.indices).(*vector).payload.(*stringPayload)
			if !reflect.DeepEqual(payload.data, data.out) {
				t.Error(fmt.Sprintf("payload.data (%v) is not equal to data.out (%v)", payload.data, data.out))
			}
			if !reflect.DeepEqual(payload.na, data.outNA) {
				t.Error(fmt.Sprintf("payload.data (%v) is not equal to data.out (%v)", payload.data, data.out))
			}
		})
	}
}

func TestStringPayload_SupportsWhicher(t *testing.T) {
	testData := []struct {
		name        string
		filter      interface{}
		isSupported bool
	}{
		{
			name:        "func(int, string, bool) bool",
			filter:      func(int, string, bool) bool { return true },
			isSupported: true,
		},
		{
			name:        "func(string, bool) bool",
			filter:      func(string, bool) bool { return true },
			isSupported: true,
		},
		{
			name:        "func(int, float64, bool) bool",
			filter:      func(int, float64, bool) bool { return true },
			isSupported: false,
		},
	}

	payload := String([]string{"one"}, nil).(*vector).payload.(Whichable)
	for _, data := range testData {
		t.Run(data.name, func(t *testing.T) {
			if payload.SupportsWhicher(data.filter) != data.isSupported {
				t.Error("Selector's support is incorrect.")
			}
		})
	}
}

func TestStringPayload_Whicher(t *testing.T) {
	testData := []struct {
		name string
		fn   interface{}
		out  []bool
	}{
		{
			name: "Odd",
			fn:   func(idx int, _ string, _ bool) bool { return idx%2 == 1 },
			out:  []bool{true, false, true, false, true, false, true, false, true, false},
		},
		{
			name: "Even",
			fn:   func(idx int, _ string, _ bool) bool { return idx%2 == 0 },
			out:  []bool{false, true, false, true, false, true, false, true, false, true},
		},
		{
			name: "func(_ int, val string, _ bool) bool {return val == 2}",
			fn:   func(_ int, val string, _ bool) bool { return val == "2" },
			out:  []bool{false, true, false, false, false, true, false, false, false, false},
		},
		{
			name: "Comparer compact",
			fn:   func(val string, _ bool) bool { return val == "39" || val == "90" },
			out:  []bool{false, false, true, false, false, false, false, true, false, false},
		},
		{
			name: "func() bool {return true}",
			fn:   func() bool { return true },
			out:  []bool{false, false, false, false, false, false, false, false, false, false},
		},
	}

	payload := String([]string{"1", "2", "39", "4", "56", "2", "45", "90", "4", "3"}, nil).(*vector).payload.(Whichable)

	for _, data := range testData {
		t.Run(data.name, func(t *testing.T) {
			result := payload.Which(data.fn)
			if !reflect.DeepEqual(result, data.out) {
				t.Error(fmt.Sprintf("Result (%v) is not equal to out (%v)", result, data.out))
			}
		})
	}
}

func TestStringPayload_SupportsApplier(t *testing.T) {
	testData := []struct {
		name        string
		applier     interface{}
		isSupported bool
	}{
		{
			name:        "func(int, string, bool) (string, bool)",
			applier:     func(int, string, bool) (string, bool) { return "", true },
			isSupported: true,
		},
		{
			name:        "func(string, bool) (string, bool)",
			applier:     func(string, bool) (string, bool) { return "", true },
			isSupported: true,
		},
		{
			name:        "func(int, string, bool) bool",
			applier:     func(int, string, bool) bool { return true },
			isSupported: false,
		},
	}

	payload := String([]string{}, nil).(*vector).payload.(Appliable)
	for _, data := range testData {
		t.Run(data.name, func(t *testing.T) {
			if payload.SupportsApplier(data.applier) != data.isSupported {
				t.Error("Applier's support is incorrect.")
			}
		})
	}
}

func TestStringPayload_Apply(t *testing.T) {
	testData := []struct {
		name        string
		applier     interface{}
		dataIn      []string
		naIn        []bool
		dataOut     []string
		naOut       []bool
		isNAPayload bool
	}{
		{
			name: "regular",
			applier: func(_ int, val string, na bool) (string, bool) {
				return fmt.Sprintf("%s.%s", val, val), na
			},
			dataIn:      []string{"1", "9", "3", "5", "7"},
			naIn:        []bool{false, true, false, true, false},
			dataOut:     []string{"1.1", "", "3.3", "", "7.7"},
			naOut:       []bool{false, true, false, true, false},
			isNAPayload: false,
		},
		{
			name: "regular compact",
			applier: func(val string, na bool) (string, bool) {
				return fmt.Sprintf("%s.%s", val, val), na
			},
			dataIn:      []string{"1", "9", "3", "5", "7"},
			naIn:        []bool{false, true, false, true, false},
			dataOut:     []string{"1.1", "", "3.3", "", "7.7"},
			naOut:       []bool{false, true, false, true, false},
			isNAPayload: false,
		},
		{
			name: "manipulate na",
			applier: func(idx int, val string, na bool) (string, bool) {
				if idx == 5 {
					return "1", true
				}
				return val, na
			},
			dataIn:      []string{"1", "2", "3", "4", "5"},
			naIn:        []bool{false, false, true, false, false},
			dataOut:     []string{"1", "2", "", "4", ""},
			naOut:       []bool{false, false, true, false, true},
			isNAPayload: false,
		},
		{
			name:        "incorrect applier",
			applier:     func(int, string, bool) bool { return true },
			dataIn:      []string{"1", "9", "3", "5", "7"},
			naIn:        []bool{false, true, false, true, false},
			dataOut:     []string{"", "", "", "", ""},
			naOut:       []bool{true, true, true, true, true},
			isNAPayload: true,
		},
	}

	for _, data := range testData {
		t.Run(data.name, func(t *testing.T) {
			payload := String(data.dataIn, data.naIn).(*vector).payload.(Appliable).Apply(data.applier)

			if !data.isNAPayload {
				payloadOut := payload.(*stringPayload)
				if !reflect.DeepEqual(data.dataOut, payloadOut.data) {
					t.Error(fmt.Sprintf("Output data (%v) does not match expected (%v)",
						payloadOut.data, data.dataOut))
				}
				if !reflect.DeepEqual(data.naOut, payloadOut.na) {
					t.Error(fmt.Sprintf("Output NA (%v) does not match expected (%v)",
						payloadOut.na, data.naOut))
				}
			} else {
				_, ok := payload.(*naPayload)
				if !ok {
					t.Error("Payload is not NA")
				}
			}
		})
	}
}

func TestStringPayload_SupportsSummarizer(t *testing.T) {
	testData := []struct {
		name        string
		summarizer  interface{}
		isSupported bool
	}{
		{
			name:        "valid",
			summarizer:  func(int, string, string, bool) (string, bool) { return "", false },
			isSupported: true,
		},
		{
			name:        "invalid",
			summarizer:  func(int, int, bool) bool { return true },
			isSupported: false,
		},
	}

	payload := String([]string{}, nil).(*vector).payload.(Summarizable)
	for _, data := range testData {
		t.Run(data.name, func(t *testing.T) {
			if payload.SupportsSummarizer(data.summarizer) != data.isSupported {
				t.Error("Summarizer's support is incorrect.")
			}
		})
	}
}

func TestStringPayload_Summarize(t *testing.T) {
	summarizer := func(idx int, prev string, cur string, na bool) (string, bool) {
		return prev + cur, na
	}

	testData := []struct {
		name        string
		summarizer  interface{}
		dataIn      []string
		naIn        []bool
		dataOut     []string
		naOut       []bool
		isNAPayload bool
	}{
		{
			name:        "true",
			summarizer:  summarizer,
			dataIn:      []string{"1", "2", "1", "6", "5"},
			naIn:        []bool{false, false, false, false, false},
			dataOut:     []string{"12165"},
			naOut:       []bool{false},
			isNAPayload: false,
		},
		{
			name:        "NA",
			summarizer:  summarizer,
			dataIn:      []string{"1", "2", "1", "6", "5"},
			naIn:        []bool{false, false, false, false, true},
			isNAPayload: true,
		},
		{
			name:        "incorrect applier",
			summarizer:  func(int, int, bool) bool { return true },
			dataIn:      []string{"1", "2", "1", "6", "5"},
			naIn:        []bool{false, true, false, true, false},
			isNAPayload: true,
		},
	}

	for _, data := range testData {
		t.Run(data.name, func(t *testing.T) {
			payload := String(data.dataIn, data.naIn).(*vector).payload.(Summarizable).Summarize(data.summarizer)

			if !data.isNAPayload {
				payloadOut := payload.(*stringPayload)
				if !reflect.DeepEqual(data.dataOut, payloadOut.data) {
					t.Error(fmt.Sprintf("Output data (%v) does not match expected (%v)",
						data.dataOut, payloadOut.data))
				}
				if !reflect.DeepEqual(data.naOut, payloadOut.na) {
					t.Error(fmt.Sprintf("Output NA (%v) does not match expected (%v)",
						data.naOut, payloadOut.na))
				}
			} else {
				naPayload, ok := payload.(*naPayload)
				if ok {
					if naPayload.length != 1 {
						t.Error("Incorrect length of NA payload (not 1)")
					}
				} else {
					t.Error("Payload is not NA")
				}
			}
		})
	}
}

func TestStringPayload_Append(t *testing.T) {
	payload := StringPayload([]string{"1", "2", "3"}, nil)

	testData := []struct {
		name    string
		vec     Vector
		outData []string
		outNA   []bool
	}{
		{
			name:    "boolean",
			vec:     Boolean([]bool{true, true}, []bool{true, false}),
			outData: []string{"1", "2", "3", "", "true"},
			outNA:   []bool{false, false, false, true, false},
		},
		{
			name:    "integer",
			vec:     Integer([]int{4, 5}, []bool{true, false}),
			outData: []string{"1", "2", "3", "", "5"},
			outNA:   []bool{false, false, false, true, false},
		},
		{
			name:    "na",
			vec:     NA(2),
			outData: []string{"1", "2", "3", "", ""},
			outNA:   []bool{false, false, false, true, true},
		},
	}

	for _, data := range testData {
		t.Run(data.name, func(t *testing.T) {
			outPayload := payload.Append(data.vec.Payload()).(*stringPayload)

			if !reflect.DeepEqual(data.outData, outPayload.data) {
				t.Error(fmt.Sprintf("Output data (%v) does not match expected (%v)",
					outPayload.data, data.outData))
			}
			if !reflect.DeepEqual(data.outNA, outPayload.na) {
				t.Error(fmt.Sprintf("Output NA (%v) does not match expected (%v)",
					outPayload.na, data.outNA))
			}
		})
	}
}

func TestStringPayload_Adjust(t *testing.T) {
	payload5 := StringPayload([]string{"1", "2", "3", "4", "5"}, nil).(*stringPayload)
	payload3 := StringPayload([]string{"1", "2", "3"}, []bool{false, false, true}).(*stringPayload)

	testData := []struct {
		name       string
		inPayload  *stringPayload
		size       int
		outPaylout *stringPayload
	}{
		{
			inPayload:  payload5,
			name:       "same",
			size:       5,
			outPaylout: StringPayload([]string{"1", "2", "3", "4", "5"}, nil).(*stringPayload),
		},
		{
			inPayload:  payload5,
			name:       "lesser",
			size:       3,
			outPaylout: StringPayload([]string{"1", "2", "3"}, nil).(*stringPayload),
		},
		{
			inPayload: payload3,
			name:      "bigger",
			size:      10,
			outPaylout: StringPayload([]string{"1", "2", "0", "1", "2", "0", "1", "2", "0", "1"},
				[]bool{false, false, true, false, false, true, false, false, true, false}).(*stringPayload),
		},
	}

	for _, data := range testData {
		t.Run(data.name, func(t *testing.T) {
			outPayload := data.inPayload.Adjust(data.size).(*stringPayload)

			if !reflect.DeepEqual(outPayload.data, data.outPaylout.data) {
				t.Error(fmt.Sprintf("Output data (%v) does not match expected (%v)",
					outPayload.data, data.outPaylout.data))
			}
			if !reflect.DeepEqual(outPayload.na, data.outPaylout.na) {
				t.Error(fmt.Sprintf("Output NA (%v) does not match expected (%v)",
					outPayload.na, data.outPaylout.na))
			}
		})
	}
}

func TestStringPayload_Find(t *testing.T) {
	payload := StringPayload([]string{"1", "2", "1", "4", "0"}, nil).(*stringPayload)

	testData := []struct {
		name   string
		needle interface{}
		pos    int
	}{
		{"existent", "4", 4},
		{"non-existent", "non", 0},
		{"incorrect type", true, 0},
	}

	for _, data := range testData {
		t.Run(data.name, func(t *testing.T) {
			pos := payload.Find(data.needle)

			if pos != data.pos {
				t.Error(fmt.Sprintf("Position (%v) does not match expected (%v)",
					pos, data.pos))
			}
		})
	}
}

func TestStringPayload_FindAll(t *testing.T) {
	payload := StringPayload([]string{"1", "2", "1", "4", "0"}, nil).(*stringPayload)

	testData := []struct {
		name   string
		needle interface{}
		pos    []int
	}{
		{"existent", "1", []int{1, 3}},
		{"non-existent", -10, []int{}},
		{"incorrect type", false, []int{}},
	}

	for _, data := range testData {
		t.Run(data.name, func(t *testing.T) {
			pos := payload.FindAll(data.needle)

			if !reflect.DeepEqual(pos, data.pos) {
				t.Error(fmt.Sprintf("Positions (%v) does not match expected (%v)",
					pos, data.pos))
			}
		})
	}
}

func TestStringPayload_Eq(t *testing.T) {
	payload := StringPayload([]string{"2", "zero", "minus", "2", "1"}, nil).(*stringPayload)

	testData := []struct {
		eq  interface{}
		cmp []bool
	}{
		{"2", []bool{true, false, false, true, false}},
		{"zero", []bool{false, true, false, false, false}},
		{2, []bool{true, false, false, true, false}},

		{int64(1), []bool{false, false, false, false, true}},
		{int32(1), []bool{false, false, false, false, true}},
		{uint64(1), []bool{false, false, false, false, true}},
		{uint32(1), []bool{false, false, false, false, true}},

		{true, []bool{false, false, false, false, false}},
	}

	for i, data := range testData {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			cmp := payload.Eq(data.eq)

			if !reflect.DeepEqual(cmp, data.cmp) {
				t.Error(fmt.Sprintf("Comparator results (%v) do not match expected (%v)",
					cmp, data.cmp))
			}
		})
	}
}

func TestStringPayload_Neq(t *testing.T) {
	payload := StringPayload([]string{"2", "zero", "minus", "2", "1"}, nil).(*stringPayload)

	testData := []struct {
		eq  interface{}
		cmp []bool
	}{
		{"2", []bool{false, true, true, false, true}},
		{2, []bool{false, true, true, false, true}},

		{int64(1), []bool{true, true, true, true, false}},
		{int32(1), []bool{true, true, true, true, false}},
		{uint64(1), []bool{true, true, true, true, false}},
		{uint32(1), []bool{true, true, true, true, false}},

		{true, []bool{true, true, true, true, true}},
	}

	for i, data := range testData {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			cmp := payload.Neq(data.eq)

			if !reflect.DeepEqual(cmp, data.cmp) {
				t.Error(fmt.Sprintf("Comparator results (%v) do not match expected (%v)",
					cmp, data.cmp))
			}
		})
	}
}

func TestStringPayload_Gt(t *testing.T) {
	payload := StringPayload([]string{"alpha", "zero", "zeroth", "betha", "gamma"}, nil).(*stringPayload)

	testData := []struct {
		val interface{}
		cmp []bool
	}{
		{"zero", []bool{false, false, true, false, false}},
		{true, []bool{false, false, false, false, false}},
	}

	for i, data := range testData {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			cmp := payload.Gt(data.val)

			if !reflect.DeepEqual(cmp, data.cmp) {
				t.Error(fmt.Sprintf("Comparator results (%v) do not match expected (%v)",
					cmp, data.cmp))
			}
		})
	}
}

func TestStringPayload_Lt(t *testing.T) {
	payload := StringPayload([]string{"alpha", "zero", "zeroth", "betha", "gamma"}, nil).(*stringPayload)

	testData := []struct {
		val interface{}
		cmp []bool
	}{
		{"zero", []bool{true, false, false, true, true}},
		{true, []bool{false, false, false, false, false}},
	}

	for i, data := range testData {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			cmp := payload.Lt(data.val)

			if !reflect.DeepEqual(cmp, data.cmp) {
				t.Error(fmt.Sprintf("Comparator results (%v) do not match expected (%v)",
					cmp, data.cmp))
			}
		})
	}
}

func TestStringPayload_Gte(t *testing.T) {
	payload := StringPayload([]string{"alpha", "zero", "zeroth", "betha", "gamma"}, nil).(*stringPayload)

	testData := []struct {
		val interface{}
		cmp []bool
	}{
		{"zero", []bool{false, true, true, false, false}},
		{true, []bool{false, false, false, false, false}},
	}

	for i, data := range testData {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			cmp := payload.Gte(data.val)

			if !reflect.DeepEqual(cmp, data.cmp) {
				t.Error(fmt.Sprintf("Comparator results (%v) do not match expected (%v)",
					cmp, data.cmp))
			}
		})
	}
}

func TestStringPayload_Lte(t *testing.T) {
	payload := StringPayload([]string{"alpha", "zero", "zeroth", "betha", "gamma"}, nil).(*stringPayload)

	testData := []struct {
		val interface{}
		cmp []bool
	}{
		{"zero", []bool{true, true, false, true, true}},
		{true, []bool{false, false, false, false, false}},
	}

	for i, data := range testData {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			cmp := payload.Lte(data.val)

			if !reflect.DeepEqual(cmp, data.cmp) {
				t.Error(fmt.Sprintf("Comparator results (%v) do not match expected (%v)",
					cmp, data.cmp))
			}
		})
	}
}
