package commands

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// confirmDelete prompts the user for confirmation before a destructive operation.
// Returns true if the user confirms, false otherwise.
// If force is true, skips the prompt and returns true.
func confirmDelete(resourceType string, resourceID interface{}, force bool) (bool, error) {
	if force {
		return true, nil
	}

	fmt.Printf("Are you sure you want to delete %s %v? [y/N]: ", resourceType, resourceID)
	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		return false, fmt.Errorf("failed to read response: %w", err)
	}

	response = strings.TrimSpace(strings.ToLower(response))
	return response == "y" || response == "yes", nil
}
