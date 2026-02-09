package main

import (
	"fmt"
	"log"
)

func main() {

	numDec := 92

	numHex, err := convertDecToHex(numDec)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("num hex: ", numHex)

}

func convertDecToHex(numDec int) (string, error) {

	if numDec < 0 {
		return "", fmt.Errorf("negative numbers not supported")
	}

	if numDec == 0 {
		return "0", nil
	}

	quotient := numDec
	remainder := 0
	res := ""

	for quotient != 0 {
		remainder = quotient % 16
		quotient /= 16
		if remainder < 10 {
			res += string(rune('0' + remainder))
		} else {
			res += string(rune('A' + (remainder - 10)))
		}
	}

	resLast := ""

	for i := len(res) - 1; i >= 0; i-- {
		resLast += string(res[i])
	}

	return resLast, nil
}
