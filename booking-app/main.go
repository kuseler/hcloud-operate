package main

import (
	"booking-app/helper"
	"fmt"
	"sync"
	"time"
)

const conferenceTickets int = 50

var conferenceName = "Go Conference"
var remainingTickets uint = 50
var bookings = make([]UserData, 0)

type UserData struct {
	firstName       string
	lastName        string
	email           string
	numberOfTickets uint
}

var wg = sync.WaitGroup{}

func main() {

	helper.GreetUsers(conferenceName, conferenceTickets, remainingTickets)

	firstName, lastName, email, userTickets := helper.GetUserInput()

	isValidName, isValidEmail, isValidTicketNumber := helper.ValidateUserInput(firstName, lastName, email, userTickets, remainingTickets)

	if isValidName && isValidEmail && isValidTicketNumber {

		bookTicket(userTickets, firstName, lastName, email)
		wg.Add(1)
		go sendTicket(userTickets, firstName, lastName, email)

		firstNames := GetFirstNames()
		fmt.Printf("%v tickets remaining for %v\n", remainingTickets, conferenceName)
		fmt.Printf("The first names of bookings are: %v\n", firstNames)
		fmt.Printf("List of bookings is %v\n", bookings)
		println("---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------")

		if remainingTickets == 0 {
			//end program
			fmt.Println("Our conference is booked out. Come back next year.")
			//break
		}
	} else {
		if !isValidName {
			fmt.Println("First Name or last Name you entered is too short.")
		}
		if !isValidEmail {
			fmt.Println("The Email Address you entered doesn't contain @ sign.")
		}
		if !isValidTicketNumber {
			fmt.Println("The Number of Tickets you entered is invalid.")
		}
	}
	wg.Wait()
}

func GetFirstNames() []string {
	firstNames := []string{}
	for _, booking := range bookings {
		firstNames = append(firstNames, booking.firstName)
	}
	return firstNames
}

func bookTicket(userTickets uint, firstName string, lastName string, email string) {
	remainingTickets = remainingTickets - userTickets

	var userData = UserData{
		firstName:       firstName,
		lastName:        lastName,
		email:           email,
		numberOfTickets: userTickets,
	}

	bookings = append(bookings, userData)

	fmt.Printf("Thank you %v %v for booking %v tickets. You will receive a confirmation email at %v\n", firstName, lastName, userTickets, email)

}

func sendTicket(userTickets uint, firstName string, lastName string, email string) {
	time.Sleep(10 * time.Second)
	var ticket = fmt.Sprintf("%v tickets for %v %v", userTickets, firstName, lastName)
	fmt.Println("#################")
	fmt.Printf("Sending ticket: %v\nto email adress: %v\n", ticket, email)
	fmt.Println("#################")
	wg.Done()
}

/* Beispiel für Switch:

city := "London"

switch city {
	case "New York":
		//ececute code for booking New York conference tickets
	case "Singapore", "London":
		//ececute code for booking Singapore and London conference tickets
	case "Berlin", "Hong Kong":
		//ececute code for booking Berlin and Hong Kong conference tickets
	default:
		fmt.Print("No valid city selected") */
