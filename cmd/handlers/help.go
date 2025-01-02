package handlers

import (
	"fmt"

	"github.com/Isaac-Fate/myst/cmd/context"
)

func ShowHelp(appContext *context.AppContext) error {
	helpText := `
🔐 MyST (My SecreTs) Help

Available Commands:
  add     Add a new secret
          - Prompts for key, value, website (optional), and notes (optional)
          - Values are encrypted using your master passphrase

  find    Search for secrets
          - Search by key, website, or notes
          - View decrypted values for found secrets

  list    List all secrets
          - Shows all stored secrets
          - Option to view decrypted values

  remove  Remove a secret
          - Remove by key
          - Confirms before deletion

  help    Show this help message

  quit    Exit the application

Tips:
  - You can type commands or use arrow keys to select
  - Use Ctrl+C to cancel any operation
  - Secret values are always encrypted before storage
  - Keep your master passphrase safe - it cannot be recovered!
`
	fmt.Println(helpText)
	return nil
}
