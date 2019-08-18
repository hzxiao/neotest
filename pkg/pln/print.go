package pln

import "github.com/fatih/color"

var Verbose bool

func Error(err error) {
	if err != nil {
		color.Red(err.Error())
	}
}

func InfoVerbose(format string, a ...interface{}) {
	if Verbose {
		color.Blue(format, a...)
	}
}

func InfoSuccess(format string, a ...interface{}) {
	color.Green(format, a...)
}

func InfoFail(format string, a ...interface{}) {
	color.Red(format, a...)
}

func Warn(format string, a ...interface{}) {
	color.Yellow(format, a...)
}
