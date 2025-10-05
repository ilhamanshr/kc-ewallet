package log_color

import "log"

func PrintRed(str string) {
	log.Println(Red(str))
}

func PrintRedf(format string, args ...interface{}) {
	log.Println(Redf(format, args...))
}

func PrintYellow(str string) {
	log.Println(Yellow(str))
}

func PrintYellowf(format string, args ...interface{}) {
	log.Println(Yellowf(format, args...))
}

func PrintMagenta(str string) {
	log.Println(Magenta(str))
}

func PrintMagentaf(format string, args ...interface{}) {
	log.Println(Magentaf(format, args...))
}

func PrintGreen(str string) {
	log.Println(Green(str))
}

func PrintGreenf(format string, args ...interface{}) {
	log.Println(Greenf(format, args...))
}
