package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"infTest.com/inf/model"
)

func GetPersonInfo(c *gin.Context) {
	personID := c.Param("person_id")

	db, err := sql.Open("mysql", "root:password@/cetec")

	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	query := `SELECT p.name, ph.number as phone_number, a.city, a.state, a.street1, a.street2, a.zip_code 
			  FROM person p
			  JOIN phone ph ON p.id = ph.person_id
			  JOIN address_join aj ON p.id = aj.person_id
			  JOIN address a ON aj.address_id = a.id
			  WHERE p.id = ?`

	var person model.Person

	err = db.QueryRow(query, personID).Scan(&person.Name, &person.PhoneNumber, &person.City, &person.State, &person.Street1, &person.Street2, &person.ZipCode)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Person not found"})
		return
	}

	c.JSON(http.StatusOK, person)
}

func CreatePerson(c *gin.Context) {
	var newPerson model.Person

	if err := c.BindJSON(&newPerson); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db, err := sql.Open("mysql", "root:password@/cetec")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	insertPersonQuery := "INSERT INTO person (name) VALUES (?)"
	result, err := db.Exec(insertPersonQuery, newPerson.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	personID, err := result.LastInsertId()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	insertAddressJoinQuery := "INSERT INTO address_join (person_id, address_id) VALUES (?, ?)"
	_, err = db.Exec(insertAddressJoinQuery, personID, 1)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Person created successfully"})
}

func main() {
	r := gin.Default()

	r.GET("/person/:person_id/info", GetPersonInfo)
	r.POST("/person/create", CreatePerson)
	fmt.Println("server starting at port 8080 ...")
	r.Run(":8080")
}
