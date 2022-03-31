package gdata

type Comparable interface {
	Compare(lhs, rhs interface{}) int
}

type CompareFunc func(lhs, rhs interface{}) int

func (f CompareFunc) Compare(lhs, rhs interface{}) int {
	return f(lhs, rhs)
}

type Int64 int64
func (val Int64) Compare(lhs, rhs interface{}) int {
	vl := lhs.(int64)
	vr := rhs.(int64)
	 if vl == vr {
	 	return 0
	 } else if vl < vr {
	 	return -1
	 } else {
	 	return 1
	 }
}

type String string
func (val String) Compare(lhs, rhs interface{}) int {
	vl := lhs.(int64)
	vr := rhs.(int64)
	if vl == vr {
		return 0
	} else if vl < vr {
		return -1
	} else {
		return 1
	}
}
