package main

import (
	"encoding/xml"
)

type PingData struct {
	XMLName     xml.Name
	EigenePings int
	AllePings   int
}
type MediaData struct {
	XMLName   xml.Name
	VideoData MISData
	AudioData MISData
}

type MediaInfo struct {
	FileName      string
	InterpretName string
	SongName      string
	StilName      string
}

type NameIndex struct {
	Name  string
	Index []int
}

type MISData struct {
	Mediainfos  []MediaInfo
	Interpreten []NameIndex
	Stile       []NameIndex
}
