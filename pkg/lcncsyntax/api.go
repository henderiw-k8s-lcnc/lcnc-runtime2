package lcncsysntax

type LcncConfig struct {
	For       LcncFor         `json:"for" yaml:"for"`
	Own       []LcncOwn       `json:"own,omitempty" yaml:"own,omitempty"`
	Watch     []LcncWatch     `json:"watch,omitempty" yaml:"watch,omitempty"`
	Vars      []*LcncVariable `json:"vars,omitempty" yaml:"vars,omitempty"`
	Resources []*LcncResource `json:"resources,omitempty" yaml:"resources,omitempty"`
	Services  []*LcncService  `json:"services,omitempty" yaml:"services,omitempty"`
}

type LcncImage struct {
	ImageName string `json:"image" yaml:"image"`
}

type LcncVariableName struct {
	VariableName string `json:"var" yaml:"var"`
}

type LcncGvk struct {
	Gvk string `json:"gvk" yaml:"gvk"`
}

type LcncFor struct {
	LcncVariableName `json:",inline" yaml:",inline"`
	LcncGvk          `json:",inline" yaml:",inline"`
}

type LcncOwn struct {
	LcncGvk `json:",inline" yaml:",inline"`
}

type LcncWatch struct {
	LcncGvk   `json:",inline" yaml:",inline"`
	LcncImage `json:",inline" yaml:",inline"`
}

type LcncVariable struct {
	LcncVariableName `json:",inline" yaml:",inline"`
	LcncQuery        `json:",inline" yaml:",inline"`
	For              LcncForLoop `json:"for,omitempty" yaml:"for,omitempty"`
	Slice            string      `json:"slice,omitempty" yaml:"string,omitempty"`
	Map              LcncMap     `json:"map,omitempty" yaml:"map,omitempty"`
}

type LcncResource struct {
	LcncImage `json:",inline" yaml:",inline"`
	For       LcncForLoop  `json:"for,omitempty" yaml:"for,omitempty"`
	Config    string       `json:"config,omitempty" yaml:"config,omitempty"`
	Input     []LcncInput  `json:"input,omitempty" yaml:"input,omitempty"`
	Output    []LcncOutput `json:"output,omitempty" yaml:"output,omitempty"`
}

type LcncService struct {
	LcncImage `json:",inline" yaml:",inline"`
}

type LcncInput struct {
	LcncVariableName `json:",inline" yaml:",inline"`
	LcncQuery        `json:",inline" yaml:",inline"`
}

type LcncOutput struct {
	LcncGvk `json:",inline" yaml:",inline"`
	Type    string
}

type LcncQuery struct {
	Query    string       `json:"query,omitempty" yaml:"query,omitempty"`
	Selector LcncSelector `json:"selector,omitempty" yaml:"selector,omitempty"`
}

type LcncSelector struct {
	Name        string            `json:"name,omitempty" yaml:"name,omitempty"`
	MatchLabels map[string]string `json:"matchLabels,omitempty" yaml:"matchLabels,omitempty"`
}

type LcncForLoop struct {
	Range     LcncForRange `json:"rane,omitempty" yaml:"range,omitempty"`
	Itarator1 string       `json:"iterator1,omitempty" yaml:"iterator1,omitempty"`
	Itarator2 string       `json:"iterator2,omitempty" yaml:"iterator2,omitempty"`
}

type LcncForRange struct {
	Gvk     string `json:"gvk,omitempty" yaml:"gvk,omitempty"`
	Variale string `json:"var,omitempty" yaml:"var,omitempty"`
}

type LcncMap struct {
	Key   string    `json:"key,omitempty" yaml:"key,omitempty"`
	Value LcncQuery `json:"value,omitempty" yaml:"value,omitempty"`
}
