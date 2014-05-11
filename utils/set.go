package utils

type Set map[interface{}]struct{}

func (set Set) Insert(item interface{}) {
	set[item] = struct{}{}
}

func (set Set) Remove(item interface{}) {
	delete(set, item)
}

func (set Set) ToSlice() []interface{} {
	slice := make([]interface{}, len(set))
	i := 0
	for val, _ := range set {
		slice[i] = val
		i++
	}
	return slice
}
