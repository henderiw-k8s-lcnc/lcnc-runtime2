package lcncsyntax

import (
	"strconv"
	"strings"
)

// this function is assumed to be executed after validation
// validate check if the for is present
func (r *LcncConfig) GetRootVertexName() string {
	for vertexName := range r.For {
		return vertexName
	}
	return ""
}

// CopyVariables copies the variable block and splits multiple entries
// in a slice to a single entry. This allows to build a generic
// resolution processor
func CopyVariables(vars []LcncVarBlock) map[string]LcncVarBlock {
	newvars := map[string]LcncVarBlock{}
	for idx, varBlock := range vars {
		for k, v := range varBlock.LcncVariables {
			newvars[strings.Join([]string{k, strconv.Itoa(idx)}, "/")] = LcncVarBlock{
				LcncBlock: varBlock.LcncBlock,
				LcncVariables: map[string]LcncVar{
					k: v,
				},
			}
		}
	}
	return newvars
}

// CopyFunctions copies the variable block and splits multiple entries
// in a slice to a single entry. This allows to build a generic
// resolution processor
func CopyFunctions(fns []LcncFunctionsBlock) map[string]LcncFunctionsBlock {
	newfns := map[string]LcncFunctionsBlock{}
	for idx, fnBlock := range fns {
		for k, v := range fnBlock.LcncFunctions {
			newfns[strings.Join([]string{k, strconv.Itoa(idx)}, "/")] = LcncFunctionsBlock{
				LcncBlock: fnBlock.LcncBlock,
				LcncFunctions: map[string]LcncFunction{
					k: v,
				},
			}
		}
	}
	return newfns
}

func GetIdxName(idxName string) (string, int) {
	split := strings.Split(idxName, "/")
	idx, _ := strconv.Atoi(split[1])
	return split[0], idx
}
