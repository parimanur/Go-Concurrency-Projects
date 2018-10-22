package main

import (
	"container/heap"
	"fmt"
	"math"
	"math/rand"
	"time"
)

/*
Data Flow
k Clients pack the value x in Request object and sends it to work channel.
Load balancer blocks on REQ channel listening to Request(s).
Load balancer chooses a worker and sends Request to one of the channels of worker WOK(i).
Worker receives Request and processes x (say calculates sin(x) lol).
Worker updates load balancer using DONE channel. LB uses this for load-balancing.
Worker writes the sin(x) value in the RESP response channel (enclosed in Request object).

Channels in play
central work channel (Type: Request)
n WOK channels (n sized worker pool, Type: Work)
k RESP channels (k clients, Type: Float)
n DONE channels (Type: Work)

|Client1|         |Load|  <-DONE-- |Worker1| processing R3 coming from Client1
|Client2| --REQ-> |Blncr| --WOK->  |Worker2| processing R8 coming from Client2
                                   |Worker3| processing R5 coming from Client1
						  <-RESP-- response of R4 to Client2

*/
// Tweak the below numbers to check the concurrency power of Go!
const nRequester = 3
const nWorker = 3
const nTotalRequests = 5

type Request struct {
	data int
	resp chan float64
}

/*
Each client sends nTotalRequests.
In that loop, it is spawning requests that are sent to the central REQ channel linked to LB.
For response, requests use a common channel (RESP) per client.
*/
func createAndRequest(req chan Request, exit chan int, workerIndex int) {
	resp := make(chan float64)
	// spawn nTotalRequests per client
	for i := 0; i < nTotalRequests; i++ {
		// wait before next request
		//time.Sleep(time.Duration(rand.Int63n(int64(time.Millisecond))))
		req <- Request{int(rand.Int31n(90)), resp}
		// read value from RESP channel
		<-resp
		fmt.Printf("Client %d Responding to req#: %d \n", workerIndex, i)
	}
	//Mark client as exited
	fmt.Printf("Exiting Client: %d \n", workerIndex)
	exit <- 1
}

/*
Each worker is a forever-running loop go-routine.
In that loop, each worker is blocked on its channel trying to get Request object and then later process it.
Worker can take multiple requests. # of pending keeps track of number of requests being executed.
pending in other words means how many requests are present/being executed in the channel of each worker.
*/
type Work struct {
	// heap index
	idx int
	// WOK channel
	wok chan Request
	// number of pending request this worker is working on
	pending int
}

func (w *Work) doWork(done chan *Work) {
	// worker works indefinitely
	for {
		// extract request from WOK channel
		req := <-w.wok
		// write to RESP channel
		req.resp <- math.Sin(float64(req.data))
		// write to DONE channel
		done <- w
	}
}

/*
The crux of Balancer is a heap (Pool) which balances based on number of pending requests.
DONE channel, is used to notify heap that worker is finished and pending counter can be decremented.
*/
type Pool []*Work

type Balancer struct {
	// a pool of workers
	pool Pool
	done chan *Work
}

func InitBalancer() *Balancer {
	done := make(chan *Work, nWorker)
	// create nWorker WOK channels
	b := &Balancer{make(Pool, 0, nWorker), done}
	for i := 0; i < nWorker; i++ {
		w := &Work{wok: make(chan Request, nRequester)}
		// put them in heap
		heap.Push(&b.pool, w)
		go w.doWork(b.done)
	}
	return b
}

// Functions to implement for the heap interface
func (p Pool) Len() int { return len(p) }

func (p Pool) Less(i, j int) bool {
	return p[i].pending < p[j].pending
}

func (p *Pool) Swap(i, j int) {
	a := *p
	a[i], a[j] = a[j], a[i]
	a[i].idx = i
	a[j].idx = j
}

func (p *Pool) Push(x interface{}) {
	n := len(*p)
	item := x.(*Work)
	item.idx = n
	*p = append(*p, item)
}

func (p *Pool) Pop() interface{} {
	old := *p
	n := len(old)
	item := old[n-1]
	item.idx = -1 // for safety
	*p = old[0 : n-1]
	return item
}

/*
If the central REQ channel has a request coming in from clients, dispatch it to least loaded Worker and update the heap.
If DONE channel reports back, the work assigned to WOK(i) has been finished.
*/
func (b *Balancer) balance(req chan Request, exit chan int) {
	numOfClientsExited := 0
	for {
		select {
		// extract request from REQ channel
		case request := <-req:
			b.dispatch(request)
		// read from DONE channel
		case w := <-b.done:
			b.completed(w)
		case <-exit:
			numOfClientsExited++
			if numOfClientsExited == nRequester {
				return
			}
		}
		//b.print()
	}
}

func (b *Balancer) dispatch(req Request) {
	// Grab least loaded worker
	w := heap.Pop(&b.pool).(*Work)
	w.wok <- req
	w.pending++
	// Put it back into heap while it is working
	heap.Push(&b.pool, w)
}

func (b *Balancer) completed(w *Work) {
	w.pending--
	// remove from heap
	heap.Remove(&b.pool, w.idx)
	// Put it back
	heap.Push(&b.pool, w)
}

func (b *Balancer) print() {
	sum := 0
	sumsq := 0
	// Print pending stats for each worker
	for _, w := range b.pool {
		fmt.Printf("%d ", w.pending)
		sum += w.pending
		sumsq += w.pending * w.pending
	}
	// Print avg for worker pool
	avg := float64(sum) / float64(len(b.pool))
	variance := float64(sumsq)/float64(len(b.pool)) - avg*avg
	fmt.Printf(" %.2f %.2f\n", avg, variance)
}
func elapsed(what string) func() {
	start := time.Now()
	return func() {
		fmt.Printf("%s took %v\n", what, time.Since(start))
	}
}

func main() {
	fmt.Printf("Starting LoadBalancer with %d workers with %d requester clients with %d resuests/client \n", nWorker, nRequester, nTotalRequests)
	work := make(chan Request)
	exit := make(chan int)
	for i := 0; i < nRequester; i++ {
		go createAndRequest(work, exit, i)
	}
	defer elapsed("LoadBalancer")()
	InitBalancer().balance(work, exit)
}
