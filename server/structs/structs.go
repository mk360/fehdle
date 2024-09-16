package structs

type UnitResponse struct {
	CargoQuery []JSONUnit `json:"cargoquery"`
}

type JSONUnit struct {
	Title struct {
		MoveType   string `json:"MoveType"`
		WeaponType string `json:"WeaponType"`
		Name       string `json:"Page"`
		IntID      string `json:"IntID"`
		GameId     string `json:"GameSort"`
		WikiName   string `json:"WikiName"`
	} `json:"title"`
}
