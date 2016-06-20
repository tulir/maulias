package main

import (
	"fmt"
	flag "maunium.net/go/mauflag"
	"os"
	"strings"
)

var wantHelp = flag.Make().LongKey("help").ShortKey("h").Bool()
var function = flag.Make().LongKey("function").ShortKey("f").Bool()
var oneliner = flag.Make().LongKey("single-line").ShortKey("s").Bool()
var withArgs = flag.Make().LongKey("with-args").ShortKey("a").Bool()
var condition = flag.Make().LongKey("if").ShortKey("i").String()
var elseCmd = flag.Make().LongKey("else").ShortKey("e").String()
var elseCmdIsLiteral = flag.Make().LongKey("with-args").ShortKey("l").Bool()

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
	err := flag.Parse()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		fmt.Fprintln(os.Stdout, help)
		os.Exit(1)
	} else if *wantHelp {
		fmt.Fprintln(os.Stdout, help)
		os.Exit(0)
	}
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
		fmt.Println("Warning: Arguments not included, as --with-args wasn't used.")
		flag.Args()[0] = strings.Split(flag.Arg(0), " ")[0]
	} else if flag.NArg() < 2 {
		fmt.Println("Usage: maulias [-afhls] [-e elseCommand] [-i condition] <ALIAS> <COMMANDS...>")
		return
	}

	if *withArgs {
		aliasWithArgs(newline, tab, elseFormat)
	} else if *function {
		aliasFunction(newline, tab, elseFormat)
	} else {
		aliasSimple()
	}
}

func aliasWithArgs(newline, tab, elseFormat string) {
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
	fmt.Print("}\n")                 // End function
}

func aliasFunction(newline, tab, elseFormat string) {
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
	fmt.Print("}\n") // End function
}

func aliasSimple() {
	fmt.Print("alias '", flag.Arg(0), "'='") // Alias beginning
	for i, arg := range flag.Args()[1:] {
		fmt.Print(replace(arg, '¤', '$')) // Commands
		if i+2 < flag.NArg() {
			fmt.Print(" && ")
		}
	}
	fmt.Print("'\n")
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
