package model

type TableDescription struct {
	Field   string
	Type    string
	Null    string
	Key     string
	Default *string
	Extra   string
}
