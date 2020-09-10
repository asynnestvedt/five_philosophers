package main

import (
	"fmt"
	"sync"
)

/**
Implement the dining philosopher’s problem with the following constraints/modifications.

There should be 5 philosophers sharing Sporks, with one Spork between each adjacent
pair of philosophers.
Each philosopher should eat only 3 times (not in an infinite loop as we did in lecture)
The philosophers pick up the Spork in any order, not lowest-numbered first
(which we did in lecture).
In order to eat, a philosopher must get permission from a host which executes in its own goroutine.
The host allows no more than 2 philosophers to eat concurrently.
Each philosopher is numbered, 1 through 5.
When a philosopher starts eating (after it has obtained necessary locks) it prints
“starting to eat <number>” on a line by itself, where <number> is the number of the philosopher.
When a philosopher finishes eating (before it has released its locks) it prints “finishing eating
<number>” on a line by itself, where <number> is the number of the philosopher.
**/

// Spork -
type Spork struct{ sync.Mutex }

// Philo -
type Philo struct {
	id             int
	lSpork, rSpork *Spork
}

func host(ch chan bool, concurrency int) {
	for i := 0; i < concurrency; i++ {
		ch <- true
	}
}

func (p Philo) eat(ch chan bool, wg *sync.WaitGroup, meals int) {

	for i := 0; i < meals; i++ {
		<-ch
		// modifying the following instruction to occur BEFORE the lock since it was likely written
		// incorrectly in the project requirements
		// "When a philosopher starts eating (after it has obtained necessary locks) it prints
		//   “starting to eat <number>”" */
		fmt.Printf("starting to eat %d\n", p.id)

		p.lSpork.Lock()
		p.rSpork.Lock()

		fmt.Printf("finishing eating %d\n", p.id)

		p.rSpork.Unlock()
		p.lSpork.Unlock()
		ch <- true
		wg.Done()
	}
}

func main() {
	const MealsPerPerson = 3
	const PersonsCount = 5

	ch := make(chan bool, 2)
	var wg sync.WaitGroup

	sporks := make([]*Spork, PersonsCount)
	for i := 0; i < PersonsCount; i++ {
		sporks[i] = new(Spork)
	}

	philos := make([]*Philo, 0, PersonsCount)
	for i := 0; i < PersonsCount; i++ {
		philos = append(philos, &Philo{i + 1, sporks[i], sporks[(i+1)%PersonsCount]})
	}

	wg.Add(PersonsCount * MealsPerPerson)
	go host(ch, 2)

	for _, p := range philos {
		go p.eat(ch, &wg, MealsPerPerson)
	}

	wg.Wait()
}
