package mysql

import (
	"blog/global"
	"blog/model"
)

func GetToolsById(id uint64) (tools *model.Tools, result int64) {
	result = global.DB.Where("tool_id=?", id).First(&tools).RowsAffected
	return tools, result
}

func QueryAllTools() (toolsArray []*model.Tools, result int64) {
	result = global.DB.Find(&toolsArray).RowsAffected
	return toolsArray, result
}

func CreateTools(tools *model.Tools) (result int64) {
	result = global.DB.Create(&tools).RowsAffected
	return result
}

func UpdateTools(toolsInfo *model.Tools) (result int64) {
	result = global.DB.Model(&toolsInfo).Where("tool_id=?", toolsInfo.ToolID).Updates(toolsInfo).RowsAffected
	return result
}

func DeleteTools(toolsId uint64) (result int64) {
	var tools *model.Tools
	result = global.DB.Delete(&tools, "tool_id=?", toolsId).RowsAffected
	return result
}
