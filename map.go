package stream

type Entry[K comparable, V any] struct {
	key   K
	value V
}

func (en Entry[K, V]) Key() K {
	return en.key
}

func (en Entry[K, V]) Value() V {
	return en.value
}

type EntrySet[K comparable, V any] []Entry[K, V]

func (set EntrySet[K, V]) Keys() []K {
	keys := make([]K, 0, len(set))
	for _, v := range set {
		keys = append(keys, v.Key())
	}
	return keys
}

func (set EntrySet[K, V]) Values() []V {
	values := make([]V, 0, len(set))
	for _, v := range set {
		values = append(values, v.Value())
	}
	return values
}

func (set EntrySet[K, V]) List() []Entry[K, V] {
	return set
}

func ConvertToEntrySet[S ~[]Entry[K, V], K comparable, V any](s S) EntrySet[K, V] {
	return EntrySet[K, V](s)
}

func GetEntrySet[M ~map[K]V, S EntrySet[K, V], K comparable, V any](m M) S {
	set := make(S, 0, len(m))
	for k, v := range m {
		entry := Entry[K, V]{
			key:   k,
			value: v,
		}
		set = append(set, entry)
	}
	return set
}


