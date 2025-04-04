package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type UserService struct {
	// not need to implement
	NotEmptyStruct bool
}
type MessageService struct {
	// not need to implement
	NotEmptyStruct bool
}

type Container struct {
	constructors map[string]func() any
	singletons   map[string]func() any
	instances    map[string]any
}

func NewContainer() *Container {
	return &Container{
		constructors: make(map[string]func() any),
		singletons:   make(map[string]func() any),
		instances:    make(map[string]any),
	}
}

func (c *Container) RegisterType(name string, constructor any) {
	f, ok := constructor.(func() any)
	if !ok {
		panic(fmt.Sprintf("invalid constructor for %s: must be a function", name))
	}
	c.constructors[name] = f
}

func (c *Container) RegisterSingletonType(name string, constructor any) {
	f, ok := constructor.(func() any)
	if !ok {
		panic(fmt.Sprintf("invalid constructor for %s: must be a function", name))
	}
	c.singletons[name] = f
}

func (c *Container) Resolve(name string) (any, error) {
	if constructor, ok := c.singletons[name]; ok {
		if instance, ok := c.instances[name]; ok {
			return instance, nil
		}

		instance := constructor()
		c.instances[name] = instance
		return instance, nil
	}

	if constructor, ok := c.constructors[name]; ok {
		instance := constructor()
		return instance, nil
	}

	return nil, fmt.Errorf("type %s not registered", name)
}

func TestDIContainer(t *testing.T) {
	container := NewContainer()
	container.RegisterType("UserService", func() interface{} {
		return &UserService{}
	})
	container.RegisterType("MessageService", func() interface{} {
		return &MessageService{}
	})

	userService1, err := container.Resolve("UserService")
	assert.NoError(t, err)
	userService2, err := container.Resolve("UserService")
	assert.NoError(t, err)

	u1 := userService1.(*UserService)
	u2 := userService2.(*UserService)
	assert.False(t, u1 == u2)

	messageService, err := container.Resolve("MessageService")
	assert.NoError(t, err)
	assert.NotNil(t, messageService)

	paymentService, err := container.Resolve("PaymentService")
	assert.Error(t, err)
	assert.Nil(t, paymentService)
}
