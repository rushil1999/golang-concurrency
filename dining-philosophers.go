/*
The Dining Philosopher Problem states that K (5 in our case) philosophers are seated around a circular table with one chopstick between each pair of philosophers. There is one chopstick between each philosopher. A philosopher may eat if he can pick up the two chopsticks adjacent to him. One chopstick may be picked up by any one of its adjacent followers but not both.
*/


package main

import (
	"fmt"
	"sync"
	"time"
)

const eatTime = 3*time.Second
const thinkTime = 2*time.Second
const hungry = 3

type Philosopher struct {
	name string
	leftFork int
	rightFork int
}

// Function attached to Philosopher struc to map real-life entity
func (p *Philosopher) sit(sitWg *sync.WaitGroup) {
	sitWg.Done()
}

// The logic is to make sure that a philosopher picks up a form with lower index value first.
// This is to prevent logical Deadlock! 
// Logical Deadlock can happen, when all the philosophers pick up the fork on their left first.
// Mutex are used to lock the forks once they are picked
func (p *Philosopher) pickupForks(forkMap map[int]*sync.Mutex) {
	pLeftFork := p.leftFork
	pRightFork := p.rightFork
	if pLeftFork < pRightFork {
		forkMap[pLeftFork].Lock()
		fmt.Printf("\t%s picked left fork\n",p.name)
		forkMap[pRightFork].Lock()
		fmt.Printf("\t%s picked right fork\n",p.name)
	} else {
		forkMap[pRightFork].Lock()
		fmt.Printf("\t%s picked Right fork\n",p.name)
		forkMap[pLeftFork].Lock()
		fmt.Printf("\t%s picked left fork\n",p.name)
	}
}

// Mutex are used to unlock the forks once they are put down
func (p *Philosopher) putdownForks(forkMap map[int]*sync.Mutex) {
	pLeftFork := p.leftFork
	pRightFork := p.rightFork

	forkMap[pLeftFork].Unlock()
	forkMap[pRightFork].Unlock()
	fmt.Printf("\t%s put down the forks\n", p.name)
}

// Function attached to Philosopher struc to map real-life entity
func (p *Philosopher) eat() {
	fmt.Printf("%s is eating\n",p.name)
	time.Sleep(eatTime)
}

// Function attached to Philosopher struc to map real-life entity
func (p *Philosopher) think() {
	fmt.Printf("%s is Thinking\n",p.name)
	time.Sleep(thinkTime)
}

// sitWg is a wait group used to execute the process once all the philisophers have seated.
// This is a requuirement of the original problem
func (p *Philosopher) dine(forkMap map[int]*sync.Mutex, sitWg *sync.WaitGroup, wg *sync.WaitGroup) {
	defer wg.Done()
	p.sit(sitWg)
	sitWg.Wait()

	for i:=0;i<hungry;i++ {
		p.pickupForks(forkMap)
		p.eat()
		p.putdownForks(forkMap)
		p.think()
	}

	fmt.Printf("%s ate and left\n", p.name)
	
}

func main() {

	philisophers := []Philosopher{
		{name: "Rushil", leftFork:4, rightFork: 0},
		{name: "Tom", leftFork:0, rightFork: 1},
		{name: "James", leftFork:1, rightFork: 2},
		{name: "Jason", leftFork:2, rightFork: 3},
		{name: "Harvey", leftFork:3, rightFork: 4},
	}

	totalPhilosophers := len(philisophers)
	forkMap := make(map[int]*sync.Mutex)


	for i:=0;i<totalPhilosophers;i++ {
		forkMap[i] = &sync.Mutex{}
	}

	// To make sure every one starts eating after everyone is seated
	sitWg := sync.WaitGroup{}
	sitWg.Add(totalPhilosophers)


	// To make sure the main end after everyone has finished eating
	wg := sync.WaitGroup{}
	wg.Add(totalPhilosophers)


	for i:=0;i<totalPhilosophers;i++ {
		p := philisophers[i]
		go p.dine(forkMap, &sitWg, &wg) //Running goroutine to start eating for each philosopher
	}


	// This is to make sure that program termniated once all philisophers have done eating
	wg.Wait()
	fmt.Println("Every ate and left....Problem solved")

}



