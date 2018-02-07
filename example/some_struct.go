package example

//go:generate flexijson-generator -p "$GOPACKAGE" -i "$GOFILE" -o "gen_flexijson_$GOFILE"

type someStruct struct {
	A      string
	B      string                 `json:"B,omitempty"`
	Sub    *SomeStruct            `json:"sub,omitempty"`
	Fields map[string]interface{} `json:"-" extrajson`
}
