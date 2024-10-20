package mysql

import (
	"blog/global"
	"blog/model"
)

// LabelIfExist 判断标签是否存在
func LabelIfExist(labelName string) (exits int64) {
	var label *model.Label
	exits = global.DB.Where("label_name=?", labelName).First(&label).RowsAffected
	return exits
}

func CreateLabel(label *model.Label) (result int64) {
	result = global.DB.Create(&label).RowsAffected
	return result
}

func QueryLabelId(labelName string) uint64 {
	var label *model.Label
	global.DB.Where("label_name=?", labelName).First(&label)
	return label.LabelID
}

func DeleteLabel(labelId uint64) {
	var count int64
	var label *model.Label
	global.DB.Table("article").Joins("JOIN label "+
		"ON article.label_id = label.label_id").Where("label.label_id = ?", labelId).Count(&count)
	if count == 0 {
		global.DB.Delete(&label, "label_id=?", labelId)
	}
}

func GetAllLabel() (labelArray []*model.Label) {
	global.DB.Find(&labelArray)
	return labelArray
}

func GetLabelNameById(id uint64) (name string) {
	var label *model.Label
	global.DB.Where("label_id=?", id).First(&label)
	name = label.LabelName
	return name
}
