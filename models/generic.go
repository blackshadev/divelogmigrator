package models

type Optional[T any] struct {
	IsFilled bool
	Value    T
}

func NewOptionalValue[T any](value *T) Optional[T] {
	if value == nil {
		return Optional[T]{
			IsFilled: false,
		}
	}

	return Optional[T]{
		IsFilled: true,
		Value:    *value,
	}
}
