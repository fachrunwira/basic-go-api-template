package files

import (
	"bufio"
	"mime/multipart"

	"github.com/fachrunwira/basic-go-api-template/lib/response"
	"github.com/labstack/echo/v4"
)

type uploadDTO struct {
	File *multipart.FileHeader `form:"file"`
}

func (h *filesHandler) UploadFile(c echo.Context) error {
	var dto uploadDTO

	if err := c.Bind(&dto); err != nil {
		h.AppLogger.Printf("FilesHandler, UploadFile: %v", err)
		return response.BadRequest(c, "invalid body request.", err)
	}

	file, err := dto.File.Open()
	if err != nil {
		h.AppLogger.Printf("FilesHandler, ReadFile: %v", err)
		return response.InternalError(c, "failed to open file", err)
	}
	defer file.Close()

	var data []string
	scanFile := bufio.NewScanner(file)
	for scanFile.Scan() {
		data = append(data, scanFile.Text())
	}

	// responseInterface := []map[string]interface{}{}
	// for _, v := range data {
	// 	val := strings.Split(v, ",")
	// 	if
	// 	responseInterface = append(responseInterface, map[string]interface{}{
	// 		"code":  val[0],
	// 		"label": val[1],
	// 	})
	// }

	// fileBytes, err := io.ReadAll(file)
	// if err != nil {
	// 	h.AppLogger.Printf("FilesHandler, ReadFile: %v", err)
	// 	return response.InternalError(c, "failed to read file", err)
	// }

	// wilayah := strings.Split(string(fileBytes), ",")
	// mapWilayah := map[string]interface{}{}
	// for _, v := range wilayah {
	// }

	// fileBytes := make([]byte, dto.File.Size)
	// _, err = file.Read(fileBytes)
	// if err != nil && err != io.EOF {
	// 	return response.InternalError(c, "failed to read file", err)
	// }

	return response.Success(c, "berhasil", nil)
}
