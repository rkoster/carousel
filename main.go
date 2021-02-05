package main

import (
	"fmt"
	"log"
	"os"

	"code.cloudfoundry.org/credhub-cli/credhub"
	"code.cloudfoundry.org/credhub-cli/credhub/auth"
	"github.com/starkandwayne/carousel/store"

	"github.com/rivo/tview"
)

func main() {
	logger := log.New(os.Stderr, "", 0)
	cfg, err := loadConfig()
	if err != nil {
		logger.Fatalf("failed to load environment configuration: %s", err)
	}

	ch, _ := credhub.New(
		cfg.Credhub.Server,
		credhub.SkipTLSValidation(true), // TODO use CA
		credhub.Auth(auth.UaaClientCredentials(cfg.Credhub.Client, cfg.Credhub.Secret)),
	)

	fmt.Println("Connected to ", ch.ApiURL)

	certs, err := ch.GetAllCertificatesMetadata()
	if err != nil {
		logger.Fatalf("failed to load certificate metadate from Credhub: %s", err)
	}

	s := store.NewStore(certs)

	root := s.Tree()

	tree := tview.NewTreeView().
		SetRoot(root).
		SetCurrentNode(root)

	// If a directory was selected, open it.
	tree.SetSelectedFunc(func(node *tview.TreeNode) {
		// reference := node.GetReference()
		// if reference == nil {
		//	return // Selecting the root node does nothing.
		// }
		node.SetExpanded(!node.IsExpanded())
	})

	if err := tview.NewApplication().SetRoot(tree, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}