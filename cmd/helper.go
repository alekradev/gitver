package cmd

import (
	"fmt"
	"github.com/spf13/viper"
	"gotver/internal/gitops"
	"gotver/internal/version"
	"log"
)

func loadConfig() {
	prints("load configuration")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}

	if err := version.ReadVersion(); err != nil {
		log.Fatal(err)
	}
	prints("load configuration success")
}

func prepareGitOperation() error {
	prints("perpare git operations")
	if err := gitops.ReadRepository(); err != nil {
		return fmt.Errorf("git repository is not initialized: %w", err)
	}

	isClean, err := gitops.IsCleanRepo()
	if err != nil {
		return err
	}

	if !isClean {
		return fmt.Errorf("git repository is not clean")
	}

	prints("prepare git operations success")
	return nil

}

func logf(format string, v ...any) {
	if verbose {
		log.Printf(format, v)
	}
}

func prints(v ...any) {
	if verbose {
		log.Print(v)
	}
}
