package tex

import (
	"fmt"
	"io"
	"time"
)

func WritePreamble(w io.Writer) {
	fmt.Fprintln(w, `\documentclass[11pt]{article}
\usepackage[margin=0.6in]{geometry}
\usepackage{tabularx}
\usepackage{array}
\usepackage{fontspec}
\setmainfont{Helvetica Neue}
\usepackage{ragged2e}
\begin{document}
\pagestyle{empty}`)
}

func TexEscape(s string) string {
	replacements := map[rune]string{
		'&': "\\&", '%': "\\%", '$': "\\$", '#': "\\#",
		'_': "\\_", '{': "\\{", '}': "\\}", '~': "\\textasciitilde{}",
		'^': "\\textasciicircum{}", '\\': "\\textbackslash{}",
	}
	out := ""
	for _, ch := range s {
		if rep, ok := replacements[ch]; ok {
			out += rep
		} else {
			out += string(ch)
		}
	}
	return out
}

func WritePostamble(w io.Writer) {
	fmt.Fprintln(w, `\end{document}`)
}

func RenderMonth(w io.Writer, year, month int, evmap map[string][]string) {
	monthName := time.Month(month).String()
	fmt.Fprintf(w, "\\section*{%s %d}\n", monthName, year)

	// weekday headers
	fmt.Fprintln(w, "\\begin{tabularx}{\\textwidth}{|*{7}{>{\\centering\\arraybackslash}X|}}\\hline")
	fmt.Fprintln(w, "Sun & Mon & Tue & Wed & Thu & Fri & Sat\\\\ \\hline")

	// compute first day
	first := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	offset := int(first.Weekday())

	// total days
	daysInMonth := time.Date(year, time.Month(month+1), 0, 0, 0, 0, 0, time.UTC).Day()

	day := 1 - offset
	for day <= daysInMonth {
		for wd := 0; wd < 7; wd++ {
			if day < 1 || day > daysInMonth {
				fmt.Fprint(w, " ~ ")
			} else {
				dateStr := fmt.Sprintf("%04d-%02d-%02d", year, month, day)
				evs := evmap[dateStr]
				fmt.Fprintf(w, "\\textbf{%d}\\\\", day)
				for _, ev := range evs {
					fmt.Fprintf(w, "%s\\\\", TexEscape(ev))
				}
			}
			if wd < 6 {
				fmt.Fprint(w, " & ")
			}
			day++
		}
		fmt.Fprintln(w, "\\\\ \\hline")
	}

	fmt.Fprintln(w, "\\end{tabularx}")
	fmt.Fprintln(w, "\\newpage")
}
