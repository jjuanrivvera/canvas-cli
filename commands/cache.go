package commands

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/canvas-cli/internal/cache"
)

var cacheCmd = &cobra.Command{
	Use:   "cache",
	Short: "Manage Canvas CLI cache",
	Long: `Manage the Canvas CLI response cache.

The cache stores API responses to reduce load on the Canvas server and improve
response times for repeated requests.

Examples:
  canvas cache stats                    # Show cache statistics
  canvas cache clear                    # Clear expired cache entries
  canvas cache clear --all              # Clear all cache entries`,
}

var cacheStatsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Show cache statistics",
	Long:  `Display statistics about the cache including size, entry counts, and hit rates.`,
	RunE:  runCacheStats,
}

var cacheClearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clear cache entries",
	Long: `Clear cached responses to free up disk space or force fresh data.

By default, only expired entries are cleared. Use --all to clear everything.

Examples:
  canvas cache clear          # Clear expired entries only
  canvas cache clear --all    # Clear all entries`,
	RunE: runCacheClear,
}

// Cache command flags
var (
	cacheClearAll bool
)

func init() {
	rootCmd.AddCommand(cacheCmd)

	// Add subcommands
	cacheCmd.AddCommand(cacheStatsCmd)
	cacheCmd.AddCommand(cacheClearCmd)

	// Flags for clear command
	cacheClearCmd.Flags().BoolVar(&cacheClearAll, "all", false, "Clear all cache entries (not just expired)")
}

func getCacheDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	return filepath.Join(home, ".canvas-cli", "cache"), nil
}

func runCacheStats(_ *cobra.Command, _ []string) error {
	cacheDir, err := getCacheDir()
	if err != nil {
		return err
	}

	// Check if cache directory exists
	if _, err := os.Stat(cacheDir); os.IsNotExist(err) {
		fmt.Println("Cache is empty (no cache directory).")
		return nil
	}

	// Get disk cache stats
	diskCache, err := cache.NewDiskCache(cacheDir, 0)
	if err != nil {
		return fmt.Errorf("failed to open cache: %w", err)
	}

	stats, err := diskCache.Stats()
	if err != nil {
		return fmt.Errorf("failed to get cache stats: %w", err)
	}

	// Calculate cache size
	size, err := getCacheDirSize(cacheDir)
	if err != nil {
		size = 0
	}

	fmt.Println("Cache Statistics")
	fmt.Println(strings.Repeat("=", 40))
	fmt.Printf("\nLocation: %s\n", cacheDir)
	fmt.Printf("Size: %s\n", formatBytes(size))
	fmt.Println()
	fmt.Printf("Total entries:   %d\n", stats.Total)
	fmt.Printf("Active entries:  %d\n", stats.Active)
	fmt.Printf("Expired entries: %d\n", stats.Expired)

	if stats.Total > 0 {
		activePercent := float64(stats.Active) / float64(stats.Total) * 100
		fmt.Printf("\nActive rate: %.1f%%\n", activePercent)
	}

	return nil
}

func runCacheClear(_ *cobra.Command, _ []string) error {
	cacheDir, err := getCacheDir()
	if err != nil {
		return err
	}

	// Check if cache directory exists
	if _, err := os.Stat(cacheDir); os.IsNotExist(err) {
		fmt.Println("Cache is already empty.")
		return nil
	}

	diskCache, err := cache.NewDiskCache(cacheDir, 0)
	if err != nil {
		return fmt.Errorf("failed to open cache: %w", err)
	}

	if cacheClearAll {
		// Confirm clearing all
		fmt.Print("Are you sure you want to clear ALL cache entries? [y/N]: ")
		reader := bufio.NewReader(os.Stdin)
		response, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read response: %w", err)
		}

		response = strings.TrimSpace(strings.ToLower(response))
		if response != "y" && response != "yes" {
			fmt.Println("Cancelled.")
			return nil
		}

		// Get count before clearing
		stats, _ := diskCache.Stats()
		totalBefore := stats.Total

		if err := diskCache.Clear(); err != nil {
			return fmt.Errorf("failed to clear cache: %w", err)
		}

		fmt.Printf("Cleared %d cache entries.\n", totalBefore)
	} else {
		// Clear only expired entries
		stats, _ := diskCache.Stats()
		expiredBefore := stats.Expired

		if expiredBefore == 0 {
			fmt.Println("No expired entries to clear.")
			return nil
		}

		// Clear expired entries by reading each file
		if err := clearExpiredEntries(cacheDir); err != nil {
			return fmt.Errorf("failed to clear expired entries: %w", err)
		}

		fmt.Printf("Cleared %d expired cache entries.\n", expiredBefore)
	}

	return nil
}

func getCacheDirSize(dir string) (int64, error) {
	var size int64

	err := filepath.Walk(dir, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return nil
	})

	return size, err
}

func formatBytes(bytes int64) string {
	const (
		KB = 1024
		MB = 1024 * KB
		GB = 1024 * MB
	)

	switch {
	case bytes >= GB:
		return fmt.Sprintf("%.2f GB", float64(bytes)/float64(GB))
	case bytes >= MB:
		return fmt.Sprintf("%.2f MB", float64(bytes)/float64(MB))
	case bytes >= KB:
		return fmt.Sprintf("%.2f KB", float64(bytes)/float64(KB))
	default:
		return fmt.Sprintf("%d bytes", bytes)
	}
}

func clearExpiredEntries(dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		path := filepath.Join(dir, entry.Name())

		// Read and check expiration
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}

		if isExpiredCacheFile(data) {
			os.Remove(path)
		}
	}

	return nil
}

func isExpiredCacheFile(data []byte) bool {
	type cacheItem struct {
		Expiration time.Time `json:"expiration"`
	}

	var item cacheItem
	if err := json.Unmarshal(data, &item); err != nil {
		return false
	}

	return time.Now().After(item.Expiration)
}
