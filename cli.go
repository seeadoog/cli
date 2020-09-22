package cli

import (
	"flag"
	"fmt"
	"os"
)

type FlagSet struct {
	Args  []string
	fg    *flag.FlagSet
	Name  string
	Usage func()
}

func (f *FlagSet) StringVal(p *string, name, val, usage string) {
	f.fg.StringVar(p, name, val, usage)
}

func (f *FlagSet) BoolVal(p *bool, name string, val bool, usage string) {
	f.fg.BoolVar(p, name, val, usage)
}

func (f *FlagSet) IntVal(p *int, name string, val int, usage string) {
	f.fg.IntVar(p, name, val, usage)
}

func (f *FlagSet) FloatVal(p *float64, name string, val float64, usage string) {
	f.fg.Float64Var(p, name, val, usage)
}

func (f *FlagSet) Parse()  {
	if err := f.fg.Parse(f.Args); err != nil {
		panic(err)
	}
	return
}

func (f *FlagSet)Arg(i int)string{
	return f.fg.Arg(i)
}

func (f *FlagSet)RunArgs()[]string{
	return f.fg.Args()
}

type Command struct {
	Usage   func()
	Runable func(flag *FlagSet)
}

type App struct {
	Args     []string
	Commands map[string]*Command
	Runable  func(flag *FlagSet)
	Usage    func()
}



func (a *App) Run() {
	if len(a.Args) == 0 {
		a.Args = os.Args
	}
	name := ""
	if len(a.Args) >= 2 {
		name = os.Args[1]
	}
	cmd := a.Commands[name]
	if cmd == nil || name == "" {
		if a.Runable != nil {
			fg := flag.NewFlagSet(a.Args[0], flag.ExitOnError)
			fg.Usage = a.Usage
			a.Runable(&FlagSet{
				Args: a.Args[1:],
				fg:   fg,
				Name: a.Args[0],
			})
			return
		} else {
			fmt.Println("commands are:")
			for k, v := range a.Commands {
				fmt.Println("\t",k)
				if v.Usage != nil{
					v.Usage()
				}
				fmt.Println()
			}
			return
		}
	}
	fg := flag.NewFlagSet(a.Args[1], flag.ExitOnError)
	fg.Usage = cmd.Usage
	cmd.Runable(&FlagSet{
		Args: a.Args[2:],
		fg:   fg,
		Name: a.Args[1],
	})
}
