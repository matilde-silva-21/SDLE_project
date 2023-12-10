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
	First  int
	Second V
}

type LexCounter[K Ordered, V Number] struct {
	Map  map[K]Pair[V]
	id K
}

func Create[K Ordered, V Number](id K) LexCounter[K, V] {
	m := make(map[K]Pair[V])
	lc := LexCounter[K, V]{Map: m, id: id}
	return lc
}

func (lexcounter LexCounter[K, V]) Inc(toSum V) {
	m := lexcounter.Map
	id := lexcounter.id
	m[id] = Pair[V]{First: m[id].First + 1, Second: m[id].Second + toSum}
}

func (lexcounter LexCounter[K, V]) Dec(toDec V) {
	m := lexcounter.Map
	id := lexcounter.id
	
	m[id] = Pair[V]{First: m[id].First + 1, Second: m[id].Second - toDec}
	
}

func (lexcounter LexCounter[K, V]) GetValue() V {
	var res V
	m := lexcounter.Map
	for _, value := range m {
		res += value.Second
	}
	return res
}

func Lexjoin[V Number](r, l Pair[V]) Pair[V] {
	r1 := r.First
	r2 := r.Second
	l1 := l.First
	l2 := l.Second

	if r1 == l1 && r2 == l2 {
		return r
	} else if l1 > r1 {
		return l
	} else if r1 > l1 {
		return r
	} else if r1 == l1 {
		res := Pair[V]{First: r1, Second: r2 + l2}
		return res
	}

	var result Pair[V]
	return result
}

func (lexcounter LexCounter[K, V]) Join(lexcounter1 LexCounter[K, V]) {
	m1 := lexcounter1.Map
	
	for key, value := range m1 {
		lexcounter.Map[key] = Lexjoin(value, lexcounter.Map[key])
	}
}

func Print[K Ordered, V Number](arguments ...LexCounter[K, V]) {
	for _, lexcounter := range arguments {
		m := lexcounter.Map
		id := lexcounter.id

		fmt.Printf("LexCounter %s: (\n", id)

		for key, value := range m {
			fmt.Printf("    %s -> %d\n", key, value)
		}

		fmt.Print(")\n\n")
	}
}

func (lexcounter *LexCounter[K, V]) SetID(ID K) {
	(*lexcounter).id = ID
}