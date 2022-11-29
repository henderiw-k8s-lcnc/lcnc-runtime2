package lcncsyntax

// lcncCfgPreHookFn processes the for, own, watch generically
type lcncCfgPreHookFn func(lcncCfg *LcncConfig)
type lcncCfgPostHookFn func(lcncCfg *LcncConfig)

// lcncGvrObjectFn processes the for, own, watch per item
type lcncGvrObjectFn func(o Origin, idx int, n string, v LcncGvrObject)

// lcncBlockFn processes the block part of the Variables and functions
type lcncBlockFn func(o Origin, idx int, v LcncBlock)

type lcncVarsPreHookFn func(v []LcncVarBlock)
type lcncVarsPostHookFn func(v []LcncVarBlock)

// lcncVarFn processes the variable in the variables section
type lcncVarFn func(o Origin, block bool, idx int, vertexName string, v LcncVar)

type lcncFunctionsPreHookFn func(v []LcncFunctionsBlock)
type lcncFunctionsPostHookFn func(v []LcncFunctionsBlock)

// lcncFunctionFn processes the function in the functions section
type lcncFunctionFn func(o Origin, block bool, idx int, vertexName string, v LcncFunction)

type lcncServicesPreHookFn func(v []LcncFunctionsBlock)
type lcncServicesPostHookFn func(v []LcncFunctionsBlock)

// lcncServiceFn processes the service in the services section
type lcncServiceFn func(o Origin, block bool, idx int, vertexName string, v LcncFunction)

type WalkConfig struct {
	lcncCfgPreHookFn        lcncCfgPreHookFn
	lcncCfgPostHookFn       lcncCfgPostHookFn
	lcncGvrObjectFn         lcncGvrObjectFn
	lcncBlockFn             lcncBlockFn
	lcncVarsPreHookFn       lcncVarsPreHookFn
	lcncVarFn               lcncVarFn
	lcncVarsPostHookFn      lcncVarsPostHookFn
	lcncFunctionsPreHookFn  lcncFunctionsPreHookFn
	lcncFunctionFn          lcncFunctionFn
	lcncFunctionsPostHookFn lcncFunctionsPostHookFn
	lcncServicesPreHookFn   lcncServicesPreHookFn
	lcncServiceFn           lcncServiceFn
	lcncServicesPostHookFn  lcncServicesPreHookFn
}

func (r *lcncparser) walkLcncConfig(fnc WalkConfig) {
	// process config entry
	if fnc.lcncCfgPreHookFn != nil {
		fnc.lcncCfgPreHookFn(r.lcncCfg)
	}

	// process for, own, watch
	if fnc.lcncGvrObjectFn != nil {
		idx := 0
		for vertexName, v := range r.lcncCfg.For {
			fnc.lcncGvrObjectFn(OriginFor, idx, vertexName, v)
			idx++
		}
		idx = 0
		for vertexName, v := range r.lcncCfg.Own {
			fnc.lcncGvrObjectFn(OriginOwn, idx, vertexName, v)
			idx++
		}
		idx = 0
		for vertexName, v := range r.lcncCfg.Watch {
			fnc.lcncGvrObjectFn(OriginWatch, idx, vertexName, v)
		}
	}

	// process variables
	if fnc.lcncVarsPreHookFn != nil {
		fnc.lcncVarsPreHookFn(r.lcncCfg.Vars)
	}
	for idx, vars := range r.lcncCfg.Vars {
		// check if there is a block
		block := false
		if vars.LcncBlock.For != nil {
			block = true
			if fnc.lcncBlockFn != nil {
				fnc.lcncBlockFn(OriginVariable, idx, vars.LcncBlock)
			}
		}
		for vertexName, v := range vars.LcncVariables {
			if fnc.lcncVarFn != nil {
				fnc.lcncVarFn(OriginVariable, block, idx, vertexName, v)
			}
		}
	}
	if fnc.lcncVarsPostHookFn != nil {
		fnc.lcncVarsPostHookFn(r.lcncCfg.Vars)
	}

	// process functions
	if fnc.lcncFunctionsPreHookFn != nil {
		fnc.lcncFunctionsPreHookFn(r.lcncCfg.Functions)
	}
	for idx, functions := range r.lcncCfg.Functions {
		// check if there is a block
		block := false
		if functions.LcncBlock.For != nil {
			block = true
			if fnc.lcncBlockFn != nil {
				fnc.lcncBlockFn(OriginFunction, idx, functions.LcncBlock)
			}
		}
		for vertexName, v := range functions.LcncFunctions {
			if fnc.lcncFunctionFn != nil {
				fnc.lcncFunctionFn(OriginFunction, block, idx, vertexName, v)
			}

		}
	}
	if fnc.lcncFunctionsPostHookFn != nil {
		fnc.lcncFunctionsPostHookFn(r.lcncCfg.Functions)
	}

	// process services
	if fnc.lcncServicesPreHookFn != nil {
		fnc.lcncServicesPreHookFn(r.lcncCfg.Services)
	}
	for idx, services := range r.lcncCfg.Services {
		for vertexName, v := range services.LcncFunctions {
			if fnc.lcncServiceFn != nil {
				fnc.lcncServiceFn(OriginService, false, idx, vertexName, v)
			}
		}
	}
	if fnc.lcncServicesPostHookFn != nil {
		fnc.lcncServicesPostHookFn(r.lcncCfg.Services)
	}

	// process config end
	if fnc.lcncCfgPostHookFn != nil {
		fnc.lcncCfgPostHookFn(r.lcncCfg)
	}
}
