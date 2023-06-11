package vector

import (
	"logarithmotechnia/embed"
	"math"
	"math/cmplx"
	"strconv"
)

type FloatPrinter struct {
	Precision int
}

type floatPayload struct {
	length  int
	data    []float64
	printer FloatPrinter
	embed.DefNAble
	DefArrangeable
}

func (p *floatPayload) Type() string {
	return "float"
}

func (p *floatPayload) Len() int {
	return p.length
}

func (p *floatPayload) Pick(idx int) any {
	return pickValueWithNA(idx, p.data, p.NA, p.length)
}

func (p *floatPayload) Data() []any {
	return dataWithNAToInterfaceArray(p.data, p.NA)
}

func (p *floatPayload) ByIndices(indices []int) Payload {
	data, na := byIndicesWithNA(indices, p.data, p.NA, math.NaN())

	return FloatPayload(data, na, p.Options()...)
}

func (p *floatPayload) SupportsWhicher(whicher any) bool {
	return supportsWhicherWithNA[float64](whicher)
}

func (p *floatPayload) Which(whicher any) []bool {
	return whichWithNA(p.data, p.NA, whicher)
}

func (p *floatPayload) Apply(applier any) Payload {
	return applyWithNA(p.data, p.NA, applier, p.Options())
}

func (p *floatPayload) Traverse(traverser any) {
	traverseWithNA(p.data, p.NA, traverser)
}

func (p *floatPayload) ApplyTo(indices []int, applier any) Payload {
	data, na := applyToWithNA(indices, p.data, p.NA, applier, math.NaN())

	if data == nil {
		return NAPayload(p.length)
	}

	return FloatPayload(data, na, p.Options()...)
}

func (p *floatPayload) SupportsSummarizer(summarizer any) bool {
	return supportsSummarizer[float64](summarizer)
}

func (p *floatPayload) Summarize(summarizer any) Payload {
	val, na := summarize(p.data, p.NA, summarizer, 0.0, math.NaN())

	return FloatPayload([]float64{val}, []bool{na}, p.Options()...)
}

func (p *floatPayload) Integers() ([]int, []bool) {
	if p.length == 0 {
		return []int{}, []bool{}
	}

	data := make([]int, p.length)
	na := make([]bool, p.Len())
	copy(na, p.NA)

	for i := 0; i < p.length; i++ {
		if p.NA[i] {
			data[i] = 0
		} else if math.IsNaN(p.data[i]) || math.IsInf(p.data[i], 1) || math.IsInf(p.data[i], -1) {
			data[i] = 0
			na[i] = true
		} else {
			data[i] = int(p.data[i])
		}
	}

	return data, na
}

func (p *floatPayload) Floats() ([]float64, []bool) {
	if p.length == 0 {
		return []float64{}, []bool{}
	}

	data := make([]float64, p.length)
	copy(data, p.data)

	na := make([]bool, p.Len())
	copy(na, p.NA)

	return data, na
}

func (p *floatPayload) Complexes() ([]complex128, []bool) {
	if p.length == 0 {
		return []complex128{}, []bool{}
	}

	data := make([]complex128, p.length)
	for i := 0; i < p.length; i++ {
		if p.NA[i] {
			data[i] = cmplx.NaN()
		} else {
			data[i] = complex(p.data[i], 0)
		}
	}

	na := make([]bool, p.Len())
	copy(na, p.NA)

	return data, na
}

func (p *floatPayload) Booleans() ([]bool, []bool) {
	if p.length == 0 {
		return []bool{}, []bool{}
	}

	data := make([]bool, p.length)

	for i := 0; i < p.length; i++ {
		if p.NA[i] {
			data[i] = false
		} else {
			data[i] = p.data[i] != 0
		}
	}

	na := make([]bool, p.length)
	copy(na, p.NA)

	return data, na
}

func (p *floatPayload) Strings() ([]string, []bool) {
	if p.length == 0 {
		return []string{}, []bool{}
	}

	data := make([]string, p.length)

	for i := 0; i < p.length; i++ {
		data[i] = p.StrForElem(i + 1)
	}

	na := make([]bool, p.Len())
	copy(na, p.NA)

	return data, na
}

func (p *floatPayload) Anies() ([]any, []bool) {
	if p.length == 0 {
		return []any{}, []bool{}
	}

	data := make([]any, p.length)
	for i := 0; i < p.length; i++ {
		if p.NA[i] {
			data[i] = nil
		} else {
			data[i] = p.data[i]
		}
	}

	na := make([]bool, p.length)
	copy(na, p.NA)

	return data, na
}

func (p *floatPayload) StrForElem(idx int) string {
	i := idx - 1

	if p.NA[i] {
		return "NA"
	}

	if math.IsInf(p.data[i], +1) {
		return "+Inf"
	}

	if math.IsInf(p.data[i], -1) {
		return "-Inf"
	}

	if math.IsNaN(p.data[i]) {
		return "NaN"
	}

	return strconv.FormatFloat(p.data[i], 'f', p.printer.Precision, 64)
}

func (p *floatPayload) Append(payload Payload) Payload {
	length := p.length + payload.Len()

	var vals []float64
	var na []bool

	if floatable, ok := payload.(Floatable); ok {
		vals, na = floatable.Floats()
	} else {
		vals, na = NAPayload(payload.Len()).(Floatable).Floats()
	}

	newVals := make([]float64, length)
	newNA := make([]bool, length)

	copy(newVals, p.data)
	copy(newVals[p.length:], vals)
	copy(newNA, p.NA)
	copy(newNA[p.length:], na)

	return FloatPayload(newVals, newNA, p.Options()...)
}

func (p *floatPayload) Adjust(size int) Payload {
	if size < p.length {
		return p.adjustToLesserSize(size)
	}

	if size > p.length {
		return p.adjustToBiggerSize(size)
	}

	return p
}

func (p *floatPayload) adjustToLesserSize(size int) Payload {
	data, na := adjustToLesserSizeWithNA(p.data, p.NA, size)

	return FloatPayload(data, na, p.Options()...)
}

func (p *floatPayload) adjustToBiggerSize(size int) Payload {
	data, na := adjustToBiggerSizeWithNA(p.data, p.NA, p.length, size)

	return FloatPayload(data, na, p.Options()...)
}

/* Finder interface */

func (p *floatPayload) Find(needle any) int {
	return find(needle, p.data, p.NA, p.convertComparator)
}

func (p *floatPayload) FindAll(needle any) []int {
	return findAll(needle, p.data, p.NA, p.convertComparator)
}

func (p *floatPayload) Eq(val any) []bool {
	return eq(val, p.data, p.NA, p.convertComparator)
}

func (p *floatPayload) Neq(val any) []bool {
	return neq(val, p.data, p.NA, p.convertComparator)
}

func (p *floatPayload) Gt(val any) []bool {
	return gt(val, p.data, p.NA, p.convertComparator)
}

func (p *floatPayload) Lt(val any) []bool {
	return lt(val, p.data, p.NA, p.convertComparator)
}

func (p *floatPayload) Gte(val any) []bool {
	return gte(val, p.data, p.NA, p.convertComparator)
}

func (p *floatPayload) Lte(val any) []bool {
	return lte(val, p.data, p.NA, p.convertComparator)
}

func (p *floatPayload) convertComparator(val any) (float64, bool) {
	var v float64
	ok := true
	switch val.(type) {
	case complex128:
		ip := imag(val.(complex128))
		if ip == 0 {
			v = real(val.(complex128))
		} else {
			ok = false
		}
	case complex64:
		ip := imag(val.(complex64))
		if ip == 0 {
			v = float64(real(val.(complex64)))
		} else {
			ok = false
		}
	case float64:
		v = val.(float64)
	case float32:
		v = float64(val.(float32))
	case int:
		v = float64(val.(int))
	case int64:
		v = float64(val.(int64))
	case int32:
		v = float64(val.(int32))
	case uint64:
		v = float64(val.(uint64))
	case uint32:
		v = float64(val.(uint32))
	default:
		ok = false
	}

	return v, ok
}

func (p *floatPayload) Groups() ([][]int, []any) {
	groups, values := groupsForData(p.data, p.NA)

	return groups, values
}

func (p *floatPayload) IsUnique() []bool {
	booleans := make([]bool, p.length)

	valuesMap := map[float64]bool{}
	wasNA := false
	wasNaN := false
	wasInfPlus := false
	wasInfMinus := false
	for i := 0; i < p.length; i++ {
		is := false

		if p.NA[i] {
			if !wasNA {
				is = true
				wasNA = true
			}
		} else if math.IsNaN(p.data[i]) {
			if !wasNaN {
				is = true
				wasNaN = true
			}
		} else if math.IsInf(p.data[i], 1) {
			if !wasInfPlus {
				is = true
				wasInfPlus = true
			}
		} else if math.IsInf(p.data[i], -1) {
			if !wasInfMinus {
				is = true
				wasInfMinus = true
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

func (p *floatPayload) Coalesce(payload Payload) Payload {
	if p.length != payload.Len() {
		payload = payload.Adjust(p.length)
	}

	var srcData []float64
	var srcNA []bool

	if same, ok := payload.(*floatPayload); ok {
		srcData = same.data
		srcNA = same.NA
	} else if floatable, ok := payload.(Floatable); ok {
		srcData, srcNA = floatable.Floats()
	} else {
		return p
	}

	dstData := make([]float64, p.length)
	dstNA := make([]bool, p.length)

	for i := 0; i < p.length; i++ {
		if p.NA[i] && !srcNA[i] {
			dstData[i] = srcData[i]
			dstNA[i] = false
		} else {
			dstData[i] = p.data[i]
			dstNA[i] = p.NA[i]
		}
	}

	return FloatPayload(dstData, dstNA, p.Options()...)
}

func (p *floatPayload) Options() []Option {
	return []Option{
		ConfOption{keyOptionPrecision, p.printer.Precision},
	}
}

func (p *floatPayload) SetOption(name string, val any) bool {
	switch name {
	case keyOptionPrecision:
		p.printer.Precision = val.(int)
	default:
		return false
	}

	return true
}

// FloatPayload creates a payload with float data.
//
// Available options are:
//   - OptionPrecision(precision int) - sets precision for printing payload's values.
func FloatPayload(data []float64, na []bool, options ...Option) Payload {
	length := len(data)
	conf := MergeOptions(options)

	vecNA := make([]bool, length)
	if len(na) > 0 {
		if len(na) == length {
			copy(vecNA, na)
		} else {
			emp := NAPayload(0)
			return emp
		}
	}

	vecData := make([]float64, length)
	for i := 0; i < length; i++ {
		if vecNA[i] {
			vecData[i] = math.NaN()
		} else {
			vecData[i] = data[i]
		}
	}

	printer := FloatPrinter{
		Precision: 3,
	}

	payload := &floatPayload{
		length:  length,
		data:    vecData,
		printer: printer,
		DefNAble: embed.DefNAble{
			NA: vecNA,
		},
	}

	conf.SetOptions(payload)

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

// FloatWithNA creates a vector with FloatPayload and allows to set NA-values.
func FloatWithNA(data []float64, na []bool, options ...Option) Vector {
	return New(FloatPayload(data, na, options...), options...)
}

// Float creates a vector with FloatPayload.
func Float(data []float64, options ...Option) Vector {
	return FloatWithNA(data, nil, options...)
}
