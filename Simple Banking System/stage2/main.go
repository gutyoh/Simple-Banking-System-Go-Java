package main

/*
[Simple Banking System - Stage 2/4: Luhn algorithm](https://hyperskill.org/projects/93/stages/516/implement)
-------------------------------------------------------------------------------
[Public and private scopes](https://hyperskill.org/learn/step/17514)
[Advanced usage of structs](https://hyperskill.org/learn/step/17498)
[Methods](https://hyperskill.org/learn/step/17739)
[Debugging Go code in GoLand](https://hyperskill.org/learn/step/23118)
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
	CardBaseDigits   = 9
	PinDigits        = 4
	LuhnAlgorithmMax = 9
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

func generateLuhnChecksumDigit(number string) int {
	sum := 0

	for i, char := range number {
		digit := int(char - '0')

		if (len(number)-i)%2 == 0 {
			digit *= 2
			if digit > LuhnAlgorithmMax {
				digit -= LuhnAlgorithmMax
			}
		}

		sum += digit
	}

	return (10 - sum%10) % 10
}

type Card struct {
	Number string
	PIN    string
}

type BankingSystem struct {
	Cards []Card
}

func (bs *BankingSystem) Start() {
	var stop bool
	for !stop {
		stop = bs.HandleMainMenuOperations()
	}
}

func (bs *BankingSystem) HandleMainMenuOperations() bool {
	for {
		bs.DisplayMainMenu()

		var choice int
		fmt.Scanln(&choice)

		switch choice {
		case 1:
			bs.CreateAccount()
		case 2:
			if bs.Login() {
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

func (*BankingSystem) DisplayMainMenu() {
	fmt.Println(MainMenuCreateAccount)
	fmt.Println(MainMenuLogin)
	fmt.Println(MenuExit)
}

func (bs *BankingSystem) CreateAccount() {
	cardNumber, pin := bs.GenerateCardNumberAndPIN()
	bs.Cards = append(bs.Cards, Card{cardNumber, pin})

	fmt.Println("\n" + CardCreatedMsg)
	fmt.Printf(CardNumberMsg, cardNumber)
	fmt.Printf(CardPINMsg, pin)
}

func (*BankingSystem) GenerateCardNumberAndPIN() (string, string) {
	cardBase := CardPrefix + generateRandomDigits(CardBaseDigits)
	checksum := generateLuhnChecksumDigit(cardBase)
	cardNumber := cardBase + fmt.Sprintf("%d", checksum)
	pin := generateRandomDigits(PinDigits)

	return cardNumber, pin
}

func generateRandomDigits(n int) string {
	maxNumber := int(math.Pow(10, float64(n)))
	return fmt.Sprintf("%0*d", n, rand.Intn(maxNumber))
}

func (bs *BankingSystem) PromptLoginCredentials() (string, string) {
	fmt.Println("\n" + CardNumberPrompt)
	var cardNumber string
	fmt.Scanln(&cardNumber)

	fmt.Println(PINPrompt)
	var pin string
	fmt.Scanln(&pin)

	return cardNumber, pin
}

func (bs *BankingSystem) Login() bool {
	cardNumber, pin := bs.PromptLoginCredentials()

	for _, c := range bs.Cards {
		if c.Number == cardNumber && c.PIN == pin {
			fmt.Println("\n" + LoggedInMsg)
			return bs.HandleAccountOperations()
		}
	}

	fmt.Println("\n" + WrongCredentialsMsg)
	return false
}

func (bs *BankingSystem) HandleAccountOperations() bool {
	for {
		bs.DisplayAccountOperationsMenu()

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

func (*BankingSystem) DisplayAccountOperationsMenu() {
	fmt.Println("\n" + AccountOperationsBalance)
	fmt.Println(AccountOperationsLogout)
	fmt.Println(MenuExit)
}

func NewBankingSystem() *BankingSystem {
	return &BankingSystem{
		Cards: make([]Card, 0),
	}
}

func main() {
	bs := NewBankingSystem()
	bs.Start()
}
