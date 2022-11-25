package lcncsyntax

type LcncConfig struct {
	// key represents the variable
	For map[string]LcncGvrObject `json:"for" yaml:"for"`
	// key represents the variable
	Own map[string]LcncGvrObject `json:"own,omitempty" yaml:"own,omitempty"`
	// key represents the variable
	Watch map[string]LcncGvrObject `json:"watch,omitempty" yaml:"watch,omitempty"`
	// key respresents the variable
	Vars      map[string]LcncVar      `json:"vars,omitempty" yaml:"vars,omitempty"`
	Resources map[string]LcncResource `json:"resources,omitempty" yaml:"resources,omitempty"`
	Services  map[string]LcncFunction `json:"services,omitempty" yaml:"services,omitempty"`
}

type LcncGvrObject struct {
	Gvr       string `json:"gvr" yaml:"gvr"`
	LcncImage `json:",inline" yaml:",inline"`
}

type LcncVar struct {
	LcncQuery `json:",inline" yaml:",inline"`
	For       LcncForLoop `json:"for,omitempty" yaml:"for,omitempty"`
}

type LcncResource struct {
	For      LcncForLoop  `json:"for,omitempty" yaml:"for,omitempty"`
	Function LcncFunction `json:"function" yaml:"function"`
}

/*
type LcncService struct {
	LcncImage `json:",inline" yaml:",inline"`
}
*/

type LcncFunction struct {
	LcncImage `json:",inline" yaml:",inline"`
	Vars      map[string]LcncVar `json:"vars,omitempty" yaml:"vars,omitempty"`
	Config    string             `json:"config,omitempty" yaml:"config,omitempty"`
	Input     map[string]string  `json:"input,omitempty" yaml:"input,omitempty"`
	// key = variableName, value is gvr format or not -> gvr format is needed for external resources
	Output map[string]string `json:"output,omitempty" yaml:"output,omitempty"`
}

type LcncImage struct {
	ImageName string `json:"image" yaml:"image"`
}

type LcncQuery struct {
	Query    *string       `json:"query,omitempty" yaml:"query,omitempty"`
	Selector *LcncSelector `json:"selector,omitempty" yaml:"selector,omitempty"`
}

type LcncSelector struct {
	Name        *string            `json:"name,omitempty" yaml:"name,omitempty"`
	MatchLabels map[string]string `json:"matchLabels,omitempty" yaml:"matchLabels,omitempty"`
}

type LcncForLoop struct {
	Range *string  `json:"range,omitempty" yaml:"range,omitempty"`
	Slice *string  `json:"slice,omitempty" yaml:"string,omitempty"`
	Map   *LcncMap `json:"map,omitempty" yaml:"map,omitempty"`
}

type LcncMap struct {
	Key   string    `json:"key,omitempty" yaml:"key,omitempty"`
	Value LcncQuery `json:"value,omitempty" yaml:"value,omitempty"`
}

/*
type LcncInput struct {
	LcncVariableName `json:",inline" yaml:",inline"`
	LcncQuery        `json:",inline" yaml:",inline"`
}
*/

/*
type LcncOutput struct {
	LcncGvr `json:",inline" yaml:",inline"`
	Type    string
}
*/

/*
type LcncForRange struct {
	Gvk     string `json:"gvk,omitempty" yaml:"gvk,omitempty"`
	Variale string `json:"var,omitempty" yaml:"var,omitempty"`
}
*/

/*
	type LcncGvr struct {
		Gvr string `json:"gvr" yaml:"gvr"`
	}

	type LcncFor struct {
		LcncVariableName `json:",inline" yaml:",inline"`
		LcncGvr          `json:",inline" yaml:",inline"`
	}

	type LcncOwn struct {
		LcncGvr `json:",inline" yaml:",inline"`
	}

	type LcncWatch struct {
		LcncGvr   `json:",inline" yaml:",inline"`
		LcncImage `json:",inline" yaml:",inline"`
	}
*/
