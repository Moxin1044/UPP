package api

import (
	"strconv"

	"upp/internal/service"

	"github.com/gin-gonic/gin"
)

type ProjectHandler struct {
	projectService *service.ProjectService
}

func NewProjectHandler(projectService *service.ProjectService) *ProjectHandler {
	return &ProjectHandler{projectService: projectService}
}

func (h *ProjectHandler) Create(c *gin.Context) {
	var req service.CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, err.Error())
		return
	}
	userID := c.GetUint("user_id")
	project, err := h.projectService.Create(userID, req)
	if err != nil {
		InternalError(c, err.Error())
		return
	}
	Success(c, project)
}

func (h *ProjectHandler) List(c *gin.Context) {
	userID := c.GetUint("user_id")
	projects, err := h.projectService.ListByUser(userID)
	if err != nil {
		InternalError(c, err.Error())
		return
	}
	Success(c, projects)
}

func (h *ProjectHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		BadRequest(c, "invalid id")
		return
	}
	project, err := h.projectService.GetByID(uint(id))
	if err != nil {
		NotFound(c, "project not found")
		return
	}
	Success(c, project)
}

func (h *ProjectHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		BadRequest(c, "invalid id")
		return
	}
	var req service.UpdateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, err.Error())
		return
	}
	project, err := h.projectService.Update(uint(id), req)
	if err != nil {
		InternalError(c, err.Error())
		return
	}
	Success(c, project)
}

func (h *ProjectHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		BadRequest(c, "invalid id")
		return
	}
	if err := h.projectService.Delete(uint(id)); err != nil {
		InternalError(c, err.Error())
		return
	}
	Success(c, nil)
}
