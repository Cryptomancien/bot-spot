package commands

import (
	"fmt"
	"testing"
)

func TestCalcAmountUSDT(t *testing.T) {

	amountUSDT := CalcAmountUSDT(200.32, "6")
	fmt.Println(amountUSDT)

	priceBTC := float32(98000.00)
	availableUSDT := float32(10.00)

	amountCycleBTC := CalcAmountBTC(availableUSDT, priceBTC)
	fmt.Println(amountCycleBTC)
}
