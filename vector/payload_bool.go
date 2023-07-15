package vector

import (
	"logarithmotechnia/embed"
	"logarithmotechnia/option"
	"math"
	"math/cmplx"
)

type booleanPayload struct {
	length int
	data   []bool
	embed.NAble
	embed.Arrangeable
}

func (p *booleanPayload) Type() string {
	return "boolean"
}

func (p *booleanPayload) Len() int {
	return p.length
}

func (p *booleanPayload) Pick(idx int) any {
	return PickValueWithNA(idx, p.data, p.NA, p.length)
}

func (p *booleanPayload) Data() []any {
	return DataWithNAToInterfaceArray(p.data, p.NA)
}

func (p *booleanPayload) ByIndices(indices []int) Payload {
	data, na := ByIndicesWithNA(indices, p.data, p.NA, false)

	return BooleanPayload(data, na, p.Options()...)
}

func (p *booleanPayload) SupportsWhicher(whicher any) bool {
	return SupportsWhicherWithNA[bool](whicher)
}

func (p *booleanPayload) Which(whicher any) []bool {
	return WhichWithNA(p.data, p.NA, whicher)
}

func (p *booleanPayload) Apply(applier any) Payload {
	return ApplyWithNA(p.data, p.NA, applier, p.Options())
}

func (p *booleanPayload) ApplyTo(indices []int, applier any) Payload {
	data, na := ApplyToWithNA(indices, p.data, p.NA, applier, false)

	if data == nil {
		return NAPayload(p.length)
	}

	return BooleanPayload(data, na, p.Options()...)
}

func (p *booleanPayload) Traverse(traverser any) {
	TraverseWithNA(p.data, p.NA, traverser)
}

func (p *booleanPayload) SupportsSummarizer(summarizer any) bool {
	return SupportsSummarizer[bool](summarizer)
}

func (p *booleanPayload) Summarize(summarizer any) Payload {
	val, na := Summarize(p.data, p.NA, summarizer, false, false)

	return BooleanPayload([]bool{val}, []bool{na}, p.Options()...)
}

func (p *booleanPayload) Integers() ([]int, []bool) {
	if p.length == 0 {
		return []int{}, []bool{}
	}

	data := make([]int, p.length)
	for i := 0; i < p.length; i++ {
		if p.NA[i] {
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
	copy(na, p.NA)

	return data, na
}

func (p *booleanPayload) Floats() ([]float64, []bool) {
	if p.length == 0 {
		return []float64{}, []bool{}
	}

	data := make([]float64, p.length)

	for i := 0; i < p.length; i++ {
		if p.NA[i] {
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
	copy(na, p.NA)

	return data, na
}

func (p *booleanPayload) Complexes() ([]complex128, []bool) {
	if p.length == 0 {
		return []complex128{}, []bool{}
	}

	data := make([]complex128, p.length)
	for i := 0; i < p.length; i++ {
		if p.NA[i] {
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
	copy(na, p.NA)

	return data, na
}

func (p *booleanPayload) Booleans() ([]bool, []bool) {
	if p.length == 0 {
		return []bool{}, []bool{}
	}

	data := make([]bool, p.length)
	copy(data, p.data)

	na := make([]bool, p.length)
	copy(na, p.NA)

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
	copy(na, p.NA)

	return data, na
}

func (p *booleanPayload) Anies() ([]any, []bool) {
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

	na := make([]bool, p.Len())
	copy(na, p.NA)

	return data, na
}

func (p *booleanPayload) Append(payload Payload) Payload {
	length := p.length + payload.Len()

	var vals []bool
	var na []bool

	if boolable, ok := payload.(Boolable); ok {
		vals, na = boolable.Booleans()
	} else {
		vals, na = NAPayload(payload.Len()).(Boolable).Booleans()
	}

	newVals := make([]bool, length)
	newNA := make([]bool, length)

	copy(newVals, p.data)
	copy(newVals[p.length:], vals)
	copy(newNA, p.NA)
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
	data, na := AdjustToLesserSizeWithNA(p.data, p.NA, size)

	return BooleanPayload(data, na)
}

func (p *booleanPayload) adjustToBiggerSize(size int) Payload {
	data, na := AdjustToBiggerSizeWithNA(p.data, p.NA, p.length, size)

	return BooleanPayload(data, na)
}

func (p *booleanPayload) Groups() ([][]int, []any) {
	groups, values := GroupsForData(p.data, p.NA)

	return groups, values
}

func (p *booleanPayload) StrForElem(idx int) string {
	if p.NA[idx-1] {
		return "NA"
	}

	if p.data[idx-1] {
		return "true"
	}

	return "false"
}

func (p *booleanPayload) Find(needle any) int {
	return Find(needle, p.data, p.NA, p.convertComparator)
}

func (p *booleanPayload) FindAll(needle any) []int {
	return FindAll(needle, p.data, p.NA, p.convertComparator)
}

func (p *booleanPayload) Eq(val any) []bool {
	return Eq(val, p.data, p.NA, p.convertComparator)
}

func (p *booleanPayload) Neq(val any) []bool {
	return Neq(val, p.data, p.NA, p.convertComparator)
}

func (p *booleanPayload) convertComparator(val any) (bool, bool) {
	v, ok := val.(bool)

	return v, ok
}

func (p *booleanPayload) IsUnique() []bool {
	booleans := make([]bool, p.length)

	valuesMap := map[bool]bool{}
	wasNA := false
	for i := 0; i < p.length; i++ {
		is := false

		if p.NA[i] {
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

func (p *booleanPayload) Options() []option.Option {
	return []option.Option{}
}

func (p *booleanPayload) SetOption(string, any) bool {
	return false
}

func (p *booleanPayload) Coalesce(payload Payload) Payload {
	if p.length != payload.Len() {
		payload = payload.Adjust(p.length)
	}

	var srcData []bool
	var srcNA []bool

	if same, ok := payload.(*booleanPayload); ok {
		srcData = same.data
		srcNA = same.NA
	} else if boolable, ok := payload.(Boolable); ok {
		srcData, srcNA = boolable.Booleans()
	} else {
		return p
	}

	dstData := make([]bool, p.length)
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

	return BooleanPayload(dstData, dstNA, p.Options()...)
}

// BooleanPayload creates a payload with boolean data.
func BooleanPayload(data []bool, na []bool, _ ...option.Option) Payload {
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

	payload := &booleanPayload{
		length: length,
		data:   vecData,
		NAble: embed.NAble{
			NA: vecNA,
		},
	}

	payload.Arrangeable = embed.Arrangeable{
		Length: payload.length,
		NAble:  payload.NAble,
		FnLess: func(i, j int) bool {
			return !payload.data[i] && payload.data[j]
		},
		FnEqual: func(i, j int) bool {
			return payload.data[i] == payload.data[j]
		},
	}

	return payload
}

// BooleanWithNA creates a vector with BooleanPayload and allows to set NA-values.
func BooleanWithNA(data []bool, na []bool, options ...option.Option) Vector {
	return New(BooleanPayload(data, na, options...), options...)
}

// Boolean creates a vector with BooleanPayload.
func Boolean(data []bool, options ...option.Option) Vector {
	return BooleanWithNA(data, nil, options...)
}
