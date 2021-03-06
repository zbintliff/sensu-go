package store

import (
	"context"
	"path"

	"github.com/sensu/sensu-go/types"
)

const keySeparator = "/"

// KeyBuilder builds multi-tenant resource keys.
type KeyBuilder struct {
	resourceName string
	namespace    Namespace
}

// NewKeyBuilder creates a new KeyBuilder.
func NewKeyBuilder(resourceName string) KeyBuilder {
	builder := KeyBuilder{resourceName: resourceName}
	return builder
}

// WithOrg adds an Organization to a key.
func (b KeyBuilder) WithOrg(org string) KeyBuilder {
	ns := Namespace{Org: org}
	return b.WithNamespace(ns)
}

// WithResource adds a resource to a key.
func (b KeyBuilder) WithResource(r types.MultitenantResource) KeyBuilder {
	b.namespace = NewNamespaceFromResource(r)
	return b
}

// WithContext adds a namespace from a context.
func (b KeyBuilder) WithContext(ctx context.Context) KeyBuilder {
	ns := NewNamespaceFromContext(ctx)
	return b.WithNamespace(ns)
}

// WithNamespace adds a Namespace to a key.
func (b KeyBuilder) WithNamespace(ns Namespace) KeyBuilder {
	b.namespace = ns
	return b
}

// Build builds a key from the components it is given.
func (b KeyBuilder) Build(keys ...string) string {
	items := append(
		[]string{
			Root,
			b.resourceName,
			b.namespace.Org,
			b.namespace.Env,
		},
		keys...,
	)
	return path.Join(items...)
}

// BuildPrefix can be used when building a key that will be used to retrieve a collection of
// records. Unlike standard build method it stops building the key when it first
// encounters a wildcard value.
func (b KeyBuilder) BuildPrefix(keys ...string) string {
	out := Root + keySeparator + b.resourceName

	keys = append([]string{b.namespace.Org, b.namespace.Env}, keys...)
	for _, key := range keys {
		// If we encounter a wildcard stop and return key
		if key == WildcardValue {
			return out
		}
		// If the key is nil we ignore and do not add
		if key == "" {
			continue
		}
		out = out + keySeparator + key
	}

	return out
}
