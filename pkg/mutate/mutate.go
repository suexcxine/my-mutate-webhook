package mutate

import (
	"encoding/json"
	"fmt"
	"log"

	v1beta1 "k8s.io/api/admission/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type patchOperation struct {
	Op    string      `json:"op"`
	Path  string      `json:"path"`
	Value interface{} `json:"value,omitempty"`
}

func patchEnvByAnnotation(pod *corev1.Pod) patchOperation {
	// TODO 考虑有多个 container 的情况
	op := "add"
	evs := make([]corev1.EnvVar, 10)
	containers := pod.Spec.Containers
	if containers != nil && len(containers) > 0 && containers[0].Env != nil {
		op = "update"
		copy(evs, containers[0].Env)
	}
	meta := pod.ObjectMeta
	annotations := meta.GetAnnotations()
	for k, v := range annotations {
		ev := corev1.EnvVar{Name: k, Value: v}
		evs = append(evs, ev)
	}
	return patchOperation{
		Op:    op,
		Path:  "/spec/containers/0/env",
		Value: evs,
	}
}

// Mutate 设置AdmissionResponse
func Mutate(body []byte) ([]byte, error) {
	log.Printf("recv: %s\n", string(body))

	admReview := v1beta1.AdmissionReview{}
	if err := json.Unmarshal(body, &admReview); err != nil {
		return nil, fmt.Errorf("unmarshal request failed with %s", err)
	}

	var err error
	var pod *corev1.Pod

	responseBody := []byte{}
	ar := admReview.Request
	resp := v1beta1.AdmissionResponse{}

	if ar != nil {

		if err := json.Unmarshal(ar.Object.Raw, &pod); err != nil {
			return nil, fmt.Errorf("unable to unmarshal pod json object %v", err)
		}

		resp.Allowed = true
		resp.UID = ar.UID
		pT := v1beta1.PatchTypeJSONPatch
		resp.PatchType = &pT

		p := patchEnvByAnnotation(pod)
		resp.Patch, err = json.Marshal(p)
		if err != nil {
			return nil, err
		}

		resp.Result = &metav1.Status{
			Status: "Success",
		}

		admReview.Response = &resp
		responseBody, err = json.Marshal(admReview)
		if err != nil {
			return nil, err
		}
	}

	log.Printf("resp: %s\n", string(responseBody))
	return responseBody, nil
}
