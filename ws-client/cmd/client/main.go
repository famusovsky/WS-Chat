package main

import (
	"bufio"
	"fmt"
	"os"
	"ws-client/internal/app"
)

func main() {
	var nickname string
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Printf("Enter your nickname:\n")
		scanner.Scan()
		nickname = scanner.Text()
		if nickname == "" {
			fmt.Println("Nickname cannot be empty")
		} else {
			break
		}
	}

	fmt.Printf("Your nickname: %s\n", nickname)

	for {
		fmt.Printf("Enter host adress:\n")
		scanner.Scan()
		err := app.Run(scanner.Text(), nickname, scanner)
		if err != nil {
			fmt.Printf("Error occured while connecting to the server:\n%s\nTrying to reconnect...\n", err)
		} else {
			break
		}
	}
}
