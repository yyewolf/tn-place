package teams

import (
	"encoding/json"
	"os"
)

// TeamFile
type TeamFile map[string]int

// Filter from integration :
// https://jsoneditoronline.org/
//
// function query (data) {
// 	let out = {};
// 	let equipes = [];
// 	_.chain(data)
// 	  .filter(item => item?.["Année"] == 1 && item?.["R1/R2/R3"] == "")
// 	  .map(item => {
// 		let e = equipes.indexOf(item?.["FIELD1"]);
// 		if (e == -1) equipes.push(item?.["FIELD1"]); e = equipes.indexOf(item?.["FIELD1"]);

// 		out[item?.["Nom"] + " " + item?.["Prénom"]] = e + 1
// 	  })
// 	  .value()
// 	return out;
//   }

func FindTeam(lastName, firstName string) int {
	// Open teams.json and find the team of the student
	d, err := os.ReadFile("teams.json")
	if err != nil {
		return -1
	}

	var teams TeamFile
	err = json.Unmarshal(d, &teams)
	if err != nil {
		return -1
	}

	team, ok := teams[lastName+" "+firstName]
	if !ok {
		return -1
	}

	return team
}
