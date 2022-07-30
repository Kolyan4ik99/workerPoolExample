package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"
)

var actions = []string{"logged in", "logged out", "created record", "deleted record", "updated account"}

type logItem struct {
	action    string
	timestamp time.Time
}

type User struct {
	id    int
	email string
	logs  []logItem
}

func (u User) getActivityInfo() string {
	output := fmt.Sprintf("UID: %d; Email: %s;\nActivity Log:\n", u.id, u.email)
	for index, item := range u.logs {
		output += fmt.Sprintf("%d. [%s] at %s\n", index, item.action, item.timestamp.Format(time.RFC3339))
	}

	return output
}

func main() {
	const countUser = 100
	const countWorkers = 10
	jobs := make(chan User, countUser)
	rand.Seed(time.Now().Unix())

	startTime := time.Now()
	wg := &sync.WaitGroup{}

	for i := 0; i < countWorkers; i++ {
		go saveUserInfo(wg, jobs)
	}

	for i := 0; i < countUser; i++ {
		tmp := i
		go func() {
			wg.Add(1)
			jobs <- generateUsers(tmp)
		}()
	}
	//close(jobs)
	wg.Wait()

	fmt.Printf("DONE! Time Elapsed: %.2f seconds\n", time.Since(startTime).Seconds())
}

func saveUserInfo(wg *sync.WaitGroup, us <-chan User) {
	for user := range us {
		fmt.Printf("WRITING FILE FOR UID %d\n", user.id)

		filename := fmt.Sprintf("users/uid%d.txt", user.id)
		file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			log.Fatal(err)
		}

		file.WriteString(user.getActivityInfo())
		time.Sleep(time.Second)
		wg.Done()
	}
}

func generateUsers(id int) User {
	users := User{}

	users = User{
		id:    id + 1,
		email: fmt.Sprintf("user%d@company.com", id+1),
		logs:  generateLogs(rand.Intn(1000)),
	}
	fmt.Printf("generated user %d\n", id+1)
	time.Sleep(time.Millisecond * 10)

	return users
}

func generateLogs(count int) []logItem {
	logs := make([]logItem, count)

	for i := 0; i < count; i++ {
		logs[i] = logItem{
			action:    actions[rand.Intn(len(actions)-1)],
			timestamp: time.Now(),
		}
	}

	return logs
}
