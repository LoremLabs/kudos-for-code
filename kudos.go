package main

import (
	"encoding/json"
	"fmt"
	"time"
)

type Kudos struct {
	TraceId     string    `json:"traceId"`     // Per project
	Id          string    `json:"id"`          // unique kudos id
	Identifier  string    `json:"identifier"`  // e.g., email
	Ts          time.Time `json:"ts"`          // generation time
	Weight      float32   `json:"weight"`      // contribution to the entire project
	Description string    `json:"description"` // contributing dependency
}

func GenerateKudos(p *Project) []Kudos {
	kudos := []Kudos{}
	traceId := NewRandomId()
	for _, d := range p.dependencies {
		for _, c := range d.contributors {
			if c.score > 0 {
				kudos = append(kudos, Kudos{
					traceId,
					NewRandomId(),
					fmt.Sprintf("did:kudos:email:%s", c.email),
					time.Now().UTC(),
					c.score,
					fmt.Sprintf("%s contribution", d.id),
				})
			}
		}
	}

	return kudos
}

func (k *Kudos) ToJSON() []byte {
	jsonData, err := json.Marshal(k)
	if err != nil {
		fmt.Println("Error marshaling to JSON:", err)
		panic(err)
	}

	return jsonData
}
