package events

import (
	"encoding/csv"
	"fmt"
	"os"
	"time"
)

// LoadCSV parses yyyy-mm-dd,text rows into a map[date][]events
func LoadCSV(path string) (map[string][]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r := csv.NewReader(f)
	_, _ = r.Read() // try skip header

	events := make(map[string][]string)
	for {
		row, err := r.Read()
		if err != nil {
			break
		}
		if len(row) < 2 {
			continue
		}
		d, err := time.Parse("2006-01-02", row[0])
		if err != nil {
			fmt.Printf("Skipping bad date %s in %s\n", row[0], path)
			continue
		}
		key := d.Format("2006-01-02")
		events[key] = append(events[key], row[1])
	}
	return events, nil
}

func Merge(base, extra map[string][]string) {
	for k, v := range extra {
		base[k] = append(base[k], v...)
	}
}
