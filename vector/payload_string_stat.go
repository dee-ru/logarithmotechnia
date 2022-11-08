package vector

func (p *stringPayload) Max() Payload {
	max, na := genMax(p.data, p.na)

	return StringPayload([]string{max}, []bool{na}, p.Options()...)
}

func (p *stringPayload) Min() Payload {
	min, na := genMin(p.data, p.na)

	return StringPayload([]string{min}, []bool{na}, p.Options()...)
}
