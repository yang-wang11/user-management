package common

import "log"

func PanicIfErrorIsNotNil(err error) {
	if err != nil {
		log.Fatalln(err.Error())
	}
}
