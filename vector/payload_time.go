package vector

import (
	"time"
)

type TimeWhicherFunc = func(int, time.Time, bool) bool
type TimeWhicherCompactFunc = func(time.Time, bool) bool
type TimeToTimeApplierFunc = func(int, time.Time, bool) (time.Time, bool)
type TimeToTimeApplierCompactFunc = func(time.Time, bool) (time.Time, bool)
type TimeSummarizerFunc = func(int, time.Time, time.Time, bool) (time.Time, bool)

type TimePrinter struct {
	Format string
}

type timePayload struct {
	length  int
	data    []time.Time
	printer TimePrinter
	DefNAble
	DefArrangeable
}

func (p *timePayload) Type() string {
	return "time"
}

func (p *timePayload) Len() int {
	return p.length
}

func (p *timePayload) Pick(idx int) interface{} {
	return pickValueWithNA(idx, p.data, p.na, p.length)
}

func (p *timePayload) Data() []interface{} {
	return dataWithNAToInterfaceArray(p.data, p.na)
}

func (p *timePayload) ByIndices(indices []int) Payload {
	data := make([]time.Time, 0, len(indices))
	na := make([]bool, 0, len(indices))

	for _, idx := range indices {
		if idx == 0 {
			data = append(data, time.Time{})
			na = append(na, true)
		} else {
			data = append(data, p.data[idx-1])
			na = append(na, p.na[idx-1])
		}
	}

	return TimePayload(data, na, p.Options()...)
}

func (p *timePayload) SupportsWhicher(whicher interface{}) bool {
	if _, ok := whicher.(TimeWhicherFunc); ok {
		return true
	}

	if _, ok := whicher.(TimeWhicherCompactFunc); ok {
		return true
	}

	return false
}

func (p *timePayload) Which(whicher interface{}) []bool {
	if byFunc, ok := whicher.(TimeWhicherFunc); ok {
		return p.selectByFunc(byFunc)
	}

	if byFunc, ok := whicher.(TimeWhicherCompactFunc); ok {
		return p.selectByCompactFunc(byFunc)
	}

	return make([]bool, p.length)
}

func (p *timePayload) selectByFunc(byFunc TimeWhicherFunc) []bool {
	booleans := make([]bool, p.length)

	for idx, val := range p.data {
		if byFunc(idx+1, val, p.na[idx]) {
			booleans[idx] = true
		}
	}

	return booleans
}

func (p *timePayload) selectByCompactFunc(byFunc TimeWhicherCompactFunc) []bool {
	booleans := make([]bool, p.length)

	for idx, val := range p.data {
		if byFunc(val, p.na[idx]) {
			booleans[idx] = true
		}
	}

	return booleans
}

func (p *timePayload) SupportsApplier(applier interface{}) bool {
	if _, ok := applier.(TimeToTimeApplierFunc); ok {
		return true
	}

	if _, ok := applier.(TimeToTimeApplierCompactFunc); ok {
		return true
	}

	return false
}

func (p *timePayload) Apply(applier interface{}) Payload {
	if applyFunc, ok := applier.(TimeToTimeApplierFunc); ok {
		return p.applyToTimeByFunc(applyFunc)
	}

	if applyFunc, ok := applier.(TimeToTimeApplierCompactFunc); ok {
		return p.applyToTimeByCompactFunc(applyFunc)
	}

	return NAPayload(p.length)
}

func (p *timePayload) applyToTimeByFunc(applyFunc TimeToTimeApplierFunc) Payload {
	data, na := applyByFunc(p.data, p.na, p.length, applyFunc, time.Time{})

	return TimePayload(data, na)
}

func (p *timePayload) applyToTimeByCompactFunc(applyFunc TimeToTimeApplierCompactFunc) Payload {
	data, na := applyByCompactFunc(p.data, p.na, p.length, applyFunc, time.Time{})

	return TimePayload(data, na)
}

func (p *timePayload) ApplyTo(whicher interface{}, applier interface{}) Payload {
	//TODO implement me
	panic("implement me")
}

func (p *timePayload) SupportsSummarizer(summarizer interface{}) bool {
	if _, ok := summarizer.(TimeSummarizerFunc); ok {
		return true
	}

	return false
}

func (p *timePayload) Summarize(summarizer interface{}) Payload {
	fn, ok := summarizer.(TimeSummarizerFunc)
	if !ok {
		return NAPayload(1)
	}

	val := time.Time{}
	na := false
	for i := 0; i < p.length; i++ {
		val, na = fn(i+1, val, p.data[i], p.na[i])

		if na {
			return NAPayload(1)
		}
	}

	return TimePayload([]time.Time{val}, nil)
}

func (p *timePayload) Strings() ([]string, []bool) {
	if p.length == 0 {
		return []string{}, []bool{}
	}

	data := make([]string, p.length)

	for i := 0; i < p.length; i++ {
		if p.na[i] {
			data[i] = ""
		} else {
			data[i] = p.StrForElem(i + 1)
		}
	}

	na := make([]bool, p.Len())
	copy(na, p.na)

	return data, na
}

func (p *timePayload) Times() ([]time.Time, []bool) {
	if p.length == 0 {
		return []time.Time{}, []bool{}
	}

	data := make([]time.Time, p.length)
	copy(data, p.data)

	na := make([]bool, p.Len())
	copy(na, p.na)

	return data, na
}

func (p *timePayload) Interfaces() ([]interface{}, []bool) {
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

func (p *timePayload) Append(payload Payload) Payload {
	length := p.length + payload.Len()

	var vals []time.Time
	var na []bool

	if timeable, ok := payload.(Timeable); ok {
		vals, na = timeable.Times()
	} else {
		vals, na = NAPayload(payload.Len()).(Timeable).Times()
	}

	newVals := make([]time.Time, length)
	newNA := make([]bool, length)

	copy(newVals, p.data)
	copy(newVals[p.length:], vals)
	copy(newNA, p.na)
	copy(newNA[p.length:], na)

	return TimePayload(newVals, newNA)
}

func (p *timePayload) Adjust(size int) Payload {
	if size < p.length {
		return p.adjustToLesserSize(size)
	}

	if size > p.length {
		return p.adjustToBiggerSize(size)
	}

	return p
}

func (p *timePayload) adjustToLesserSize(size int) Payload {
	data := make([]time.Time, size)
	na := make([]bool, size)

	copy(data, p.data)
	copy(na, p.na)

	return TimePayload(data, na)
}

func (p *timePayload) adjustToBiggerSize(size int) Payload {
	cycles := size / p.length
	if size%p.length > 0 {
		cycles++
	}

	data := make([]time.Time, cycles*p.length)
	na := make([]bool, cycles*p.length)

	for i := 0; i < cycles; i++ {
		copy(data[i*p.length:], p.data)
		copy(na[i*p.length:], p.na)
	}

	data = data[:size]
	na = na[:size]

	return TimePayload(data, na)
}

func (p *timePayload) StrForElem(idx int) string {
	return p.data[idx-1].Format(p.printer.Format)
}

/* Finder interface */

func (p *timePayload) Find(needle interface{}) int {
	val, ok := needle.(time.Time)
	if !ok {
		return 0
	}

	for i, datum := range p.data {
		if !p.na[i] && val.Equal(datum) {
			return i + 1
		}
	}

	return 0
}

func (p *timePayload) FindAll(needle interface{}) []int {
	val, ok := needle.(time.Time)
	if !ok {
		return []int{}
	}

	found := []int{}
	for i, datum := range p.data {
		if !p.na[i] && val.Equal(datum) {
			found = append(found, i+1)
		}
	}

	return found
}

/* Comparable interface */

func (p *timePayload) Eq(val interface{}) []bool {
	cmp := make([]bool, p.length)

	v, ok := val.(time.Time)
	if !ok {
		return cmp
	}

	for i, datum := range p.data {
		if p.na[i] {
			cmp[i] = false
		} else {
			cmp[i] = datum.Equal(v)
		}
	}

	return cmp
}

func (p *timePayload) Neq(val interface{}) []bool {
	cmp := make([]bool, p.length)

	v, ok := val.(time.Time)
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
			cmp[i] = !datum.Equal(v)
		}
	}

	return cmp
}

func (p *timePayload) Gt(val interface{}) []bool {
	cmp := make([]bool, p.length)

	v, ok := val.(time.Time)
	if !ok {
		return cmp
	}

	for i, datum := range p.data {
		if p.na[i] {
			cmp[i] = false
		} else {
			if p.na[i] {
				cmp[i] = false
			} else {
				cmp[i] = datum.After(v)
			}
		}
	}

	return cmp
}

func (p *timePayload) Lt(val interface{}) []bool {
	cmp := make([]bool, p.length)

	v, ok := val.(time.Time)
	if !ok {
		return cmp
	}

	for i, datum := range p.data {
		if p.na[i] {
			cmp[i] = false
		} else {
			if p.na[i] {
				cmp[i] = false
			} else {
				cmp[i] = datum.Before(v)
			}
		}
	}

	return cmp
}

func (p *timePayload) Gte(val interface{}) []bool {
	cmp := make([]bool, p.length)

	v, ok := val.(time.Time)
	if !ok {
		return cmp
	}

	for i, datum := range p.data {
		if p.na[i] {
			cmp[i] = false
		} else {
			cmp[i] = datum.After(v) || datum.Equal(v)
		}
	}

	return cmp
}

func (p *timePayload) Lte(val interface{}) []bool {
	cmp := make([]bool, p.length)

	v, ok := val.(time.Time)
	if !ok {
		return cmp
	}

	for i, datum := range p.data {
		if p.na[i] {
			cmp[i] = false
		} else {
			cmp[i] = datum.Before(v) || datum.Equal(v)
		}
	}

	return cmp
}

func (p *timePayload) Groups() ([][]int, []interface{}) {
	groupMap := map[time.Time][]int{}
	ordered := []time.Time{}
	na := []int{}

	for i, val := range p.data {
		idx := i + 1

		if p.na[i] {
			na = append(na, idx)
			continue
		}

		if _, ok := groupMap[val]; !ok {
			groupMap[val] = []int{}
			ordered = append(ordered, val)
		}

		groupMap[val] = append(groupMap[val], idx)
	}

	groups := make([][]int, len(ordered))
	for i, val := range ordered {
		groups[i] = groupMap[val]
	}

	if len(na) > 0 {
		groups = append(groups, na)
	}

	values := make([]interface{}, len(groups))
	for i, val := range ordered {
		values[i] = interface{}(val)
	}
	if len(na) > 0 {
		values[len(values)-1] = nil
	}

	return groups, values
}

func (p *timePayload) IsUnique() []bool {
	booleans := make([]bool, p.length)

	valuesMap := map[string]bool{}
	wasNA := false
	for i := 0; i < p.length; i++ {
		is := false

		if p.na[i] {
			if !wasNA {
				is = true
				wasNA = true
			}
		} else {
			strTime := p.data[i].Format(p.printer.Format)

			if _, ok := valuesMap[strTime]; !ok {
				is = true
				valuesMap[strTime] = true
			}
		}

		booleans[i] = is
	}

	return booleans
}

func (p *timePayload) Options() []Option {
	return []Option{
		OptionTimeFormat(p.printer.Format),
	}
}

func TimePayload(data []time.Time, na []bool, options ...Option) Payload {
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

	vecData := make([]time.Time, length)
	for i := 0; i < length; i++ {
		if vecNA[i] {
			vecData[i] = time.Time{}
		} else {
			vecData[i] = data[i]
		}
	}

	printer := TimePrinter{Format: time.RFC3339}
	if conf.HasOption(KeyOptionTimeFormat) {
		printer.Format = conf.Value(KeyOptionTimeFormat).(string)
	}

	payload := &timePayload{
		length:  length,
		data:    vecData,
		printer: printer,
		DefNAble: DefNAble{
			na: vecNA,
		},
	}

	payload.DefArrangeable = DefArrangeable{
		Length:   payload.length,
		DefNAble: payload.DefNAble,
		FnLess: func(i, j int) bool {
			return payload.data[i].Before(payload.data[j])
		},
		FnEqual: func(i, j int) bool {
			return payload.data[i].Equal(payload.data[j])
		},
	}

	return payload
}

func (p *timePayload) Coalesce(payload Payload) Payload {
	if p.length != payload.Len() {
		payload = payload.Adjust(p.length)
	}

	var srcData []time.Time
	var srcNA []bool

	if same, ok := payload.(*timePayload); ok {
		srcData = same.data
		srcNA = same.na
	} else if timeable, ok := payload.(Timeable); ok {
		srcData, srcNA = timeable.Times()
	} else {
		return p
	}

	dstData := make([]time.Time, p.length)
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

	return TimePayload(dstData, dstNA, p.Options()...)
}

func TimeWithNA(data []time.Time, na []bool, options ...Option) Vector {
	return New(TimePayload(data, na, options...), options...)
}

func Time(data []time.Time, options ...Option) Vector {
	return TimeWithNA(data, nil, options...)
}
