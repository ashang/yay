package main

import (
	"fmt"
	"os"

	"github.com/jguer/yay"
)

func usage() {
	fmt.Println(`usage:  yay <operation> [...]
    operations:
    yay {-h --help}
    yay {-V --version}
    yay {-D --database} <options> <package(s)>
    yay {-F --files}    [options] [package(s)]
    yay {-Q --query}    [options] [package(s)]
    yay {-R --remove}   [options] <package(s)>
    yay {-S --sync}     [options] [package(s)]
    yay {-T --deptest}  [options] [package(s)]
    yay {-U --upgrade}  [options] <file(s)>

    New operations:
    yay -Qstats   displays system information

    New options:
    --topdown     shows repository's packages first and then aur's
    --downtop     shows aur's packages first and then repository's
    --noconfirm   skip user input on package install
`)
}

var version = "1.82"

func parser() (op string, options []string, packages []string, err error) {
	if len(os.Args) < 2 {
		err = fmt.Errorf("no operation specified")
		return
	}

	for _, arg := range os.Args[1:] {
		if arg[0] == '-' && arg[1] != '-' {
			op = arg
		}

		if arg[0] == '-' && arg[1] == '-' {
			if arg == "--help" {
				op = arg
			} else if arg == "--topdown" {
				yay.SortMode = yay.TopDown
			} else if arg == "--downtop" {
				yay.SortMode = yay.DownTop
			} else if arg == "--noconfirm" {
				yay.NoConfirm = true
				options = append(options, arg)
			} else {
				options = append(options, arg)
			}
		}

		if arg[0] != '-' {
			packages = append(packages, arg)
		}
	}

	if op == "" {
		op = "yogurt"
	}

	return
}

func main() {
	op, options, pkgs, err := parser()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	yay.Config()

	switch op {
	case "-Cd":
		err = yay.CleanDependencies(pkgs)
	case "-Qstats":
		err = yay.LocalStatistics(version)
	case "-Ss":
		for _, pkg := range pkgs {
			err = yay.Search(pkg)
		}
	case "-S":
		err = yay.Install(pkgs, options)
	case "-Syu", "-Suy":
		err = yay.Upgrade(options)
	case "-Si":
		err = yay.SingleSearch(pkgs, options)
	case "yogurt":
		for _, pkg := range pkgs {
			err = yay.NumberMenu(pkg, options)
			break
		}
	case "--help", "-h":
		usage()
	default:
		err = yay.PassToPacman(op, pkgs, options)
	}

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
