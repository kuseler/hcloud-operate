package helper

import (
	"fmt"
	"strings"
)

func GreetUsers(conferenceName string, conferenceTickets int, remainingTickets uint) {
	fmt.Printf("Welcome to %v booking application!\n", conferenceName)
	fmt.Printf("We have total of %v tickets and %v are still available.\n", conferenceTickets, remainingTickets)
	fmt.Println("Get your tickets here to attend.")
}

func GetUserInput() (string, string, string, uint) {
	var firstName string
	var lastName string
	var email string
	var userTickets uint

	//ask user for their name
	fmt.Println("Enter you first Name:")
	fmt.Scan(&firstName)

	fmt.Println("Enter you last Name:")
	fmt.Scan(&lastName)

	fmt.Println("Enter your email adress:")
	fmt.Scan(&email)

	fmt.Println("Enter number of tickets:")
	fmt.Scan(&userTickets)

	return firstName, lastName, email, userTickets
}

func ValidateUserInput(firstName string, lastName string, email string, userTickets uint, remainingTickets uint) (bool, bool, bool) {
	isValidName := len(firstName) >= 2 && len(lastName) >= 2
	isValidEmail := strings.Contains(email, "@")
	isValidTicketNumber := userTickets > 0 && userTickets <= remainingTickets
	return isValidName, isValidEmail, isValidTicketNumber
}
