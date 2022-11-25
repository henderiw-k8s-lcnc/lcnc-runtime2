package main

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/henderiw-k8s-lcnc/lcnc-runtime2/pkg/lcnc"
	"github.com/henderiw-k8s-lcnc/lcnc-runtime2/pkg/lcncsyntax"
	"github.com/yndd/ndd-runtime/pkg/logging"
	"gopkg.in/yaml.v2"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

const dir = "./examples"

func main() {
	debug := true
	zlog := zap.New(zap.UseDevMode(debug), zap.JSONEncoder())
	ctrl.SetLogger(zlog)
	logger := logging.NewLogrLogger(zlog.WithName("lcnc runtime"))

	files, err := os.ReadDir(dir)
	if err != nil {
		logger.Debug("cannot read directory", "error", err)
		os.Exit(1)
	}
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".yaml") {

			logger.Debug("file", "filename", f.Name())
			b, err := os.ReadFile(filepath.Join(dir, f.Name()))
			if err != nil {
				logger.Debug("cannot read file", "error", err)
				os.Exit(1)
			}

			lcncCfg := &lcncsyntax.LcncConfig{}
			if err := yaml.Unmarshal(b, lcncCfg); err != nil {
				logger.Debug("cannot unmarshal", "error", err)
				os.Exit(1)
			}

			l, err := lcnc.New(lcncCfg)
			if err != nil {
				logger.Debug("cannot create lcnc", "error", err)
				os.Exit(1)
			}

			extRes, err := l.GetExternalResources()
			if err != nil {
				logger.Debug("cannot get external resources", "error", err)
				os.Exit(1)
			}
			logger.Debug("external resources", "for", f.Name(), "external resoruces", extRes)

			if err := l.Transform(); err != nil {
				logger.Debug("cannot transformm resources", "error", err)
				os.Exit(1)
			}	
		}
	}

	// Parse config map
	/*
		mgr, err := manager.New(ctrl.GetConfigOrDie(), manager.Options{
			Namespace: os.Getenv("POD_NAMESPACE"),
		})
		if err != nil {
			logger.Debug("unable to create manager", "error", err)
			os.Exit(1)
		}

		var1 := lcncsyntax.LcncVariable{
			LcncVariableName: lcncsyntax.LcncVariableName{VariableName: "x"},
			LcncQuery: lcncsyntax.LcncQuery{
				Query: "topo.yndd.io/v1alpha1/templates",
				Selector: lcncsyntax.LcncSelector{
					MatchLabels: map[string]string{
						"yndd.io/topology": "$infra.spec.topology",
						"yndd.io/linktype": "$infra.spec.topology",
					},
				},
			},
		}

		value, err := lcncsyntax.GetValue(var1.LcncQuery.Query)
		if err != nil {
			logger.Debug("Cannot get value", "error", err)
			os.Exit(1)
		}

		if value.Kind == lcncsyntax.GVRKind {
			gvk, err := mgr.GetRESTMapper().KindFor(*value.Gvr)
			if err != nil {
				logger.Debug("Cannot get gvk", "error", err)
				os.Exit(1)
			}

			//if len(v.LcncQuery.Selector.MatchLabels) != 0 {
			opts := []client.ListOption{
				client.MatchingLabels(var1.LcncQuery.Selector.MatchLabels),
			}
			l := getUnstructuredObj(gvk)
			if err := mgr.GetClient().List(context.TODO(), l, opts...); err != nil {
				logger.Debug("Cannot get value", "error", err)
				os.Exit(1)
			}
		}

		var2 := lcncsyntax.LcncVariable{
			LcncVariableName: lcncsyntax.LcncVariableName{VariableName: "x"},
			For: lcncsyntax.LcncForLoop{
				Range:     "topo.yndd.io/v1alpha1/templates",
				Iterator2: "parentTemplateName",
			},
		}
	*/

}

func getUnstructuredObj(gvk schema.GroupVersionKind) *unstructured.UnstructuredList {
	var u unstructured.UnstructuredList
	u.SetAPIVersion(gvk.GroupVersion().String())
	u.SetKind(gvk.Kind)
	uCopy := u.DeepCopy()
	return uCopy
}
