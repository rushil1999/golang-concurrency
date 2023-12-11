/*
The Sleeping Barber Problem is a classical synchronization problem in which a barber shop with one barber, a waiting room, and a number of customers is simulated. The problem involves coordinating the access to the waiting room and the barber chair so that only one customer is in the chair at a time and the barber is always working on a customer if there is one in the chair, otherwise the barber is sleeping until a customer arrives
*/

package main

import (
	"fmt"
	"sync"
	"time"
	"math/rand"
)

// Barber interface to specify methods that need to be caled
type barberInterface interface {
	operate(<-chan Customer)
	sleep()
}

type Customer struct {
	name string
	id int
}

type Barber struct {
	name string
}

const haircutTime = 2*time.Second // Time need to cut hair of a customer

// This function will generate customers with some random delay
func generateCustomer(id int) Customer{
	customerName := fmt.Sprintf("customer_%d", id+1)
	newCustomer := Customer{
		name: customerName,
		id: id+1,
	}
	delay := rand.Intn(10)
	time.Sleep(time.Duration(delay)*time.Second)
	return newCustomer
}

// This function is to push the generated customers to designated channel
func generateCustomers(waitLine chan<- Customer) {
	customerId :=0 
	for {
		newCustomer := generateCustomer(customerId) 

		select {
		case waitLine <- newCustomer: // For non blocking send
			fmt.Printf("%s is waiting\n", newCustomer.name)
		}

		customerId = newCustomer.id
	}
}

// Function attached to Barber struct to map real-life entity
func (b *Barber)sleep(){
	fmt.Printf("Barber is sleeping !!\n")
}

// Function to start the barber go routine
func (b *Barber)operate(waitLine <-chan Customer, wg *sync.WaitGroup) {
	for {
		b.sleep()
		select {
			case newCustomer, more := <- waitLine:
				fmt.Printf("Cutting hair of customer %s\n", newCustomer.name)
				time.Sleep(haircutTime) 
				fmt.Printf("FINISHED Cutting hair of customer %s!\n", newCustomer.name)
				wg.Done()
				if !more {
					fmt.Printf("Shop closed!!!")
					return
				} 
		}
	}
}




func main(){

	barber := &Barber{name: "La Fashion"}
	waitLine := make(chan Customer, 5) //Only a certain number of customers can wait in line

	wg := sync.WaitGroup{}
	wg.Add(10) // Will have in total 10 customer before the program terminates
	go generateCustomers(waitLine)
	go barber.operate(waitLine, &wg)


	wg.Wait()

}