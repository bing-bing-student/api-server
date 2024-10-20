package controller

import (
	"blog/service/mysql"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type GetAllLabelResponse struct {
	LabelID   string `json:"label_id"`
	LabelName string `json:"label_name"`
}

func GetAllLabel(c *gin.Context) {
	var labelArrayResponse []GetAllLabelResponse
	labelArray := mysql.GetAllLabel()
	for _, label := range labelArray {
		labelId := strconv.FormatUint(label.LabelID, 10)
		labelArrayResponse = append(labelArrayResponse, GetAllLabelResponse{
			LabelID:   labelId,
			LabelName: label.LabelName,
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"labelArray": labelArrayResponse,
	})
}
