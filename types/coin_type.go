package types

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

const Iota = "iota"
const Lua = "lua"

type Unit struct {
	Denom   string `json:"denom"`
	Decimal int    `json:"decimal"`
	IsMin   bool   `json:"is_min"`
}

type Units = []Unit

type CoinType struct {
	Name  string `json:"name"`
	Units Units  `json:"units"`
}

type CoinTypeSet struct {
	CoinTypes []CoinType `json:"coin_type"`
}

func NewCoinTypeSet() CoinTypeSet {
	var typs []CoinType
	return CoinTypeSet{
		typs,
	}
}

func (cts *CoinTypeSet) Add(ct CoinType) {
	for _, typ := range cts.CoinTypes {
		if typ.Name == ct.Name {
			return
		}
	}
	cts.CoinTypes = append(cts.CoinTypes, ct)
}

func (ct CoinType) Convert(orgCoinStr string, denom string) (destCoinStr string, err error) {
	orgDenom, orgAmt, _ := GetCoin(orgCoinStr)
	var destUint Unit
	if destUint,err = ct.GetUnit(denom); err != nil {
		return destCoinStr,errors.New("not exist unit " + orgDenom)
	}

	if orgUnit, ok := ct.GetUnit(orgDenom); ok == nil {
		numerator := NewIntWithDecimal(1, destUint.Decimal)
		denominator := NewIntWithDecimal(1, orgUnit.Decimal)
		rat := NewRatFromInt(numerator, denominator)
		amount, _ := NewRatFromDecimal(orgAmt, destUint.Decimal)
		amt := amount.Mul(rat).DecimalString(orgUnit.Decimal)
		destCoinStr = fmt.Sprintf("%s%s",amt,destUint.Denom)
		return destCoinStr, nil
	}
	return destCoinStr, errors.New("not exist unit " + orgDenom)
}

func (ct CoinType) ConvertToIota(coinStr string) (coin Coin, err error) {
	minUint := ct.GetMinUnit()

	if destCoinStr,err := ct.Convert(coinStr,minUint.Denom);err == nil {
		coin,err = ParseCoin(destCoinStr)
		return coin,err
	}

	return coin, errors.New("convert error")
}

func (ct CoinType) GetUnit(denom string) (u Unit, err error) {
	for _, unit := range ct.Units {
		if denom == unit.Denom {
			return unit, nil
		}
	}
	return u, errors.New("not find unit " + denom)
}

func (ct CoinType) GetMinUnit() (unit Unit) {
	for _, unit := range ct.Units {
		if unit.IsMin {
			return unit
		}
	}
	return unit
}

func (ct CoinType) GetMaxUnit() (unit Unit) {
	unit ,_ = ct.GetUnit(ct.Name)
	return unit
}


func (ct CoinType) String() string {
	bz, _ := json.Marshal(ct)
	return string(bz)
}

func NewDefaultCoinType(name string) CoinType {
	org := Unit{
		Denom:   name,
		Decimal: 0,
		IsMin:   false,
	}

	iota := Unit{
		Denom:   name + "_" + Iota,
		Decimal: 18,
		IsMin:   true,
	}

	lua := Unit{
		Denom:   name + "_" + Lua,
		Decimal: 9,
		IsMin:   false,
	}

	units := make(Units, 3)
	units[0] = iota
	units[1] = lua
	units[2] = org

	return CoinType{
		Name:  name,
		Units: units,
	}
}

func CoinTypeKey(coinName string) string {
	return fmt.Sprintf("%s/%s","coin_types",coinName)
}

func GetCoin(coinStr string) (denom, amount string, err error) {
	var (
		reDnm  = `[[:alpha:]][[:word:]]{2,15}`
		reAmt  = `[0-9]+[.]?[0-9]*`
		reSpc  = `[[:space:]]*`
		reCoin = regexp.MustCompile(fmt.Sprintf(`^(%s)%s(%s)$`, reAmt, reSpc, reDnm))
	)

	coinStr = strings.TrimSpace(coinStr)

	matches := reCoin.FindStringSubmatch(coinStr)
	if matches == nil {
		err = fmt.Errorf("invalid coin expression: %s", coinStr)
		return
	}
	denom, amount = matches[2], matches[1]
	return
}
