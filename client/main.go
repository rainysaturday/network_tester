package main

import (
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

func get_data() int {
	resp, err := http.Get("http://localhost:8080?size=100")
	if err != nil {
		fmt.Println(err)
		return 0
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return 0
	}

	return len(body)
}

func call_many(num_requests int) {
	t1 := time.Now()
	for i := 0; i < num_requests; i++ {
		get_data()
	}

	taken := time.Since(t1).Milliseconds()
	regs_p_s := (num_requests * 1000) / int(taken)
	fmt.Printf("Finished sequential connections in %d ms = %d reqs/s\n", taken, regs_p_s)
}

func call_many_par(num_requests int, max_par int) {
	var wg sync.WaitGroup
	t1 := time.Now()
	for i := 0; i < max_par; i++ {
		wg.Add(1)
		go func(id int) {
			for i := 0; i < (num_requests / max_par); i++ {
				get_data()
			}
			wg.Done()
		}(i)
	}
	wg.Wait()

	taken := time.Since(t1).Milliseconds()
	regs_p_s := (num_requests * 1000) / int(taken)
	fmt.Printf("Finished parallel %d connections in %d ms = %d reqs/s\n", max_par, taken, regs_p_s)
}

func burst_test() {
	call_many(10000)
	for i := 1; i < 20; i++ {
		call_many_par(10000, i)
	}
}

func main() {
	burst_test()
}
