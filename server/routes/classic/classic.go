package classic

import (
	"encoding/json"
	"fehdle/routes/common"
	"fehdle/structs"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	cron "github.com/robfig/cron/v3"
)

type MainUnit struct {
	MovementType string
	WeaponType   string
	Name         string
	Properties   string
	Entries      string
	IntID        string
	GameId       int
}

type GuessingResult struct {
	WeaponTypeDifference uint8  `json:"wpnDiff"` // should return 0 if the match is exact, 1 if either color or weapon is correct, 2 if nothing is correct
	WeaponTypeData       string `json:"wpn"`
	GameIdDiffDirection  int    `json:"gameIdDiff"` // -1 if today's hero was released before the choice, 0 if it's the same game, or 1 if they were released after
	Name                 string `json:"name"`
	GameId               int    `json:"gameId"`
}

var mainUnit MainUnit = MainUnit{
	MovementType: "",
	WeaponType:   "",
	Name:         "",
	IntID:        "",
	GameId:       0,
}

func Guess(writer http.ResponseWriter, request *http.Request) {
	byteIntId, _ := io.ReadAll(request.Body)
	match, e := regexp.Match("[0-9]", byteIntId)
	if !match || e != nil {
		writer.Write([]byte("Invalid payload format: expected a valid IntID, received " + string(byteIntId)))
		writer.WriteHeader(400)
	}

	var intIdString = string(byteIntId)
	guessed, e := common.FindHero(intIdString)

	if e != nil {
		writer.WriteHeader(404)
	} else {
		var compared = compareWithMainUnit(guessed)
		byteResponse, _ := json.Marshal(compared)
		fmt.Println(string(byteResponse))
		writer.Write(byteResponse)
	}
}

func compareWithMainUnit(chosenPick structs.JSONUnit) GuessingResult {
	convertedResultGameId, _ := strconv.Atoi(chosenPick.Title.GameId)

	if chosenPick.Title.IntID == mainUnit.IntID {
		var correctResult GuessingResult = GuessingResult{
			Name:                 chosenPick.Title.Name,
			WeaponTypeDifference: 0,
			GameIdDiffDirection:  0,
			WeaponTypeData:       chosenPick.Title.WeaponType,
			GameId:               convertedResultGameId,
		}

		return correctResult
	}

	var wrongResult GuessingResult = GuessingResult{
		Name:           chosenPick.Title.Name,
		WeaponTypeData: chosenPick.Title.WeaponType,
	}

	var splitCorrectWeaponType = strings.Split(mainUnit.WeaponType, " ")
	var splitAnswerWeaponType = strings.Split(wrongResult.WeaponTypeData, " ")
	var weaponTypeDifference uint8 = 0
	if splitAnswerWeaponType[0] != splitCorrectWeaponType[0] {
		weaponTypeDifference++
	}

	if splitAnswerWeaponType[1] != splitCorrectWeaponType[1] {
		weaponTypeDifference++
	}

	wrongResult.WeaponTypeDifference = weaponTypeDifference
	wrongResult.GameId = convertedResultGameId

	if convertedResultGameId > mainUnit.GameId {
		wrongResult.GameIdDiffDirection = 1
	} else {
		wrongResult.GameIdDiffDirection = -1
	}

	return wrongResult
}

func UpdateGoroutine() {
	updateCron := cron.New()
	updateCron.AddFunc("* * * * *", func() {
		var unit = common.UpdateMainUnit().CargoQuery[0].Title
		mainUnit.MovementType = unit.MoveType
		mainUnit.WeaponType = unit.WeaponType
		mainUnit.IntID = unit.IntID
		mainUnit.Name = unit.Name
		gameIdInt, _ := strconv.Atoi(unit.GameId)
		mainUnit.GameId = gameIdInt
	})

	updateCron.Start()

	select {}
}
