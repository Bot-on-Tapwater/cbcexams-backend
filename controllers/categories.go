package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

/* GetCategories returns a list of education levels */
// GetCategories handles the HTTP request to retrieve categorized educational data.
// It responds with a JSON object containing the following:
// - levels: A list of educational levels such as "Grade 9", "Grade 8", etc.
// - education levels: A mapping of broader education categories (e.g., "Pre-Primary", "High School")
//   to their respective levels.
// - resource types: A mapping of education levels or roles (e.g., "High School", "Teacher")
//   to the types of resources available (e.g., "KCSE Papers", "Lesson Plans").
// - subjects: A mapping of education levels (e.g., "High School", "Pre-Primary") to the subjects
//   taught at those levels (e.g., "Mathematics", "English").
//
// Parameters:
// - c: The Gin context, which provides request and response handling.
//
// Response:
// - HTTP 200 OK: A JSON object containing the categorized educational data.
func GetCategories(c *gin.Context) {
	levels := []string{
		"Grade 9",
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

	educationLevels := map[string][]string{
		"Pre-Primary":   {"Playgroup", "PP1", "PP2"},
		"Lower Primary": {"Grade 1", "Grade 2", "Grade 3"},
		"Upper Primary": {"Grade 4", "Grade 5", "Grade 6"},
		"Junior School": {"Grade 7", "Grade 8", "Grade 9"},
		"Senior School": {"Grade 10", "Grade 11", "Grade 12"},
		"High School":   {"Form 1", "Form 2", "Form 3", "Form 4"},
	}

	resourceTypesEducationLevels := map[string][]string{
		"All Education Levels": {"Opener Exam", "Mid Term Exam", "End Term Exam", "Schemes of Work", "Lesson Plan", "Notes", "Assignment", "Topic-tests ", "Lesson Plans", "Syllabus", "Schemes of Work", "Study Guide", "Marking Scheme", "Design-Material"},
		"High School":          {"KCSE", "Mock"},
		"Upper Primary":        {"KPSEA"},
		"Teacher":              {"Lesson Plans", "Syllabus", "Schemes of Work", "Study Guide", "Marking Scheme", "Design-Material"},
		"Misc":                 {"Assessment Book", "Record of Work", "CBC Assessment Rubric"},
	}

	resourceTypeCategories := map[string][]string{
		"Exams and Past Papers": {
			"Opener Exam",
			"Mid Term Exam",
			"End Term Exam",
			"KCSE",
			"Mock",
			"KPSEA",
			"Topic-tests",
		},
		"Teacher's Resources": {
			"Schemes of Work",
			"Lesson Plan",
			"Syllabus",
			"Study Guide",
			"Marking Scheme",
			"Design-Material",
			"Record of Work",
			"CBC Assessment Rubric",
		},
		"Notes": {
			"Notes",
			"Assignment",
			"Study Guide", // included again for relevance
		},
		"Other": {
			"Assessment Book",
			"Design-Material", // repeated where fitting
		},
	}

	subjects := map[string][]string{
		"High School":      {"Mathematics", "English", "Kiswahili", "Biology", "Chemistry", "Physics", "History & Government", "Geography", "Christian Religious Education", "Islamic Religious Education", "Hindu Religious Education", "Business Studies", "Agriculture", "Computer Studies", "Home Science", "Art & Design", "Music", "French", "German", "Arabic", "Aviation Technology", "Woodwork", "Metalwork"},
		"Pre-Primary":      {"Language Activities", "English", "Kiswahili", "Mathematical Activities", "Environmental Activities", "Psychomotor & Creative Activities", "Art", "Music", "Movement", "Christian Religious Education", "Islamic Religious Education", "Hindu Religious Education", "Pastoral Instruction"},
		"Lower Primary":    {"English", "Kiswahili", "Mathematics", "Environmental Activities", "Hygiene & Nutrition", "Christian Religious Education", "Islamic Religious Education", "Hindu Religious Education", "Movement & Creative Arts", "Music", "Art", "Physical Education"},
		"Upper-Primary":    {"English", "Kiswahili", "Mathematics", "Science & Technology", "Social Studies", "History", "Geography", "Citizenship", "Christian Religious Education", "Islamic Religious Education", "Hindu Religious Education"},
		"Junior-Secondary": {"English", "Kiswahili", "Mathematics", "Integrated Science", "Health Education", "Pre-Technical Studies", "Social Studies", "History", "Geography", "Civics", "Business Studies", "Christian Religious Education", "Islamic Religious Education", "Hindu Religious Education", "Agriculuture", "Life Skills", "Computer Science", "Performing Arts", "Music", "Drama", "Visual Arts", "Art & Design", "French", "German", "Arabic", "Kenyan Sign Language"},
	}

	c.JSON(http.StatusOK, gin.H{"levels": levels, "education levels": educationLevels, "resource types education level": resourceTypesEducationLevels, "resource types categories": resourceTypeCategories, "subjects": subjects})
}
