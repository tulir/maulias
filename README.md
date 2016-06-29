# maulias
[![Build Status](http://img.shields.io/travis/tulir293/maulias.svg?style=flat-square)](https://travis-ci.org/tulir293/maulias)
[![License](http://img.shields.io/:license-mit-blue.svg?style=flat-square)](http://tulir293.mit-license.org)

A command to generate bash functions that act like advanced aliases.

```
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
  ```
