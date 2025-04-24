package main

import (
	"strings"
)

func ValidString(str string) string {
	//what replaces bad words
	const replacement = "****"

	//Words not allowed
	badWords := []string{"kerfuffle", "sharbert", "fornax"}

	strList := strings.Split(str, " ")
	for index, val := range strList {
		for _, bad := range badWords {
			if strings.EqualFold(val, bad) {
				strList[index] = replacement
				break
			}
		}
	}
	ans := strings.Join(strList, " ")

	return ans
}
