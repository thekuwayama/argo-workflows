package controller

import (
	"k8s.io/utils/env"

	wfv1 "github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	wfretry "github.com/argoproj/argo-workflows/v3/workflow/util/retry"
)

// RetryTweak is a 2nd order function interface for tweaking the retry
type RetryTweak = func(retryStrategy wfv1.RetryStrategy, nodes wfv1.Nodes, tmpl *wfv1.Template)

// RetryOnDifferentHost append affinity with fail host to template
func RetryOnDifferentHost(retryNodeName string) RetryTweak {
	return func(retryStrategy wfv1.RetryStrategy, nodes wfv1.Nodes, tmpl *wfv1.Template) {
		if retryStrategy.Affinity == nil {
			return
		}
		hostNames := wfretry.GetFailHosts(nodes, retryNodeName)
		hostLabel := env.GetString("RETRY_HOST_NAME_LABEL_KEY", "kubernetes.io/hostname")
		if hostLabel != "" && len(hostNames) > 0 {
			tmpl.Affinity = wfretry.AddHostnamesToAffinity(hostLabel, hostNames, tmpl.Affinity)
		}
	}
}
