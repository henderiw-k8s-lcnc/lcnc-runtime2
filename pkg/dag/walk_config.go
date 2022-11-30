package dag

type WalkInitFn func()

type WalkEntryFn func(from string, depth int)

type WalkConfig struct {
	Dep         bool
	WalkInitFn  WalkInitFn
	WalkEntryFn WalkEntryFn
}
