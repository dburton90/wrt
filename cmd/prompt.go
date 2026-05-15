package cmd

import (
	"fmt"
	"strings"

	"github.com/manifoldco/promptui"
)

func promptString(label, defaultVal string, required bool) (string, error) {
	validate := func(s string) error {
		if required && strings.TrimSpace(s) == "" {
			return fmt.Errorf("required")
		}
		return nil
	}
	p := promptui.Prompt{
		Label:    label,
		Default:  defaultVal,
		Validate: validate,
	}
	val, err := p.Run()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(val), nil
}

func promptOptional(label string) (string, error) {
	p := promptui.Prompt{
		Label: label + " (optional, press Enter to skip)",
	}
	val, err := p.Run()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(val), nil
}

// promptBackportBranches prompts the user to enter version=branch pairs until empty input.
func promptBackportBranches() (map[string]string, error) {
	m := map[string]string{}
	fmt.Println("Add backport branch mappings (format: version=branch, e.g. 1.1.1=release/1.1.1).")
	fmt.Println("Press Enter with empty input when done.")
	for {
		p := promptui.Prompt{Label: "Backport mapping"}
		val, err := p.Run()
		if err != nil {
			return m, nil
		}
		val = strings.TrimSpace(val)
		if val == "" {
			break
		}
		parts := strings.SplitN(val, "=", 2)
		if len(parts) != 2 {
			fmt.Println("  Invalid format. Use version=branch (e.g. 1.1.1=release/1.1.1)")
			continue
		}
		m[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
	}
	return m, nil
}
