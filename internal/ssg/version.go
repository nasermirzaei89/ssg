package ssg

import "fmt"

func PrintVersion() error {
	_, err := fmt.Println("Version: ???")
	if err != nil {
		return fmt.Errorf("error on println: %w", err)
	}

	return nil
}
