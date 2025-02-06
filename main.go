package main

import (
	"fmt"
)

func frontend(p_x int, p_y int) {

	for i := 0; i < 10; i++ {

		for j := 0; j < 10; j++ {

			if i == p_y && j == p_x {
				fmt.Printf("R  ")
			} else {
				fmt.Printf("-  ")
			}
			if j == 9 {
				fmt.Printf("\n")
			}
		}
	}
}
func main() {
	var x, y int
	fmt.Printf("Enter the player's position (x y):\n")
	fmt.Scanln(&x, &y)
	frontend(x, y)
}