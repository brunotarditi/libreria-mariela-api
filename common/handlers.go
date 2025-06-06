package common

import (
	"libreria/requests"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Get[T any](ops Operations[T]) gin.HandlerFunc {
	return func(c *gin.Context) {
		models, err := ops.FindAll()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, models)
	}
}

func GetByID[T any](ops Operations[T]) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		model, err := ops.FindByID(id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Entidad no encontrada"})
			return
		}
		c.JSON(http.StatusOK, model)
	}
}

func Create[T any, R requests.MapperRequest[T]](ops Operations[T]) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request R
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if validatable, ok := any(request).(requests.ValidateRequest); ok {
			if err := validatable.Validate(ops.(*GormOperations[T]).db); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
		}

		model, err := mapper(request)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := ops.Create(model); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, model)
	}
}

func Update[T any, R requests.MapperRequest[T]](ops Operations[T]) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		existing, err := ops.FindByID(id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Entidad no encontrada"})
			return
		}

		var request R
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if validatable, ok := any(request).(requests.ValidateRequest); ok {
			if err := validatable.Validate(ops.(*GormOperations[T]).db); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
		}

		model, err := mapper(request, existing)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := ops.Update(model); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, model)
	}
}

func Delete[T any](ops Operations[T]) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if err := ops.Delete(id); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Entidad no encontrada"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Entidad eliminada"})
	}
}

func mapper[T any, R requests.MapperRequest[T]](req R, existing ...T) (T, error) {
	if len(existing) > 0 {
		return req.UpdateModel(existing[0])
	}
	return req.ToModel()
}
