package vector

import (
	"strconv"
)

const maxIntPrint = 5

// integer is a structure, subsisting integer vectors
type integer struct {
	length int
	data []int
	na []bool
}

func (p *integer) Length() int {
	return p.length
}

func (p *integer) NA() []bool {
	return p.na
}

func (p *integer) ByIndices(indices []int) Payload {
	data := make([]int, 1, len(indices)+1)
	na := make([]bool, 1, len(indices)+1))

	for _, idx := range indices {
		data = append(data, p.data[idx])
		na = append(na, p.na[idx])
	}

	return &integer{
		length: len(data)-1,
		data:   data,
		na:     na,
	}
}

func (p *integer) SupportsFilter(selector interface{}) bool {
	if _, ok := selector.([]int); ok {
		return true
	}

	if _, ok := selector.(func (int, int, bool) bool); ok {
		return true
	}

	return false
}

func (p *integer) Filter(filter interface{}) []bool {
	if byFunc, ok := filter.(func (int, int, bool) bool); ok {
		return p.selectByFunc(byFunc)
	}

	return make([]bool, p.length)
}

func (p *integer) selectIndices(indices []int) []bool {
	booleans := make([]bool, p.length)

	for _, idx := range indices {
		if idx >= 1 && idx <= p.length {
			booleans[idx-1] = true
		}
	}

	return booleans
}

func (p *integer) selectByFunc(byFunc func (int, int, bool) bool) []bool {
	booleans := make([]bool, p.length)

	for idx, val := range p.data {
		if byFunc(idx, val, p.na[idx]) {
			booleans[idx-1] = true
		}
	}

	return booleans
}

func (p *integer) Integers() ([]int, []bool) {
	if p.length == 0 {
		return []int{}, []bool{}
	}

	data := make([]int, p.length)
	copy(data, p.data[1:])

	na := make([]bool, p.Length())
	copy(na, p.na[1:])

	return data, na
}

func (p *integer) Floats() ([]float64, []bool) {
	if p.length == 0 {
		return []float64{}, nil
	}

	data := make([]float64, p.length)

	for i := 1; i <= p.length; i++ {
		data[i-1] = float64(p.data[i])
	}

	na := make([]bool, p.Length())
	copy(na, p.na[1:])

	return data, na
}

func (p *integer) Booleans() ([]bool, []bool) {
	if p.length == 0 {
		return []bool{}, nil
	}

	data := make([]bool, p.length)

	for i := 1; i <= p.length; i++ {
		data[i-1] = p.data[i] != 0
	}

	na := make([]bool, p.Length())
	copy(na, p.na[1:])

	return data, na
}

func (p *integer) Strings() ([]string, []bool) {
	if p.length == 0 {
		return []string{}, nil
	}

	data := make([]string, p.length)

	for i := 1; i <= p.length; i++ {
		data[i-1] = strconv.Itoa(p.data[i])
	}

	na := make([]bool, p.Length())
	copy(na, p.na[1:])

	return data, na
}

func (p *integer) String() string {
	str := "["

	if p.length > 0 {
		str += p.strForElem(1)
	}
	if p.length > 1 {
		for i := 2; i <= p.length; i++ {
			if i <= maxIntPrint {
				str += ", " + p.strForElem(i)
			} else {
				str += ", ..."
				break
			}
		}
	}

	str += "]"

	return str
}

func (p *integer) strForElem(idx int) string {
	str := strconv.Itoa(p.data[idx])
	if p.DefNameable.HasNameFor(idx) {
		str += " (" + p.DefNameable.Name(idx) + ")"
	}
	return str
}

// NewIntegerPayload creates a new integer vector
func NewIntegerPayload(data []int, na []bool) Payload {
	length := len(data)

	vecData := make([]int, length+1)
	if length > 0 {
		copy(vecData[1:], data)
	}

	vecNA := make([]bool, length+1)
	if len(na) == length {
		copy(vecNA[1:], na)
	}

	payload := &integer{
		length: length,
		data:   vecData,
		na: vecNA,
	}

	return payload
}
