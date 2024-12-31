package main

import (
	"fmt"
	"math/rand"
	"sync"
)

type Lift struct {
	GF    int // Ground Floor
	TF    int // Top Floor
	CF    int // Current Floor
	Dir   int // Direction: 1 = Up, 0 = Down
	q     map[int]chan User
	Users []User
	lock  sync.Mutex
}

type User struct {
	ID  string
	In  int
	Out int
	Dir int
}

var UsersNames = []string{"John", "Jane", "Doe", "Smith", "Williams", "Brown", "Jones", "Garcia", "Miller", "Davis"}

func getRandomUser() User {
	index := rand.Intn(len(UsersNames))

	id := rand.Intn(1000000)
	if id == 0 {
		id = 1
	}

	in := rand.Intn(10)
	out := rand.Intn(10)

	dir := 0
	if in > out {
		dir = 0
	} else {
		dir = 1
	}

	return User{
		ID:  fmt.Sprintf("%d-%s", id, UsersNames[index]),
		In:  in,
		Out: out,
		Dir: dir,
	}
}

func (l *Lift) MovieLift() {
	if l.Dir == 1 && l.CF < l.TF {
		l.CF++
	} else if l.Dir == 0 && l.CF > l.GF {
		l.CF--
	}

	if l.CF == l.TF {
		l.Dir = 0
	} else if l.CF == l.GF {
		l.Dir = 1
	}
}

func (l *Lift) AddUser() {
	l.lock.Lock()
	ch := l.q[l.CF]

	totalUsers := len(ch)
	if totalUsers == 0 {
		fmt.Printf("No users in queue on floor %d\n", l.CF)
		l.lock.Unlock()
		return
	}

	for user := range ch {
		if user.In == l.CF {
			l.Users = append(l.Users, user)
			fmt.Printf("%s added into lift\n", user.ID)
		}
		if len(ch) == 0 {
			break
		}
	}
	l.lock.Unlock()
}

func (l *Lift) RemoveUser() {
	l.lock.Lock()
	for i, user := range l.Users {
		if user.Out == l.CF {
			fmt.Printf("%d %d \n", i, len(l.Users))
			if i < len(l.Users) {
				l.Users = l.Users[:i]
			} else {
				l.Users = append(l.Users[:i], l.Users[i+1:]...)
			}
			fmt.Printf("%s removed from lift\n", user.ID)
		}
	}
	l.lock.Unlock()
}

func (l *Lift) AddUserToQueue(user User) {
	l.lock.Lock()
	ch := l.q[user.In]
	// fmt.Printf("ok: %t, user.In: %d\n", ok, user.In)
	totalUsers := len(ch)
	if totalUsers == 3 {
		fmt.Printf("Queue on floor %d is full\n", user.In)
		l.lock.Unlock()
		return
	}
	l.q[user.In] <- user
	l.lock.Unlock()
	fmt.Printf("%s added to queue on floor %d and exit floor %d\n", user.ID, user.In, user.Out)
}
