package commands

import (
	"github.com/fatih/color"
	"main/database"
	"os"
	"strconv"
)

func makeRange(min, max int) []int {
	a := make([]int, max-min+1)
	for i := range a {
		a[i] = min + i
	}
	return a
}

func Clear() {
	args := os.Args[2:]

	if len(args) != 2 {
		color.Red("Start and End required")
		color.Cyan("Example: go run . -cl 12 35")
		return
	}

	startStr := args[0]
	endStr := args[1]

	start, _ := strconv.Atoi(startStr)
	end, _ := strconv.Atoi(endStr)

	r := makeRange(start, end)

	for i := range r {
		color.White("Deleting %d", r[i])
		database.DeleteByIdInt(int32(r[i]))
	}

	color.Green("Range successfully deleted")
}
