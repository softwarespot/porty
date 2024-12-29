package helpers

import "time"

// FormatAsDateTime formats time to a "YYYY-MM-DD HH:MM:SS" string in the current location
func FormatAsDateTime(t time.Time) string {
	if t.IsZero() {
		return "0000-00-00 00:00:00"
	}
	return t.Format("2006-01-02 15:04:05")
}
