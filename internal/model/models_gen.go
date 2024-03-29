// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"fmt"
	"io"
	"strconv"
)

type NewWeapon struct {
	Name string `json:"name"`
}

type Weapon struct {
	Name        string          `json:"name"`
	Type        string          `json:"type"`
	Phy         int             `json:"phy"`
	Mag         int             `json:"mag"`
	Fir         int             `json:"fir"`
	Lit         int             `json:"lit"`
	Hol         int             `json:"hol"`
	Cri         int             `json:"cri"`
	Sta         int             `json:"sta"`
	Str         AttributeScales `json:"str"`
	Dex         AttributeScales `json:"dex"`
	Int         AttributeScales `json:"int"`
	Fai         AttributeScales `json:"fai"`
	Arc         AttributeScales `json:"arc"`
	Any         string          `json:"any"`
	Phyb        int             `json:"phyb"`
	Magb        int             `json:"magb"`
	Firb        int             `json:"firb"`
	Litb        int             `json:"litb"`
	Holb        int             `json:"holb"`
	Bst         string          `json:"bst"`
	Rst         string          `json:"Rst"`
	Wgt         string          `json:"wgt"`
	Upgrade     string          `json:"upgrade"`
	ID          string          `json:"id"`
	Custom      bool            `json:"custom"`
	LastUpdated string          `json:"lastUpdated"`
}

type AttributeScales string

const (
	AttributeScalesA AttributeScales = "A"
	AttributeScalesB AttributeScales = "B"
	AttributeScalesC AttributeScales = "C"
	AttributeScalesD AttributeScales = "D"
	AttributeScalesE AttributeScales = "E"
	AttributeScales_ AttributeScales = "_"
)

var AllAttributeScales = []AttributeScales{
	AttributeScalesA,
	AttributeScalesB,
	AttributeScalesC,
	AttributeScalesD,
	AttributeScalesE,
	AttributeScales_,
}

func (e AttributeScales) IsValid() bool {
	switch e {
	case AttributeScalesA, AttributeScalesB, AttributeScalesC, AttributeScalesD, AttributeScalesE, AttributeScales_:
		return true
	}
	return false
}

func (e AttributeScales) String() string {
	return string(e)
}

func (e *AttributeScales) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = AttributeScales(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid AttributeScales", str)
	}
	return nil
}

func (e AttributeScales) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type Attributes string

const (
	AttributesStr Attributes = "STR"
	AttributesDex Attributes = "DEX"
	AttributesInt Attributes = "INT"
	AttributesFai Attributes = "FAI"
	AttributesArc Attributes = "ARC"
)

var AllAttributes = []Attributes{
	AttributesStr,
	AttributesDex,
	AttributesInt,
	AttributesFai,
	AttributesArc,
}

func (e Attributes) IsValid() bool {
	switch e {
	case AttributesStr, AttributesDex, AttributesInt, AttributesFai, AttributesArc:
		return true
	}
	return false
}

func (e Attributes) String() string {
	return string(e)
}

func (e *Attributes) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = Attributes(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Attributes", str)
	}
	return nil
}

func (e Attributes) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
