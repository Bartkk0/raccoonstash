package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"flag"
	"fmt"
	"raccoonstash/internal"
)

var cliConfig struct {
	register   string
	unregister string
	list       bool
	stats      bool
}

func main() {
	internal.AddGlobalFlags()
	flag.StringVar(&cliConfig.register, "register", "", "Register a token (rand to generate random token)")
	flag.StringVar(&cliConfig.unregister, "unregister", "", "Unregister a token (all to unregister all)")
	flag.BoolVar(&cliConfig.list, "tokens", false, "List all tokens")
	flag.BoolVar(&cliConfig.stats, "stats", false, "Show stats")
	flag.Parse()

	internal.InitializeDatabase()

	if cliConfig.register != "" {
		if cliConfig.register == "rand" {
			b := make([]byte, 32)
			_, err := rand.Read(b)
			if err != nil {
				println(err)
				return
			}
			cliConfig.register = hex.EncodeToString(b)
		}

		err := internal.Queries.RegisterToken(context.Background(), cliConfig.register)
		if err != nil {
			println(err.Error())
			return
		}
		println("Registered token", cliConfig.register)
		return
	}
	if cliConfig.unregister != "" {
		if cliConfig.unregister != "all" {
			err := internal.Queries.UnregisterToken(context.Background(), cliConfig.unregister)
			if err != nil {
				println(err.Error())
				return
			}
			println("Unregistered token", cliConfig.unregister)
		} else {
			err := internal.Queries.UnregisterAllTokens(context.Background())
			if err != nil {
				println(err.Error())
				return
			}
			println("Unregistered all tokens")
		}
		return
	}

	if cliConfig.list {
		tokens, _ := internal.Queries.GetAllTokens(context.Background())
		println("Token\tCreated at")
		for _, token := range tokens {
			println(token.Token, "\t", token.CreatedAt.String())
		}
		return
	}

	if cliConfig.stats {
		stats, _ := internal.Queries.GetStats(context.Background())
		fmt.Printf("%d files %d bytes\n", stats.NFiles, stats.FilesSize)
		fmt.Printf("%d pastes %d bytes\n", stats.NPastes, stats.PastesSize)
	}
}
