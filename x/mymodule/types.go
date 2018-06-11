package mymodule

// Description - description fields for a candidate
type ValueNum struct {
	Num int64 `json:"num"`
}

func NewValueNum(num int64)ValueNum {
	return ValueNum{
		Num:num,
	}
}