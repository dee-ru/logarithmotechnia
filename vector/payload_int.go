package vector

import (
	"math"
	"math/cmplx"
	"strconv"
)

const maxIntPrint = 15

type IntegerWhicherFunc = func(int, int, bool) bool
type IntegerWhicherCompactFunc = func(int, bool) bool
type IntegerApplierFunc = func(int, int, bool) (int, bool)
type IntegerApplierCompactFunc = func(int, bool) (int, bool)
type IntegerSummarizerFunc = func(int, int, int, bool) (int, bool)

// integerPayload is a structure, subsisting Integer vectors
type integerPayload struct {
	length int
	data   []int
	DefNAble
	DefArrangeable
}

func (p *integerPayload) Type() string {
	return "integer"
}

func (p *integerPayload) Len() int {
	return p.length
}

func (p *integerPayload) Pick(idx int) interface{} {
	return pickValueWithNA(idx, p.data, p.na, p.length)
}

func (p *integerPayload) Data() []interface{} {
	return dataWithNAToInterfaceArray(p.data, p.na)
}

func (p *integerPayload) ByIndices(indices []int) Payload {
	data, na := byIndices(indices, p.data, p.na, 0)

	return IntegerPayload(data, na, p.Options()...)
}

func (p *integerPayload) SupportsWhicher(whicher any) bool {
	return supportsWhicher[int](whicher)
}

func (p *integerPayload) Which(whicher interface{}) []bool {
	if byFunc, ok := whicher.(IntegerWhicherFunc); ok {
		return selectByFunc(p.data, p.na, byFunc)
	}

	if byFunc, ok := whicher.(IntegerWhicherCompactFunc); ok {
		return selectByCompactFunc(p.data, p.na, byFunc)
	}

	return make([]bool, p.length)
}

func (p *integerPayload) SupportsApplier(applier any) bool {
	return supportApplier[int](applier)
}

func (p *integerPayload) Apply(applier interface{}) Payload {
	if applyFunc, ok := applier.(IntegerApplierFunc); ok {
		data, na := applyByFunc(p.data, p.na, p.length, applyFunc, 0)

		return IntegerPayload(data, na, p.Options()...)
	}

	if applyFunc, ok := applier.(IntegerApplierCompactFunc); ok {
		data, na := applyByCompactFunc(p.data, p.na, p.length, applyFunc, 0)

		return IntegerPayload(data, na, p.Options()...)
	}

	return NAPayload(p.length)

}

func (p *integerPayload) ApplyTo(indices []int, applier interface{}) Payload {
	//TODO implement me
	panic("implement me")
}

func (p *integerPayload) SupportsSummarizer(summarizer interface{}) bool {
	if _, ok := summarizer.(IntegerSummarizerFunc); ok {
		return true
	}

	return false
}

func (p *integerPayload) Summarize(summarizer interface{}) Payload {
	fn, ok := summarizer.(IntegerSummarizerFunc)
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

	return IntegerPayload([]int{val}, nil, p.Options()...)
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

func (p *integerPayload) Append(payload Payload) Payload {
	length := p.length + payload.Len()

	var vals []int
	var na []bool

	if intable, ok := payload.(Intable); ok {
		vals, na = intable.Integers()
	} else {
		vals, na = NAPayload(payload.Len()).(Intable).Integers()
	}

	newVals := make([]int, length)
	newNA := make([]bool, length)

	copy(newVals, p.data)
	copy(newVals[p.length:], vals)
	copy(newNA, p.na)
	copy(newNA[p.length:], na)

	return IntegerPayload(newVals, newNA, p.Options()...)
}

func (p *integerPayload) Adjust(size int) Payload {
	if size < p.length {
		return p.adjustToLesserSize(size)
	}

	if size > p.length {
		return p.adjustToBiggerSize(size)
	}

	return p
}

func (p *integerPayload) adjustToLesserSize(size int) Payload {
	data, na := adjustToLesserSizeWithNA(p.data, p.na, size)

	return IntegerPayload(data, na, p.Options()...)
}

func (p *integerPayload) adjustToBiggerSize(size int) Payload {
	data, na := adjustToBiggerSizeWithNA(p.data, p.na, p.length, size)

	return IntegerPayload(data, na, p.Options()...)
}

func (p *integerPayload) Groups() ([][]int, []interface{}) {
	groups, values := groupsForData(p.data, p.na)

	return groups, values
}

func (p *integerPayload) StrForElem(idx int) string {
	if p.na[idx-1] {
		return "NA"
	}

	return strconv.Itoa(p.data[idx-1])
}

func (p *integerPayload) Find(needle interface{}) int {
	val, ok := p.convertComparator(needle)
	if !ok {
		return 0
	}

	for i, datum := range p.data {
		if !p.na[i] && val == datum {
			return i + 1
		}
	}

	return 0
}

/* Finder interface */

func (p *integerPayload) FindAll(needle interface{}) []int {
	val, ok := p.convertComparator(needle)
	if !ok {
		return []int{}
	}

	found := []int{}
	for i, datum := range p.data {
		if !p.na[i] && val == datum {
			found = append(found, i+1)
		}
	}

	return found
}

func (p *integerPayload) Eq(val interface{}) []bool {
	cmp := make([]bool, p.length)

	v, ok := p.convertComparator(val)
	if !ok {
		return cmp
	}

	for i, datum := range p.data {
		if p.na[i] {
			cmp[i] = false
		} else {
			cmp[i] = datum == v
		}
	}

	return cmp
}

/* Comparable interface */

func (p *integerPayload) Neq(val interface{}) []bool {
	cmp := make([]bool, p.length)

	v, ok := p.convertComparator(val)
	if !ok {
		for i := range p.data {
			cmp[i] = true
		}

		return cmp
	}

	for i, datum := range p.data {
		if p.na[i] {
			cmp[i] = true
		} else {
			cmp[i] = datum != v
		}
	}

	return cmp
}

func (p *integerPayload) Gt(val interface{}) []bool {
	cmp := make([]bool, p.length)

	v, ok := p.convertComparator(val)
	if !ok {
		return cmp
	}

	for i, datum := range p.data {
		if p.na[i] {
			cmp[i] = false
		} else {
			cmp[i] = datum > v
		}
	}

	return cmp
}

func (p *integerPayload) Lt(val interface{}) []bool {
	cmp := make([]bool, p.length)

	v, ok := p.convertComparator(val)
	if !ok {
		return cmp
	}

	for i, datum := range p.data {
		if p.na[i] {
			cmp[i] = false
		} else {
			cmp[i] = datum < v
		}
	}

	return cmp
}

func (p *integerPayload) Gte(val interface{}) []bool {
	cmp := make([]bool, p.length)

	v, ok := p.convertComparator(val)
	if !ok {
		return cmp
	}

	for i, datum := range p.data {
		if p.na[i] {
			cmp[i] = false
		} else {
			cmp[i] = datum >= v
		}
	}

	return cmp
}

func (p *integerPayload) Lte(val interface{}) []bool {
	cmp := make([]bool, p.length)

	v, ok := p.convertComparator(val)
	if !ok {
		return cmp
	}

	for i, datum := range p.data {
		if p.na[i] {
			cmp[i] = false
		} else {
			cmp[i] = datum <= v
		}
	}

	return cmp
}

func (p *integerPayload) convertComparator(val interface{}) (int, bool) {
	var v int
	ok := true
	switch val.(type) {
	case int:
		v = val.(int)
	case int64:
		v = int(val.(int64))
	case int32:
		v = int(val.(int32))
	case uint64:
		v = int(val.(uint64))
	case uint32:
		v = int(val.(uint32))
	case complex128:
		ip := imag(val.(complex128))
		rp, fp := math.Modf(real(val.(complex128)))
		if ip == 0 && fp == 0 {
			v = int(rp)
		} else {
			ok = false
		}
	case complex64:
		ip := imag(val.(complex64))
		rp, fp := math.Modf(float64(real(val.(complex64))))
		if ip == 0 && fp == 0 {
			v = int(rp)
		} else {
			ok = false
		}
	case float64:
		rp, fp := math.Modf(val.(float64))
		if fp == 0 {
			v = int(rp)
		} else {
			ok = false
		}
	case float32:
		rp, fp := math.Modf(float64(val.(float32)))
		if fp == 0 {
			v = int(rp)
		} else {
			ok = false
		}
	default:
		ok = false
	}

	return v, ok
}

func (p *integerPayload) IsUnique() []bool {
	booleans := make([]bool, p.length)

	valuesMap := map[int]bool{}
	wasNA := false
	for i := 0; i < p.length; i++ {
		is := false

		if p.na[i] {
			if !wasNA {
				is = true
				wasNA = true
			}
		} else {
			if _, ok := valuesMap[p.data[i]]; !ok {
				is = true
				valuesMap[p.data[i]] = true
			}
		}

		booleans[i] = is
	}

	return booleans
}

func (p *integerPayload) Coalesce(payload Payload) Payload {
	if p.length != payload.Len() {
		payload = payload.Adjust(p.length)
	}

	var srcData []int
	var srcNA []bool

	if same, ok := payload.(*integerPayload); ok {
		srcData = same.data
		srcNA = same.na
	} else if intable, ok := payload.(Intable); ok {
		srcData, srcNA = intable.Integers()
	} else {
		return p
	}

	dstData := make([]int, p.length)
	dstNA := make([]bool, p.length)

	for i := 0; i < p.length; i++ {
		if p.na[i] && !srcNA[i] {
			dstData[i] = srcData[i]
			dstNA[i] = false
		} else {
			dstData[i] = p.data[i]
			dstNA[i] = p.na[i]
		}
	}

	return IntegerPayload(dstData, dstNA, p.Options()...)
}

func (p *integerPayload) Options() []Option {
	return []Option{}
}

func IntegerPayload(data []int, na []bool, _ ...Option) Payload {
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

	payload := &integerPayload{
		length: length,
		data:   vecData,
		DefNAble: DefNAble{
			na: vecNA,
		},
	}

	payload.DefArrangeable = DefArrangeable{
		Length:   payload.length,
		DefNAble: payload.DefNAble,
		FnLess: func(i, j int) bool {
			return payload.data[i] < payload.data[j]
		},
		FnEqual: func(i, j int) bool {
			return payload.data[i] == payload.data[j]
		},
	}

	return payload
}

func IntegerWithNA(data []int, na []bool, options ...Option) Vector {
	return New(IntegerPayload(data, na, options...), options...)
}

func Integer(data []int, options ...Option) Vector {
	return IntegerWithNA(data, nil, options...)
}
