# workerPoolExample

## Example

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