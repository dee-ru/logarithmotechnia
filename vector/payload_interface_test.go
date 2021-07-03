package vector

import (
	"fmt"
	"logarithmotechnia.com/logarithmotechnia/util"
	"math"
	"reflect"
	"strconv"
	"testing"
)

func TestInterface(t *testing.T) {
	testInterfaceEmpty(t)

	emptyNA := []bool{false, false, false, false, false}

	testData := []struct {
		name          string
		data          []interface{}
		na            []bool
		outData       []interface{}
		names         map[string]int
		expectedNames map[string]int
	}{
		{
			name:    "normal + false na",
			data:    []interface{}{1, 2, 3, 4, 5},
			na:      []bool{false, false, false, false, false},
			outData: []interface{}{1, 2, 3, 4, 5},
			names:   nil,
		},
		{
			name:    "normal + empty na",
			data:    []interface{}{1, 2, 3, 4, 5},
			na:      []bool{},
			outData: []interface{}{1, 2, 3, 4, 5},
			names:   nil,
		},
		{
			name:    "normal + nil na",
			data:    []interface{}{1, 2, 3, 4, 5},
			na:      nil,
			outData: []interface{}{1, 2, 3, 4, 5},
			names:   nil,
		},
		{
			name:    "normal + mixed na",
			data:    []interface{}{1, 2, 3, 4, 5},
			na:      []bool{false, true, true, true, false},
			outData: []interface{}{1, nil, nil, nil, 5},
			names:   nil,
		},
		{
			name:          "normal + names",
			data:          []interface{}{1, 2, 3, 4, 5},
			na:            []bool{false, false, false, false, false},
			outData:       []interface{}{1, 2, 3, 4, 5},
			names:         map[string]int{"one": 1, "three": 3, "five": 5},
			expectedNames: map[string]int{"one": 1, "three": 3, "five": 5},
		},
		{
			name:          "normal + incorrect names",
			data:          []interface{}{1, 2, 3, 4, 5},
			na:            []bool{false, false, false, false, false},
			outData:       []interface{}{1, 2, 3, 4, 5},
			names:         map[string]int{"zero": 0, "one": 1, "three": 3, "five": 5, "seven": 7},
			expectedNames: map[string]int{"one": 1, "three": 3, "five": 5},
		},
	}

	for _, data := range testData {
		t.Run(data.name, func(t *testing.T) {
			config := Config{NamesMap: data.names}
			v := Interface(data.data, data.na, config).(*vector)

			length := len(data.data)
			if v.length != length {
				t.Error(fmt.Sprintf("Vector length (%d) is not equal to data length (%d)\n", v.length, length))
			}

			payload, ok := v.payload.(*interfacePayload)
			if !ok {
				t.Error("Payload is not interfacePayload")
			} else {
				if !reflect.DeepEqual(payload.data, data.outData) {
					t.Error(fmt.Sprintf("Payload data (%v) is not equal to correct data (%v)\n",
						payload.data, data.data))
				}

				if v.length != v.DefNameable.length || v.length != payload.length {
					t.Error(fmt.Sprintf("Lengths are different: (vv.length - %d, "+
						"vv.DefNameable.length - %d, payload.length - %d, ",
						v.length, v.DefNameable.length, payload.length))
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

			if data.names != nil {
				if !reflect.DeepEqual(v.names, data.expectedNames) {
					t.Error(fmt.Sprintf("Vector names (%v) is not equal to out names (%v)",
						v.names, data.expectedNames))
				}
			}
		})
	}
}

func testInterfaceEmpty(t *testing.T) {
	vec := Interface([]interface{}{1, 2, 3, 4, 5}, []bool{false, false, true, false})
	_, ok := vec.(*vector).payload.(*emptyPayload)
	if !ok {
		t.Error("Vector's payload is not empty")
	}
}

func TestInterfacePayload_Type(t *testing.T) {
	vec := Interface([]interface{}{}, nil)
	if vec.Type() != "interface" {
		t.Error("Type is incorrect.")
	}
}

func TestInterfacePayload_Len(t *testing.T) {
	testData := []struct {
		in        []interface{}
		outLength int
	}{
		{[]interface{}{1, 2, 3, 4, 5}, 5},
		{[]interface{}{1, 2, 3}, 3},
		{[]interface{}{}, 0},
		{nil, 0},
	}

	for i, data := range testData {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			payload := Interface(data.in, nil).(*vector).payload
			if payload.Len() != data.outLength {
				t.Error(fmt.Sprintf("Payloads's length (%d) is not equal to out (%d)",
					payload.Len(), data.outLength))
			}
		})
	}
}

func TestInterfacePayload_SupportsWhicher(t *testing.T) {
	testData := []struct {
		name        string
		filter      interface{}
		isSupported bool
	}{
		{
			name:        "func(int, interface{}, bool) bool",
			filter:      func(int, interface{}, bool) bool { return true },
			isSupported: true,
		},
		{
			name:        "func(int, float64, bool) bool",
			filter:      func(int, float64, bool) bool { return true },
			isSupported: false,
		},
	}

	payload := Interface([]interface{}{1}, nil).(*vector).payload.(Whichable)
	for _, data := range testData {
		t.Run(data.name, func(t *testing.T) {
			if payload.SupportsWhicher(data.filter) != data.isSupported {
				t.Error("Selector's support is incorrect.")
			}
		})
	}
}

func TestInterfacePayload_Whicher(t *testing.T) {
	testData := []struct {
		name string
		fn   interface{}
		out  []bool
	}{
		{
			name: "Odd",
			fn:   func(idx int, _ interface{}, _ bool) bool { return idx%2 == 1 },
			out:  []bool{true, false, true, false, true, false, true, false, true, false},
		},
		{
			name: "Even",
			fn:   func(idx int, _ interface{}, _ bool) bool { return idx%2 == 0 },
			out:  []bool{false, true, false, true, false, true, false, true, false, true},
		},
		{
			name: "Nth(3)",
			fn:   func(idx int, _ interface{}, _ bool) bool { return idx%3 == 0 },
			out:  []bool{false, false, true, false, false, true, false, false, true, false},
		},
		{
			name: "func() bool {return true}",
			fn:   func() bool { return true },
			out:  []bool{false, false, false, false, false, false, false, false, false, false},
		},
	}

	payload := Interface([]interface{}{true, false, true, false, true, false, true, false, true, false}, nil).(*vector).payload.(Whichable)

	for _, data := range testData {
		t.Run(data.name, func(t *testing.T) {
			result := payload.Which(data.fn)
			if !reflect.DeepEqual(result, data.out) {
				t.Error(fmt.Sprintf("Result (%v) is not equal to out (%v)", result, data.out))
			}
		})
	}
}

func TestInterfacePayload_SupportsApplier(t *testing.T) {
	testData := []struct {
		name        string
		applier     interface{}
		isSupported bool
	}{
		{
			name:        "func(int, interface{}, bool) (bool, bool)",
			applier:     func(int, interface{}, bool) (interface{}, bool) { return 1, true },
			isSupported: true,
		},
		{
			name:        "func(int, float64, bool) bool",
			applier:     func(int, int, bool) bool { return true },
			isSupported: false,
		},
	}

	payload := Interface([]interface{}{}, nil).(*vector).payload.(Appliable)
	for _, data := range testData {
		t.Run(data.name, func(t *testing.T) {
			if payload.SupportsApplier(data.applier) != data.isSupported {
				t.Error("Applier's support is incorrect.")
			}
		})
	}
}

func TestInterfacePayload_Apply(t *testing.T) {
	testData := []struct {
		name    string
		applier interface{}
		dataIn  []interface{}
		naIn    []bool
		dataOut []interface{}
		naOut   []bool
	}{
		{
			name: "regular",
			applier: func(idx int, val interface{}, na bool) (interface{}, bool) {
				if idx == 5 {
					return 5, na
				}
				return val, na
			},
			dataIn:  []interface{}{true, true, true, false, false},
			naIn:    []bool{false, true, false, true, false},
			dataOut: []interface{}{true, nil, true, nil, 5},
			naOut:   []bool{false, true, false, true, false},
		},
		{
			name: "manipulate na",
			applier: func(idx int, val interface{}, na bool) (interface{}, bool) {
				newNA := na
				if idx == 5 {
					newNA = true
				}
				return val, newNA
			},
			dataIn:  []interface{}{true, true, false, false, true},
			naIn:    []bool{false, true, false, true, false},
			dataOut: []interface{}{true, nil, false, nil, nil},
			naOut:   []bool{false, true, false, true, true},
		},
		{
			name:    "incorrect applier",
			applier: func(int, int, bool) bool { return true },
			dataIn:  []interface{}{true, true, false, false, true},
			naIn:    []bool{false, true, false, true, false},
			dataOut: []interface{}{nil, nil, nil, nil, nil},
			naOut:   []bool{true, true, true, true, true},
		},
	}

	for _, data := range testData {
		t.Run(data.name, func(t *testing.T) {
			payload := Interface(data.dataIn, data.naIn).(*vector).payload
			payloadOut := payload.(Appliable).Apply(data.applier).(*interfacePayload)

			if !reflect.DeepEqual(data.dataOut, payloadOut.data) {
				t.Error(fmt.Sprintf("Output data (%v) does not match expected (%v)",
					payloadOut.data, data.dataOut))
			}
			if !reflect.DeepEqual(data.naOut, payloadOut.na) {
				t.Error(fmt.Sprintf("Output NA (%v) does not match expected (%v)",
					payloadOut.na, data.naOut))
			}
		})
	}
}

func TestInterfacePayload_Integers(t *testing.T) {
	convertor := func(idx int, val interface{}, na bool) (int, bool) {
		if na {
			return 0, true
		}

		switch v := val.(type) {
		case float64:
			return int(v), false
		case int:
			return v, false
		default:
			return 0, true
		}
	}

	testData := []struct {
		name      string
		dataIn    []interface{}
		naIn      []bool
		convertor func(idx int, val interface{}, na bool) (int, bool)
		dataOut   []int
		naOut     []bool
	}{
		{
			name:      "regular",
			dataIn:    []interface{}{1, 2.5, "three", 4 + 3i, 5, 0},
			naIn:      []bool{false, false, false, false, true, false},
			convertor: convertor,
			dataOut:   []int{1, 2, 0, 0, 0, 0},
			naOut:     []bool{false, false, true, true, true, false},
		},
		{
			name:      "without converter",
			dataIn:    []interface{}{1, 2.5, "three", 4 + 3i, 5, 0},
			naIn:      []bool{false, false, false, false, true, false},
			convertor: nil,
			dataOut:   []int{0, 0, 0, 0, 0, 0},
			naOut:     []bool{true, true, true, true, true, true},
		},
		{
			name:      "empty",
			dataIn:    []interface{}{},
			naIn:      []bool{},
			convertor: convertor,
			dataOut:   []int{},
			naOut:     []bool{},
		},
	}

	for _, data := range testData {
		t.Run(data.name, func(t *testing.T) {
			vec := Interface(data.dataIn, data.naIn, OptionConvertors(InterfaceConvertors{Intabler: data.convertor}))
			payload := vec.(*vector).payload.(*interfacePayload)

			integers, na := payload.Integers()
			if !reflect.DeepEqual(integers, data.dataOut) {
				t.Error(fmt.Sprintf("Result data (%v) is not equal to expected (%v)", integers, data.dataOut))
			}
			if !reflect.DeepEqual(na, data.naOut) {
				t.Error(fmt.Sprintf("Result na (%v) is not equal to expected (%v)", na, data.naOut))
			}
		})
	}
}

func TestInterfacePayload_Floats(t *testing.T) {
	convertor := func(idx int, val interface{}, na bool) (float64, bool) {
		if na {
			return math.NaN(), true
		}

		switch v := val.(type) {
		case float64:
			return v, false
		case int:
			return float64(v), false
		default:
			return math.NaN(), true
		}
	}

	testData := []struct {
		name      string
		dataIn    []interface{}
		naIn      []bool
		convertor func(idx int, val interface{}, na bool) (float64, bool)
		dataOut   []float64
		naOut     []bool
	}{
		{
			name:      "regular",
			dataIn:    []interface{}{1, 2.5, "three", 4 + 3i, 5, 0},
			naIn:      []bool{false, false, false, false, true, false},
			convertor: convertor,
			dataOut:   []float64{1, 2.5, math.NaN(), math.NaN(), math.NaN(), 0},
			naOut:     []bool{false, false, true, true, true, false},
		},
		{
			name:      "without converter",
			dataIn:    []interface{}{1, 2.5, "three", 4 + 3i, 5, 0},
			naIn:      []bool{false, false, false, false, true, false},
			convertor: nil,
			dataOut:   []float64{math.NaN(), math.NaN(), math.NaN(), math.NaN(), math.NaN(), math.NaN()},
			naOut:     []bool{true, true, true, true, true, true},
		},
		{
			name:      "empty",
			dataIn:    []interface{}{},
			naIn:      []bool{},
			convertor: convertor,
			dataOut:   []float64{},
			naOut:     []bool{},
		},
	}

	for _, data := range testData {
		t.Run(data.name, func(t *testing.T) {
			vec := Interface(data.dataIn, data.naIn, OptionConvertors(InterfaceConvertors{Floatabler: data.convertor}))
			payload := vec.(*vector).payload.(*interfacePayload)

			floats, na := payload.Floats()
			if !util.EqualFloatArrays(floats, data.dataOut) {
				t.Error(fmt.Sprintf("Result data (%v) is not equal to expected (%v)", floats, data.dataOut))
			}
			if !reflect.DeepEqual(na, data.naOut) {
				t.Error(fmt.Sprintf("Result na (%v) is not equal to expected (%v)", na, data.naOut))
			}
		})
	}
}
