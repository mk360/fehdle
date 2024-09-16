package common

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type ImageResponse struct {
	Cargoquery []struct {
		Page string `json:"Page"`
	} `json:"cargoquery"`
}

var redirectlessHttpClient = http.Client{
	CheckRedirect: func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	},
}

func convertImageType(imgType string) string {
	switch imgType {
	case "thumbnail":
		return "_Face_FC"
	case "portrait":
		return "_Face"
	}

	return ""
}

func GetImage(heroName string, imgType string) []byte {
	var cacheRequest, _ = http.NewRequest("GET", "https://feheroes.fandom.com/wiki/Special:Redirect/file/"+url.QueryEscape(strings.Replace(heroName, " ", "_", -1))+convertImageType(imgType)+".webp", nil)
	redirectResponse, _ := redirectlessHttpClient.Do(cacheRequest)

	var location = strings.Replace(redirectResponse.Header.Get("Location"), "/revision/latest", "/revision/latest/scale-to-width-down/300", 1)
	imageCDNLocation, _ := http.Get(location)

	defer imageCDNLocation.Body.Close()
	imageByteData, _ := io.ReadAll(imageCDNLocation.Body)

	return imageByteData
}

func ImgRoute(response http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	if len(request.Form["id"]) == 0 {
		response.WriteHeader(400)
		return
	}

	var imgType = request.Form.Get("imgType")

	var query = url.Values{}
	query.Add("action", "cargoquery")
	query.Add("format", "json")
	query.Add("tables", "Units")
	query.Add("fields", "Units.WikiName=Page")
	query.Add("where", "Properties holds not \"enemy\" and IntID = "+request.Form.Get("id"))

	resp, e := http.Get("https://feheroes.fandom.com/api.php?" + query.Encode())

	if e != nil {
		response.Write([]byte(""))
		return
	}

	defer resp.Body.Close()

	data, _ := io.ReadAll(resp.Body)
	var unmarshaled ImageResponse = ImageResponse{}
	json.Unmarshal(data, &unmarshaled)
	if len(unmarshaled.Cargoquery) == 0 {
		fmt.Println("Searched for a missing id, " + request.Form.Get("id"))
		response.WriteHeader(404)
		return
	}
	var imageByteData = GetImage(unmarshaled.Cargoquery[0].Page, imgType)
	response.Write(imageByteData)
}
