package controller

import (
	"net/http"
	"users/model"
	"users/service"

	"github.com/gin-gonic/gin"
)

func Approval(c *gin.Context) {
	var input struct {
		UserID       int                  `json:"user_id"`
		ApprovalType service.ApprovalType `json:"approval_type"`
	}

	if err := c.ShouldBind(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, &model.GlobalResponse{
			Success: false,
			Message: err.Error(),
		})
	}

	s := service.GetService()
	defer func() {
		if r := recover(); r != nil {
			err := s.ErrorCheck()
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, &model.GlobalResponse{
					Success: false,
					Message: err.Error(),
				})
				return
			}
		}
	}()

	err := s.SellerUpdateApproval(c.Request.Context(), input.UserID, input.ApprovalType)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, &model.GlobalResponse{
			Success: false,
			Message: err.Error(),
		})
	}

	c.JSON(http.StatusOK, &model.GlobalResponse{
		Success: true,
		Message: "Approval status updated successfully",
	})
}
