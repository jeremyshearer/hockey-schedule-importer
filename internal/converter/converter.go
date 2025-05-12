package converter

import (
	"github.com/jeremyshearer/hockey-schedule-importer/internal/benchapp"
	"github.com/jeremyshearer/hockey-schedule-importer/internal/gamesheets"
)

func Convert(gamesheetsGames []gamesheets.Game) ([]benchapp.Game, error) {
	var benchappGames []benchapp.Game
	for _, game := range gamesheetsGames {
		benchappGame := benchapp.Game{
			Type:     "GAME",
			GameType: "REGULAR",
			Title:    "",
			Away:     game.Visitor,
			Home:     game.Home,
			Date:     game.Date.Format("02/01/2006"),
			Time:     game.Date.Format("3:04 PM"),
			Duration: "1:30",
			Location: game.Location,
			Address:  "",
			Notes:    "",
		}
		benchappGames = append(benchappGames, benchappGame)
	}
	return benchappGames, nil
}
