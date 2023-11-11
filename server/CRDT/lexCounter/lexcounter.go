package lexCounter

import "fmt"

type Ordered interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64 | ~string
}

type Number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64
}

type Pair[V Number] struct {
	first  int
	second V
}

type LexCounter[K Ordered, V Number] struct {
	m  map[K]Pair[V]
	id K
}

func Create[K Ordered, V Number](id K) LexCounter[K, V] {
	m := make(map[K]Pair[V])
	lc := LexCounter[K, V]{m: m, id: id}
	return lc
}

func (lexcounter LexCounter[K, V]) Inc(toSum V) {
	m := lexcounter.m
	id := lexcounter.id
	m[id] = Pair[V]{first: m[id].first + 1, second: m[id].second + toSum}
}

func (lexcounter LexCounter[K, V]) Dec(toDec V) {
	m := lexcounter.m
	id := lexcounter.id
	if m[id].second <= toDec {
		m[id] = Pair[V]{first: m[id].first + 1, second: 0}
	} else {
		m[id] = Pair[V]{first: m[id].first + 1, second: m[id].second - toDec}
	}
}

func (lexcounter LexCounter[K, V]) GetValue() V {
	var res V
	m := lexcounter.m
	for _, value := range m {
		res += value.second
	}
	return res
}

func Lexjoin[V Number](r, l Pair[V]) Pair[V] {
	r1 := r.first
	r2 := r.second
	l1 := l.first
	l2 := l.second

	if r1 == l1 && r2 == l2 {
		return r
	} else if l1 > r1 {
		return l
	} else if r1 > l1 {
		return r
	} else if r1 == l1 {
		res := Pair[V]{first: r1, second: r2 + l2}
		return res
	}

	var result Pair[V]
	return result
}

func (lexcounter LexCounter[K, V]) Join(lexcounter1 LexCounter[K, V]) {
	m1 := lexcounter1.m
	
	for key, value := range m1 {
		lexcounter.m[key] = Lexjoin(value, lexcounter.m[key])
	}
}

func Print[K Ordered, V Number](arguments ...LexCounter[K, V]) {
	for _, lexcounter := range arguments {
		m := lexcounter.m
		id := lexcounter.id

		fmt.Printf("LexCounter %s: (\n", id)

		for key, value := range m {
			fmt.Printf("    %s -> %d\n", key, value)
		}

		fmt.Print(")\n\n")
	}
}
