package main

import (
	"log"
	"os"
	"fmt"
	"strings"
	"math/rand"
	"github.com/logrusorgru/aurora/v4"
)

type Game struct {
	Words map[string]bool
	Word string
	Guesses []string
	Attempts int
}

func (g *Game) ReadGuess() {
	var guess string
	for {
		fmt.Scanln(&guess)
		guess = strings.ToUpper(guess)
		if g.Words[guess] {
			break
		} else {
			fmt.Print(aurora.BrightRed("Not a valid word! Try again: "))
		}
	}
	g.Guesses = append(g.Guesses, guess)
	g.Attempts--
}

func (g *Game) RenderGuess() bool {
	guess := g.Guesses[len(g.Guesses)-1]
	letterCount := make(map[string]int)
	for _, letter := range g.Word {
		letterCount[string(letter)]++
	}
	for i, letter := range guess {
		if letterCount[string(letter)] > 0 {
			letterCount[string(letter)]--
			if string(letter) == string(g.Word[i]) {
				fmt.Print(aurora.BrightGreen(string(letter) + " "))
			} else {
				fmt.Print(aurora.BrightYellow(string(letter) + " "))
			}
		} else {
			fmt.Print(aurora.BrightRed(string(letter) + " "))
		}
	}
	fmt.Println()
	if guess == g.Word {
		return true
	} else {
		return false
	}
}

func main() {
	content, err := os.ReadFile("words.txt")
    if err != nil {
        log.Fatal(err)
    }
    wordOptions := strings.Split(string(content), "\n")
	wordsMap := make(map[string]bool)
	for _, word := range wordOptions {
		wordsMap[word] = true
	}
	wordIndex := rand.Intn(len(wordOptions))
	chosen := wordOptions[wordIndex]
	fmt.Println(aurora.BgBlack(aurora.White("Guess the word!")))

	for _ = range(len(chosen)) {
		fmt.Print(aurora.BrightRed("_ "))
	}

	game := Game{
		Words: wordsMap,
		Word: chosen,
		Guesses: []string{},
		Attempts: 6,
	}

	fmt.Println()

	won := false
	for game.Attempts > 0 && !won {
		game.ReadGuess()
		if game.RenderGuess() {
			won = true
		}
	}
	if !won {
		fmt.Print(aurora.BrightRed("You lost! The word was: "))
		fmt.Print(aurora.BrightGreen(game.Word))
	} else {
		fmt.Println(aurora.Sprintf(aurora.BrightGreen("You won with %s attempts left!\n"), aurora.BrightRed(fmt.Sprint(game.Attempts))))
	}
}