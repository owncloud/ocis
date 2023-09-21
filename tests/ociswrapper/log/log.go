package log

import "log"

func Println(message string) {
	log.Println("[ociswrapper]", message)
}

func Panic(err error) {
	log.Panic("[ociswrapper]", err.Error())
}

func Fatalln(err error) {
	log.Fatalln("[ociswrapper]", err.Error())
}
