package context

import (
	"github.com/Isaac-Fate/myst/internal/config"
	"github.com/Isaac-Fate/myst/internal/manager"
)

type AppContext struct {
	Config        config.Config
	Passphrase    string
	SecretManager *manager.SecretManager
}
