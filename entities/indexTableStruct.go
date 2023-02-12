package entities

type IndexTableRow struct {
	UID        int `json:"uid"`
	NumInArray int `json:"numInArray"`
}

type IndexTable struct {
	Uid  int             `json:"previousMaxUID"`
	Rows []IndexTableRow `json:"rows"`
}
