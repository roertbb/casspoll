package main

import (
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/manifoldco/promptui"
)

const layout = "2006-01-02T15:04:05"

func selectOption(label string, options []string) (int, string, error) {
	prompt := promptui.Select{
		Label: label,
		Items: options,
	}

	idx, result, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return -1, "", err
	}

	return idx, result, nil
}

func promptString(label string) (string, error) {
	prompt := promptui.Prompt{
		Label: label,
	}

	result, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return "", err
	}

	return result, nil
}

func promptNonEmptyString(label string) (string, error) {
	prompt := promptui.Prompt{
		Label: label,
		Validate: func(input string) error {
			if len(input) == 0 {
				return errors.New("Invalid value")
			}
			return nil
		},
	}

	result, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return "", err
	}

	return result, nil
}

func promptFutureDate(label string) (time.Time, error) {
	prompt := promptui.Prompt{
		Label: label,
		Validate: func(input string) error {
			date, err := time.Parse(layout, input)
			if date.Before(time.Now()) {
				return errors.New("Date need to be in the future")
			}
			if err != nil {
				return errors.New("Failed to parse date")
			}
			return nil
		},
	}

	result, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return time.Now(), err
	}

	date, _ := time.Parse(layout, result)
	return date, nil
}

func clearScreen() {
	print("\033[H\033[2J")
}

func reverse(s interface{}) {
	size := reflect.ValueOf(s).Len()
	swap := reflect.Swapper(s)
	for i, j := 0, size-1; i < j; i, j = i+1, j-1 {
		swap(i, j)
	}
}
