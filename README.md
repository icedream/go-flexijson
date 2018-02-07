# FlexiJSON

This helper allows a developer to generate a wrapping struct that decodes and
encodes arbitrary extra JSON fields that are not part of the actual defined
Go struct.

This should be combined with the usage of `go generate` for instant awesomeness!

## Installation

`go get -v github.com/icedream/go-flexijson/...`

## Usage

`flexijson-generator -p <package name> -i <input go file> -o <output go file>`

The output .go file will contain the wrapping struct code and should not need
any additional modifications. The wrapping struct type name is the same as the
original struct type name but with the first letter in uppercase, so make sure
your original type name starts with a lowercase letter.

This code generator works even if the original struct code uses types that have
not yet been declared, even the generated wrapping struct itself. This allows
for recursive usage.

## Example

An example is available in the [example](example/) directory.

```go
package main

//go:generate flexijson-generator -p "$GOPACKAGE" -i "$GOFILE" -o "gen_flexijson_$GOFILE"

// someStruct will be wrapped by the generated SomeStruct type.
type someStruct struct {
	A      string                 `json:"A,omitempty"`
	B      string                 `json:"B,omitempty"`
	Sub    *SomeStruct            `json:"sub,omitempty"`
	Fields map[string]interface{} `json:"-,extra"`
}
```
