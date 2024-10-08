package common

import (
	"encoding/json"
	"errors"
	"fehdle/structs"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

var source = rand.NewSource(time.Now().UnixNano())

type TotalRosterResponse struct {
	CargoQuery []struct {
		Title struct {
			Count string `json:"Count"`
		} `json:"title"`
	} `json:"cargoquery"`
}

func GetRosterSize() int32 {
	var query = url.Values{
		"action": {"cargoquery"},
		"format": {"json"},
		"tables": {"Units"},
		"fields": {"COUNT(*)=Count"},
		"where":  {"Properties holds not \"enemy\""},
	}
	resp, e := http.Get("https://feheroes.fandom.com/api.php?" + query.Encode())
	if e != nil {
		log.Fatalln(e)
	}
	defer resp.Body.Close()
	data, _ := io.ReadAll(resp.Body)
	var countResponse TotalRosterResponse = TotalRosterResponse{}
	json.Unmarshal(data, &countResponse)
	conv, _ := strconv.ParseInt(countResponse.CargoQuery[0].Title.Count, 10, 64)
	return int32(conv)
}

func UpdateMainUnit() structs.UnitResponse {
	r := rand.New(source)

	var totalUnits = GetRosterSize()
	var query = url.Values{
		"action": {"cargoquery"},
		"format": {"json"},
		"tables": {"Units"},
		"limit":  {"1"},
		"where":  {"Properties holds not \"enemy\" "},
		"fields": {"MoveType, WeaponType, _pageName=Page, Properties, Entries, IntID"},
	}

	var randomOffset = r.Intn(int(totalUnits))
	query.Set("offset", strconv.FormatInt(int64(randomOffset), 10))
	resp, e := http.Get("https://feheroes.fandom.com/api.php?" + query.Encode())
	if e != nil {
		log.Fatalln(e)
	}
	defer resp.Body.Close()
	data, _ := io.ReadAll(resp.Body)
	var responseStruct structs.UnitResponse = structs.UnitResponse{}
	json.Unmarshal(data, &responseStruct)
	fmt.Println(responseStruct)
	return responseStruct
}

func FindHero(intId string) (structs.JSONUnit, error) {
	var query = url.Values{
		"action": {"cargoquery"},
		"format": {"json"},
		"tables": {"Units"},
		"fields": {"MoveType, WeaponType, _pageName=Page, WikiName, GameSort, IntID"},
		"where":  {"IntID = " + intId},
	}
	resp, e := http.Get("https://feheroes.fandom.com/api.php?" + query.Encode())
	if e != nil {
		log.Fatalln(e)
	}
	defer resp.Body.Close()
	data, _ := io.ReadAll(resp.Body)
	var foundUnit structs.UnitResponse = structs.UnitResponse{}
	json.Unmarshal(data, &foundUnit)
	fmt.Println(foundUnit)
	var jsonUnit structs.JSONUnit = structs.JSONUnit{}
	if len(foundUnit.CargoQuery) > 0 {
		jsonUnit = foundUnit.CargoQuery[0]
		return jsonUnit, nil
	}
	var notFoundError = errors.New("Could not find unit with intId " + intId)
	return jsonUnit, notFoundError
}
