package utils

import ("log")

func CheckFatalError(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}