package common

import (
    "fmt"
    "time"
)

// FormatDate mengubah time.Time menjadi string dengan format yyyy-mm-dd
func FormatDate(t time.Time) string {
    return t.Format("2006-01-02")
}

// ParseDate mengubah string dengan format yyyy-mm-dd menjadi time.Time
func ParseDate(dateStr string) (time.Time, error) {
    return time.Parse("2006-01-02", dateStr)
}

// GenerateID membuat ID sederhana dengan prefix dan timestamp
func GenerateID(prefix string) string {
    return fmt.Sprintf("%s-%d", prefix, time.Now().UnixNano())
}

// ContainsString cek apakah slice string mengandung sebuah nilai
func ContainsString(slice []string, str string) bool {
    for _, v := range slice {
        if v == str {
            return true
        }
    }
    return false
}
