package common

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
)

type Kudos struct {
	TraceId     string    `json:"traceId"`       // Per project
	Id          string    `json:"id"`            // unique kudos id
	Identifier  string    `json:"identifier"`    // e.g., email
	Name        string    `json:"name"`          // author name
	Ts          time.Time `json:"ts"`            // generation time
	Weight      float64   `json:"weight,string"` // contribution to the entire project
	PackageName string    `json:"packageName"`   // package name
	Description string    `json:"description"`   // contributing dependency
	Type        string    `json:"type"`          // code
}

func GenerateKudos(p *Project) []Kudos {
	subjectsPrefix := "email"
	kudos := []Kudos{}
	traceId := NewRandomId()
	for _, d := range p.dependencies {
		packageName := ExtractPackageName(d.id)
		for _, c := range d.contributors {
			if c.score > 0 {
				kudos = append(kudos, Kudos{
					traceId,
					NewRandomId(),
					fmt.Sprintf("%s:%s", subjectsPrefix, c.email),
					c.name,
					time.Now().UTC().Truncate(time.Second),
					ToFixed(c.score, 6),
					packageName,
					packageName,
					"code",
				})
			}
		}
	}

	return kudos
}

func (k *Kudos) ToJSON() []byte {
	jsonData, err := json.Marshal(k)
	if err != nil {
		log.Println("Error marshaling to JSON:", err)
		panic(err)
	}

	return jsonData
}
