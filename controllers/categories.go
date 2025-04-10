package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

/* GetEducationLevels returns a list of education levels */
func GetEducationLevels(c *gin.Context) {
	levels := []string{
		"Grade 8",
        "Grade 7",
        "Grade 6",
        "Grade 5",
        "Grade 4",
        "Grade 3",
        "Grade 2",
        "Grade 1",
        "Playgroup",
        "PP1",
        "PP2",
        "Form 1",
        "Form 2",
        "Form 3",
        "Form 4",
	}

	c.JSON(http.StatusOK, gin.H{"levels": levels})
}