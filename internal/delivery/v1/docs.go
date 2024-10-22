package v1

import (
	"errors"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sixojke/test-astral/domain"
	"github.com/sixojke/test-astral/pkg/logger"
)

type uploadDocumentInpMeta struct {
	Name         string   `form:"name"`
	IsFile       bool     `form:"is_file"`
	IsPublic     bool     `form:"public"`
	Mime         string   `form:"mime"`
	Grants       []string `form:"grant[]"`
	DocumentData string   `form:"json"`
}

func (u *uploadDocumentInpMeta) validate() error {
	if u.Name == "" {
		return domain.ErrNameIsEmpty
	}

	return nil
}

type uploadDocumentData struct {
	DocumentData interface{} `json:"json"`
	File         string      `json:"file"`
}

// @Summary Upload document
// @Security UsersAuth
// @Tags docs
// @Description Upload document
// @ModuleID uploadDocument
// @Accept multipart/form-data
// @Produce json
// @Param name formData string false "Document name"
// @Param is_file formData bool false "Is file"
// @Param public formData bool false "Is public"
// @Param mime formData string false "Document mime type"
// @Param grant[] formData string false "Grant array"
// @Param json formData string false "Document data"
// @Param file formData file false "Document file"
// @Success 200 {object} swagData{data=uploadDocumentData} "Document uploaded successfully"
// @Failure 400 {object} swagError "Bad Request"
// @Failure 500 {object} swagError "Internal Server Error"
// @Router /docs [post]
func (h *Handler) uploadDocument(c *gin.Context) {
	var inp uploadDocumentInpMeta
	if err := c.ShouldBind(&inp); err != nil {
		errResponse(c, http.StatusBadRequest, err.Error(), domain.ErrInvalidMetaData.Error())

		return
	}
	logger.Debugf("%v", inp)

	if err := inp.validate(); err != nil {
		errResponse(c, http.StatusBadRequest, err.Error(), err.Error())

		return
	}

	userId := getUserIdByContext(c)

	var fileName string
	var filePath string
	if inp.IsFile {
		file, err := c.FormFile("file")
		if err != nil {
			errResponse(c, http.StatusBadRequest, err.Error(), domain.ErrFileNotFound.Error())

			return
		}

		if file.Size > h.config.HTTPServer.MaxFileSizeMb<<20 {
			errResponse(c, http.StatusBadRequest, domain.ErrFileIsTooLarge.Error(), domain.ErrFileIsTooLarge.Error())

			return
		}

		filePath = h.filePathGenerator(userId, file.Filename)

		if fileExists(filePath) {
			errResponse(c, http.StatusBadRequest, domain.ErrFileThisNameIsAlready.Error(), domain.ErrFileThisNameIsAlready.Error())

			return
		}

		if err := c.SaveUploadedFile(file, filePath); err != nil {
			errResponse(c, http.StatusBadRequest, err.Error(), domain.ErrInternalServerError.Error())

			return
		}

		fileName = file.Filename
	}

	if err := h.service.Document.Create(&domain.Document{
		Name:         inp.Name,
		Mime:         inp.Mime,
		FilePath:     filePath,
		IsFile:       inp.IsFile,
		IsPublic:     inp.IsPublic,
		DocumentData: inp.DocumentData,
		Grants:       inp.Grants,
	}, userId); err != nil {
		errResponse(c, http.StatusInternalServerError, err.Error(), domain.ErrInternalServerError.Error())

		if inp.IsFile {
			if err = os.Remove(filePath); err != nil {
				logger.Errorf("failed to delete file: %v", err)
			}
		}

		return
	}

	newResponse(c, http.StatusOK, uploadDocumentData{
		DocumentData: inp.DocumentData,
		File:         fileName,
	}, nil)
}

type getDocumentsData struct {
	Documents *[]domain.Document `json:"docs"`
}

// @Summary Get documents
// @Security UsersAuth
// @Tags docs
// @Description Get documents by user
// @ModuleID getDocuments
// @Accept json
// @Produce json
// @Param login query string false "User login"
// @Param key query string false "Key for filter"
// @Param value query string false "Value for filter"
// @Param limit query int false "Limit for pagination"
// @Param page query int false "Page for pagination"
// @Success 200 {object} swagData{data=getDocumentsData} "Documents list"
// @Failure 400 {object} swagError "Bad Request"
// @Failure 404 {object} swagError "User not found"
// @Failure 500 {object} swagError "Internal Server Error"
// @Router /docs [get]
func (h *Handler) getDocuments(c *gin.Context) {
	filterParams := domain.PrepareFillterParams(c.Query("key"), c.Query("value"), c.Query("limit"), c.Query("page"))

	documents, err := h.service.Document.GetByUser(c.Query("login"), getUserIdByContext(c), filterParams)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			errResponse(c, http.StatusNotFound, err.Error(), err.Error())
		} else {
			errResponse(c, http.StatusInternalServerError, err.Error(), domain.ErrInternalServerError.Error())
		}

		return
	}

	newResponse(c, http.StatusOK, getDocumentsData{
		Documents: documents,
	}, nil)
}

// @Summary Get document by ID
// @Security UsersAuth
// @Tags docs
// @Description Get document by ID
// @ModuleID getDocument
// @Accept json
// @Produce json
// @Param id path string true "Document ID"
// @Success 200 {object} swagData{data=domain.Document} "Document"
// @Failure 400 {object} swagError "Bad Request"
// @Failure 404 {object} swagError "Document not found"
// @Failure 500 {object} swagError "Internal Server Error"
// @Router /docs/{id} [get]
func (h *Handler) getDocument(c *gin.Context) {
	documentId := c.Param("id")

	if documentId == "" {
		errResponse(c, http.StatusBadRequest, domain.ErrParameterIsEmpty.Error(), domain.ErrParameterIsEmpty.Error())

		return
	}

	document, err := h.service.Document.GetById(documentId, getUserIdByContext(c))
	if err != nil {
		if errors.Is(err, domain.ErrDocumentNotFound) {
			errResponse(c, http.StatusNotFound, err.Error(), err.Error())
		} else {
			errResponse(c, http.StatusInternalServerError, err.Error(), domain.ErrInternalServerError.Error())
		}

		return
	}

	if !document.IsFile {
		newResponse(c, http.StatusOK, document, nil)

		return
	}

	if !fileExists(document.FilePath) {
		errResponse(c, http.StatusNotFound, domain.ErrFileIsDamagedOrNotFound.Error(), domain.ErrFileIsDamagedOrNotFound.Error())

		return
	}

	c.File(document.FilePath)
}

// @Summary Check document by ID
// @Security UsersAuth
// @Tags docs
// @Description Check document by ID
// @ModuleID checkDocument
// @Accept json
// @Produce json
// @Param id path string true "Document ID"
// @Success 200 {object} nil "Success"
// @Failure 400 {object} swagError "Bad Request"
// @Failure 404 {object} swagError "Document not found"
// @Failure 500 {object} swagError "Internal Server Error"
// @Router /docs/{id} [head]
func (h *Handler) checkDocument(c *gin.Context) {
	documentId := c.Param("id")

	if documentId == "" {
		errResponse(c, http.StatusBadRequest, domain.ErrParameterIsEmpty.Error(), domain.ErrParameterIsEmpty.Error())

		return
	}

	exists, err := h.service.Document.CheckById(documentId, getUserIdByContext(c))
	if err != nil {
		if errors.Is(err, domain.ErrDocumentNotFound) {
			errResponse(c, http.StatusNotFound, err.Error(), err.Error())
		} else {
			errResponse(c, http.StatusInternalServerError, err.Error(), domain.ErrInternalServerError.Error())
		}

		return
	}

	if !exists {
		errResponse(c, http.StatusNotFound, domain.ErrDocumentNotFound.Error(), domain.ErrDocumentNotFound.Error())

		return
	}

	newResponse(c, http.StatusOK, nil, nil)
}

// @Summary Delete document by ID
// @Security UsersAuth
// @Tags docs
// @Description Delete document by ID
// @ModuleID deleteDocument
// @Accept json
// @Produce json
// @Param id path string true "Document ID"
// @Success 200 {object} swagResponse{response=map[string]bool} "Success"
// @Failure 400 {object} swagError "Bad Request"
// @Failure 404 {object} swagError "Document not found"
// @Failure 500 {object} swagError "Internal Server Error"
// @Router /docs/{id} [delete]
func (h *Handler) deleteDocument(c *gin.Context) {
	documentId := c.Param("id")

	if documentId == "" {
		errResponse(c, http.StatusBadRequest, domain.ErrParameterIsEmpty.Error(), domain.ErrParameterIsEmpty.Error())

		return
	}

	if err := h.service.Document.Delete(documentId, getUserIdByContext(c)); err != nil {
		if errors.Is(err, domain.ErrDocumentNotFound) {
			errResponse(c, http.StatusBadRequest, err.Error(), err.Error())
		} else {
			errResponse(c, http.StatusInternalServerError, err.Error(), domain.ErrInternalServerError.Error())
		}

		return
	}

	newResponse(c, http.StatusOK, nil, map[string]bool{
		documentId: true,
	})
}
