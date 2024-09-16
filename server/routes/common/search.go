package common

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type HeroData struct {
	Page       string `json:"Page"`
	MoveType   string `json:"MoveType"`
	WeaponType string `json:"WeaponType"`
	Properties string `json:"Properties"`
	IntID      int    `json:"IntID"`
}

type HeroSearchQuery struct {
	Cargoquery []struct {
		Title HeroData `json:"title"`
	} `json:"cargoquery"`
}

func searchHero(query string) []HeroData {
	var urlQuery = url.Values{}
	urlQuery.Add("format", "json")
	urlQuery.Add("action", "cargoquery")
	urlQuery.Add("tables", "Units")
	urlQuery.Add("fields", "_pageName=Page, MoveType, WeaponType, Properties, WikiName, IntID")
	urlQuery.Add("where", "lower(_pageName) like \"%"+strings.ToLower(query)+"%\"")
	res, _ := http.Get("https://feheroes.fandom.com/api.php?" + urlQuery.Encode())
	data, _ := io.ReadAll(res.Body)
	var response HeroSearchQuery = HeroSearchQuery{}
	var heroes []HeroData = make([]HeroData, len(response.Cargoquery))
	json.Unmarshal(data, &response)
	for i, _ := range response.Cargoquery {
		heroes[i] = response.Cargoquery[i].Title
	}

	return heroes
}

func SearchRoute(responseWriter http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	var query = request.Form["q"]
	if len(query) > 0 && len(query[0]) > 2 {
		var resp = searchHero(query[0])
		var marshaled, _ = json.Marshal(resp)
		responseWriter.Write(marshaled)
	} else {
		responseWriter.Write([]byte("[]"))
	}
}
