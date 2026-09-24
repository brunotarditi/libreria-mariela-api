package common

import (
	"fmt"
	"libreria/requests"
	"libreria/responses"
	"libreria/utils"
	"net/http"
	"strconv"

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

func Paginated[T any](ops Operations[T]) gin.HandlerFunc {
	return func(c *gin.Context) {
		page, _ := strconv.Atoi(c.DefaultQuery("page", "0"))
		size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

		options := QueryOptions{
			Search:    c.Query("search"),
			Sort:      c.DefaultQuery("sort", "id"),
			Direction: c.DefaultQuery("direction", "asc"),
		}

		offset := page * size

		total, err := ops.Count(options)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		models, err := ops.Paginated(offset, size, options)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data":  models,
			"total": total,
			"page":  page,
			"size":  size,
		})
	}
}

func Create[T any, R requests.MapperRequest[T]](ops Operations[T]) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request R
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if validator, ok := any(request).(requests.ValidateRequest); ok {
			if err := validator.Validate(ops.(*GormOperations[T]).db); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
		}

		model, err := mapper(request)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		createdModel, err := ops.Create(model)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, createdModel)
	}
}

func CreateMany[T any, R interface {
	requests.MapperArrayRequest[T]
	requests.ValidateRequest
}](ops Operations[T]) gin.HandlerFunc {
	return func(c *gin.Context) {
		var requests R
		if err := c.ShouldBindJSON(&requests); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := requests.Validate(ops.(*GormOperations[T]).db); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		models, err := mapperArray(requests)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := ops.CreateMany(models); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, models)
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

		if validator, ok := any(request).(requests.ValidateOnUpdate[T]); ok {
			if err := validator.ValidateUpdate(ops.(*GormOperations[T]).db, existing); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
		}

		model, err := mapper(request, existing)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		updatedModel, err := ops.Update(model)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, updatedModel)
	}
}

func Delete[T any](ops Operations[T]) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if err := ops.Delete(id); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "No encontrado"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Eliminado con éxito"})
	}
}

func BulkDelete[T any](ops Operations[T]) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req requests.BulkDeleteRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos: " + err.Error()})
			return
		}

		if err := req.Validate(); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		deletedCount, err := ops.DeleteMany(req.IDs)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		var message string
		if deletedCount == 1 {
			message = "Se eliminó 1 registro con éxito"
		} else {
			message = fmt.Sprintf("Se eliminaron %d registros con éxito", deletedCount)
		}

		c.JSON(http.StatusOK, responses.BulkDeleteResponse{
			Success:      true,
			Message:      message,
			DeletedCount: deletedCount,
		})
	}
}

type ExcelExportable interface {
	ExcelHeaders() []string
	ExcelRow() []interface{}
}

func Export[T ExcelExportable](ops Operations[T], entityName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		items, err := ops.FindAll()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		var headers []string
		var rows [][]interface{}

		if len(items) > 0 {
			headers = items[0].ExcelHeaders()
			rows = make([][]interface{}, len(items))
			for i, item := range items {
				rows[i] = item.ExcelRow()
			}
		} else {
			var dummy T
			headers = dummy.ExcelHeaders()
			rows = [][]interface{}{}
		}

		sheetName := utils.CapitalizeFirst(entityName)
		file, err := utils.BuildExcel(sheetName, headers, rows)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error al generar el archivo Excel: " + err.Error()})
			return
		}

		filename := fmt.Sprintf("%s_export.xlsx", entityName)
		if err := utils.StreamExcel(c, file, filename); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error al transmitir el archivo Excel: " + err.Error()})
			return
		}
	}
}

func mapper[T any, R requests.MapperRequest[T]](req R, existing ...T) (T, error) {
	if len(existing) > 0 {
		return req.UpdateModel(existing[0])
	}
	return req.ToModel()
}

func mapperArray[T any, R requests.MapperArrayRequest[T]](req R) ([]T, error) {
	return req.ToArrayModel()
}
