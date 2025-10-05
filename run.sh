# Capture only stdout from the Go program (leave logs/warnings on stderr).
# Use tail to get the last non-empty line in case the program prints other info.
OUTFILE=$(go run . --year 2025 --birthdays data/birthdays.csv --holidays US 2>/dev/null | sed '/^\s*$/d' | tail -n 1)
if [ -z "$OUTFILE" ]; then
	echo "No output file generated" >&2
	exit 1
fi

# Ensure xelatex runs in the output directory so generated aux/log files are written there.
OUTDIR=$(dirname "$OUTFILE")
BASENAME=$(basename "$OUTFILE")
if [ -n "$OUTDIR" ] && [ "$OUTDIR" != "." ]; then
	cd "$OUTDIR" || { echo "Failed to cd to $OUTDIR" >&2; exit 1; }
fi

# Run xelatex non-interactively on the generated file name (in the output dir).
xelatex -interaction=nonstopmode "$BASENAME"