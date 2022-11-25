package lcncsyntax

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"

	"k8s.io/apimachinery/pkg/runtime/schema"
)

type Externalresources struct {
	resources []schema.GroupVersionResource
}

func NewExternalResources() *Externalresources {
	return &Externalresources{
		resources: []schema.GroupVersionResource{},
	}
}

func (r *Externalresources) Add(er schema.GroupVersionResource) {
	found := false
	for _, resource := range r.resources {
		if resource.Group == er.Group &&
			resource.Version == er.Version &&
			resource.Resource == er.Resource {
			return
		}
	}
	if !found {
		r.resources = append(r.resources, er)
	}
}

func (r *Externalresources) Get() []schema.GroupVersionResource {
	return r.resources
}

func (r *LcncConfig) GetRootVertexName() (string, error) {
	if len(r.For) != 1 {
		return "", fmt.Errorf("lcnc config must have just 1 for statement, got: %v", r.For)
	}
	for vertexName := range r.For {
		return vertexName, nil
	}
	return "", fmt.Errorf("we should never come here, got: %v", r.For)
}

func (r *LcncConfig) GetExternalResources() ([]schema.GroupVersionResource, error) {
	er := NewExternalResources()
	// for must exist
	if len(r.For) != 1 {
		return nil, fmt.Errorf("for statement only accepts 1 value: %s", r.For)
	}
	for _, f := range r.For {
		gvr, err := GetGVR(f.Gvr)
		if err != nil {
			return nil, errors.Wrap(err, "for")
		}
		er.Add(*gvr)
	}
	for _, own := range r.Own {
		gvr, err := GetGVR(own.Gvr)
		if err != nil {
			return nil, errors.Wrap(err, "own")
		}
		er.Add(*gvr)
	}
	for _, watch := range r.Watch {
		gvr, err := GetGVR(watch.Gvr)
		if err != nil {
			return nil, errors.Wrap(err, "watch")
		}
		er.Add(*gvr)
	}
	for _, variable := range r.Vars {
		if variable.LcncQuery.Query != nil {
			value, err := GetValue(*variable.LcncQuery.Query)
			if err != nil {
				return nil, errors.Wrap(err, "variable")
			}
			if value.Kind == GVRKind {
				er.Add(*value.Gvr)
			}
		}
		if variable.For.Range != nil {
			value, err := GetValue(*variable.For.Range)
			if err != nil {
				return nil, err
			}
			if value.Kind == GVRKind {
				er.Add(*value.Gvr)
			}
		}
	}
	for _, resource := range r.Resources {
		for _, output := range resource.Function.Output {
			value, err := GetValue(output)
			if err != nil {
				return nil, err
			}

			if value.Kind == GVRKind {
				er.Add(*value.Gvr)
			}
		}
	}
	return er.Get(), nil
}

type ValueKind string

const (
	VariableKind               ValueKind = "variable"
	KeyVariableReferenceKind   ValueKind = "keyVariableReference" // used for $KEY, $VALUE
	RootVariableReferenceKind  ValueKind = "rootVariableReference"
	ChildVariableReferenceKind ValueKind = "childVariableRefrence"
	GVRKind                    ValueKind = "gvr"
)

type Value struct {
	Kind     ValueKind
	Variable []string
	Gvr      *schema.GroupVersionResource
}

func GetValue(s string) (*Value, error) {
	if len(s) <= 1 {
		return nil, fmt.Errorf("input value should have an input string with len > 1, got: %s", s)
	}
	// check if this is a variable or a gvr
	if string(s[0:1]) == "$" {
		varKind, varValue := GetVariable(s)
		return &Value{
			Kind:     varKind,
			Variable: varValue,
		}, nil
	}
	if len(strings.Split(s, "/")) == 1 {
		// this is a regular variable w/o a reference ($)
		return &Value{
			Kind:     VariableKind,
			Variable: []string{s},
		}, nil
	}
	// this is a gvr
	gvr, err := GetGVR(s)
	if err != nil {
		return nil, err
	}
	return &Value{
		Kind: GVRKind,
		Gvr:  gvr,
	}, nil

}

func GetVariable(s string) (ValueKind, []string) {
	// remove the first char from the string
	// split the string with the . delineator
	split := strings.Split(s[1:], ".")
	if split[0] == "VALUE" || split[0] == "KEY" {
		return KeyVariableReferenceKind, split
	}
	if len(split) > 1 {
		return ChildVariableReferenceKind, split
	}
	return RootVariableReferenceKind, split
}

func GetGVR(s string) (*schema.GroupVersionResource, error) {
	split := strings.Split(s, "/")
	if len(split) != 3 {
		return nil, fmt.Errorf("expecting a GVR in format <group>/<version>/<resource>, got: %s", s)
	}
	return &schema.GroupVersionResource{
		Group:    split[0],
		Version:  split[1],
		Resource: split[2],
	}, nil
}
