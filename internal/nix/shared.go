package nix

import (
	"fmt"
	"strconv"
	"strings"

	"srn/internal/shared"
)

func runNixCommand(command string, args ...string) {
	shared.CheckErr(shared.ExecuteCommand(command, args...))
}

// parseGenerations takes a slice of strings, which can contain numbers or ranges (e.g., "1-5"),
// and returns a slice of strings with all the individual numbers.
func parseGenerations(args []string) ([]string, error) {
	var generations []string
	for _, arg := range args {
		// Check if the argument contains a hyphen, indicating a range.
		if strings.Contains(arg, "-") {
			parts := strings.Split(arg, "-")
			if len(parts) != 2 {
				return nil, fmt.Errorf("invalid range format: %s", arg)
			}
			// Convert the start and end of the range to integers.
			start, err1 := strconv.Atoi(parts[0])
			end, err2 := strconv.Atoi(parts[1])
			if err1 != nil || err2 != nil {
				return nil, fmt.Errorf("invalid range numbers: %s", arg)
			}
			// Ensure the start of the range is not greater than the end.
			if start > end {
				return nil, fmt.Errorf("start of range cannot be greater than end: %s", arg)
			}
			// Append all numbers in the range to the generations slice.
			for i := start; i <= end; i++ {
				generations = append(generations, strconv.Itoa(i))
			}
		} else {
			// If it's not a range, just append the number.
			generations = append(generations, arg)
		}
	}
	return generations, nil
}
