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
		return -1
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return -1
	}

	return len(body)
}

func call_many(num_requests int) {
	t1 := time.Now()
	for i := 0; i < num_requests; i++ {
		if get_data() == -1 {
			fmt.Printf("Failed sequential connections after %d requests \n", i)
		}
	}

	taken := time.Since(t1).Milliseconds()
	regs_p_s := (num_requests * 1000) / int(taken)
	fmt.Printf("Finished sequential connections in %d ms = %d reqs/s\n", taken, regs_p_s)
}

func call_many_delay(num_requests int, delay_micros int) bool {
	t1 := time.Now()
	delay := time.Duration(delay_micros) * time.Microsecond
	for i := 0; i < num_requests; i++ {
		if delay_micros != 0 {
			time.Sleep(delay)
		}
		if get_data() == -1 {
			fmt.Printf("Failed sequential connections w/delay in %d ms after %d requests \n", delay, i)
			return false
		}
	}

	taken := time.Since(t1).Milliseconds()
	regs_p_s := (num_requests * 1000) / int(taken)
	fmt.Printf("Finished sequential connections w/delay in %d ms = %d reqs/s\n", taken, regs_p_s)
	return true
}

func call_many_par(num_requests int, max_par int) {
	var wg sync.WaitGroup
	t1 := time.Now()
	for i := 0; i < max_par; i++ {
		wg.Add(1)
		go func(id int) {
			for i := 0; i < (num_requests / max_par); i++ {
				if get_data() == -1 {
					fmt.Printf("Failed parallel connections after %d requests \n", i)
				}
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

func rate_limit_test() {
	for i := 20000; i >= 0; i -= 100 {
		call_many_delay(100, i)
	}
}

func main() {
	burst_test()
	rate_limit_test()
}
