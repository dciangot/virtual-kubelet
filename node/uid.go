package node

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
)

type podUIDKey struct {
	Namespace string
	Name      string
	UID       types.UID
}

func newPodUIDKey(pod *corev1.Pod) podUIDKey {
	return podUIDKey{
		Namespace: pod.Namespace,
		Name:      pod.Name,
		UID:       pod.UID,
	}
}

func objectName(namespace, name string) string {
	if len(namespace) > 0 {
		return namespace + "/" + name
	}
	return name
}

func (k *podUIDKey) ObjectName() string {
	return objectName(k.Namespace, k.Name)
}

func (k *podUIDKey) String() string {
	return k.ObjectName() + "/" + string(k.UID)
}

// uidProviderWrapper wraps a legacy PodLifecycleHandler to handle GetPodByUID/GetPodStatusByUID,
// by ignoring the pod UID and calling GetPod/GetPodStatus using only the pod namespace/name.
type uidProviderWrapper struct {
	PodLifecycleHandler
}

var _ PodUIDLifecycleHandler = (*uidProviderWrapper)(nil)

func (p *uidProviderWrapper) GetPodByUID(ctx context.Context, namespace, name string, uid types.UID) (*corev1.Pod, error) {
	return p.GetPod(ctx, namespace, name)
}

func (p *uidProviderWrapper) GetPodStatusByUID(ctx context.Context, namespace, name string, uid types.UID) (*corev1.PodStatus, error) {
	return p.GetPodStatus(ctx, namespace, name)
}

// NotifyPods forwards the notification callback to the underlying handler if it implements
// PodNotifier. This ensures that legacy providers which also implement PodNotifier are
// correctly recognised as an asyncProvider, so their notifier callback is wired up and
// CreatePod/UpdatePod/DeletePod can call it without panicking on a nil function.
func (p *uidProviderWrapper) NotifyPods(ctx context.Context, f func(*corev1.Pod)) {
	if notifier, ok := p.PodLifecycleHandler.(PodNotifier); ok {
		notifier.NotifyPods(ctx, f)
	}
}

// uidBasedProvider wraps a PodUIDLifecycleHandler to implement PodLifecycleHandler, by stubbing
// the legacy GetPod/GetPodStatus methods. These methods are never called if the provider implements
// PodUIDLifecycleHandler, but are required to be implemented as PodControllerConfig takes a
// PodLifecycleHandler. In the future, PodControllerConfig will take an asyncProvider, and this will
// no longer be necessary.
type uidBasedProvider struct {
	PodUIDLifecycleHandler
}

// UIDBasedProvider wraps a PodUIDLifecycleHandler into a PodLifecycleHandler that can be used in
// PodControllerConfig.
func UIDBasedProvider(p PodUIDLifecycleHandler) PodLifecycleHandler {
	return &uidBasedProvider{PodUIDLifecycleHandler: p}
}

func (p *uidBasedProvider) GetPod(ctx context.Context, namespace, name string) (*corev1.Pod, error) {
	panic("GetPod should never be called when GetPodByUID is implemented")
}

func (p *uidBasedProvider) GetPodStatus(ctx context.Context, namespace, name string) (*corev1.PodStatus, error) {
	panic("GetPodStatus should never be called when GetPodStatusByUID is implemented")
}
