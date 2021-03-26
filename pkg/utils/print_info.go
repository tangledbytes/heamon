package utils

import "fmt"

// PrintInfo prints version, commit and date info
// on the terminal
func PrintInfo(version, commit, date string) {
	fmt.Println("🔥 Heamon version:", version)
	fmt.Println("🛠️ Commit:", commit)
	fmt.Println("📅 Release Date:", date)
}
