// This work is subject to the CC0 1.0 Universal (CC0 1.0) Public Domain Dedication
// license. Its contents can be found at:
// http://creativecommons.org/publicdomain/zero/1.0/ and
// http://creativecommons.org/publicdomain/zero/1.0/legalcode

package optarg

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Option struct {
	Name        string
	ShortName   string
	Description string
	defaultval  interface{}
	value       string
}

var (
	options     = make([]*Option, 0)
	Remainder   = make([]string, 0)
	ShortSwitch = "-"
	LongSwitch  = "--"
	UsageInfo   = fmt.Sprintf("Usage: %s [options]:", os.Args[0])
	HeaderFmt   = "\n[%s]"
)

const headerName = "__hdr"

// Returns usage information in a neatly formatted string.
func UsageString() string {
	var desc, str, format string
	var ok bool
	var lines, output []string
	var i int

	offset := 0

	// Find the largest length of the option name list. Needed to align
	// the description blocks consistently.
	for _, v := range options {
		if v.ShortName == headerName {
			continue
		}
		// If there's no short name, skip it in the description
		if len(v.ShortName) == 0 {
			str = fmt.Sprintf("%s%s:", LongSwitch, v.Name)
		} else {
			str = fmt.Sprintf("%s%s, %s%s: ", LongSwitch, v.Name, ShortSwitch, v.ShortName)
		}
		if len(str) > offset {
			offset = len(str)
		}
	}

	offset++ // add margin.

	output = append(output, UsageInfo)

	for _, v := range options {
		if v.ShortName == headerName {
			output = append(output, fmt.Sprintf(HeaderFmt, v.Description))
			continue
		}

		// Print namelist. right-align it based on the maximum width
		// found in previous loop.
		if len(v.ShortName) == 0 {
			str = fmt.Sprintf("%s%s: ", LongSwitch, v.Name)
		} else {
			str = fmt.Sprintf("%s%s, %s%s: ", LongSwitch, v.Name, ShortSwitch, v.ShortName)
		}
		format = fmt.Sprintf("%%%ds", offset)
		str = fmt.Sprintf(format, str)
		desc = v.Description

		// boolean flags need no 'default value' description. They are either set or not.
		if _, ok = v.defaultval.(bool); !ok {
			if fmt.Sprintf("%v", v.defaultval) != "" {
				desc = fmt.Sprintf("%s (defaults to: %v)", desc, v.defaultval)
			}
		}

		// Format and print left-aligned, word-wrapped description with
		// a @margin left margin size using my super duper patented
		// multi-line string wrap routine (see string.go). Assume
		// maximum of 80 characters screen width. Which makes block
		// width equal to 80 - @offset. I would prefer to use
		// ALIGN_JUSTIFY for added sexy, but it looks a little funky for
		// short descriptions. So we'll stick with the establish left-
		// aligned text.
		lines = multilineWrap(desc, 80, offset, 0, _ALIGN_LEFT)

		// First line needs to be appended to where we left off.
		output = append(output, fmt.Sprintf("%s%s", str, strings.TrimSpace(lines[0])))

		// Print the rest as-is (properly indented).
		for i = 1; i < len(lines); i++ {
			output = append(output, lines[i])
		}
	}
	return strings.Join(output, "\n") + "\n"
}

// Prints usage information in a neatly formatted overview.
func Usage() {
	fmt.Print(UsageString())
	return
}

// Parse os.Args using the previously added Options.
func Parse() <-chan *Option {
	c := make(chan *Option)
	Remainder = make([]string, 0)
	go processArgs(c)
	return c
}

func processArgs(c chan *Option) {
	var opt *Option
	var ok bool
	var tok, v string

	for i := range os.Args {
		if i == 0 {
			continue
		} // skip app name

		if v = strings.TrimSpace(os.Args[i]); len(v) == 0 {
			continue
		}

		if len(v) >= 3 && v[0:2] == LongSwitch {
			if v = strings.TrimSpace(v[2:]); len(v) == 0 {
				Remainder = append(Remainder, LongSwitch)
			} else {
				if opt = findOption(v); opt == nil {
					fmt.Fprintf(os.Stderr, "Unknown option '--%s' specified.\n", v)
					Usage()
					os.Exit(1)
				}

				if _, ok = opt.defaultval.(bool); ok {
					opt.value = "true"
					c <- opt
					opt = nil
				}
			}

		} else if len(v) >= 2 && v[0:1] == ShortSwitch {
			if v = strings.TrimSpace(v[1:]); len(v) == 0 {
				Remainder = append(Remainder, ShortSwitch)
			} else {
				for i = range v {
					tok = v[i : i+1]
					opt = findOption(tok)
					if opt == nil {

						fmt.Fprintf(os.Stderr, "Unknown option '-%s' specified.\n", tok)
						Usage()
						os.Exit(1)
					}

					if _, ok := opt.defaultval.(bool); ok {
						opt.value = "true"
						c <- opt
						opt = nil
					}
				}
			}

		} else {
			if opt == nil {
				Remainder = append(Remainder, v)
			} else {
				opt.value = v
				c <- opt
				opt = nil
			}
		}
	}
	close(c)
}

// Adds a section header. Useful to separate options into discrete groups.
func Header(title string) {
	options = append(options, &Option{
		ShortName:   headerName,
		Description: title,
	})
}

// Add a new command line option to check for.
func Add(shortname, name, description string, defaultvalue interface{}) {
	options = append(options, &Option{
		ShortName:   shortname,
		Name:        name,
		Description: description,
		defaultval:  defaultvalue,
	})
}

func findOption(name string) *Option {
	for i := range options {
		if options[i].Name == name || (len(options[i].ShortName) > 0 && options[i].ShortName == name) {
			return options[i]
		}
	}
	return nil
}

func (this *Option) String() string { return this.value }

func (this *Option) Bool() bool {
	if b, err := strconv.ParseBool(this.value); err == nil {
		return b
	}
	return false
}

func (this *Option) Int() int {
	if v, err := strconv.Atoi(this.value); err == nil {
		return v
	}
	return this.defaultval.(int)
}
func (this *Option) Int8() int8   { return int8(this.Int()) }
func (this *Option) Int16() int16 { return int16(this.Int()) }
func (this *Option) Int32() int32 { return int32(this.Int()) }
func (this *Option) Int64() int64 {
	if v, err := strconv.ParseInt(this.value, 10, 64); err == nil {
		return v
	}
	return this.defaultval.(int64)
}

func (this *Option) Uint() uint {
	if v, err := strconv.ParseUint(this.value, 10, 0); err == nil {
		return uint(v)
	}
	return this.defaultval.(uint)
}
func (this *Option) Uint8() uint8   { return uint8(this.Int()) }
func (this *Option) Uint16() uint16 { return uint16(this.Int()) }
func (this *Option) Uint32() uint32 { return uint32(this.Int()) }
func (this *Option) Uint64() uint64 {
	if v, err := strconv.ParseUint(this.value, 10, 64); err == nil {
		return v
	}
	return this.defaultval.(uint64)
}

func (this *Option) Float32() float32 {
	if v, err := strconv.ParseFloat(this.value, 32); err == nil {
		return float32(v)
	}
	return this.defaultval.(float32)
}

func (this *Option) Float64() float64 {
	if v, err := strconv.ParseFloat(this.value, 64); err == nil {
		return v
	}
	return this.defaultval.(float64)
}
