package main

import "time"

func test_go_routine(){

	for i := 0; i < 21000; i ++ {
		go func() {

			for {
				time.Sleep(0)
			}
		}()
	}

	for {
		time.Sleep(0)
	}
}
