package model

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var reSRV = regexp.MustCompile(`^(?P<weight>\d+) (?P<port>\d+) (?P<target>\S+)$`)

// Type is a DNS type.
type Type uint16

//supported types
const (
	TypeNone  = Type(0)
	TypeA     = Type(1)
	TypeNS    = Type(2)
	TypeCNAME = Type(5)
	TypePTR   = Type(12)
	TypeMX    = Type(15)
	TypeTXT   = Type(16)
	TypeAAAA  = Type(28)
	TypeSRV   = Type(33)
)

func (t Type) String() string {
	var enumVal string

	switch t {
	case TypeA:
		enumVal = "A"
	case TypeNS:
		enumVal = "NS"
	case TypeCNAME:
		enumVal = "CNAME"
	case TypePTR:
		enumVal = "PTR"
	case TypeMX:
		enumVal = "MX"
	case TypeTXT:
		enumVal = "TXT"
	case TypeAAAA:
		enumVal = "AAA"
	case TypeSRV:
		enumVal = "SRV"
	}

	return enumVal
}

// MarshalJSON marshals Type into json.
func (t Type) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, t.String())), nil
}

// UnmarshalJSON unmarshals Level from json.
func (t *Type) UnmarshalJSON(b []byte) error {

	switch string(bytes.Trim(b, `"`)) {
	case "A":
		*t = TypeA

	case "NS":
		*t = TypeNS

	case "CNAME":
		*t = TypeCNAME

	case "PTR":
		*t = TypePTR

	case "MX":
		*t = TypeMX

	case "TXT":
		*t = TypeTXT

	case "AAA":
		*t = TypeAAAA

	case "SRV":
		*t = TypeSRV

	default:
		return errors.New("invalid Type")
	}

	return nil
}

//Record api data model
type Record struct {
	Type     Type   `json:"type"`
	TTL      uint32 `json:"ttl,omitempty"`
	Priority int    `json:"priority,omitempty"`
	Name     string `json:"name"`
	Content  string `json:"content"`
	Key      string `json:"key,omitempty"`
	Path     string `json:"-"`
}

func (r Record) ToEtcd() (*Etcd, error) {
	var e = new(Etcd)
	if r.TTL != 0 {
		e.TTL = r.TTL
	}
	if r.Priority != 0 {
		e.Priority = r.Priority
	}
	switch r.Type {
	case TypeMX:
		e.Mail = true
		fallthrough
	case TypeA, TypeCNAME, TypePTR, TypeAAAA:
		e.Host = r.Content
	// case TypeNS:

	case TypeTXT:
		e.Text = r.Content
	case TypeSRV:

		matchs := reSRV.FindStringSubmatch(r.Content)
		if len(matchs) != 4 {
			return nil, errors.New("Content field format is incorrect")
		}
		e.Weight, _ = strconv.Atoi(matchs[1])
		e.Port, _ = strconv.Atoi(matchs[2])
		e.Host = matchs[3]
	default:
		return nil, errors.New("Type field invalid")
	}

	if !strings.HasSuffix(r.Path, "/") {
		r.Path += "/"
	}
	if r.Type == TypePTR {
		r.Path += "arpa/in-addr/"
	}
	keys := strings.Split(r.Name, ".")
	for i, j := 0, len(keys)-1; i < j; i, j = i+1, j-1 {
		keys[i], keys[j] = keys[j], keys[i]
	}
	e.Key = r.Path + strings.Join(keys, "/")
	return e, nil
}
