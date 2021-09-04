package vector

import (
	"math"
	"math/cmplx"
)

type BooleanWhicherFunc = func(int, bool, bool) bool
type BooleanWhicherCompactFunc = func(bool, bool) bool
type BooleanToBooleanApplierFunc = func(int, bool, bool) (bool, bool)
type BooleanToBooleanApplierCompactFunc = func(bool, bool) (bool, bool)
type BooleanSummarizerFunc = func(int, bool, bool, bool) (bool, bool)

type booleanPayload struct {
	length int
	data   []bool
	DefNAble
}

func (p *booleanPayload) Type() string {
	return "boolean"
}

func (p *booleanPayload) Len() int {
	return p.length
}

func (p *booleanPayload) ByIndices(indices []int) Payload {
	data := make([]bool, 0, len(indices))
	na := make([]bool, 0, len(indices))

	for _, idx := range indices {
		data = append(data, p.data[idx-1])
		na = append(na, p.na[idx-1])
	}

	return &booleanPayload{
		length: len(data),
		data:   data,
		DefNAble: DefNAble{
			na: na,
		},
	}
}

func (p *booleanPayload) SupportsWhicher(whicher interface{}) bool {
	if _, ok := whicher.(BooleanWhicherFunc); ok {
		return true
	}

	if _, ok := whicher.(BooleanWhicherCompactFunc); ok {
		return true
	}

	return false
}

func (p *booleanPayload) Which(whicher interface{}) []bool {
	if byFunc, ok := whicher.(BooleanWhicherFunc); ok {
		return p.selectByFunc(byFunc)
	}

	if byFunc, ok := whicher.(BooleanWhicherCompactFunc); ok {
		return p.selectByCompactFunc(byFunc)
	}

	return make([]bool, p.length)
}

func (p *booleanPayload) selectByFunc(byFunc BooleanWhicherFunc) []bool {
	booleans := make([]bool, p.length)

	for idx, val := range p.data {
		if byFunc(idx+1, val, p.na[idx]) {
			booleans[idx] = true
		}
	}

	return booleans
}

func (p *booleanPayload) selectByCompactFunc(byFunc BooleanWhicherCompactFunc) []bool {
	booleans := make([]bool, p.length)

	for idx, val := range p.data {
		if byFunc(val, p.na[idx]) {
			booleans[idx] = true
		}
	}

	return booleans
}

func (p *booleanPayload) SupportsApplier(applier interface{}) bool {
	if _, ok := applier.(BooleanToBooleanApplierFunc); ok {
		return true
	}

	if _, ok := applier.(BooleanToBooleanApplierCompactFunc); ok {
		return true
	}

	return false
}

func (p *booleanPayload) Apply(applier interface{}) Payload {
	if applyFunc, ok := applier.(BooleanToBooleanApplierFunc); ok {
		return p.applyToBooleanByFunc(applyFunc)
	}

	if applyFunc, ok := applier.(BooleanToBooleanApplierCompactFunc); ok {
		return p.applyToBooleanByCompactFunc(applyFunc)
	}

	return NAPayload(p.length)
}

func (p *booleanPayload) applyToBooleanByFunc(applyFunc BooleanToBooleanApplierFunc) Payload {
	data := make([]bool, p.length)
	na := make([]bool, p.length)

	for i := 0; i < p.length; i++ {
		dataVal, naVal := applyFunc(i+1, p.data[i], p.na[i])
		if naVal {
			dataVal = false
		}
		data[i] = dataVal
		na[i] = naVal
	}

	return BooleanPayload(data, na)
}

func (p *booleanPayload) applyToBooleanByCompactFunc(applyFunc BooleanToBooleanApplierCompactFunc) Payload {
	data := make([]bool, p.length)
	na := make([]bool, p.length)

	for i := 0; i < p.length; i++ {
		dataVal, naVal := applyFunc(p.data[i], p.na[i])
		if naVal {
			dataVal = false
		}
		data[i] = dataVal
		na[i] = naVal
	}

	return BooleanPayload(data, na)
}

func (p *booleanPayload) SupportsSummarizer(summarizer interface{}) bool {
	if _, ok := summarizer.(BooleanSummarizerFunc); ok {
		return true
	}

	return false
}

func (p *booleanPayload) Summarize(summarizer interface{}) Payload {
	fn, ok := summarizer.(BooleanSummarizerFunc)
	if !ok {
		return NAPayload(1)
	}

	val := false
	na := false
	for i := 0; i < p.length; i++ {
		val, na = fn(i+1, val, p.data[i], p.na[i])
		if na {
			return NAPayload(1)
		}
	}

	return BooleanPayload([]bool{val}, nil)
}

func (p *booleanPayload) Integers() ([]int, []bool) {
	if p.length == 0 {
		return []int{}, []bool{}
	}

	data := make([]int, p.length)
	for i := 0; i < p.length; i++ {
		if p.na[i] {
			data[i] = 0
		} else {
			if p.data[i] {
				data[i] = 1
			} else {
				data[i] = 0
			}
		}
	}

	na := make([]bool, p.Len())
	copy(na, p.na)

	return data, na
}

func (p *booleanPayload) Floats() ([]float64, []bool) {
	if p.length == 0 {
		return []float64{}, []bool{}
	}

	data := make([]float64, p.length)

	for i := 0; i < p.length; i++ {
		if p.na[i] {
			data[i] = math.NaN()
		} else {
			if p.data[i] {
				data[i] = 1
			} else {
				data[i] = 0
			}
		}
	}

	na := make([]bool, p.Len())
	copy(na, p.na)

	return data, na
}

func (p *booleanPayload) Complexes() ([]complex128, []bool) {
	if p.length == 0 {
		return []complex128{}, []bool{}
	}

	data := make([]complex128, p.length)
	for i := 0; i < p.length; i++ {
		if p.na[i] {
			data[i] = cmplx.NaN()
		} else {
			if p.data[i] {
				data[i] = 1
			} else {
				data[i] = 0
			}
		}
	}

	na := make([]bool, p.Len())
	copy(na, p.na)

	return data, na
}

func (p *booleanPayload) Booleans() ([]bool, []bool) {
	if p.length == 0 {
		return []bool{}, []bool{}
	}

	data := make([]bool, p.length)
	copy(data, p.data)

	na := make([]bool, p.length)
	copy(na, p.na)

	return data, na
}

func (p *booleanPayload) Strings() ([]string, []bool) {
	if p.length == 0 {
		return []string{}, []bool{}
	}

	data := make([]string, p.length)

	for i := 0; i < p.length; i++ {
		data[i] = p.StrForElem(i + 1)
	}

	na := make([]bool, p.Len())
	copy(na, p.na)

	return data, na
}

func (p *booleanPayload) Interfaces() ([]interface{}, []bool) {
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

	na := make([]bool, p.Len())
	copy(na, p.na)

	return data, na
}

func (p *booleanPayload) Append(vec Vector) Payload {
	length := p.length + vec.Len()

	vals, na := vec.Booleans()

	newVals := make([]bool, length)
	newNA := make([]bool, length)

	copy(newVals, p.data)
	copy(newVals[p.length:], vals)
	copy(newNA, p.na)
	copy(newNA[p.length:], na)

	return BooleanPayload(newVals, newNA)
}

func (p *booleanPayload) Adjust(size int) Payload {
	if size < p.length {
		return p.adjustToLesserSize(size)
	}

	if size > p.length {
		return p.adjustToBiggerSize(size)
	}

	return p
}

func (p *booleanPayload) adjustToLesserSize(size int) Payload {
	data := make([]bool, size)
	na := make([]bool, size)

	copy(data, p.data)
	copy(na, p.na)

	return BooleanPayload(data, na)
}

func (p *booleanPayload) adjustToBiggerSize(size int) Payload {
	cycles := size / p.length
	if size%p.length > 0 {
		cycles++
	}

	data := make([]bool, cycles*p.length)
	na := make([]bool, cycles*p.length)

	for i := 0; i < cycles; i++ {
		copy(data[i*p.length:], p.data)
		copy(na[i*p.length:], p.na)
	}

	data = data[:size]
	na = na[:size]

	return BooleanPayload(data, na)
}

func (p *booleanPayload) StrForElem(idx int) string {
	if p.na[idx-1] {
		return "NA"
	}

	if p.data[idx-1] {
		return "true"
	}

	return "false"
}

/* Finder interface */

func (p *booleanPayload) Find(needle interface{}) int {
	val, ok := needle.(bool)
	if !ok {
		return 0
	}

	for i, datum := range p.data {
		if val == datum {
			return i + 1
		}
	}

	return 0
}

func (p *booleanPayload) FindAll(needle interface{}) []int {
	val, ok := needle.(bool)
	if !ok {
		return []int{}
	}

	found := []int{}
	for i, datum := range p.data {
		if val == datum {
			found = append(found, i+1)
		}
	}

	return found
}

func BooleanPayload(data []bool, na []bool) Payload {
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

	vecData := make([]bool, length)
	for i := 0; i < length; i++ {
		if vecNA[i] {
			vecData[i] = false
		} else {
			vecData[i] = data[i]
		}
	}

	return &booleanPayload{
		length: length,
		data:   vecData,
		DefNAble: DefNAble{
			na: vecNA,
		},
	}
}

func Boolean(data []bool, na []bool) Vector {
	return New(BooleanPayload(data, na))
}
