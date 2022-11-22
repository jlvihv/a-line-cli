package service

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_GetTemplates(t *testing.T) {
	templateService := NewTemplateService()
	data := templateService.GetTemplates()
	t.Log(data)
	assert.NotEmpty(t, data)
}

func Test_GetTemplateDetail(t *testing.T) {
	templateService := NewTemplateService()
	detailData, _ := templateService.GetTemplateDetail(4)
	assert.NotNil(t, detailData)
}
