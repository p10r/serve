package main

import (
	"context"
	"github.com/p10r/serve/db"
	"github.com/p10r/serve/discord"
	"github.com/p10r/serve/domain"
	"github.com/p10r/serve/flashscore"
	"log"
	"os"
	"time"
)

var (
	favouriteLeagues = []string{
		"Italy: SuperLega",
		"Italy: SuperLega - Play Offs",
		"Italy: Coppa Italia A1",
		"Italy: Coppa Italia A1 Women",
		"Italy: Serie A1 Women",
		"Italy: Serie A1 Women - Playoffs",
		"Poland: PlusLiga",
		"Poland: PlusLiga - Play Offs",
		"France: Ligue A - Play Offs",
		"France: Ligue A",
		"Russia: Super League - Play Offs",
		"Russia: Super League",
		"Russia: Russia Cup",
		"World: Nations League",
		"World: Nations League - Play Offs",
		"World: Nations League Women",
		"World: Nations League Women - Play Offs",
		"World: Pan-American Cup",
		"World: World Championship - First round",
		"World: World Championship - Second round",
		"World: World Championship - Play Offs",
		"World: World Championship Women - First round",
		"Germany: VBL Supercup",
		"Germany: 1. Bundesliga",
		"Germany: 1. Bundesliga - Losers stage",
		"Germany: 1. Bundesliga - Winners stage",
		"Germany: 1. Bundesliga - Play Offs",
		"Germany: DVV Cup",
		"Turkey: Sultanlar Ligi Women",
		"Turkey: Sultanlar Ligi Women - Play Offs",
		"Turkey: Efeler Ligi",
		"TURKEY: Efeler Ligi - Play Offs",
		"Turkey: Efeler Ligi - 5th-8th places",
		"Europe: Champions League",
		"Europe: Champions League Women",
		"Europe: Champions League Women - Play Offs",
		"Europe: Champions League - Play Offs",
		"Europe: CEV Cup",
		"Europe: European Championships Women",
		"Europe: European Championships",
	}
	flashscoreUri = "https://flashscore.p.rapidapi.com"
	apiKey        = os.Getenv("API_KEY")
	discordUri    = os.Getenv("DISCORD_URI")
	dsn           = os.Getenv("DSN")
)

func main() {
	log.Println("Starting serve")
	workflow()
}

func workflow() {
	if apiKey == "" {
		log.Fatal("API_KEY has not been set")
	}

	if discordUri == "" {
		log.Fatal("DISCORD_URI has not been set")
	}

	conn := db.NewDB(dsn)
	err := conn.Open()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("DSN is set to %v", dsn)

	store := db.NewMatchStore(conn)
	flashscoreClient := flashscore.NewClient(flashscoreUri, apiKey)
	discordClient := discord.NewClient(discordUri)

	now := func() time.Time {
		return time.Now()
	}

	importer := domain.NewMatchImporter(store, flashscoreClient, discordClient, favouriteLeagues, now)

	_, err = importer.ImportScheduledMatches(context.Background())
	if err != nil {
		log.Println("Error when trying to import scheduled matches: ", err)
	}

	return
}
