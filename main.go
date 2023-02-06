package main

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var out io.Writer = os.Stdout
var unittest = false // trigger for unit test injection

func main() {
	args := os.Args
	lena := len(args)
	if lena == 3 {
		good(args[1:])
	}
}

// returns slice of runes from incoming string
func xr(str string) []rune { return []rune(str) }

// slice of strings from slice of runes
func xrxs(xr []rune) (xs []string) {
	for _, r := range xr {
		xs = append(xs, string(r))
	}
	return
}

// incoming string to slice of strings
func sxs(s string) []string {
	return xrxs(xr(s))
}

// incoming slice of strings to string
func xss(xs []string) string {
	return strings.Join(xs, "")
}

// cut broken rune at the end of string, was detected before, probably no needed at the moment
func cutLastRune(str string) (cutted string) {
	runes := xr(str)
	runes = runes[:len(runes)-1]
	cutted = string(runes)
	return
}

// read file from "fpath" location and return (string, bool)
func rfile(fpath string) (str string, ok bool) {
	bytes, err := os.ReadFile(fpath)
	if err == nil {
		// str = cutLastRune(string(bytes))
		str = string(bytes)
		ok = true
	} else {
		str = ""
		ok = false
	}
	return
}

// write file to "fpath" location
func wfile(fpath string, data string) {
	os.WriteFile(fpath, []byte(data), 0666)
}

// punctuation list
var plist = [...]string{",", ".", ";", ":", "!", "?", "'"}

// check the given string is in punctuation list
func isPunctuation(s string) bool {
	for _, p := range plist {
		if s == p {
			return true
		}
	}
	return false
}

// add spaces before and after the given string
func addSpaces(s string) string {
	return " " + s + " "
}

func cutMultiSpaces(s string) string {
	reg, err := regexp.Compile("[ ]{2,}") // more than one space
	if err == nil {
		s = reg.ReplaceAllString(s, " ")
	}
	return s
}

// add spaces before and after every punctuation symbols ", . ; : ! ?"
func preformatPunctuationSpaces(str string) string {
	runes := xr(str) // split string to runes slice
	xs := xrxs(runes)
	for ind, s := range xs {
		if isPunctuation(s) {
			xs[ind] = addSpaces(s)
		}
	}
	return strings.Join(xs, "")
}

// add spaces around punctuation and before and after string, to manage case of format rules placed f.e. at the end of string
func preformatSpaces(str string) string {
	return cutMultiSpaces(addSpaces(preformatPunctuationSpaces(str)))
}

// hex rule implementation
func hexManager(s string) string {
	xs := regspace.Split(s, -1)
	lena := len(xs) - 1
	inum, err := strconv.ParseInt(xs[lena], 16, 0)
	if err == nil {
		xs[lena] = strconv.FormatInt(inum, 10)
	}
	return strings.Join(xs, " ")
}

// bin rule implementation
func binManager(s string) string {
	xs := regspace.Split(s, -1)
	lena := len(xs) - 1
	inum, err := strconv.ParseInt(xs[lena], 2, 0)
	if err == nil {
		xs[lena] = strconv.FormatInt(inum, 10)
	}
	return strings.Join(xs, " ")
}

// capitalise all letters of word
func upMaker(s string) string {
	return strings.ToUpper(s)
}

// up rule implementation. upper case all letters
func upManager(s string, xnum int) string {
	// myprint("UP manager EXECUTED", true)
	xs := regspace.Split(s, -1)
	lena := len(xs)
	for i := 0; i < xnum; i++ {
		if i-lena < 0 {
			xs[lena-i-1] = upMaker(xs[lena-i-1])
		} else {
			break
		}
	}
	return strings.Join(xs, " ")
}

func lowMaker(s string) string {
	return strings.ToLower(s)
}

// low rule implementation. lower case all letters
func lowManager(s string, xnum int) string {
	// myprint("LOW manager EXECUTED", true)
	xs := regspace.Split(s, -1)
	lena := len(xs)
	for i := 0; i < xnum; i++ {
		if i-lena < 0 {
			xs[lena-i-1] = lowMaker(xs[lena-i-1])
		} else {
			break
		}
	}
	return strings.Join(xs, " ")
}

// capitalise first letter of word
func capMaker(s string) string {
	// myprint("CAP maker EXECUTED", true)
	xs := sxs(s)
	if len(xs) > 0 {
		// myprint("xs[0] "+xs[0], true)
		// fmt.Println(len(xs[0]))
		xs[0] = upMaker(xs[0])
		// fmt.Println(xs[0])
	}
	return xss(xs)
}

// cap rule implementation. capitalise the first letter
func capManager(s string, xnum int) string {
	// myprint("CAP manager EXECUTED", true)
	xs := regspace.Split(s, -1)
	lena := len(xs)
	for i := 0; i < xnum; i++ {
		if i-lena < 0 {
			xs[lena-i-1] = capMaker(xs[lena-i-1])
		} else {
			break
		}
	}
	return strings.Join(xs, " ")
}

// regular expression to find number
var regnum, _ = regexp.Compile(`\d+`)

// regular expression to find space
var regspace, _ = regexp.Compile(`\s`)

// regular expressions to find every rule name
var reghex, _ = regexp.Compile(`hex`)
var regbin, _ = regexp.Compile(`bin`)
var regup, _ = regexp.Compile(`up`)
var reglow, _ = regexp.Compile(`low`)
var regcap, _ = regexp.Compile(`cap`)

func ruleMaker(rule, s string) string {
	xnum := 1 // default one word before rule will be changed
	// calculate the number of words before rule to modify
	if regnum.MatchString(rule) { // rule include number
		// convert string to integer using standard library
		num, err := strconv.Atoi(regnum.FindString(rule))
		if err == nil { // no error, now set new repeat number
			xnum = num
		}
	}
	// managing every case of rules
	if reghex.MatchString(rule) {
		s = hexManager(s)
	} else if regbin.MatchString(rule) {
		s = binManager(s)
	} else if regup.MatchString(rule) {
		s = upManager(s, xnum)
	} else if reglow.MatchString(rule) {
		s = lowManager(s, xnum)
	} else if regcap.MatchString(rule) {
		s = capManager(s, xnum)
	}
	return s
}

// string modify according to rules up,low,cap,hex,bin
func ruleManager(s string) string {
	ruleregraw := `(\s\((up|low|cap)\s,\s\d+\)\s|\s\((up|low|cap|hex|bin)\)\s)`
	rreg, _ := regexp.Compile(ruleregraw)
	haverule := true // some rule still was found inside string
	// looped trying to find the rule inside string
	for haverule {
		if rreg.MatchString(s) {
			rule := rreg.FindString(s) // get the rule
			faceback := rreg.Split(s, 2)
			s = ruleMaker(rule, faceback[0])
			s += " " + faceback[1]
			// haverule = false
		} else {
			haverule = false
		}
	}
	return s
}

// list of punctuations need to be formatted separately
var x6 = plist[:6]

// check the given string is in x6 list
func inx6(s string) bool {
	for _, p := range x6 {
		if s == p {
			return true
		}
	}
	return false
}

// regular expression for required punctuation and left space before
var reg0, _ = regexp.Compile(`(\s\,)`)
var reg1, _ = regexp.Compile(`(\s\.)`)
var reg2, _ = regexp.Compile(`(\s\;)`)
var reg3, _ = regexp.Compile(`(\s\:)`)
var reg4, _ = regexp.Compile(`(\s\!)`)
var reg5, _ = regexp.Compile(`(\s\?)`)

// slice of regular expressions to find punctuations with left space
var reg = []*regexp.Regexp{reg0, reg1, reg2, reg3, reg4, reg5}
var rep = strings.Split(",.;:!?", "") // slice of replacements

// remove spaces before every instance of next punctuation
func x6join(s string) string {
	for i := range reg {
		s = reg[i].ReplaceAllString(s, rep[i])
	}
	return s
}

// regular expression for removing internal spaces inside paired '
var regnos, _ = regexp.Compile(`(\s'\s.*\s'\s)`)

// left internal space from beginning of string
var regleft, _ = regexp.Compile(`^\s'\s`)

// right internal space from beginning of string
var regright, _ = regexp.Compile(`\s'\s$`)

// remove internal spaces of paired ' bla hla ' matches
func noSpacePairs(s string) string {
	havematch := true // some rule still was found inside string
	// looped trying to find the rule inside string
	for havematch {
		if regnos.MatchString(s) {
			match := regnos.FindString(s) // get the match
			// formatted version of match for replacement
			cutted := regleft.ReplaceAllString(match, " '")
			cutted = regright.ReplaceAllString(cutted, "' ")
			s = strings.Replace(s, match, cutted, 1)
			// havematch = false
		} else {
			havematch = false
		}
	}
	return s
}

// add n after founded a
func aan(s string) string {
	xs := sxs(s)
	nxs := []string{}
	nxs = append(nxs, xs[0:2]...)
	nxs = append(nxs, "n")
	nxs = append(nxs, xs[2:]...)
	return strings.Join(nxs, "")
}

// regular expression to find 'a vowel'
var regan, _ = regexp.Compile(`\s(a|A)\s(a|A|e|E|i|I|o|O|u|U|å|Å|ö|Ö|ä|Ä|h|H)`)

// very strange rule for h
func anHorse(s string) string {
	havematch := true // default value of for
	for havematch {
		if regan.MatchString(s) { // there is still present substring inside
			match := regan.FindString(s)
			s = strings.Replace(s, match, aan(match), 1)
		} else {
			havematch = false
		}
	}
	return s
}

// string modification according to the task rules
func smod(s string) string {
	// change text formatting according to rules with statements
	s = ruleManager(s)

	// remove x6 punctuation left space before, to join them to the word
	s = x6join(s)

	// remove paired punctuation mark internal spaces
	s = noSpacePairs(s)

	// 'a ' to 'an ' before words which begins from vowels and h (an house, an horse, an hiphop ... okaaay)
	s = anHorse(s)

	// trim spaces, added earlier at start and end of string
	s = strings.TrimSpace(s)

	return s
}

// print every array slice element on new string
func myarprint(ars []string) {
	fmt.Println(len(ars))
	for _, s := range ars {
		fmt.Println(s)
	}
}

func myprint(str string, ok bool) {
	if ok {
		fmt.Println(str)
	}
}

// input data looks good, continue calculation
func good(ar []string) {
	// read input data from file
	raw, _ := rfile(ar[0])
	// myprint("\nraw\n"+raw+"\n", ok)

	pfs := preformatSpaces(raw)
	// myprint(pfs, true)

	//modify data
	mods := smod(pfs)
	// myprint(mods, true)
	if unittest { // write for 'testing'
		out.Write([]byte(mods))
	}

	// write to output file
	wfile(ar[1], mods)
}
