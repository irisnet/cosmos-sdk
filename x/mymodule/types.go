package mymodule

// Description - description fields for a candidate
type ValueNum struct {
	num int64 `json:"num"`
}

func NewValueNum(num int64)ValueNum {
	return ValueNum{
		num:num,
	}
}