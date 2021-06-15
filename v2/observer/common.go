package main

import (
	"fmt"
	"time"

	myk8s "github.com/brigadecore/brigade/v2/internal/kubernetes"
	corev1 "k8s.io/api/core/v1"
)

const apiRequestTimeout = 30 * time.Second

// syncDeletedPod only fires when a pod deletion is COMPLETE. i.e. The pod is
// completely gone.
func (o *observer) syncDeletedPod(obj interface{}) {
	o.syncMu.Lock()
	defer o.syncMu.Unlock()
	pod := obj.(*corev1.Pod)
	// Remove this pod from the set of pods we were tracking for deletion.
	// Managing this set is essential to not leaking memory.
	delete(o.deletingPodsSet, namespacedPodName(pod.Namespace, pod.Name))
}

func (o *observer) getPodTimeoutDuration(
	pod *corev1.Pod,
	max time.Duration,
) time.Duration {
	rawDuration := pod.Annotations[myk8s.AnnotationTimeoutDuration]
	if rawDuration == "" {
		return max
	}

	// Attempt to set the timeout per the annotation on the pod itself
	timeout, err := time.ParseDuration(rawDuration)
	// Fallback to the max if we are unable to parse timeout value
	if err != nil {
		o.errFn(
			fmt.Errorf(
				"unable to parse timeout duration %q for pod %q; "+
					"using configured maximum of %q",
				rawDuration,
				pod.Name,
				max,
			),
		)
		return max
	}
	// ... or if the parsed duration exceeds the max
	if timeout > max {
		o.errFn(
			fmt.Errorf(
				"timeout duration %q for pod %q exceeds the configured maximum; "+
					"using configured maximum of %q",
				timeout,
				pod.Name,
				max,
			),
		)
		return max
	}
	return timeout
}
