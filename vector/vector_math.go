package vector

type SummerV interface {
	Sum() Vector
}

type SummerP interface {
	Sum() Payload
}

func (v *vector) Sum() Vector {
	if v.IsGrouped() {

		vectors := v.GroupVectors()
		outValues := make([]Vector, len(vectors))
		for i := 0; i < len(vectors); i++ {
			outValues[i] = vectors[i].Sum()
		}

		return Combine(outValues...)
	}

	if summer, ok := v.payload.(SummerP); ok {
		return New(summer.Sum(), v.Options()...)
	}

	return NA(1)
}
