package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ddol/go-family-calendar/events"
	"github.com/ddol/go-family-calendar/tex"
)

func main() {
	year := flag.Int("year", 2025, "Calendar year")
	birthdays := flag.String("birthdays", "", "CSV file of birthdays")
	holidays := flag.String("holidays", "", "Comma-separated country codes for holiday CSVs (e.g. US,IE,NZ)")
	// Default output file uses a timestamp so multiple runs don't overwrite each other.
	defaultOut := filepath.Join("output", time.Now().Format("20060102-15.04.05")+"_calendar.tex")
	outfile := flag.String("out", defaultOut, "Output LaTeX file")
	flag.Parse()

	// Load events
	allEvents := make(map[string][]string) // yyyy-mm-dd -> [events]

	if *birthdays != "" {
		evs, err := events.LoadCSV(*birthdays)
		if err != nil {
			log.Fatalf("Failed to load birthdays: %v", err)
		}
		events.Merge(allEvents, evs)
	}

	if *holidays != "" {
		for _, cc := range strings.Split(*holidays, ",") {
			// Holiday files are now named with the year prefix: data/holidays/YYYY-CC.csv
			path := fmt.Sprintf("data/holidays/%d-%s.csv", *year, strings.ToUpper(strings.TrimSpace(cc)))
			evs, err := events.LoadCSV(path)
			if err != nil {
				// Fail gracefully: log a warning and continue compilation with other events
				log.Printf("Warning: could not load holiday file %s: %v", path, err)
				continue
			}
			events.Merge(allEvents, evs)
		}
	}

	// Ensure output directory exists, then open output file
	outDir := filepath.Dir(*outfile)
	if outDir != "" && outDir != "." {
		if err := os.MkdirAll(outDir, 0o755); err != nil {
			log.Fatalf("Failed to create output directory %s: %v", outDir, err)
		}
	}

	out, err := os.Create(*outfile)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	// Generate .tex
	tex.WritePreamble(out)
	for month := 1; month <= 12; month++ {
		tex.RenderMonth(out, *year, month, allEvents)
	}
	tex.WritePostamble(out)

	// Print the filename on stdout so wrapper scripts can capture it.
	fmt.Println(*outfile)
}
