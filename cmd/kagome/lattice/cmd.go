// Copyright 2015 ikawaha
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// 	You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package lattice

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/ikawaha/kagome/tokenizer"
)

// subcommand property
var (
	CommandName  = "lattice"
	Description  = `lattice viewer`
	UsageMessage = "%s [-udic userdic_file] [-output output_file] [-v] sentence\n"
	ErrorWriter  = os.Stderr
)

// options
type option struct {
	udic    string
	output  string
	verbose bool
	input   string
	flagSet *flag.FlagSet
}

func newOption() (o *option) {
	o = &option{
		flagSet: flag.NewFlagSet(CommandName, flag.ContinueOnError),
	}
	// option settings
	o.flagSet.StringVar(&o.udic, "udic", "", "user dic")
	o.flagSet.StringVar(&o.output, "output", "", "output file")
	o.flagSet.BoolVar(&o.verbose, "v", false, "verbose mode")
	return
}

func (o *option) parse(args []string) (err error) {
	if err = o.flagSet.Parse(args); err != nil {
		return
	}
	// validations
	if o.flagSet.NArg() == 0 {
		return fmt.Errorf("input is empty")
	}
	o.input = strings.Join(o.flagSet.Args(), " ")
	return
}

// command main
func command(opt *option) error {
	t := tokenizer.New()
	var out = os.Stdout
	if opt.output != "" {
		var err error
		out, err = os.OpenFile(opt.output, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0666)
		if err != nil {
			fmt.Fprintln(ErrorWriter, err)
			os.Exit(1)
		}
		defer out.Close()
	}
	var udic tokenizer.UserDic
	if opt.udic != "" {
		var err error
		udic, err = tokenizer.NewUserDic(opt.udic)
		if err != nil {
			return err
		}
		t.SetUserDic(udic)
	}
	if opt.udic != "" {
		if udic, err := tokenizer.NewUserDic(opt.udic); err != nil {
			fmt.Fprintln(ErrorWriter, err)
			os.Exit(1)
		} else {
			t.SetUserDic(udic)
		}
	}

	tokens := t.Dot(opt.input, out)
	if opt.verbose {
		for i, size := 1, len(tokens); i < size; i++ {
			tok := tokens[i]
			f := tok.Features()
			if tok.Class == tokenizer.DUMMY {
				fmt.Fprintf(ErrorWriter, "%s\n", tok.Surface)
			} else {

				fmt.Fprintf(ErrorWriter, "%s\t%v\n", tok.Surface, strings.Join(f, ","))
			}
		}
	}
	return nil
}

func Run(args []string) error {
	opt := newOption()
	if e := opt.parse(args); e != nil {
		Usage()
		PrintDefaults()
		return e
	}
	return command(opt)
}

func Usage() {
	fmt.Fprintf(ErrorWriter, UsageMessage, CommandName)
}

func PrintDefaults() {
	o := newOption()
	o.flagSet.PrintDefaults()
}
