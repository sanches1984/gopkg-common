package types

import "reflect"

func SliceStringIntersection(a, b []string) (c []string) {
	m := make(map[string]bool)

	for _, item := range a {
		m[item] = true
	}

	for _, item := range b {
		if _, ok := m[item]; ok {
			c = append(c, item)
		}
	}
	return
}

func SliceStringContains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func SliceIntContains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func SliceInt32Contains(s []int32, e int32) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func SliceInt64Contains(s []int64, e int64) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func IsIntersect(a interface{}, b interface{}) bool {
	hash := make(map[interface{}]bool)
	av := reflect.ValueOf(a)
	bv := reflect.ValueOf(b)

	for i := 0; i < av.Len(); i++ {
		el := av.Index(i).Interface()
		hash[el] = true
	}

	for i := 0; i < bv.Len(); i++ {
		el := bv.Index(i).Interface()
		if _, found := hash[el]; found {
			return true
		}
	}

	return false
}

func SliceUniqueInt(intSlice []int) []int {
	keys := make(map[int]struct{}, len(intSlice))
	ret := make([]int, 0, len(intSlice))
	for _, item := range intSlice {
		if _, ok := keys[item]; !ok {
			keys[item] = struct{}{}
			ret = append(ret, item)
		}
	}
	return ret
}

func SliceUniqueInt64(intSlice []int64) []int64 {
	keys := make(map[int64]struct{}, len(intSlice))
	ret := make([]int64, 0, len(intSlice))
	for _, item := range intSlice {
		if _, ok := keys[item]; !ok {
			keys[item] = struct{}{}
			ret = append(ret, item)
		}
	}
	return ret
}

func SliceNotZeroInt64(id int64) []int64 {
	if id == 0 {
		return []int64{}
	}
	return []int64{id}
}
