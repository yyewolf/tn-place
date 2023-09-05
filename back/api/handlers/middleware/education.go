package middlewares

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/markbates/goth"
	"github.com/yyewolf/tn-place/back/internal/education"
	"golang.org/x/oauth2"
	admin "google.golang.org/api/admin/directory/v1"
	"google.golang.org/api/option"
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

func SetEducation(ctx *gin.Context, user goth.User) {
	providerI, _ := ctx.Get("provider")
	provider := providerI.(goth.Provider)

	c := context.Background()
	oconfig := &oauth2.Config{}
	token, err := provider.RefreshToken(user.RefreshToken)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error refreshing token"})
		return
	}
	adminService, err := admin.NewService(c, option.WithTokenSource(oconfig.TokenSource(c, token)))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	t, err := adminService.Users.Get(user.UserID).Projection("custom").CustomFieldMask("Education").ViewType("domain_public").Do()
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	edc := &education.Education{}
	err = json.Unmarshal(t.CustomSchemas["Education"], edc)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	edc.Team = FindTeam(user.LastName, user.FirstName)

	ctx.Set("education", edc)
}
