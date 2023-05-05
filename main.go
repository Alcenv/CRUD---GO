package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
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

	r := gin.Default()
	db, err := gorm.Open(mysql.Open("root:@tcp(127.0.0.1:3306)/sistema"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&Empleado{})

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

	db.AutoMigrate(&Empleado{}, &ModelWithoutCreatedAt{})
	r.LoadHTMLGlob("plantillas/*")
	r.Run(":8080")
}
