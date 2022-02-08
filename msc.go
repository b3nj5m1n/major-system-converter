package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wordwrap"
	"github.com/muesli/reflow/wrap"
)

const (
	heightResults = 8
)

var styleError = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#FAFAFA")).
	Background(lipgloss.Color("#d14774"))

type arrayFlags []string

func (i *arrayFlags) String() string {
	return ""
}
func (i *arrayFlags) Set(value string) error {
	*i = append(*i, strings.TrimSpace(value))
	return nil
}

var (
	datasetFiles arrayFlags
)

type Word struct {
	word      string
	number    string
	pos       string
	phonemes  string
	frequency int
}

type WordList []Word

type AllWords WordList
type ResultWords WordList
type SelectedWords WordList

func (w WordList) Len() int {
	return len(w)
}
func (w WordList) Less(i, j int) bool {
	return w[i].frequency > w[j].frequency
}
func (w WordList) Swap(i, j int) {
	w[i], w[j] = w[j], w[i]
}

type Result struct {
	array []string
	lines int
}

type Model struct {
	allWords      AllWords
	resultWords   ResultWords
	selectedWords SelectedWords
	loading       bool
	input         textinput.Model
	spinner       spinner.Model
	width         int
	scroll        int
	result        Result
}

func initialModel() Model {
	ti := textinput.New()
	ti.Placeholder = "Number to convert"
	ti.Focus()
	ti.CharLimit = 256
	ti.Width = 20

	s := spinner.New()
	spinners := []spinner.Spinner{spinner.Line, spinner.Dot, spinner.MiniDot, spinner.Jump, spinner.Pulse, spinner.MiniDot, spinner.Jump, spinner.Pulse, spinner.MiniDot, spinner.Jump, spinner.Pulse, spinner.Points, spinner.Globe, spinner.Moon, spinner.Monkey, spinner.Globe, spinner.Moon, spinner.Monkey, spinner.Globe, spinner.Moon, spinner.Monkey, spinner.Globe, spinner.Moon, spinner.Monkey}
	rand.Seed(time.Now().UnixNano())
	index := rand.Intn(len(spinners))
	s.Spinner = spinners[index]

	return Model{
		allWords:      nil,
		resultWords:   nil,
		selectedWords: nil,
		loading:       true,
		input:         ti,
		spinner:       s,
		result:        Result{array: nil, lines: 0},
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, loadWords)
}

func loadWords() tea.Msg {
	var wordList AllWords

	for _, path := range datasetFiles {
		f, err := os.Open(path)
		if err != nil {
			fmt.Println(styleError.Render("Error reading file."))
			return tea.Quit
		}

		defer f.Close()

		r := csv.NewReader(f)
		for {
			record, err := r.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				fmt.Println(styleError.Render("Error reading file."))
				return tea.Quit
			}
			frequency, convErr := strconv.Atoi(record[4])
			if convErr != nil {
				frequency = 0
			}
			word := Word{
				word:      record[0],
				number:    record[1],
				pos:       record[2],
				phonemes:  record[3],
				frequency: frequency,
			}
			wordList = append(wordList, word)
		}
	}
	return wordList
}

func findWords(number string, searchList WordList) tea.Cmd {
	return func() tea.Msg {
		var result ResultWords
		for _, word := range searchList {
			if word.number == number {
				result = append(result, word)
			}
		}
		sort.Sort(WordList(result))
		return result
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" || msg.String() == "q" {
			return m, tea.Quit
		}
		if m.loading {
			return m, nil
		}
		switch msg.String() {
		case "up", "k":
			if m.scroll > 0 {
				m.scroll--
			}
		case "down", "j":
			m.scroll++
		case "0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "backspace":
			m.input, _ = m.input.Update(msg)
			if m.input.Value() == "" {
				m.resultWords = nil
				return m, nil
			} else {
				return m, findWords(m.input.Value(), WordList(m.allWords))
			}
		default:
			return m, nil
		}
		return m, formatWordList(m)
	case spinner.TickMsg:
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case AllWords:
		m.allWords = msg
		m.loading = false
		return m, nil
	case ResultWords:
		m.resultWords = msg
		m.scroll = 0
		return m, formatWordList(m)
	case tea.WindowSizeMsg:
		m.width = msg.Width
		return m, nil
	case Result:
		m.result = msg
		if m.scroll < 0 {
			m.scroll = 0
		}
		if (m.result.lines - (heightResults / 2)) < m.scroll {
			m.scroll = m.result.lines - (heightResults / 2)
		}
		return m, nil
	default:
		return m, nil
	}
}

func frequencyStyle(w Word) lipgloss.Style {
	// 500 most common words
	if w.frequency >= 470945 {
		return lipgloss.NewStyle().
			Italic(true).
			Underline(true)
	}
	// 1000 most common words
	if w.frequency >= 251904 {
		return lipgloss.NewStyle().
			Underline(true)
	}
	// 2500 most common words
	if w.frequency >= 100742 {
		return lipgloss.NewStyle().
			Italic(true)
	}
	// 10000 most common words
	if w.frequency >= 16752 {
		return lipgloss.NewStyle()
	}
	return lipgloss.NewStyle().
		Faint(true)
}

func formatWord(w Word) string {
	style := frequencyStyle(w)
	word := w.word
	switch w.pos {
	case "JJ":
		return style.
			Bold(true).
			Foreground(lipgloss.Color("#ffffff")).
			Background(lipgloss.Color("#00d5ff")).
			Render(word)
	case "JJR":
		return style.
			Foreground(lipgloss.Color("#00ffea")).
			Render(word)
	case "JJS":
		return style.
			Foreground(lipgloss.Color("#00ffbf")).
			Render(word)
	case "NN":
		return style.
			Bold(true).
			Foreground(lipgloss.Color("#ffffff")).
			Background(lipgloss.Color("#ff4eff")).
			Render(word)
	case "NNS":
		return style.
			Foreground(lipgloss.Color("#ff89ff")).
			Render(word)
	case "NNP":
		return style.
			Bold(true).
			Foreground(lipgloss.Color("#ffffff")).
			Background(lipgloss.Color("#ff89ff")).
			Render(word)
	case "NNPS":
		return style.
			Foreground(lipgloss.Color("#ff89ff")).
			Render(word)
	case "VB":
		return style.
			Bold(true).
			Foreground(lipgloss.Color("#ffffff")).
			Background(lipgloss.Color("#ffb327")).
			Render(word)
	case "VBD", "VBG", "VBN", "VBP", "VBZ":
		return style.
			Foreground(lipgloss.Color("#ffd35a")).
			Render(word)
	case "FW":
		return style.
			Foreground(lipgloss.Color("#717171")).
			Render(word)
	}
	return style.
		Foreground(lipgloss.Color("#778899")).
		Render(word)
}

func formatWordList(m Model) tea.Cmd {
	return func() tea.Msg {
		result := ""
		for _, word := range m.resultWords {
			result += formatWord(word) + ", "
		}
		result = wrap.String(wordwrap.String(result, m.width*2/3), m.width*2/3)
		resultSplit := strings.Split(result, "\n")
		resultLen := len(resultSplit)
		return Result{array: resultSplit, lines: resultLen}
	}
}

func scrollResult(result Result, scroll int) string {
	resultStr := ""
	for i := scroll; i < scroll+heightResults; i++ {
		if i >= result.lines {
			resultStr += "\n"
		} else {
			resultStr += result.array[i] + "\n"
		}
	}
	return resultStr
}

func (m Model) View() string {
	style := lipgloss.NewStyle().
		MaxWidth(m.width - (3 * 2))
	styleContainer := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#874BFD")).
		Padding(1, 1).
		MarginTop(2).
		Width(m.width / 2)
	styleContainerResults := lipgloss.NewStyle().
		Width(m.width * 2 / 3).
		MaxWidth(m.width * 2 / 3).
		MaxHeight(heightResults)

	s := ""

	inputStr := ""
	if m.loading {
		inputStr += m.spinner.View()
		inputStr += " Loading dataset"
	} else {
		inputStr += m.input.View()
	}
	input := lipgloss.Place(m.width, 0,
		lipgloss.Center, lipgloss.Center,
		styleContainer.Render(inputStr),
	)
	s += input

	results := ""
	if m.resultWords != nil {
		results = styleContainerResults.Render(scrollResult(m.result, m.scroll))
	}
	s += lipgloss.PlaceHorizontal(m.width, lipgloss.Center, results)
	s += "\n"

	var styleHelp = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#708090")).
		MarginTop(2).
		UnsetBackground()
	help := styleHelp.Render(wordwrap.String("q to quit, j or up arrow to go up, k or down arrow to go down", m.width))
	s += lipgloss.PlaceHorizontal(m.width, lipgloss.Center, help)

	return style.Render(s)
}

func main() {
	flag.Var(&datasetFiles, "d", "Dataset file(s) to use.")
	flag.Parse()

	if len(datasetFiles) < 1 {
		fmt.Println(styleError.Render("Specify at least one dataset."))
		os.Exit(1)
	}

	p := tea.NewProgram(initialModel())
	if err := p.Start(); err != nil {
		fmt.Println(styleError.Render("Alas, there's been an error."))
		os.Exit(1)
	}
}
