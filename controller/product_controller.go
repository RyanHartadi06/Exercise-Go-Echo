package controller

import (
	"Go-Echo/config"
	"Go-Echo/constants"
	"Go-Echo/helper"
	"Go-Echo/model"
	"io"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
)

func GetProductController(e echo.Context) error {
	var products []model.Product

	err := config.DB.Find(&products).Error

	if err != nil {
		return e.JSON(http.StatusInternalServerError, model.Response{
			Message: "Error",
			Data:    nil,
		})
	}
	response := constants.Response{
		Message: "Success",
		Data:    products,
	}


	return e.JSON(http.StatusOK, response)
}

func DeleteProductController(e echo.Context) error {
	err := config.DB.Delete(&model.Product{Name: e.Param("name")}).Error

	if err != nil {
		return e.JSON(http.StatusInternalServerError, model.Response{
			Message: "Error",
			Data:    nil,
		})
	}
	response := constants.Response{
		Message: "Success",
	}


	return e.JSON(http.StatusOK, response)

}
func StoreProductController(e echo.Context) error {
	file, err := e.FormFile("file")
	if err != nil {
		return err
	}

	// Open the file
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	generateRandom := helper.GenerateRandomFileName()
	// Destination file path
	dstPath := "uploads/" + generateRandom

	// Create the destination file
	dst, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	// Copy the file content to the destination
	if _, err = io.Copy(dst, src); err != nil {
		return err
	}
	//return e.String(http.StatusOK, "File uploaded successfully.")
	product := model.Product{
		Name:  e.FormValue("name"),
		Stock: e.FormValue("stock"),
		Image: dstPath,
	}

	e.Bind(&product)

	errSave := config.DB.Save(&product).Error
	if errSave != nil {
		return e.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": err.Error(),
		})
	}
	response := constants.Response{
		Message: "Success",
		Data:    product,
	}

	return e.JSON(http.StatusOK, response)
}
