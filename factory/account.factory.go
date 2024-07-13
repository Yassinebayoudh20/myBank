package factory

import (
	"math/rand"

	"github.com/Yassinebayoudh20/my_bank/util"
)

func RandomOwner() string {
	return util.RandomString(6)
}

func RandomMoney() int64 {
	return util.RandomInt(0, 1000)
}

func RandomCurrency() string {
	currencies := []string{"EUR", "USD", "TND"}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}
