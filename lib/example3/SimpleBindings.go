package example3

type SimpleBinding map[string]SimpleTerm

func (b SimpleBinding) Merge(b2 SimpleBinding) SimpleBinding {

	result := SimpleBinding{}

	for k, v := range b {
		result[k] = v
	}

	for k, v := range b2 {
		result[k] = v
	}

	return result
}

func (b SimpleBinding) String() string {

	s, sep := "", ""

	for k, v := range b {
		s += sep + k + "=" + v.String()
		sep = ", "
	}

	return "{" + s + "}"
}