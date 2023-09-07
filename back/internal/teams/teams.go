package teams

import (
	"encoding/json"
	"image/color"
	"os"
)

// TeamFile
type TeamFile map[string]string

// Filter from integration :
// https://jsoneditoronline.org/
//
// function query (data) {
// 	let out = {};
// 	_.chain(data)
// 	  .filter(item => item?.["Rang"] == "")
// 	  .map(item => {
// 		out[item?.["Nom"] + " " + item?.["Prénom"]] = item?.["Équipe"]
// 	  })
// 	  .value()
// 	return out;
//   }

var Colors = map[string]color.Color{
	"Jaune or":        color.NRGBA{255, 215, 0, 255},
	"Rouge écarlate":  color.NRGBA{237, 0, 0, 255},
	"Noir anthracite": color.NRGBA{0, 0, 0, 255},
	"Vert amande":     color.NRGBA{193, 217, 188, 255},
	"Bleu ciel":       color.NRGBA{135, 206, 235, 255},
	"Rose fuchsia":    color.NRGBA{252, 64, 138, 255},
}

func FindTeam(lastName, firstName string) string {
	// Open teams.json and find the team of the student
	d, err := os.ReadFile("teams.json")
	if err != nil {
		return ""
	}

	var teams TeamFile
	err = json.Unmarshal(d, &teams)
	if err != nil {
		return ""
	}

	team, ok := teams[lastName+" "+firstName]
	if !ok {
		return ""
	}

	return team
}
