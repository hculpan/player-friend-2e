package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
)

type Spell struct {
	Level        string   `json:"level"`
	Title        string   `json:"title"`
	Schools      []string `json:"schools"`
	Reversible   bool     `json:"reversible"`
	Range        string   `json:"range"`
	Components   string   `json:"components"`
	Duration     string   `json:"duration"`
	CastingTime  string   `json:"castingTime"`
	AreaOfEffect string   `json:"areaOfEffect"`
	SavingThrow  string   `json:"savingThrow"`
	Description  []string `json:"description"`
}

func main() {
	file, err := os.Open("WizardSpells.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	spellLevelRegex, err := regexp.Compile("-Level Spells")
	if err != nil {
		log.Fatal(err)
	}

	r := bufio.NewReader(file)

	var spells []Spell
	var currentSpell Spell = Spell{}
	var readyForNewSpell bool = true
	currentSpellLevel := "First-Level"

	fmt.Println("Starting read of file")
	for {
		line, err := r.ReadString('\r')
		if err != nil {
			if err == io.EOF {
				break
			} else {
				log.Fatal(err)
			}
		}

		if len(strings.TrimSpace(line)) == 0 {
			continue
		}

		trimmedLine := strings.TrimSpace(line)

		if readyForNewSpell {
			if currentSpell.Title != "" {
				spells = append(spells, currentSpell)
			}
			currentSpell = Spell{
				Level: currentSpellLevel,
				Title: trimmedLine,
			}
			readyForNewSpell = false
		} else if strings.HasPrefix(trimmedLine, "(") {
			schools := strings.Split(trimmedLine[1:len(trimmedLine)-1], ", ")
			currentSpell.Schools = schools
		} else if strings.HasPrefix(trimmedLine, "Reversible") {
			currentSpell.Reversible = true
		} else if strings.HasPrefix(trimmedLine, "Range:") {
			lines := strings.Split(trimmedLine, "Component")
			if len(lines) != 2 {
				log.Fatal("Error parsing range at ", currentSpell.Title, " ", lines)
			}
			currentSpell.Range = strings.TrimSpace(strings.TrimPrefix(lines[0], "Range: "))
			currentSpell.Components = strings.TrimSpace(strings.TrimPrefix(lines[1], "s: "))
		} else if strings.HasPrefix(trimmedLine, "Duration:") {
			lines := strings.Split(trimmedLine, "Casting Time")
			if len(lines) != 2 {
				log.Fatal("Error parsing duration at ", currentSpell.Title, " ", lines)
			}
			currentSpell.Duration = strings.TrimSpace(strings.TrimPrefix(lines[0], "Duration: "))
			currentSpell.CastingTime = strings.TrimSpace(strings.TrimPrefix(lines[1], ": "))
		} else if strings.HasPrefix(trimmedLine, "Area of Effect:") {
			lines := strings.Split(trimmedLine, "Saving Throw")
			if len(lines) != 2 {
				log.Fatal("Error parsing Area of Effect at ", currentSpell.Title, " ", lines)
			}
			currentSpell.AreaOfEffect = strings.TrimSpace(strings.TrimPrefix(lines[0], "Area of Effect: "))
			currentSpell.SavingThrow = strings.TrimSpace(strings.TrimPrefix(lines[1], ": "))
		} else if line[1] == '\t' {
			desc, newTitle := readDescription(trimmedLine, r)
			currentSpell.Description = desc
			spells = append(spells, currentSpell)

			if spellLevelRegex.MatchString(newTitle) {
				currentSpellLevel = strings.TrimSuffix(newTitle, " Spells")
				currentSpell = Spell{}
				readyForNewSpell = true
				continue
			}

			currentSpell = Spell{
				Title: newTitle,
				Level: currentSpellLevel,
			}
		}
	}
	fmt.Println("Finished read of file")

	if currentSpell.Title != "" {
		spells = append(spells, currentSpell)
	}

	jsonFile, _ := os.Create("wizard-spells.json")
	defer jsonFile.Close()

	jsonWriter := io.Writer(jsonFile)
	encoder := json.NewEncoder(jsonWriter)
	encoder.SetIndent("", "    ")
	encoder.Encode(spells)
}

func readDescription(line string, r *bufio.Reader) ([]string, string) {
	result := []string{line}
	newTitle := ""

	for {
		line, err := r.ReadString('\r')
		if err != nil {
			if err == io.EOF {
				break
			} else {
				log.Fatal(err)
			}
		}

		trimmedLine := strings.TrimSpace(line)
		if line[1] == '\t' || len(trimmedLine) == 0 {
			result = append(result, trimmedLine)
		} else {
			newTitle = trimmedLine
			break
		}
	}

	if result[len(result)-1] == "" {
		result = result[:len(result)-1]
	}

	return result, newTitle
}
