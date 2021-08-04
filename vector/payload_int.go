package vector

import (
	"math"
	"math/cmplx"
	"strconv"
)

const maxIntPrint = 5

// integerPayload is a structure, subsisting Integer vectors
type integerPayload struct {
	length int
	data   []int
	DefNAble
}

func (p *integerPayload) Type() string {
	return "integer"
}

func (p *integerPayload) Len() int {
	return p.length
}

func (p *integerPayload) ByIndices(indices []int) Payload {
	data := make([]int, 0, len(indices))
	na := make([]bool, 0, len(indices))

	for _, idx := range indices {
		data = append(data, p.data[idx-1])
		na = append(na, p.na[idx-1])
	}

	return &integerPayload{
		length: len(data),
		data:   data,
		DefNAble: DefNAble{
			na: na,
		},
	}
}

func (p *integerPayload) SupportsWhicher(whicher interface{}) bool {
	if _, ok := whicher.(func(int, int, bool) bool); ok {
		return true
	}

	return false
}

func (p *integerPayload) Which(whicher interface{}) []bool {
	if byFunc, ok := whicher.(func(int, int, bool) bool); ok {
		return p.selectByFunc(byFunc)
	}

	return make([]bool, p.length)
}

func (p *integerPayload) selectByFunc(byFunc func(int, int, bool) bool) []bool {
	booleans := make([]bool, p.length)

	for idx, val := range p.data {
		if byFunc(idx+1, val, p.na[idx]) {
			booleans[idx] = true
		}
	}

	return booleans
}

func (p *integerPayload) SupportsApplier(applier interface{}) bool {
	if _, ok := applier.(func(int, int, bool) (int, bool)); ok {
		return true
	}

	return false
}

func (p *integerPayload) Apply(applier interface{}) Payload {
	var data []int
	var na []bool

	if applyFunc, ok := applier.(func(int, int, bool) (int, bool)); ok {
		data, na = p.applyByFunc(applyFunc)
	} else {
		return NAPayload(p.length)
	}

	return &integerPayload{
		length: p.length,
		data:   data,
		DefNAble: DefNAble{
			na: na,
		},
	}
}

func (p *integerPayload) applyByFunc(applyFunc func(int, int, bool) (int, bool)) ([]int, []bool) {
	data := make([]int, p.length)
	na := make([]bool, p.length)

	for i := 0; i < p.length; i++ {
		dataVal, naVal := applyFunc(i+1, p.data[i], p.na[i])
		if naVal {
			dataVal = 0
		}
		data[i] = dataVal
		na[i] = naVal
	}

	return data, na
}

func (p *integerPayload) SupportsSummarizer(summarizer interface{}) bool {
	if _, ok := summarizer.(func(int, int, int, bool) (int, bool)); ok {
		return true
	}

	return false
}

func (p *integerPayload) Summarize(summarizer interface{}) Payload {
	fn, ok := summarizer.(func(int, int, int, bool) (int, bool))
	if !ok {
		return NAPayload(1)
	}

	val := 0
	na := false
	for i := 0; i < p.length; i++ {
		val, na = fn(i+1, val, p.data[i], p.na[i])

		if na {
			return NAPayload(1)
		}
	}

	return IntegerPayload([]int{val}, nil)
}

func (p *integerPayload) Integers() ([]int, []bool) {
	if p.length == 0 {
		return []int{}, []bool{}
	}

	data := make([]int, p.length)
	copy(data, p.data)

	na := make([]bool, p.Len())
	copy(na, p.na)

	return data, na
}

func (p *integerPayload) Floats() ([]float64, []bool) {
	if p.length == 0 {
		return []float64{}, []bool{}
	}

	data := make([]float64, p.length)

	for i := 0; i < p.length; i++ {
		if p.na[i] {
			data[i] = math.NaN()
		} else {
			data[i] = float64(p.data[i])
		}
	}

	na := make([]bool, p.Len())
	copy(na, p.na)

	return data, na
}

func (p *integerPayload) Complexes() ([]complex128, []bool) {
	if p.length == 0 {
		return []complex128{}, []bool{}
	}

	data := make([]complex128, p.length)
	for i := 0; i < p.length; i++ {
		if p.na[i] {
			data[i] = cmplx.NaN()
		} else {
			data[i] = complex(float64(p.data[i]), 0)
		}
	}

	na := make([]bool, p.Len())
	copy(na, p.na)

	return data, na
}

func (p *integerPayload) Booleans() ([]bool, []bool) {
	if p.length == 0 {
		return []bool{}, []bool{}
	}

	data := make([]bool, p.length)

	for i := 0; i < p.length; i++ {
		if p.na[i] {
			data[i] = false
		} else {
			data[i] = p.data[i] != 0
		}
	}

	na := make([]bool, p.length)
	copy(na, p.na)

	return data, na
}

func (p *integerPayload) Strings() ([]string, []bool) {
	if p.length == 0 {
		return []string{}, []bool{}
	}

	data := make([]string, p.length)

	for i := 0; i < p.length; i++ {
		if p.na[i] {
			data[i] = ""
		} else {
			data[i] = strconv.Itoa(p.data[i])
		}
	}

	na := make([]bool, p.Len())
	copy(na, p.na)

	return data, na
}

func (p *integerPayload) Interfaces() ([]interface{}, []bool) {
	if p.length == 0 {
		return []interface{}{}, []bool{}
	}

	data := make([]interface{}, p.length)
	for i := 0; i < p.length; i++ {
		if p.na[i] {
			data[i] = nil
		} else {
			data[i] = p.data[i]
		}
	}

	na := make([]bool, p.length)
	copy(na, p.na)

	return data, na
}

func (p *integerPayload) Append(vec Vector) Payload {
	length := p.length + vec.Len()

	vals, na := vec.Integers()

	newVals := make([]int, length)
	newNA := make([]bool, length)

	copy(newVals, p.data)
	copy(newVals[p.length:], vals)
	copy(newNA, p.na)
	copy(newNA[p.length:], na)

	return IntegerPayload(newVals, newNA)
}

func (p *integerPayload) StrForElem(idx int) string {
	if p.na[idx-1] {
		return "NA"
	}

	return strconv.Itoa(p.data[idx-1])
}

func IntegerPayload(data []int, na []bool) Payload {
	length := len(data)

	vecNA := make([]bool, length)
	if len(na) > 0 {
		if len(na) == length {
			copy(vecNA, na)
		} else {
			emp := NAPayload(0)
			return emp
		}
	}

	vecData := make([]int, length)
	for i := 0; i < length; i++ {
		if vecNA[i] {
			vecData[i] = 0
		} else {
			vecData[i] = data[i]
		}
	}

	return &integerPayload{
		length: length,
		data:   vecData,
		DefNAble: DefNAble{
			na: vecNA,
		},
	}
}

func Integer(data []int, na []bool, options ...Config) Vector {
	config := mergeConfigs(options)

	return New(IntegerPayload(data, na), config)
}
