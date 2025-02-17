package commands

import exchange "main/services"

func New() {
	exchange.CheckConnection()

	exchange.CheckBalanceUSDT()
}
