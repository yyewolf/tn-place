package teams

import (
	"encoding/json"
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
