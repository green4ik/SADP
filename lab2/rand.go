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
		return nil, errors.New("модуль має бути більшим за нуль")
	}
	if multiplier < 0 && multiplier >= modulus {
		return nil, errors.New("множник має бути не меншим за нуль і не більшим за модуль")
	}
	if increment < 0 && increment >= modulus {
		return nil, errors.New("приріст має бути не меншим за нуль і не більшим за модуль")
	}
	if seed < 0 && increment >= modulus {
		return nil, errors.New("початкове значення має бути не меншим за нуль і не більшим за модуль")
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
