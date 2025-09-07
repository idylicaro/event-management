package mappers

import (
	"reflect"
)

// Mapper interface for generic mapping operations
type Mapper[S any, T any] interface {
	Map(source S) T
}

// FieldMapper provides generic field mapping between structs
type FieldMapper[S any, T any] struct {
	customMappings map[string]func(S) any
}

// NewFieldMapper creates a new generic field mapper
func NewFieldMapper[S any, T any]() *FieldMapper[S, T] {
	return &FieldMapper[S, T]{
		customMappings: make(map[string]func(S) any),
	}
}

// WithCustomMapping adds a custom mapping for a specific field
func (fm *FieldMapper[S, T]) WithCustomMapping(fieldName string, mapper func(S) any) *FieldMapper[S, T] {
	fm.customMappings[fieldName] = mapper
	return fm
}

// Map performs the mapping from source to target using reflection
func (fm *FieldMapper[S, T]) Map(source S) T {
	var target T

	sourceVal := reflect.ValueOf(source)
	targetVal := reflect.ValueOf(&target).Elem()

	// Handle pointer types for source
	if sourceVal.Kind() == reflect.Ptr {
		if sourceVal.IsNil() {
			return target
		}
		sourceVal = sourceVal.Elem()
	}

	// Handle pointer types for target
	targetType := targetVal.Type()
	if targetType.Kind() == reflect.Ptr {
		// Create new instance of the struct that the pointer points to
		targetStructType := targetType.Elem()
		newTarget := reflect.New(targetStructType)
		targetVal.Set(newTarget)
		targetVal = newTarget.Elem()
		targetType = targetStructType
	}

	sourceType := sourceVal.Type()

	for i := 0; i < targetType.NumField(); i++ {
		targetField := targetType.Field(i)
		targetFieldVal := targetVal.Field(i)

		if !targetFieldVal.CanSet() {
			continue
		}

		// Check for custom mapping first
		if customMapper, exists := fm.customMappings[targetField.Name]; exists {
			customVal := customMapper(source)
			if customVal != nil {
				targetFieldVal.Set(reflect.ValueOf(customVal))
			}
			continue
		}

		// Find matching field in source
		if sourceField, found := sourceType.FieldByName(targetField.Name); found {
			sourceFieldVal := sourceVal.FieldByName(sourceField.Name)

			if sourceFieldVal.IsValid() && sourceFieldVal.Type().AssignableTo(targetField.Type) {
				targetFieldVal.Set(sourceFieldVal)
			}
		}
	}

	return target
}

// BuilderMapper provides a fluent interface for mapping
type BuilderMapper[S any, T any] struct {
	mapper func(S) T
}

// NewBuilderMapper creates a new builder mapper
func NewBuilderMapper[S any, T any](mapper func(S) T) *BuilderMapper[S, T] {
	return &BuilderMapper[S, T]{mapper: mapper}
}

// Map executes the mapping
func (bm *BuilderMapper[S, T]) Map(source S) T {
	return bm.mapper(source)
}
