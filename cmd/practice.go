package main

// importing statements
import (
	"fmt"
	"math/rand"
	"os"
	"sort"
	"sync"
	"sync/atomic"
	"time"
)

func workertesting(done chan bool) {
	fmt.Print("working..\n")
	time.Sleep(5 * time.Second)
	done <- true
}

func pong(pings <-chan string, pongs chan<- string) {
	// take from channel and put it in another channel
	msg := <-pings
	pongs <- msg
}

func pinger() {
	pings := make(chan string, 1)
	pongs := make(chan string, 1)
	go func() {
		pings <- "message from another channel"
	}()
	pong(pings, pongs)
	fmt.Println(<-pongs)
}

func selecTester() {
	c1 := make(chan string)
	c2 := make(chan string)
	go func() {
		time.Sleep(1 * time.Second)
		c1 <- "one"
	}()
	go func() {
		time.Sleep(2 * time.Second)
		c2 <- "two"
	}()

	for i := 0; i < 2; i++ {
		select {
		case msg1 := <-c1: // waiting from message from thread 1
			fmt.Println("received", msg1)
		case msg2 := <-c2: // waiting from message from thread 2
			fmt.Println("receieved", msg2)
		}
	}
}

func timeout() {
	c1 := make(chan string, 1) // buffered channel
	go func() {
		time.Sleep(2 * time.Second) // sleep for 2 second
		c1 <- "result 1"
	}()

	// select always choose one value out of all the values
	select {
	case res := <-c1:
		fmt.Println(res)
	case <-time.After(1 * time.Second): // after 1 second print the timeout
		fmt.Println("timeout 1")
	}

	c2 := make(chan string, 1)
	go func() {
		time.Sleep(2 * time.Second)
		c2 <- "result 2"
	}()
	select {
	case res := <-c2:
		fmt.Println(res)
	case <-time.After(3 * time.Second):
		fmt.Println("timeout 2")
	}
}

func nonblocking() {
	messages := make(chan string)
	signals := make(chan bool)

	select {
	case msg := <-messages:
		fmt.Println("recevied messages", msg)
	default:
		fmt.Println("no message received")
	}

	msg := "hi"
	select {
	case messages <- msg:
		fmt.Println("sent message", msg)
	default:
		fmt.Println("no message sent")

	}

	select {
	case msg := <-messages:
		fmt.Println("received message", msg)
	case sig := <-signals:
		fmt.Println("received signal", sig)
	default:
		fmt.Println("no activity ")

	}

}

func closing() {
	jobs := make(chan int, 5)
	done := make(chan bool)

	go func() {
		for {
			j, isChannelClosed := <-jobs
			if isChannelClosed {
				fmt.Println("received jobs", j)
			} else {
				fmt.Println("recieved all jobs... channel closed")
				done <- true // workers are done processing jobs
				return
			}
		}
	}()

	for j := 1; j <= 3; j++ {
		jobs <- j
		fmt.Println("sent jobs", j)
	}
	fmt.Println("closing job channel")
	close(jobs)
	fmt.Println("all jobs Sent ... channel  closed")
	<-done // wait for the worker goroutine to finish
}

func timers() {
	timer1 := time.NewTimer(2 * time.Second)
	<-timer1.C
	fmt.Println("Timer 1 fired")
	timer2 := time.NewTimer(time.Second)
	go func() {
		<-timer2.C
		fmt.Println("Timer 2 Fired")
	}()

	stop2 := timer2.Stop()
	if stop2 {
		fmt.Println("timer 2 stopped")
	}
	time.Sleep(2 * time.Second)
}

// executes cron jobs
func tickers() {
	ticker := time.NewTicker(500 * time.Millisecond)
	done := make(chan bool)
	go func() {
		for {
			select {
			case <-done:
				return
			case t := <-ticker.C:
				fmt.Println("Tick at ", t)
			}
		}
	}()
	time.Sleep(1600 * time.Millisecond)
	ticker.Stop() // not necessary but good habit
	done <- true  // signal go routine that ticker is stopped
	fmt.Println("Ticker Stopped")
}

func worker(id int, jobs <-chan int, results chan<- int) {
	for j := range jobs {
		fmt.Println("worker", id, "started job", j)
		time.Sleep(time.Second)
		fmt.Println("woker", id, "finished job", j)
		results <- j * 2 // send the result via channel
	}
}

func pools() {
	const numJobs = 5
	jobs := make(chan int, numJobs)
	results := make(chan int, numJobs)

	for w := 1; w <= 3; w++ {
		go worker(w, jobs, results) // create 3 worker threads
	}

	for j := 1; j <= numJobs; j++ {
		jobs <- j // populate jobs channel
	}
	close(jobs)
	for a := 1; a <= numJobs; a++ {
		<-results // wait for all the result to be calculated
	}
}

func waitWorker(id int) {
	fmt.Printf("worker %d starting \n", id)
	time.Sleep(time.Second)
	fmt.Printf("worker %d done\n", id)
}

func waitGroups() {
	var wg sync.WaitGroup
	for i := 1; i <= 5; i++ {
		wg.Add(1)
		i := i // make a i copy to be used in different goroutines
		go func() {

			defer wg.Done()
			waitWorker(i)
		}()
	}
	wg.Wait() // Wait for all goroutine to finish
}

// rate limiting methods
func rateLimiting() {
	requests := make(chan int, 5) // buffered requests
	for i := 1; i <= 5; i++ {
		requests <- i
	}
	close(requests)
	limiter := time.Tick(200 * time.Millisecond) // ticker
	for req := range requests {
		<-limiter
		fmt.Println("request", req, time.Now())
	}
	burstyLimiter := make(chan time.Time, 3)
	for i := 0; i < 3; i++ {
		burstyLimiter <- time.Now()
	}

	// Fill Burst limiter with 1 Values every 200 uptill it's buffer capacity
	go func() {
		for t := range time.Tick(200 * time.Millisecond) {
			burstyLimiter <- t
		}
	}()

	burstyRequests := make(chan int, 5)
	for i := 1; i <= 5; i++ {
		burstyRequests <- i
	}

	close(burstyRequests)
	fmt.Println("burst requests")
	for req := range burstyRequests {
		<-burstyLimiter
		fmt.Println("request", req, time.Now())
	}
}

func atomiccounter() {
	var ops uint64 // +ve int
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			for c := 0; c < 1000; c++ {
				atomic.AddUint64(&ops, 1) // add to ops value
			}
			wg.Done()
		}()
	}

	wg.Wait() // wait for all the goroutines to finish
	fmt.Println("ops:", ops)
}

// accessing data safely accorss goroutines
type Container struct {
	mu       sync.Mutex
	counters map[string]int
}

func (c *Container) inc(name string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.counters[name]++
}

func mutexUse() {
	c := Container{
		counters: map[string]int{"a": 0, "b": 0},
	}

	var wg sync.WaitGroup
	//closure function ( similar to lambda function in python)
	doIncrement := func(name string, n int) {
		for i := 0; i < n; i++ {
			c.inc(name)
		}
		wg.Done()
	}
	wg.Add(3)
	// Increment Values Concurrently
	go doIncrement("a", 10000)
	go doIncrement("a", 10000)
	go doIncrement("b", 10000)
	wg.Wait()
	fmt.Println(c.counters)
}

// stateful goroutines
type readOp struct {
	key  int
	resp chan int // reponse channel
}

type writeOp struct {
	key  int
	val  int
	resp chan bool // response channel
}

func states() {
	var readOps uint64
	var writeOps uint64

	reads := make(chan readOp)
	writes := make(chan writeOp)

	go func() {
		var state = make(map[int]int)
		for {
			select {
			case read := <-reads:
				read.resp <- state[read.key]
			case write := <-writes:
				state[write.key] = write.val
				write.resp <- true
			}
		}
	}()

	// Goroutines for reading Operations
	for r := 0; r < 100; r++ {
		go func() {
			// infinte loop
			for {
				read := readOp{
					key:  rand.Intn(5),
					resp: make(chan int)}
				reads <- read
				<-read.resp
				atomic.AddUint64(&readOps, 1)
				time.Sleep(time.Millisecond)
			}
		}()
	}

	//  goroutines for writing operation
	for w := 0; w < 10; w++ {
		go func() {
			for {
				write := writeOp{
					key:  rand.Intn(5),
					val:  rand.Intn(100),
					resp: make(chan bool)}
				writes <- write
				<-write.resp
				atomic.AddUint64(&writeOps, 1)
				time.Sleep(time.Millisecond)
			}
		}()
	}

	time.Sleep(time.Second)
	readOpsFinal := atomic.LoadUint64(&readOps)
	fmt.Println("readOps:", readOpsFinal)
	writeOpsFinal := atomic.LoadUint64(&writeOps)
	fmt.Println("writeOps:", writeOpsFinal)
	fmt.Println("total operation", readOpsFinal+writeOpsFinal)
}

func sorting() {
	ints := []int{1, 2, 5, 3, 9, 6}
	sort.Ints(ints)
	fmt.Println("Ints : sorting", ints)
	fmt.Println("checking if ints are sorted or not ( with builtint functions)")
	s := sort.IntsAreSorted(ints)
	fmt.Println("Sorted ? ", s)
}

type byLength []string

func (s byLength) Len() int {
	return len(s)
}

func (s byLength) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s byLength) Less(i, j int) bool {
	return len(s[i]) < len(s[j])
}

func sorting2() {
	fruits := []string{"peach", "banana", "kiwi"}
	fmt.Println("Before sorting,", fruits)
	sort.Sort(byLength(fruits))
	fmt.Println("After sorting,", fruits)
}

// returns pointer to the file descriptor
func createFile(p string) *os.File {
	fmt.Println("creating", p)
	f, err := os.Create(p)
	if err != nil {
		panic(err)
	}
	return f
}

func writeFile(f *os.File) {
	fmt.Println("writing to file")
	fmt.Fprintf(f, "data")
}
func closeFile(f *os.File) {
	fmt.Println("closing file ")
	err := f.Close()
	if err != nil {
		// printing to StdErr output stream
		fmt.Fprintf(os.Stderr, "error: &v\n", err)
		os.Exit(1)
	}
}

func fileoperation() {
	f := createFile("defer.txt")
	defer closeFile(f)
	writeFile(f)

}

// Main Method
func main() {
	fileoperation()
}
