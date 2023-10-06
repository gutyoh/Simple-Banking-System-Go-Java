package main

/*
[Simple Banking System - Stage 4/4: Advanced system](https://hyperskill.org/projects/93/stages/518/implement)
-------------------------------------------------------------------------------
[CRUD Operations — Update](https://hyperskill.org/learn/step/33258)
[CRUD Operations — Delete](https://hyperskill.org/learn/step/31914)
[Transactions](https://hyperskill.org/learn/step/35352)
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
	AccountOperationsBalance      = "1. Balance"
	AccountOperationsAddIncome    = "2. Add income"
	AccountOperationsDoTransfer   = "3. Do transfer"
	AccountOperationsCloseAccount = "4. Close account"
	AccountOperationsLogout       = "5. Log out"
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

	CardPrefix      = "400000"
	CardCreatedMsg  = "Your card has been created"
	CardNumberMsg   = "Your card number:\n%s\n"
	CardPINMsg      = "Your card PIN:\n%s\n\n"
	BalanceMsg      = "Balance: %d"
	IncomePrompt    = "Enter income:"
	TransferPrompt  = "Transfer\nEnter card number:"
	CloseAccountMsg = "The account has been closed!"

	IncomeAddedMsg = "Income was added!"

	CardNotFoundMsg = "Such a card does not exist."

	NotEnoughMoneyMsg = "Not enough money!"

	TransferSuccessfulMsg = "Transfer successful!"

	TransferFailedMsg = "Transfer failed."

	TransferToSameAccountMsg = "You can't transfer money to the same account!"

	TransferToInvalidAccountMsg = "Probably you made a mistake in the card number. Please try again!"

	TransferAmountPrompt = "Enter how much money you want to transfer:"
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

	tx := bs.db.Begin()
	result := tx.Create(&card)
	if result.Error != nil {
		log.Printf("cannot create card: %v\n", result.Error)
		tx.Rollback()
		return
	}
	tx.Commit()

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
	exit := bs.HandleAccountOperations(&card)

	return exit
}

func (bs *BankingSystem) HandleAccountOperations(card *Card) bool {
	for {
		bs.DisplayAccountOperationsMenu()

		var choice int
		fmt.Scanln(&choice)

		switch choice {
		case 1:
			fmt.Printf("\n"+BalanceMsg+"\n", card.Balance)
		case 2:
			bs.AddIncome(card)
		case 3:
			bs.InitiateTransfer(card)
		case 4:
			bs.CloseAccount(card)
		case 5:
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
	fmt.Println(AccountOperationsAddIncome)
	fmt.Println(AccountOperationsDoTransfer)
	fmt.Println(AccountOperationsCloseAccount)
	fmt.Println(AccountOperationsLogout)
	fmt.Println(MenuExit)
}

func (bs *BankingSystem) AddIncome(card *Card) {
	fmt.Println(IncomePrompt)
	var income int
	fmt.Scanln(&income)

	card.Balance += income

	tx := bs.db.Begin()
	result := tx.Save(&card)
	if result.Error != nil {
		log.Printf("cannot update balance: %v\n", result.Error)
		tx.Rollback()
		return
	}
	tx.Commit()

	fmt.Println(IncomeAddedMsg)
}

func (bs *BankingSystem) InitiateTransfer(senderCard *Card) {
	recipientCardNumber := bs.PromptForRecipientCardNumber()

	if !bs.CanTransferBetweenCards(senderCard, recipientCardNumber) {
		return
	}

	recipientCard, err := bs.GetCard(recipientCardNumber)
	if err != nil {
		fmt.Println(CardNotFoundMsg)
		return
	}

	transferAmount := bs.PromptForTransferAmount()
	if transferAmount <= 0 || senderCard.Balance < transferAmount {
		fmt.Println(NotEnoughMoneyMsg)
		return
	}

	if bs.ExecuteTransfer(senderCard, recipientCard, transferAmount) {
		fmt.Println(TransferSuccessfulMsg)
	} else {
		fmt.Println(TransferFailedMsg)
	}
}

func (*BankingSystem) PromptForRecipientCardNumber() string {
	fmt.Println(TransferPrompt)
	var recipientCardNumber string
	fmt.Scanln(&recipientCardNumber)
	return recipientCardNumber
}

func (*BankingSystem) CanTransferBetweenCards(senderCard *Card, recipientCardNumber string) bool {
	if senderCard.Number == recipientCardNumber {
		fmt.Println(TransferToSameAccountMsg)
		return false
	}

	if !luhnAlgorithm(recipientCardNumber) {
		fmt.Println(TransferToInvalidAccountMsg)
		return false
	}

	return true
}

func (bs *BankingSystem) ExecuteTransfer(sender *Card, recipient *Card, amount int) bool {
	tx := bs.db.Begin()
	result := tx.Model(sender).Update("balance", gorm.Expr("balance - ?", amount))
	if result.Error != nil {
		log.Printf("cannot update sender balance: %v\n", result.Error)
		tx.Rollback()
		return false
	}

	result = tx.Model(recipient).Update("balance", gorm.Expr("balance + ?", amount))
	if result.Error != nil {
		log.Printf("cannot update recipient balance: %v\n", result.Error)
		tx.Rollback()
		return false
	}
	tx.Commit()

	return true
}

func (bs *BankingSystem) GetCard(cardNumber string) (*Card, error) {
	var card Card
	result := bs.db.Where("number = ?", cardNumber).First(&card)
	if result.Error != nil {
		return nil, result.Error
	}
	return &card, nil
}

func (*BankingSystem) PromptForTransferAmount() int {
	fmt.Println(TransferAmountPrompt)
	var amount int
	fmt.Scanln(&amount)
	return amount
}

func (bs *BankingSystem) CloseAccount(card *Card) {
	tx := bs.db.Begin()
	result := tx.Delete(&card)
	if result.Error != nil {
		fmt.Printf("cannot delete card: %v\n", result.Error)
		tx.Rollback()
		return
	}
	tx.Commit()

	fmt.Println(CloseAccountMsg)
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
