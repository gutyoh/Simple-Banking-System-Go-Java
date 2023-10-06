package main

/*
[Simple Banking System - Stage 1/4: Card anatomy](https://hyperskill.org/projects/93/stages/515/implement)
-------------------------------------------------------------------------------
[String formatting](https://hyperskill.org/learn/step/16860)
[Control statements](https://hyperskill.org/learn/step/16235)
[Math package](https://hyperskill.org/learn/step/18431)
[Type conversion and overflow](https://hyperskill.org/learn/step/18710)
[Generating random numbers and strings](https://hyperskill.org/learn/step/28497)
[Working with slices](https://hyperservices.herokuapp.com/graph/prereq?topic=1701)
[Operations with maps](https://hyperservices.herokuapp.com/graph/prereq?topic=1850)
[Functional decomposition](https://hyperskill.org/learn/step/17506)
[Debugging Go code](https://hyperskill.org/learn/step/23076)
*/

import (
	"fmt"
	"math"
	"math/rand"
)

// Main menu options
const (
	MainMenuCreateAccount = "1. Create an account"
	MainMenuLogin         = "2. Log into account"
	MenuExit              = "0. Exit"
)

// Account operations options
const (
	AccountOperationsBalance = "1. Balance"
	AccountOperationsLogout  = "2. Log out"
)

// Banking system prompts
const (
	CardNumberPrompt = "Enter your card number:"
	PINPrompt        = "Enter your PIN:"
)

// Digit constants
const (
	CardBaseDigits = 10
	PinDigits      = 4
)

// Banking system messages
const (
	WrongCredentialsMsg = "Wrong card number or PIN"

	WrongOptionMsg = "Wrong option!"
	LoggedInMsg    = "You have successfully logged in!"
	LoggedOutMsg   = "You have successfully logged out!"
	GoodbyeMsg     = "Bye!"

	CardPrefix     = "400000"
	CardCreatedMsg = "Your card has been created"
	CardNumberMsg  = "Your card number:\n%s\n"
	CardPINMsg     = "Your card PIN:\n%s\n\n"
	BalanceMsg     = "Balance: %d"
)

func startBankingSystem(cards map[string]string) {
	var stop bool
	for !stop {
		stop = handleMainMenuOperations(cards)
	}
}

func handleMainMenuOperations(cards map[string]string) bool {
	for {
		displayMainMenu()

		var choice int
		fmt.Scanln(&choice)

		switch choice {
		case 1:
			createAccount(cards)
		case 2:
			if login(cards) {
				fmt.Println("\n" + GoodbyeMsg)
				return true
			}
		case 0:
			fmt.Println("\n" + GoodbyeMsg)
			return true
		default:
			fmt.Println("\n" + WrongOptionMsg)
		}
	}
}

func displayMainMenu() {
	fmt.Println(MainMenuCreateAccount)
	fmt.Println(MainMenuLogin)
	fmt.Println(MenuExit)
}

func createAccount(cards map[string]string) {
	card, pin := generateCardNumberAndPIN()
	cards[card] = pin
	fmt.Println("\n" + CardCreatedMsg)
	fmt.Printf(CardNumberMsg, card)
	fmt.Printf(CardPINMsg, pin)
}

func generateCardNumberAndPIN() (string, string) {
	cardNumber := CardPrefix + generateRandomDigits(CardBaseDigits)
	pin := generateRandomDigits(PinDigits)

	return cardNumber, pin
}

func generateRandomDigits(n int) string {
	maxNumber := int(math.Pow(10, float64(n)))
	return fmt.Sprintf("%0*d", n, rand.Intn(maxNumber))
}

func promptLoginCredentials() (string, string) {
	fmt.Println("\n" + CardNumberPrompt)
	var cardNumber string
	fmt.Scanln(&cardNumber)

	fmt.Println(PINPrompt)
	var pin string
	fmt.Scanln(&pin)

	return cardNumber, pin
}

func login(cards map[string]string) bool {
	cardNumber, pin := promptLoginCredentials()

	if storedPin, exists := cards[cardNumber]; exists && storedPin == pin {
		fmt.Println("\n" + LoggedInMsg)
		return handleAccountOperations()
	}

	fmt.Println("\n" + WrongCredentialsMsg)
	return false
}

func handleAccountOperations() bool {
	for {
		displayAccountOperationsMenu()

		var choice int
		fmt.Scanln(&choice)

		switch choice {
		case 1:
			fmt.Printf("\n"+BalanceMsg+"\n", 0)
		case 2:
			fmt.Println("\n" + LoggedOutMsg)
			return false
		case 0:
			return true
		default:
			fmt.Println("\n" + WrongOptionMsg)
		}
	}
}

func displayAccountOperationsMenu() {
	fmt.Println("\n" + AccountOperationsBalance)
	fmt.Println(AccountOperationsLogout)
	fmt.Println(MenuExit)
}

func main() {
	cards := make(map[string]string)
	startBankingSystem(cards)
}
