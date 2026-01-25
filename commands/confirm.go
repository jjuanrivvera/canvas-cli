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
// If dryRun is true (global flag), shows what would be deleted and returns false.
func confirmDelete(resourceType string, resourceID interface{}, force bool) (bool, error) {
	// In dry-run mode, show what would be deleted but don't proceed
	if dryRun {
		fmt.Printf("DRY RUN: Would delete %s %v\n", resourceType, resourceID)
		fmt.Println("No changes were made.")
		return false, nil
	}

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

// confirmDeleteWithDetails prompts for confirmation and shows resource details.
// This provides a preview of what will be deleted.
func confirmDeleteWithDetails(resourceType string, resourceID interface{}, details map[string]interface{}, force bool) (bool, error) {
	// In dry-run mode, show detailed preview
	if dryRun {
		fmt.Printf("DRY RUN: Would delete %s %v\n", resourceType, resourceID)
		if len(details) > 0 {
			fmt.Println("Resource details:")
			for key, value := range details {
				fmt.Printf("  %s: %v\n", key, value)
			}
		}
		fmt.Println("\nNo changes were made.")
		return false, nil
	}

	// Show details before asking for confirmation
	if len(details) > 0 && !force {
		fmt.Printf("About to delete %s %v:\n", resourceType, resourceID)
		for key, value := range details {
			fmt.Printf("  %s: %v\n", key, value)
		}
		fmt.Println()
	}

	if force {
		return true, nil
	}

	fmt.Print("Are you sure you want to delete this? [y/N]: ")
	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		return false, fmt.Errorf("failed to read response: %w", err)
	}

	response = strings.TrimSpace(strings.ToLower(response))
	return response == "y" || response == "yes", nil
}

// confirmUpdate shows what will be changed in dry-run mode
func confirmUpdateDryRun(resourceType string, resourceID interface{}, changes map[string]interface{}) bool {
	if !dryRun {
		return false
	}

	fmt.Printf("DRY RUN: Would update %s %v\n", resourceType, resourceID)
	if len(changes) > 0 {
		fmt.Println("Changes:")
		for key, value := range changes {
			fmt.Printf("  %s: %v\n", key, value)
		}
	}
	fmt.Println("\nNo changes were made.")
	return true
}
