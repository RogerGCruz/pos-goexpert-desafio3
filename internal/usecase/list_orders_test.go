package usecase

import (
	"errors"
	"testing"

	"github.com/devfullcycle/20-CleanArch/internal/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type orderRepositoryStub struct {
	orders []entity.Order
	err    error
}

func (r *orderRepositoryStub) Save(order *entity.Order) error {
	_ = order
	return nil
}

func (r *orderRepositoryStub) GetAll() ([]entity.Order, error) {
	if r.err != nil {
		return nil, r.err
	}

	return r.orders, nil
}

func TestListOrderUseCaseExecute(t *testing.T) {
	repository := &orderRepositoryStub{
		orders: []entity.Order{
			{ID: "1", Price: 10, Tax: 2, FinalPrice: 12},
			{ID: "2", Price: 20, Tax: 4, FinalPrice: 24},
		},
	}

	useCase := NewListOrderUseCase(repository)
	output, err := useCase.Execute()

	require.NoError(t, err)
	require.Len(t, output, 2)
	assert.Equal(t, "1", output[0].ID)
	assert.Equal(t, 10.0, output[0].Price)
	assert.Equal(t, 2.0, output[0].Tax)
	assert.Equal(t, 12.0, output[0].FinalPrice)
	assert.Equal(t, "2", output[1].ID)
	assert.Equal(t, 20.0, output[1].Price)
	assert.Equal(t, 4.0, output[1].Tax)
	assert.Equal(t, 24.0, output[1].FinalPrice)
}

func TestListOrderUseCaseExecuteReturnsRepositoryError(t *testing.T) {
	repository := &orderRepositoryStub{err: errors.New("repository error")}

	useCase := NewListOrderUseCase(repository)
	output, err := useCase.Execute()

	assert.Nil(t, output)
	assert.EqualError(t, err, "repository error")
}
