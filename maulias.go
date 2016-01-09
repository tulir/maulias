package main

import (
	"fmt"
	flag "github.com/ogier/pflag"
	"strings"
)

var function = flag.BoolP("function", "f", false, "")
var oneliner = flag.BoolP("single-line", "s", false, "")
var withArgs = flag.BoolP("with-args", "a", false, "")
var condition = flag.StringP("if", "i", "", "")
var elseCmd = flag.StringP("else", "e", "", "")
var elseCmdIsLiteral = flag.BoolP("else-literal", "l", false, "")

var help = `Generates bash functions or alias commands using given values. Intended as an advanced aliasing system.

Usage:
  maulias [-afls] [-e elseCommand] [-i condition] <ALIAS> <COMMAND...>

Help options:
  -h, --help               Show this help page

Application options:
  -f, --function           Generate a function rather than an alias command.
  -s, --single-line        Generate a single-line function. Doesn't affect anything without --function
  -a, --with-args          Generate a function that also checks arguments. When used, --function is ignored
  -i, --if=CONDITION       If set, this will be used as the condition inside an if. Only works with --function. Doesn't work with --with-args
  -e, --else=CMD           The command to use in the else case (used with --with-args and --if)
  -l, --else-literal       Don't add \"$@\"; to the else statement.
`

func init() {
	flag.Usage = func() {
		print(help)
	}
	flag.Parse()
}

func main() {
	var newline, tab, elseFormat string

	if *oneliner {
		newline = " "
		tab = ""
	} else {
		newline = "\n"
		tab = "	"
	}
	if *elseCmdIsLiteral {
		elseFormat = tab + tab + "%[1]s" + newline
	} else {
		elseFormat = tab + tab + "%[1]s \"$@\";" + newline
	}

	if strings.Contains(flag.Arg(0), " ") && !*withArgs {
		println("Warning: Arguments not included, as --with-args wasn't used.")
		flag.Args()[0] = strings.Split(flag.Arg(0), " ")[0]
	} else if flag.NArg() < 2 {
		println("Usage: maulias [-afls] [-e elseCommand] [-i condition] <ALIAS> <COMMANDS...>")
		return
	}

	if *withArgs {
		sargs := strings.Split(flag.Arg(0), " ")
		fmt.Print(sargs[0], "() {", newline) // Start function
		sargs = sargs[1:]
		fmt.Print(tab, "if [ ") // Start if. Condition bracket opening
		for i, sarg := range sargs {
			fmt.Printf("\"$%[1]d\" = \"%[2]s\"", i+1, sarg) // If conditions
			if i+1 < len(sargs) {
				fmt.Printf(" -a ") // Bash equivalent to "&&" (logical AND)
			}
		}
		fmt.Print(" ];", newline)       // Condition bracket closing
		fmt.Print(tab, "then", newline) // Then
		for _, arg := range flag.Args()[1:] {
			fmt.Printf(tab+tab+"%[1]s;"+newline, replace(arg, '¤', '$')) // Commands
		}
		fmt.Print(tab, "else", newline)  // "Else"
		fmt.Printf(elseFormat, *elseCmd) // Else command
		fmt.Print(tab, "fi", newline)    // End if
		fmt.Print("}")                   // End function
	} else if *function {
		fmt.Print(flag.Arg(0), "() {", newline) // Start function
		if len(*condition) != 0 {
			fmt.Printf(tab+"if [ %[1]s ]; ", *condition)
			fmt.Print(tab, "then", newline)
			for _, arg := range flag.Args()[1:] {
				fmt.Printf(tab+tab+"%[1]s;"+newline, replace(arg, '¤', '$')) // Commands
			}
			if len(*elseCmd) != 0 {
				fmt.Print(tab, "else", newline)  // "Else"
				fmt.Printf(elseFormat, *elseCmd) // Else command
			}
			fmt.Print(tab, "fi", "newline") // End if
		} else {
			for _, arg := range flag.Args()[1:] {
				fmt.Printf(tab+tab+"%[1]s;"+newline, replace(arg, '¤', '$')) // Commands
			}
		}
		fmt.Print("}") // End function
	} else {
		fmt.Print("alias '", flag.Arg(0), "'='") // Alias beginning
		for i, arg := range flag.Args()[1:] {
			fmt.Print(replace(arg, '¤', '$')) // Commands
			if i+2 < flag.NArg() {
				fmt.Print(" && ")
			}
		}
		fmt.Print("'")
	}
	fmt.Print("\n")
}

func replace(str string, toRepl, replWith rune) string {
	strune := []rune(str)
	for i, char := range strune {
		if i != 0 && strune[i-1] == '\\' {
			continue
		}
		if char == toRepl {
			strune[i] = replWith
		}
	}
	return string(strune)
}
