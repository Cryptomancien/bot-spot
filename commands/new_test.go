package commands

import (
	"fmt"
	"testing"
)

func TestCalcAmountUSDT(t *testing.T) {

	amountUSDT := CalcAmountUSDT(200.32, "6")
	fmt.Println(amountUSDT)

	priceBTC := 98000.00
	availableUSDT := 10.00

	amountCycleBTC := CalcAmountBTC(availableUSDT, priceBTC)
	fmt.Println(amountCycleBTC)
}
