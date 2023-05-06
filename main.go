package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Empleado struct {
	Id     uuid.UUID `json:"id" gorm:"type:uuid;default:uuid_generate_v4()"`
	Nombre string    `gorm:"not null"`
	Correo string    `gorm:"not null"`
}

func main() {
	DB_URL := os.Getenv("DB_URL")
	if DB_URL == "" {
		DB_URL = "postgresql://root@localhost"
	}

	DB_PORT := os.Getenv("DB_PORT")
	if DB_PORT == "" {
		DB_PORT = "26257"
	}

	DB_INIT := os.Getenv("DB_INIT")
	if DB_INIT == "" {
		DB_INIT = "true"
	}

	r := gin.Default()
	dbURL := fmt.Sprintf("%s:%s/defaultdb?sslmode=disable&application_name=$ demos_golang", DB_URL, DB_PORT)
	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	if DB_INIT != "false" {
		db.AutoMigrate(&Empleado{})
	}

	r.GET("/", func(c *gin.Context) {
		var empleados []Empleado
		db.Find(&empleados)
		fmt.Println(empleados)
		c.HTML(http.StatusOK, "inicio.html", gin.H{
			"empleados": empleados,
		})
	})

	r.GET("/crear", func(c *gin.Context) {
		c.HTML(http.StatusOK, "crear.html", nil)
	})

	r.POST("/insertar", func(c *gin.Context) {
		nombre := c.PostForm("nombre")
		correo := c.PostForm("correo")
		empleado := Empleado{Nombre: nombre, Correo: correo}
		db.Create(&empleado)
		c.Redirect(http.StatusMovedPermanently, "/")
	})

	r.GET("/borrar/:id", func(c *gin.Context) {
		id := c.Param("id")
		_uuid := uuid.MustParse(id)
		fmt.Println(_uuid)
		empleado := Empleado{Id: _uuid}
		db.First(&empleado)
		db.Delete(&empleado)
		c.Redirect(http.StatusMovedPermanently, "/")
	})

	r.GET("/editar/:id", func(c *gin.Context) {
		id := c.Param("id")
		fmt.Println(id)
		empleado := Empleado{Id: uuid.MustParse(id)}
		db.First(&empleado)
		c.HTML(http.StatusOK, "editar.html", gin.H{
			"empleado": empleado,
		})
	})

	r.POST("/actualizar", func(c *gin.Context) {
		id := c.PostForm("id")
		empleado := Empleado{Id: uuid.MustParse(id)}
		db.First(&empleado)
		if empleado.Id == uuid.Nil {
			fmt.Printf("No se encontró el empleado con el ID %d", empleado.Id) // Imprime un mensaje en la consola
			c.String(http.StatusBadRequest, "No se encontró el empleado con el ID especificado")
			return
		}
		nombre := c.PostForm("nombre")
		correo := c.PostForm("correo")
		empleado.Nombre = nombre
		empleado.Correo = correo
		db.Save(&empleado)
		c.Redirect(http.StatusMovedPermanently, "/")
	})

	r.LoadHTMLGlob("plantillas/*")
	r.Run(":3000")
}
