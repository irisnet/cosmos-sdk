package types

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

const (
	//1 iris = 10^3 iris-milli
	Milli = "milli"

	//1 iris = 10^6 iris-micro
	Micro = "micro"

	//1 iris = 10^9 iris-nano
	Nano = "nano"

	//1 iris = 10^12 iris-pico
	Pico = "pico"

	//1 iris = 10^15 iris-femto
	Femto = "femto"

	//1 iris = 10^18 iris-atto
	Atto = "atto"
)

var (
	MainUnit = func(coinName string) Unit {
		return NewUnit(coinName, 0)
	}

	MilliUnit = func(coinName string) Unit {
		denom := fmt.Sprintf("%s-%s", coinName, Milli)
		return NewUnit(denom, 3)
	}

	MicroUnit = func(coinName string) Unit {
		denom := fmt.Sprintf("%s-%s", coinName, Micro)
		return NewUnit(denom, 6)
	}

	NanoUnit = func(coinName string) Unit {
		denom := fmt.Sprintf("%s-%s", coinName, Nano)
		return NewUnit(denom, 9)
	}

	PicoUnit = func(coinName string) Unit {
		denom := fmt.Sprintf("%s-%s", coinName, Pico)
		return NewUnit(denom, 12)
	}

	FemtoUnit = func(coinName string) Unit {
		denom := fmt.Sprintf("%s-%s", coinName, Femto)
		return NewUnit(denom, 15)
	}

	AttoUnit = func(coinName string) Unit {
		denom := fmt.Sprintf("%s-%s", coinName, Atto)
		return NewUnit(denom, 18)
	}
)

type Unit struct {
	Denom   string `json:"denom"`
	Decimal int    `json:"decimal"`
}

func NewUnit(denom string, decimal int) Unit {
	return Unit{
		Denom:   denom,
		Decimal: decimal,
	}
}
func (u Unit) GetPrecision() Int {
	return NewIntWithDecimal(1, u.Decimal)
}

type Units = []Unit

type CoinType struct {
	Name         string `json:"name"`
	MinUnitDenom string `json:"minUnitDenom"`
	Units        Units  `json:"units"`
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
	orgDenom, orgAmt, err := GetCoin(orgCoinStr)
	if err != nil {
		return destCoinStr, err
	}
	var destUint Unit
	if destUint, err = ct.GetUnit(denom); err != nil {
		return destCoinStr, errors.New("not exist unit " + orgDenom)
	}
	// 目标Coin = 原金额 * (10^目标精度 / 10^原精度)
	if orgUnit, err := ct.GetUnit(orgDenom); err == nil {
		rat := NewRatFromInt(destUint.GetPrecision(), orgUnit.GetPrecision())
		amount, _ := NewRatFromDecimal(orgAmt, destUint.Decimal)//将原金额按照目标精度转化
		amt := amount.Mul(rat).DecimalString(destUint.Decimal)
		destCoinStr = fmt.Sprintf("%s%s", amt, destUint.Denom)
		return destCoinStr, nil
	}
	return destCoinStr, errors.New("not exist unit " + orgDenom)
}

func (ct CoinType) ConvertToMinCoin(coinStr string) (coin Coin, err error) {
	minUint := ct.GetMinUnit()

	if destCoinStr, err := ct.Convert(coinStr, minUint.Denom); err == nil {
		coin, err = ParseCoin(destCoinStr)
		return coin, err
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
		if unit.Denom == ct.MinUnitDenom {
			return unit
		}
	}
	return unit
}

func (ct CoinType) GetMaxUnit() (unit Unit) {
	unit, _ = ct.GetUnit(ct.Name)
	return unit
}

func (ct CoinType) String() string {
	bz, _ := json.Marshal(ct)
	return string(bz)
}

func NewDefaultCoinType(name string) CoinType {
	units := GetDefaultUnits(name)
	return CoinType{
		Name:         name,
		Units:        units,
		MinUnitDenom: units[6].Denom,
	}
}

func CoinTypeKey(coinName string) string {
	return fmt.Sprintf("%s/%s/%s", "global", "coin_types", coinName)
}

func GetDefaultUnits(coin string) Units {
	units := make(Units, 7)
	units[0] = MainUnit(coin)
	units[1] = MilliUnit(coin)
	units[2] = MicroUnit(coin)
	units[3] = NanoUnit(coin)
	units[4] = PicoUnit(coin)
	units[5] = FemtoUnit(coin)
	units[6] = AttoUnit(coin)
	return units
}

func GetCoin(coinStr string) (denom, amount string, err error) {
	var (
		reDnm  = `[A-Za-z\-]{2,15}`
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

func GetMainCoinDenom(coinStr string) (coinName string, err error) {
	denom, _, err := GetCoin(coinStr)
	if err != nil {
		return coinName, err
	}
	coinName = strings.Split(denom, "-")[0]
	return coinName, nil
}
