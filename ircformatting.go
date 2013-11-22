// apply irc style colour/formatting codes
// eg, used can provide a string like:
//  Roses are %REDred%NORMAL, violets are %BLUEblue%NORMAL.
package main

import "strings"

var formatmap = map[string]string{
    "NORMAL":       "\u000f",
    "BOLD":         "\u0002",
    "UNDERLINE":    "\u001f",
    "REVERSE":      "\u0016",
    "WHITE":        "\u000300",
    "BLACK":        "\u000301",
    "DARK_BLUE":    "\u000302",
    "DARK_GREEN":   "\u000303",
    "RED":          "\u000304",
    "BROWN":        "\u000305",
    "PURPLE":       "\u000306",
    "OLIVE":        "\u000307",
    "YELLOW":       "\u000308",
    "GREEN":        "\u000309",
    "TEAL":         "\u000310",
    "CYAN":         "\u000311",
    "BLUE":         "\u000312",
    "MAGENTA":      "\u000313",
    "DARK_GRAY":    "\u000314",
    "LIGHT_GRAY":   "\u000315",
}

func applyFormatting(str string) string {
    for k, v := range formatmap {
        str = strings.Replace(str, "#"+k, v, -1)
        str = strings.Replace(str, "%"+k, v, -1)
    }
    return str
}
