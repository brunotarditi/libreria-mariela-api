package requests

import "errors"

type BulkDeleteRequest struct {
	IDs []uint `json:"ids" binding:"required,min=1"`
}

func (r BulkDeleteRequest) Validate() error {
	if len(r.IDs) == 0 {
		return errors.New("debe proporcionar al menos un ID para eliminar")
	}
	for _, id := range r.IDs {
		if id == 0 {
			return errors.New("los IDs deben ser números enteros positivos mayores a cero")
		}
	}
	return nil
}
