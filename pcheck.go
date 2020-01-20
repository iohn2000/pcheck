package main

import (
	"bufio"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func main() {
	fmt.Printf("Pwned Passwords\nPwned Passwords are 555,278,657 real world passwords previously exposed in data breaches.\nLast update July 2019\n")

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter password to check if exposed: ")
	text, _ := reader.ReadString('\r')
	text = strings.Replace(text, "\r", "", -1)

	h := sha1.New()
	io.WriteString(h, text)

	bs := h.Sum(nil)

	hashString := hex.EncodeToString(bs[:])

	firstFive := hashString[0:5]
	theRest := strings.ToUpper(hashString[5:])

	resp, getError := http.Get("https://api.pwnedpasswords.com/range/" + firstFive)
	if getError != nil {
		panic(getError)
	}
	defer resp.Body.Close()
	body, readError := ioutil.ReadAll(resp.Body)
	if readError != nil {
		panic(readError)
	}

	response := string(body)

	splitResponse := strings.Split(response, "\r\n")

	found := false
	for _, value := range splitResponse {
		idx := strings.Index(value, theRest)
		if idx > -1 {
			hashcount := strings.Split(value, ":")
			fmt.Printf("\npassword '%s' exposed %v times! Change it!\n\n", text, hashcount[1])
			found = true
			break
		}
	}

	if !found {
		fmt.Println("Lucky you, this password has not been found.")
	}
}
