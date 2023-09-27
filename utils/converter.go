package utils

type Int interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

func ConvertIntSlice[U, T Int](in []T) (out []U) {
	out = make([]U, len(in))
	for i := range in {
		out[i] = U(in[i])
	}
	return
}
