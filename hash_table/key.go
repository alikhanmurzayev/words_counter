package hash_table

type Key interface {
	Hash() int
	Interface() interface{}
}
