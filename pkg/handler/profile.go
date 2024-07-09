package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gordejka179/t-bmstu/pkg/database"
)

func (h *Handler) profileMainPage(c *gin.Context) {
	profile, err := database.GetInfoForProfilePage(c.GetString("username"))

	if err != nil {
		// TODO return error
		return
	}

	c.HTML(http.StatusOK, "profile.tmpl", gin.H{
		"NickName": profile.Username,
		"Surname":  profile.LastName,
		"Name":     profile.FirstName,
		"Email":    profile.Email,
	})
}
