package fake

import "k8s.io/apimachinery/pkg/runtime"

type Recorder struct {
}

func (fr *Recorder) Event(object runtime.Object, eventtype, reason, message string) {

}

// Eventf is just like Event, but with Sprintf for the message field.
func (fr *Recorder) Eventf(object runtime.Object, eventtype, reason, messageFmt string, args ...interface{}) {
}

// AnnotatedEventf is just like eventf, but with annotations attached
func (fr *Recorder) AnnotatedEventf(object runtime.Object, annotations map[string]string, eventtype, reason, messageFmt string, args ...interface{}) {
}
