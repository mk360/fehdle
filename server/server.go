package main

import (
	"fehdle/routes"
	"net/http"
)

// func compareWithMainUnit(chosenPick JSONUnit) GuessingResult {
// 	if chosenPick.Title.IntID == mainUnit.IntID {
// 		var correctResult GuessingResult = GuessingResult{
// 			Name:                chosenPick.Title.Name,
// 			WeaponType:          0,
// 			GameIdDiffDirection: 0,
// 			WeaponTypeData:      chosenPick.Title.WeaponType,
// 			WikiName:            chosenPick.Title.WikiName,
// 		}

// 		return correctResult
// 	}

// 	var wrongResult GuessingResult = GuessingResult{
// 		Name:           chosenPick.Title.Name,
// 		WeaponTypeData: chosenPick.Title.WeaponType,
// 		WikiName:       chosenPick.Title.WikiName,
// 	}

// 	var splitCorrectWeaponType = strings.Split(mainUnit.WeaponType, " ")
// 	var splitAnswerWeaponType = strings.Split(wrongResult.WeaponTypeData, " ")
// 	var weaponTypeDifference uint = 0
// 	if splitAnswerWeaponType[0] != splitCorrectWeaponType[0] {
// 		weaponTypeDifference++
// 	}

// 	if splitAnswerWeaponType[1] != splitCorrectWeaponType[1] {
// 		weaponTypeDifference++
// 	}

// 	wrongResult.WeaponType = weaponTypeDifference

// 	convertedResultGameId, _ := strconv.Atoi(chosenPick.Title.GameId)

// 	if convertedResultGameId > mainUnit.GameId {
// 		wrongResult.GameIdDiffDirection = 1
// 	} else {
// 		wrongResult.GameIdDiffDirection = -1
// 	}

// 	return wrongResult
// }

func corsMiddleware(next http.Handler) http.Handler {
	var corsMiddleware = http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		// writer.Header().Add("Access-Control-Allow-Origin", os.Getenv("CORS_DOMAIN"))
		writer.Header().Add("Access-Control-Allow-Origin", "*")
		writer.Header().Add("Access-Control-Allow-Methods", "POST")
		next.ServeHTTP(writer, request)
	})

	return corsMiddleware
}

func main() {
	mux := http.NewServeMux()

	go routes.UpdateGoroutine()

	var classicGuess = http.HandlerFunc(routes.ClassicRoute)
	var shadowGuess = http.HandlerFunc(routes.ShadowRoute)

	mux.Handle("/classic", corsMiddleware(classicGuess))
	mux.Handle("/shadow", corsMiddleware(shadowGuess))

	http.ListenAndServe(":4444", mux)
}
