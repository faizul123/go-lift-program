package main

import (
	"fmt"
	"net/http"
	"time"
)

func main() {
	fmt.Println("Lift Program started")
	var queue map[int]chan User
	queue = make(map[int]chan User)
	for i := 0; i < 10; i++ {
		queue[i] = make(chan User, 3)
		defer close(queue[i])
	}

	lift := Lift{
		GF:    0,
		TF:    10,
		CF:    0,
		Dir:   1,
		Users: []User{},
		q:     queue,
	}

	userTicker := time.NewTicker(400 * time.Millisecond)
	defer userTicker.Stop()

	liftTicker := time.NewTicker(1 * time.Second)
	defer liftTicker.Stop()

	go func() {
		for {
			select {
			case <-userTicker.C:
				user := getRandomUser()
				if user.In != user.Out {
					lift.AddUserToQueue(user)
					fmt.Printf("%s user generated on floor %d and exit floor %d\n", user.ID, user.In, user.Out)
				}
			case <-liftTicker.C:
				lift.AddUser()
				lift.MovieLift()
				lift.RemoveUser()
				fmt.Printf("Lift moved to floor %d\n", lift.CF)
			}
		}
	}()

	http.ListenAndServe("localhost:8080", nil)
	// time.Sleep(10 * time.Second)
	fmt.Println("Lift Program ended")
}
