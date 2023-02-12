package entities

type IndexTableRow struct {
	UID             uint64
	FileStartOffset int
	FileEndOffset   int
}

type IndexTable struct {
	PreviousMaxUID uint64
	OverallOffset  int
	Rows           []IndexTableRow
}
