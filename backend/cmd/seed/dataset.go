package main

import (
	"log"

	"zone/internal/seeddata"
)

// printReport prints the manual command's credentials table. The server never
// calls this because it includes the shared demo password. The dataset itself
// lives in internal/seeddata so both entry points use the same rows.
func printReport(result *seeddata.Result) {
	log.Printf("\nlog in as any of these (password: %s):", seeddata.DemoPassword)
	for _, u := range result.Users {
		kind := "public"
		if u.IsPrivate {
			kind = "private"
		}
		log.Printf("  %-24s %-18s %-10s id=%d", u.Email, u.Nickname, kind, u.ID)
	}

	log.Println("\nsession tokens (cookie name: session_token):")
	for _, u := range result.Users {
		log.Printf("  %-24s %s", u.Nickname, result.Sessions[u.ID])
	}
}
