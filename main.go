package main

import (
	"fmt"
	"github.com/rivo/tview"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

func getBitwardenPassphrase(keyFile string) (string, error) {
	// Execute rbw get command to retrieve the passphrase
	cmd := exec.Command("rbw", "get", "-f", "passphrase", keyFile)
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get passphrase from Bitwarden: %v", err)
	}
	return strings.TrimSpace(string(output)), nil
}

func showHostVerificationPrompt(prompt string) (bool, string) {
	app := tview.NewApplication()

	pages := tview.NewPages()

	// Initial modal
	modal := tview.NewModal().
		SetText(prompt).
		AddButtons([]string{"yes", "fingerprint", "no"})

	var userChoice bool
	var customFingerprint string

	// Create fingerprint input form
	fingerprintForm := tview.NewForm()
	fingerprintForm.
		AddInputField("Fingerprint:", "", 60, nil, func(text string) {
			customFingerprint = text
		}).
		AddButton("Confirm", func() {
			userChoice = true
			app.Stop()
		}).
		AddButton("Cancel", func() {
			pages.SwitchToPage("modal")
		})
	fingerprintForm.SetBorder(true).SetTitle("Enter Host Fingerprint")

	// Add both pages
	pages.AddPage("modal", modal, true, true)
	pages.AddPage("fingerprint", fingerprintForm, true, false)

	modal.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		switch buttonLabel {
		case "yes":
			userChoice = true
			app.Stop()
		case "fingerprint":
			pages.SwitchToPage("fingerprint")
		case "no":
			userChoice = false
			app.Stop()
		}
	})

	if err := app.SetRoot(pages, true).EnableMouse(true).Run(); err != nil {
		return false, ""
	}

	return userChoice, customFingerprint
}

func showPrompt(prompt string) string {
	app := tview.NewApplication()

	var response string

	promptForm := tview.NewForm()
	promptForm.
		AddInputField(prompt, "", 60, nil, func(text string) {
			response = text
		}).
		AddButton("Confirm", func() {
			app.Stop()
		})
	if err := app.SetRoot(promptForm, true).SetFocus(promptForm).Run(); err != nil {
		panic(err)
	}
	if err := app.SetRoot(promptForm, true).EnableMouse(true).Run(); err != nil {
		return ""
	}

	return response
}

func main() {
	arguments := os.Args[1:]
	// Get the prompt from arguments
	prompt := strings.Join(arguments, " ")
	// If the prompt is empty, print a message and exit
	if prompt == "" {
		fmt.Println("I'm a ssh askpass helper and I'm not designed to work standalone")
		return
	}
	// Check if it's a host verification prompt
	if strings.Contains(strings.ToLower(prompt), "the authenticity of host") {
		trust, fingerprint := showHostVerificationPrompt(prompt)
		if trust {
			if fingerprint != "" {
				fmt.Println(fingerprint)
			} else {
				fmt.Println("yes")
			}
		} else {
			fmt.Println("no")
		}
		return
	}

	// If it's a passphrase prompt, get the key file name from the prompt
	// Assuming prompt format: "Enter passphrase for /path/to/key:"
	if strings.Contains(strings.ToLower(prompt), "passphrase") {
		var keyFile string
		keyFileRegex, _ := regexp.Compile(` for(?: key)? '?(.+?)'?: `)
		keyFileMatch := keyFileRegex.FindString(prompt)
		keyFilePath := strings.TrimSuffix(strings.TrimPrefix(keyFileMatch, " for '"), "': ")
		keyFile = strings.Split(keyFilePath, "/")[len(strings.Split(keyFileMatch, "/"))-1]
		_, err := fmt.Fprintf(os.Stderr, "Using item: %v\n", keyFile)
		if err != nil {
			return
		}

		if keyFile != "" {
			passphrase, e := getBitwardenPassphrase(keyFile)
			if e != nil {
				_, er := fmt.Fprintf(os.Stderr, "Error: %v\n", e)
				if er != nil {
					return
				}
				os.Exit(1)
			}
			fmt.Print(passphrase)
			return
		}
	}

	if strings.Contains(strings.ToLower(prompt), "password for 'https://") {
		_, err := fmt.Fprintf(os.Stderr, "Askpass for Git HTTP(S) is not supported (yet).\nConsider configuring git to use the native one provided by rbw instead:\nhttps://git-scm.com/docs/gitcredentials\n")
		// TODO: another tool for git credentials manage
		// for now use https://github.com/doy/rbw/blob/main/bin/git-credential-rbw
		if err != nil {
			return
		}
		os.Exit(1)
	}

	// If we don't recognize the prompt, just ask the user
	response := showPrompt(prompt + " ")
	fmt.Print(strings.TrimSpace(response))
}
