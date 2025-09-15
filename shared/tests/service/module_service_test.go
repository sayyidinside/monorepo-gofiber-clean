package service_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/sayyidinside/monorepo-gofiber-clean/shared/domain/entity"
	"github.com/sayyidinside/monorepo-gofiber-clean/shared/domain/service"
	mocks "github.com/sayyidinside/monorepo-gofiber-clean/shared/mocks/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_ModuleService_GetByID_Success(t *testing.T) {
	moduleRepository := new(mocks.ModuleRepository)

	service := service.NewModuleService(moduleRepository)

	fakeModule := &entity.Module{ID: 1, Name: "test", UUID: uuid.New()}
	moduleRepository.On("FindByID", mock.Anything, uint(1)).Return(fakeModule, nil)

	response := service.GetByID(context.Background(), 1)

	assert.Equal(t, 200, response.Status)
	assert.True(t, response.Success)
	assert.NotNil(t, response.Data)
	assert.Nil(t, response.Errors)
}

func Test_ModuleService_GetByID_NotFound(t *testing.T) {
	moduleRepository := new(mocks.ModuleRepository)

	service := service.NewModuleService(moduleRepository)

	moduleRepository.On("FindByID", mock.Anything, uint(1)).Return(nil, nil)

	response := service.GetByID(context.Background(), 1)

	assert.Equal(t, 404, response.Status)
	assert.False(t, response.Success)
	assert.Nil(t, response.Data)
	assert.Equal(t, "Module not found", response.Message)
}
