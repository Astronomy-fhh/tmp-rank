package store

type Element interface {
	Compare(interface{}) int
	GetKey() int64
}

type Store[T Element] interface {
	Add(T)
	GetSort(uid int64, rangeNum int) ([]T, int64)
}
