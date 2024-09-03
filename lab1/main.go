package main

import (
	"errors"
	"fmt"
	"os"
)

type LCG struct {
	modulus    int
	multiplier int
	increment  int
	seed       int
}

func NewLCG(modulus int, multiplier int, increment int, seed int) (*LCG, error) {
	if modulus <= 0 {
		return nil, errors.New("Модуль має бути більшим за нуль")
	}
	if multiplier < 0 && multiplier >= modulus {
		return nil, errors.New("Множник має бути не меншим за нуль і не більшим за модуль")
	}
	if increment < 0 && increment >= modulus {
		return nil, errors.New("Приріст має бути не меншим за нуль і не більшим за модуль")
	}
	if seed < 0 && increment >= modulus {
		return nil, errors.New("Початкове значення має бути не меншим за нуль і не більшим за модуль")
	}
	return &LCG{
		modulus:    modulus,
		multiplier: multiplier,
		increment:  increment,
		seed:       seed,
	}, nil
}

func (lcg *LCG) Next() int {
	lcg.seed = (lcg.multiplier*lcg.seed + lcg.increment) % lcg.modulus
	return lcg.seed
}
func toFile(lcg LCG, n int) {
	fileLcg := lcg
	file, err := os.OpenFile("file.txt", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Помилка відкриття файлу : %v\n", err)
		return
	}
	defer file.Close()
	for i := 0; i < n; i++ {
		lcgseed := fileLcg.Next()
		fmt.Fprintf(file, "%d, ", lcgseed)
	}
	fmt.Fprintln(file, "=============================================================================================================================================================")
}
func main() {
	fmt.Println("Введіть кількість виведених чисел:")
	var iters int
	_, err := fmt.Fscanf(os.Stdin, "%d", &iters)
	if err != nil {
		fmt.Println("Введіть ціле число:", err)
		return
	}
	if iters <= 0 {
		fmt.Println("Число має бути більше нуля")
		return
	}
	lcg, err := NewLCG((2<<11)-1, (2<<4)*(2<<4), 2, 8)
	if err != nil {
		fmt.Println("Умови параметрів незадовільні")
		return
	}
	toFile(*lcg, iters)
	fmt.Printf("Модуль : %d \nМножник : %d \nПриріст : %d \nПочаткове число : %d\n",
		lcg.modulus, lcg.multiplier, lcg.increment, lcg.seed)
	fmt.Println("===============================")
	arr := make([]int, 0)
	for i := 0; i < iters; i++ {
		lcgseed := lcg.Next()
		arr = append(arr, lcgseed)
		fmt.Println(lcgseed)
		for j, val := range arr {
			if val == lcgseed && i != j {
				defer fmt.Printf("Елемент %d повторюється з елементом %d", j, i)
				return
			}

		}

	}
	fmt.Println("===============================")
}
