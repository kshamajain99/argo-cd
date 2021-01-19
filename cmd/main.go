package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	appcontroller "github.com/argoproj/argo-cd/cmd/argocd-application-controller/commands"
	dex "github.com/argoproj/argo-cd/cmd/argocd-dex/commands"
	reposerver "github.com/argoproj/argo-cd/cmd/argocd-repo-server/commands"
	apiserver "github.com/argoproj/argo-cd/cmd/argocd-server/commands"
	util "github.com/argoproj/argo-cd/cmd/argocd-util/commands"
	cli "github.com/argoproj/argo-cd/cmd/argocd/commands"
	cmdutil "github.com/argoproj/argo-cd/cmd/util"
)

const (
	binaryNameEnv = "ARGOCD_BINARY_NAME"
)

func main() {
	var command *cobra.Command

	binaryName := filepath.Base(os.Args[0])
	if val := os.Getenv(binaryNameEnv); val != "" {
		binaryName = val
	}
	switch binaryName {
	case "argocd", "argocd-linux-amd64", "argocd-darwin-amd64", "argocd-windows-amd64.exe":
		command = cli.NewCommand()
	case "argocd-util", "argocd-util-linux-amd64", "argocd-util-darwin-amd64", "argocd-util-windows-amd64.exe":
		command = util.NewCommand()
	case "argocd-server":
		command = apiserver.NewCommand()
	case "argocd-application-controller":
		command = appcontroller.NewCommand()
	case "argocd-repo-server":
		command = reposerver.NewCommand()
	case "argocd-dex":
		command = dex.NewCommand()
	default:
		if suggestionsString := cmdutil.FindSuggestions(binaryName); suggestionsString != "" {
			fmt.Printf("unknown command %q for %s\n", binaryName, suggestionsString)
			os.Exit(1)
		}
		// default is argocd cli
		command = cli.NewCommand()
	}

	if err := command.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
