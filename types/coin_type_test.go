package types

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

const name  = "iris"

func makeCoinType() CoinType{
	return NewDefaultCoinType(name)
}

func TestCoinType_Convert(t *testing.T) {
	coinType := makeCoinType()
	coin := NewCoin(name,1)

	coinStr,err := coinType.Convert(coin.String(),"iris-milli")
	assert.Nil(t,err)
	assert.Equal(t,NewCoin("iris-milli",pow10(3)).String(),coinStr)

	coinStr,err = coinType.Convert(coin.String(),"iris-micro")
	assert.Nil(t,err)
	assert.Equal(t,NewCoin("iris-micro",pow10(6)).String(),coinStr)

	coinStr,err = coinType.Convert(coin.String(),"iris-nano")
	assert.Nil(t,err)
	assert.Equal(t,NewCoin("iris-nano",pow10(9)).String(),coinStr)

	coinStr,err = coinType.Convert(coin.String(),"iris-pico")
	assert.Nil(t,err)
	assert.Equal(t,NewCoin("iris-pico",pow10(12)).String(),coinStr)

	coinStr,err = coinType.Convert(coin.String(),"iris-femto")
	assert.Nil(t,err)
	assert.Equal(t,NewCoin("iris-femto",pow10(15)).String(),coinStr)

	coinStr,err = coinType.Convert(coin.String(),"iris-atto")
	assert.Nil(t,err)
	assert.Equal(t,NewCoin("iris-atto",pow10(18)).String(),coinStr)

}

func TestCoinType_ConvertToMinCoin(t *testing.T) {
	coinType := makeCoinType()
	coinIota,err := coinType.ConvertToMinCoin("1.1iris")
	assert.Nil(t,err)
	assert.Equal(t,NewCoin("iris-atto",11 * pow10(17) ),coinIota)

	coinAtto := "1 iris-atto"
	coinfemto,_ := coinType.Convert(coinAtto,"iris-femto")
	assert.Equal(t,"0.001iris-femto",coinfemto)


	coinfemto,_ = coinType.Convert(coinAtto,"iris-pico")
	assert.Equal(t,"0.000001iris-pico",coinfemto)

}


func pow10(n int) int64 {
	return NewIntWithDecimal(1, n).Int64()
}