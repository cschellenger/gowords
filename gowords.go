package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Game struct {
	Words    map[string]bool
	Word     string
	Guesses  []string
	Attempts int
	Won      bool
}

func (g *Game) RenderGuess() (string, error) {
	guess := g.Guesses[len(g.Guesses)-1]
	if g.Words[guess] == false {
		return "", fmt.Errorf("Invalid word: %s", guess)
	}
	if g.Attempts <= 0 {
		return "", fmt.Errorf("No attempts left")
	}
	g.Attempts--
	letterCount := make(map[string]int)
	for _, letter := range g.Word {
		letterCount[string(letter)]++
	}
	// find the exact matches first
	for i, letter := range guess {
		if string(letter) == string(g.Word[i]) {
			letterCount[string(letter)]--
		}
	}
	// print the word with colors
	rendered := ""
	for i, letter := range guess {
		if string(letter) == string(g.Word[i]) {
			rendered += "[green]" + string(letter) + " "
		} else if letterCount[string(letter)] > 0 {
			letterCount[string(letter)]--
			rendered += "[yellow]" + string(letter) + " "
		} else {
			rendered += "[red]" + string(letter) + " "
		}
	}
	if g.Word == guess {
		g.Won = true
	}
	return rendered, nil
}

func (g *Game) RenderLetters() string {
	// print the letters with colors
	rendered := ""
	letters := make(map[string]string)
	for _, guess := range g.Guesses {
		for i, guessLetter := range guess {
			if string(g.Word[i]) == string(guessLetter) {
				letters[string(guessLetter)] = "[green]"
			}
			for _, wordLetter := range g.Word {
				if string(guessLetter) == string(wordLetter) {
					if letters[string(guessLetter)] == "" {
						letters[string(guessLetter)] = "[yellow]"
					}
				}
			}
			if letters[string(guessLetter)] == "" {
				letters[string(guessLetter)] = "[red]"
			}
		}
	}
	for _, l := range "QWERTYUIOP\nASDFGHJKL\nZXCVBNM\n" {
		if l == '\n' {
			rendered += "\n"
		} else {
			if letters[string(l)] == "" {
				letters[string(l)] = "[white]"
			}
			rendered += letters[string(l)] + string(l) + " "
		}
	}
	return rendered
}

func main() {
	content, err := os.ReadFile("words.txt")
	if err != nil {
		log.Fatal(err)
	}
	wordOptions := strings.Split(strings.ToUpper(string(content)), "\n")
	wordsMap := make(map[string]bool)
	for _, word := range wordOptions {
		wordsMap[word] = true
	}
	wordIndex := rand.Intn(len(wordOptions))
	chosen := wordOptions[wordIndex]
	game := Game{
		Words:    wordsMap,
		Word:     chosen,
		Guesses:  []string{},
		Attempts: 6,
	}

	app := tview.NewApplication()
	inputField := tview.NewInputField().SetFieldWidth(10)
	inputField.SetAcceptanceFunc(func(text string, lastChar rune) bool {
		if game.Attempts <= 0 || game.Won {
			return false
		}
		return len(text) <= len(game.Word)
	})
	guessesView := tview.NewTextView().
		SetDynamicColors(true).
		SetSize(6, 20).
		SetTextAlign(tview.AlignLeft)

	lettersView := tview.NewTextView().
		SetDynamicColors(true).
		SetSize(6, 20).
		SetTextAlign(tview.AlignCenter).
		SetText(game.RenderLetters())

	messageView := tview.NewTextView().
		SetDynamicColors(true).
		SetSize(1, 40).
		SetTextAlign(tview.AlignLeft)

	quitButton := tview.NewButton("Quit")
	quitButton.SetSelectedFunc(func() {
		app.Stop()
	})

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(tview.NewFlex().
			SetDirection(tview.FlexColumn).
			AddItem(guessesView, 0, 1, false).
			AddItem(lettersView, 0, 1, false), 0, 4, false).
		AddItem(inputField, 0, 1, true).
		AddItem(messageView, 0, 1, false).
		AddItem(quitButton, 1, 1, false)

	inputField.SetDoneFunc(func(key tcell.Key) {
		switch key {
		case tcell.KeyEnter:
			guess := strings.ToUpper(inputField.GetText())
			game.Guesses = append(game.Guesses, guess)
			rendered, err := game.RenderGuess()
			inputField.SetText("")
			if err == nil {
				guessesView.SetText(guessesView.GetText(false) + "\n" + rendered)
				messageView.SetText("")
				lettersView.SetText(game.RenderLetters())
			} else {
				messageView.SetText(err.Error())
			}
			if game.Won {
				inputField.SetDisabled(true)
				messageView.SetText("[green]You won! The word was: " + game.Word)
				app.SetFocus(quitButton)
			} else if game.Attempts <= 0 {
				inputField.SetDisabled(true)
				messageView.SetText("[red]You lost! The word was: " + game.Word)
				app.SetFocus(quitButton)
			} else {
				app.SetFocus(inputField)
			}
		case tcell.KeyEscape, tcell.KeyTab:
			app.SetFocus(quitButton)
		}
	})

	if err := app.SetRoot(flex, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}
