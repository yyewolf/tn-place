package middlewares

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/markbates/goth"
	"github.com/yyewolf/tn-place/back/internal/education"
	"github.com/yyewolf/tn-place/back/internal/teams"
	"golang.org/x/oauth2"
	admin "google.golang.org/api/admin/directory/v1"
	"google.golang.org/api/option"
)

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

	edc.Team = teams.FindTeam(user.LastName, user.FirstName)

	ctx.Set("education", edc)
}
