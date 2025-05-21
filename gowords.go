package main

import (
	"log"
	"os"
	"fmt"
	"strings"
	"math/rand"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Game struct {
	Words map[string]bool
	Word string
	Guesses []string
	Attempts int
	Won bool
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
	game := Game{
		Words: wordsMap,
		Word: chosen,
		Guesses: []string{},
		Attempts: 6,
	}

	app := tview.NewApplication()
	form := tview.NewForm().
		AddTextView("", "", 20, 6, true, false).
		AddInputField("", "", 20, nil, nil).
		AddButton("Quit", func() {
			app.Stop()
		})
	textView := form.GetFormItem(0).(*tview.TextView)
	inputField := form.GetFormItem(1).(*tview.InputField)
	inputField.SetAcceptanceFunc(func(text string, lastChar rune) bool {
		if game.Attempts <= 0 || game.Won {
			return false
		}
		return len(text) <= len(game.Word)
	})
	inputField.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			guess := strings.ToUpper(inputField.GetText())
			game.Guesses = append(game.Guesses, guess)
			rendered, err := game.RenderGuess()
			inputField.SetText("")
			if err == nil {
				textView.SetText(textView.GetText(false) + "\n" + rendered)
			}
			if game.Won {
				inputField.SetDisabled(true)
				form.AddTextView("", "You won!", 20, 1, true, false)
			} else if game.Attempts <= 0 {
				inputField.SetDisabled(true)
				form.AddTextView("", "You lost! The word was: " + game.Word, 40, 1, true, false)
			} else {	
				go func() {
					app.QueueUpdateDraw(func() {
						app.SetFocus(inputField)
					})
				}()
			}
		}
	})

	form.SetBorder(true).SetTitle("Go Words").SetTitleAlign(tview.AlignLeft)
	if err := app.SetRoot(form, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}