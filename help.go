package ssg

import "fmt"

func PrintHelp() error {
	_, err := fmt.Println("This is a help")
	if err != nil {
		return fmt.Errorf("error on println: %w", err)
	}

	return nil
}
