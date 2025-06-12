package controllers

import (
	"net/http"
	"os"

	"github.com/agl/fio/internal/domain/entities"
	"github.com/agl/fio/internal/domain/interfaces"
	"github.com/agl/fio/internal/presentation/extractors"
	. "github.com/agl/fio/pkg/logger"
	"github.com/gin-gonic/gin"
)

type PersonHandler struct {
	PORT    string
	service interfaces.PersonService
}

func NewPersonHandler(service interfaces.PersonService) *PersonHandler {
	PORT := os.Getenv("PORT")

	return &PersonHandler{PORT: PORT, service: service}
}

// @title People Library API
// @version 1.0.0
// @description API for managing people information
// @host localhost:6060
// @BasePath /
func (p *PersonHandler) StartApi() {
	r := gin.Default()

	r.GET("/person/:id", p.getPerson)
	r.GET("/person/filter", p.getPerson_Filter)

	r.DELETE("/person/:id", p.deletePerson)

	r.PATCH("/person/:id", p.updatePerson)

	r.POST("/person", p.createPerson)

	r.Run(":" + p.PORT)
}


// updatePerson godoc
// @Summary Partially update an existing person
// @Description Update person's information by ID. Only provided fields will be updated.
// @Tags People
// @Accept json
// @Produce json
// @Param id path int true "Person ID"
// @Param person body entities.Person true "Person object with fields to update"
// @Success 200 {object} map[string]responses.ResponseMessage
// @Failure 400 {object} map[string]string{error,details}
// @Failure 409 {object} map[string]string{error,details}
// @Router /person/{id} [patch]
func (p *PersonHandler) updatePerson(ctx *gin.Context) {
	id := ctx.Param("id")

	Log.Debug("Received request to update person", "id", id)

	var person entities.Person
	if err := ctx.ShouldBindJSON(&person); err != nil {
		Log.Info("Failed to bind JSON for updating person", "id", id, "error", err)

		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid input",
			"details": err.Error(),
		})

		return
	}

	if err := p.service.UpdatePersonByID(person, id); err != nil {
		Log.Info("Failed to update person in database", "id", id, "error", err)
		ctx.JSON(http.StatusConflict, gin.H{
			"error":   "Failed to update person",
			"details": err.Error(),
		})
		return
	}

	Log.Info("Successfully updated person", "id", id)

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Person updated successfully",
		"id":      id,
	})
}

// createPerson godoc
// @Summary Create a new person
// @Description Create a person entity and enrich it with age, gender, and nationality by name
// @Tags People
// @Accept json
// @Produce json
// @Param person body entities.Person true "Person object"
// @Success 201 {object} map[string]responses.ResponseMessage
// @Failure 400 {object} map[string]string{error,details}
// @Failure 409 {object} map[string]string{error,details}
// @Router /person [post]
func (p *PersonHandler) createPerson(ctx *gin.Context) {
	var received_person entities.ReceivedPerson

	Log.Debug("Received request to create person", "request", received_person)

	if err := ctx.ShouldBindJSON(&received_person); err != nil {
		Log.Info("Failed to bind JSON", "error", err)

		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid input",
			"details": err.Error(),
		})

		return
	}

	person, err := extractors.GetExtraUserInfoByName(received_person.Name, received_person)

	if err != nil {
		Log.Info("Failed to retrieve extra person data", "name", person.Name, "error", err)

		ctx.JSON(http.StatusConflict, gin.H{
			"error":   "Failed to retrieve extra person data",
			"details": err.Error(),
		})

		return
	}

	Log.Debug("Person enriched with extra data", "name", person.Name, "age", person.Age, "gender", person.Gender, "nationality", person.Nationality)

	id, err := p.service.CreatePerson(person)
	if err != nil {
		Log.Info("Failed to create person", "name", person.Name, "error", err)

		ctx.JSON(http.StatusConflict, gin.H{
			"error":   "Failed to create person",
			"details": err.Error(),
		})

		return
	}

	Log.Info("Person created successfully", "id", id)

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Person created successfully",
		"id":      id,
	})
}

// getPerson godoc
// @Summary Get person by ID
// @Description Get a single person by their ID
// @Tags People
// @Produce json
// @Param id path int true "Person ID"
// @Success 200 {object} map[string]responses.FoundPerson
// @Failure 400 {object} map[string]string{error,details}
// @Router /person/{id} [get]
func (p *PersonHandler) getPerson(ctx *gin.Context) {
	id := ctx.Param("id")
	Log.Debug("Received request for getPerson", "id", id)

	person, err := p.service.GetPersonByID(id)
	if err != nil {
		Log.Info("Failed to get person by ID", "id", id, "error", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Failed to get person by ID",
			"details": err.Error(),
		})
		return
	}

	Log.Info("Person successfully fetched", "id", id)
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Person received successfully",
		"data":    person,
	})
}

// getPerson_Filter godoc
// @Summary Get people by filter
// @Description Get a list of people filtered by parameters with pagination
// @Tags People
// @Produce json
// @Param name query string false "Name to filter by"
// @Param surname query string false "Surname to filter by"
// @Param patronymic query string false "Patronymic to filter by"
// @Param age query int false "Age to filter by"
// @Param gender query string false "Gender to filter by"
// @Param nationality query string false "Nationality to filter by"
// @Param page query int false "Page number"
// @Param limit query int false "Number of results per page"
// @Success 200 {object} map[string]responses.FoundPerson
// @Failure 400 {object} map[string]string{error,details}
// @Router /person/filter [get]
func (p *PersonHandler) getPerson_Filter(ctx *gin.Context) {
	page := ctx.Query("page")
	limit := ctx.Query("limit")

	name := ctx.DefaultQuery("name", "")
	surname := ctx.DefaultQuery("surname", "")
	patronymic := ctx.DefaultQuery("patronymic", "")
	age := ctx.DefaultQuery("age", "0")
	gender := ctx.DefaultQuery("gender", "")
	nationality := ctx.DefaultQuery("nationality", "")

	Log.Debug("Received request for getPerson_Age", "age", age, "page", page, "limit", limit)

	people, err := p.service.GetPeopleByFilter(name, surname, age, gender, nationality, page, limit, &patronymic)
	if err != nil {
		Log.Info("Failed to get people by age", "age", age, "error", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Failed to get people by age",
			"details": err.Error(),
		})
		return
	}

	Log.Info("Successfully fetched people by age", "age", age, "count", len(people))
	ctx.JSON(http.StatusOK, gin.H{
		"message": "People received successfully",
		"data":    people,
	})
}

// deletePerson godoc
// @Summary Delete person by ID
// @Description Delete a single person from the database by their ID
// @Tags People
// @Produce json
// @Param id path int true "Person ID"
// @Success 200 {object} map[string]responses.ResponseMessage
// @Failure 409 {object} map[string]string{error,details}
// @Router /person/{id} [delete]
func (p *PersonHandler) deletePerson(ctx *gin.Context) {
	id := ctx.Param("id")

	Log.Debug("Received request to delete person", "id", id)

	err := p.service.DeletePersonByID(id)
	if err != nil {
		Log.Info("Failed to delete person", "id", id, "error", err)

		ctx.JSON(http.StatusConflict, gin.H{
			"error":   "Failed to delete person",
			"details": err.Error(),
		})

		return
	}

	Log.Info("Successfully deleted person", "id", id)

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Person deleted successfully",
		"id":      id,
	})
}