package env_test

import (
	"testing"

	"github.com/patrickhuber/go-cross/env"
)

type conformance struct {
	env env.Environment
}

type Conformance interface {
	CanRoundtrip(t *testing.T)
	CanDelete(t *testing.T)
	CanLookup(t *testing.T)
}

func NewConformance(e env.Environment) Conformance {
	return &conformance{
		env: e,
	}
}

func (c *conformance) CanRoundtrip(t *testing.T) {
	key := "TEST"
	value := "VALUE"

	e := c.env

	err := e.Set(key, value)
	if err != nil {
		t.Fatal(err)
	}

	v := e.Get(key)
	if v != value {
		t.Fatalf("expected %s=%s but found %s=%s", key, value, key, v)
	}
}

func (c *conformance) CanDelete(t *testing.T) {
	key := "TEST"
	value := "VALUE"

	e := c.env

	err := e.Set(key, value)
	if err != nil {
		t.Fatal(err)
	}

	err = e.Delete(key)
	if err != nil {
		t.Fatal(err)
	}

	_, ok := e.Lookup(key)
	if ok {
		t.Fatalf("expected %s to be deleted but it was found", key)
	}
}

func (c *conformance) CanLookup(t *testing.T) {
	key := "TEST"
	value := "VALUE"

	e := c.env

	err := e.Set(key, value)
	if err != nil {
		t.Fatal(err)
	}
	_, ok := e.Lookup(key)
	if !ok {
		t.Fatalf("key %s was not found", key)
	}
}

func TestMemory(t *testing.T) {
	t.Run("can roundtrip", func(t *testing.T) {
		NewConformance(env.NewMemory()).
			CanRoundtrip(t)
	})
	t.Run("can delete", func(t *testing.T) {
		NewConformance(env.NewMemory()).
			CanDelete(t)
	})
	t.Run("can lookup", func(t *testing.T) {
		NewConformance(env.NewMemory()).
			CanLookup(t)
	})
}

func TestStdLib(t *testing.T) {
	t.Run("can roundtrip", func(t *testing.T) {
		NewConformance(env.New()).
			CanRoundtrip(t)
	})
	t.Run("can delete", func(t *testing.T) {
		NewConformance(env.NewMemory()).
			CanDelete(t)
	})
	t.Run("can lookup", func(t *testing.T) {
		NewConformance(env.NewMemory()).
			CanLookup(t)
	})
}
