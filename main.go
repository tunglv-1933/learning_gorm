package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type course struct {
	gorm.Model
	ID    int `gorm:"primaryKey;autoIncrement"`
	Title string
}

type courseContent struct {
	gorm.Model
	ID          int `gorm:"primaryKey;autoIncrement"`
	Title       string
	Description string
}

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		panic("Error loading .env file")
	}
	mysqlUserName := os.Getenv("MYSQL_USERNAME")
	mysqlPassword := os.Getenv("MYSQL_PASSWORD")
	databaseName := os.Getenv("DATABASE_NAME")
	dsn := fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local", mysqlUserName, mysqlPassword, databaseName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&course{}, &courseContent{})

	router := gin.Default()
	router.GET("/courses", func(c *gin.Context) { getCourses(c, db) })
	router.POST("/courses", func(c *gin.Context) { createCourse(c, db) })
	router.GET("/course/:id", func(c *gin.Context) { getCourse(c, db) })
	router.PUT("/course/:id", func(c *gin.Context) { updateCourse(c, db) })
	router.DELETE("/course/:id", func(c *gin.Context) { deleteCourse(c, db) })

	router.Run(":8000")
}

func getCourses(c *gin.Context, db *gorm.DB) {
	var courses []course
	db.Find(&courses)
	c.JSON(200, courses)
}

func createCourse(c *gin.Context, db *gorm.DB) {
	var newCourse course

	if err := c.ShouldBind(&newCourse); err == nil {
		db.Create(&newCourse)
		c.JSON(200, newCourse)
	} else {
		c.JSON(200, gin.H{"message": "error!"})
	}
	return
}

func getCourse(c *gin.Context, db *gorm.DB) {
	var findCourse course
	db.Where("id = ?", c.Param("id")).Find(&findCourse)

	if findCourse.ID > 0 {
		c.JSON(200, findCourse)
	} else {
		c.JSON(200, gin.H{"message": "not found!"})
	}
	return
}

func updateCourse(c *gin.Context, db *gorm.DB) {
	var findCourse course
	db.Where("id = ?", c.Param("id")).Find(&findCourse)

	if findCourse.ID > 0 {
		c.Bind(&findCourse)
		db.Save(&findCourse)
		c.JSON(200, findCourse)
	} else {
		c.JSON(200, gin.H{"message": "not found!"})
	}

	return
}

func deleteCourse(c *gin.Context, db *gorm.DB) {

	var findCourse course
	db.Where("id = ?", c.Param("id")).Find(&findCourse)

	if findCourse.ID > 0 {
		db.Delete(&findCourse)
		c.JSON(200, gin.H{"message": "delete successfull"})
	} else {
		c.JSON(200, gin.H{"message": "not found!"})
	}

	return
}
