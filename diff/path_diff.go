package diff

import (
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"
)

// PathDiff is a diff between path item objects: https://swagger.io/specification/#path-item-object
type PathDiff struct {
	SummaryDiff     *ValueDiff      `json:"summary,omitempty" yaml:"summary,omitempty"`
	DescriptionDiff *ValueDiff      `json:"description,omitempty" yaml:"description,omitempty"`
	OperationsDiff  *OperationsDiff `json:"operations,omitempty" yaml:"operations,omitempty"`
	ServersDiff     *ServersDiff    `json:"servers,omitempty" yaml:"servers,omitempty"`
	ParametersDiff  *ParametersDiff `json:"parameters,omitempty" yaml:"parameters,omitempty"`
}

func newPathDiff() *PathDiff {
	return &PathDiff{}
}

// Empty indicates whether a change was found in this element
func (pathDiff *PathDiff) Empty() bool {
	if pathDiff == nil {
		return true
	}

	return pathDiff == nil || *pathDiff == *newPathDiff()
}

func getPathDiff(config *Config, pathItem1, pathItem2 *openapi3.PathItem) (*PathDiff, error) {

	diff, err := getPathDiffInternal(config, pathItem1, pathItem2)
	if err != nil {
		return nil, err
	}
	if diff.Empty() {
		return nil, nil
	}
	return diff, nil
}

func getPathDiffInternal(config *Config, pathItem1, pathItem2 *openapi3.PathItem) (*PathDiff, error) {

	if pathItem1 == nil || pathItem2 == nil {
		return nil, fmt.Errorf("path item is nil")
	}

	result := newPathDiff()
	var err error

	result.SummaryDiff = getValueDiff(pathItem1.Summary, pathItem2.Summary)
	result.DescriptionDiff = getValueDiff(pathItem1.Description, pathItem2.Description)

	result.OperationsDiff, err = getOperationsDiff(config, pathItem1, pathItem2)
	if err != nil {
		return nil, err
	}

	result.ServersDiff = getServersDiff(config, &pathItem1.Servers, &pathItem2.Servers)
	result.ParametersDiff, err = getParametersDiff(config, pathItem1.Parameters, pathItem2.Parameters)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Apply applies the diff
func (pathDiff *PathDiff) Patch(pathItem *openapi3.PathItem) error {

	if pathDiff.Empty() {
		return nil
	}

	err := pathDiff.OperationsDiff.Patch(pathItem.Operations())
	if err != nil {
		return err
	}

	return err
}
