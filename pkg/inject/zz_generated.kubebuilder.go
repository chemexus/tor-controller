package inject

import (
	torv1alpha1 "github.com/kragniz/tor-controller/pkg/apis/tor/v1alpha1"
	rscheme "github.com/kragniz/tor-controller/pkg/client/clientset/versioned/scheme"
	"github.com/kragniz/tor-controller/pkg/controller/onionservice"
	"github.com/kragniz/tor-controller/pkg/inject/args"
	"github.com/kubernetes-sigs/kubebuilder/pkg/inject/run"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes/scheme"
)

func init() {
	rscheme.AddToScheme(scheme.Scheme)

	// Inject Informers
	Inject = append(Inject, func(arguments args.InjectArgs) error {
		Injector.ControllerManager = arguments.ControllerManager

		if err := arguments.ControllerManager.AddInformerProvider(&torv1alpha1.OnionService{}, arguments.Informers.Tor().V1alpha1().OnionServices()); err != nil {
			return err
		}

		// Add Kubernetes informers

		if c, err := onionservice.ProvideController(arguments); err != nil {
			return err
		} else {
			arguments.ControllerManager.AddController(c)
		}
		return nil
	})

	// Inject CRDs
	Injector.CRDs = append(Injector.CRDs, &torv1alpha1.OnionServiceCRD)
	// Inject PolicyRules
	Injector.PolicyRules = append(Injector.PolicyRules, rbacv1.PolicyRule{
		APIGroups: []string{"tor.k8s.io"},
		Resources: []string{"*"},
		Verbs:     []string{"*"},
	})
	// Inject GroupVersions
	Injector.GroupVersions = append(Injector.GroupVersions, schema.GroupVersion{
		Group:   "tor.k8s.io",
		Version: "v1alpha1",
	})
	Injector.RunFns = append(Injector.RunFns, func(arguments run.RunArguments) error {
		Injector.ControllerManager.RunInformersAndControllers(arguments)
		return nil
	})
}
