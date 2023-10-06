package main

/*
[Simple Banking System - Stage 3/4: I'm so lite](https://hyperskill.org/projects/93/stages/517/implement)
-------------------------------------------------------------------------------
[Errors](https://hyperskill.org/learn/step/16774)
[Declaring GORM Models](https://hyperskill.org/learn/step/28639)
[Migrations](https://hyperskill.org/learn/step/22043)
[CRUD Operations — Create](https://hyperskill.org/learn/step/22859)
[CRUD Operations — Read](https://hyperskill.org/learn/step/24151)
*/

import (
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"math"
	"math/rand"
)

// DatabaseName is the name of the database file
const DatabaseName = "card.s3db"

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
	WrongOptionMsg      = "Wrong option!"
	LoggedInMsg         = "You have successfully logged in!"
	LoggedOutMsg        = "You have successfully logged out!"
	GoodbyeMsg          = "Bye!"

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
	ID      uint   `gorm:"primaryKey"`
	Number  string `gorm:"unique;not null"`
	PIN     string
	Balance int `gorm:"default:0"`
}

func (Card) TableName() string {
	return "card"
}

type BankingSystem struct {
	db *gorm.DB
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
	card := Card{Number: cardNumber, PIN: pin}
	result := bs.db.Create(&card)
	if result.Error != nil {
		log.Printf("cannot create card: %v\n", result.Error)
		return
	}

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

	var card Card
	result := bs.db.Where("number = ? AND pin = ?", cardNumber, pin).First(&card)
	if result.Error != nil {
		fmt.Println("\n" + WrongCredentialsMsg)
		return false
	}

	fmt.Println("\n" + LoggedInMsg)
	exit := bs.HandleAccountOperations()

	return exit
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

func NewBankingSystem(db *gorm.DB) (*BankingSystem, error) {
	if !db.Migrator().HasTable(&Card{}) {
		err := db.Migrator().CreateTable(&Card{})
		if err != nil {
			return nil, fmt.Errorf("failed to create `card` table: %w", err)
		}
	}

	return &BankingSystem{
		db: db,
	}, nil
}

func main() {
	db, err := gorm.Open(sqlite.Open(DatabaseName), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to open %s: %v", DatabaseName, err)
	}

	bs, err := NewBankingSystem(db)
	if err != nil {
		log.Fatalf("failed to initialize the Banking System application: %v", err)
	}

	bs.Start()
}
