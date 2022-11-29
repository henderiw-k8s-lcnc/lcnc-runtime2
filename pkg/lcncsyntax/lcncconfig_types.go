package lcncsyntax

type LcncConfig struct {
	// key represents the variable
	For map[string]LcncGvrObject `json:"for" yaml:"for"`
	// key represents the variable
	Own map[string]LcncGvrObject `json:"own,omitempty" yaml:"own,omitempty"`
	// key represents the variable
	Watch map[string]LcncGvrObject `json:"watch,omitempty" yaml:"watch,omitempty"`
	// key respresents the variable
	Vars      []LcncVarBlock       `json:"vars,omitempty" yaml:"vars,omitempty"`
	Functions []LcncFunctionsBlock `json:"fucntions,omitempty" yaml:"functions,omitempty"`
	Services  []LcncFunctionsBlock `json:"services,omitempty" yaml:"services,omitempty"`
	//Services map[string]LcncFunction `json:"services,omitempty" yaml:"services,omitempty"`
}

type LcncGvrObject struct {
	Gvr       string `json:"gvr" yaml:"gvr"`
	LcncImage `json:",inline" yaml:",inline"`
}

type LcncVarBlock struct {
	LcncBlock    `json:",inline" yaml:",inline"`
	LcncVariables map[string]LcncVar `json:",inline" yaml:",inline"`
}

type LcncBlock struct {
	For *LcncFor `json:"for,omitempty" yaml:"for,omitempty"`
	// TODO add IF statement block as standalone and within the if statement
}

type LcncFor struct {
	Range *string `json:"range,omitempty" yaml:"range,omitempty"`
}

type LcncVar struct {
	Slice *LcncSlice `json:"slice,omitempty" yaml:"slice,omitempty"`
	Map   *LcncMap   `json:"map,omitempty" yaml:"map,omitempty"`
}

type LcncSlice struct {
	LcncValue `json:"value,omitempty" yaml:"value,omitempty"`
}

type LcncMap struct {
	Key       *string `json:"key,omitempty" yaml:"key,omitempty"`
	LcncValue `json:"value,omitempty" yaml:"value,omitempty"`
}

type LcncValue struct {
	LcncQuery `json:",inline" yaml:",inline"`
	String    *string `json:"string,omitempty" yaml:"string,omitempty"`
}

type LcncFunctionsBlock struct {
	LcncBlock    `json:",inline" yaml:",inline"`
	LcncFunctions map[string]LcncFunction `json:",inline" yaml:",inline"`
}

type LcncFunction struct {
	LcncImage `json:",inline" yaml:",inline"`
	//Vars      []LcncVarBlock    `json:"vars,omitempty" yaml:"vars,omitempty"`
	Vars   map[string]LcncVar `json:"vars,omitempty" yaml:"vars,omitempty"`
	Config string             `json:"config,omitempty" yaml:"config,omitempty"`
	// input is always a GVK of some sort
	Input map[string]string `json:"input,omitempty" yaml:"input,omitempty"`
	// key = variableName, value is gvr format or not -> gvr format is needed for external resources
	Output map[string]string `json:"output,omitempty" yaml:"output,omitempty"`
}

type LcncImage struct {
	ImageName *string `json:"image" yaml:"image"`
}

type LcncQuery struct {
	Query    *string       `json:"query,omitempty" yaml:"query,omitempty"`
	Selector *LcncSelector `json:"selector,omitempty" yaml:"selector,omitempty"`
}

type LcncSelector struct {
	Name        *string           `json:"name,omitempty" yaml:"name,omitempty"`
	MatchLabels map[string]string `json:"matchLabels,omitempty" yaml:"matchLabels,omitempty"`
}

/*
type LcncForLoop struct {
	Range *string  `json:"range,omitempty" yaml:"range,omitempty"`
	Slice *string  `json:"slice,omitempty" yaml:"string,omitempty"`
	Map   *LcncMap `json:"map,omitempty" yaml:"map,omitempty"`
}
*/
