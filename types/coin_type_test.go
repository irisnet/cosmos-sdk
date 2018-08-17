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

	coinLua,err:=coinType.Convert(coin.String(),"iris_lua")
	assert.Nil(t,err)
	assert.Equal(t,NewCoin("iris_lua",1000000000).String(),coinLua)

	coinIris,err:=coinType.Convert(coin.String(),"iris")
	assert.Nil(t,err)
	assert.Equal(t,NewCoin("iris",1).String(),coinIris)

	coinIota,err:= coinType.Convert(coin.String(),"iris_iota")
	assert.Nil(t,err)
	assert.Equal(t,NewCoin("iris_iota",1000000000000000000).String(),coinIota)

	coin = NewCoin("iris_lua",10000000)
	coinIris2,err := coinType.Convert(coin.String(),"iris")
	assert.Nil(t,err)
	assert.Equal(t,"0.01iris",coinIris2)
}

func TestCoinType_ConvertToIota(t *testing.T) {
	coinType := makeCoinType()
	coinIota,err := coinType.ConvertToIota("1.1iris")
	assert.Nil(t,err)
	assert.Equal(t,NewCoin("iris_iota",1100000000000000000),coinIota)
}

func TestCoinType_ConvertFromCoinString(t *testing.T) {
	coinType := makeCoinType()
	coinStr := "1.2iris"
	coin,err := coinType.ConvertToIota(coinStr)
	assert.Nil(t,err)
	assert.Equal(t,NewCoin("iris_iota",1200000000000000000),coin)
}