package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Empleado struct {
	Id        uint      `gorm:"primaryKey"`
	Nombre    string    `gorm:"not null"`
	Correo    string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime:milli"`
}

type ModelWithoutCreatedAt struct {
	Id        uint `gorm:"primaryKey"`
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type ActualizarEmpleado struct {
	Nombre string `form:"nombre"`
	Correo string `form:"correo"`
	Id     uint   `form:"id"`
}

func (ModelWithoutCreatedAt) TableName() string {
	return "empleados"
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
		db.AutoMigrate(&Empleado{}, &ModelWithoutCreatedAt{})
	}

	r.GET("/", func(c *gin.Context) {
		var empleados []Empleado
		db.Find(&empleados)
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
		var empleado Empleado
		db.First(&empleado, id)
		db.Delete(&empleado)
		c.Redirect(http.StatusMovedPermanently, "/")
	})

	r.GET("/editar/:id", func(c *gin.Context) {
		id := c.Param("id")
		var empleado Empleado
		db.First(&empleado, id)
		c.HTML(http.StatusOK, "editar.html", gin.H{
			"empleado": empleado,
		})
	})

	r.POST("/actualizar", func(c *gin.Context) {
		var empleadoActualizar ActualizarEmpleado
		if err := c.ShouldBind(&empleadoActualizar); err != nil {
			fmt.Println(err) // Imprime el error en la consola
			c.String(http.StatusBadRequest, "Error al obtener los datos del empleado a actualizar")
			return
		}

		var empleado Empleado
		db.First(&empleado, empleadoActualizar.Id)
		if empleado.Id == 0 {
			fmt.Printf("No se encontró el empleado con el ID %d", empleadoActualizar.Id) // Imprime un mensaje en la consola
			c.String(http.StatusBadRequest, "No se encontró el empleado con el ID especificado")
			return
		}

		empleado.Nombre = empleadoActualizar.Nombre
		empleado.Correo = empleadoActualizar.Correo
		db.Save(&empleado)
		c.Redirect(http.StatusMovedPermanently, "/")
	})

	r.LoadHTMLGlob("plantillas/*")
	r.Run(":3000")
}
