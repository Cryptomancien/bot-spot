package commands

import (
	"fmt"
	"testing"
)

func TestCalcAmountUSDT(t *testing.T) {

	amountUSDT := CalcAmountUSDT(5000, "6")
	fmt.Println(amountUSDT)

}
