/*
Copyright 2019 The Kruise Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package validating

import (
	"context"
	"net/http"

	appsv1alpha1 "github.com/openkruise/kruise/pkg/apis/apps/v1alpha1"
	admissionv1beta1 "k8s.io/api/admission/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/runtime/inject"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission/types"
)

func init() {
	webhookName := "validating-create-update-uniteddeployment"
	if HandlerMap[webhookName] == nil {
		HandlerMap[webhookName] = []admission.Handler{}
	}
	HandlerMap[webhookName] = append(HandlerMap[webhookName], &UnitedDeploymentCreateUpdateHandler{})
}

// UnitedDeploymentCreateUpdateHandler handles UnitedDeployment
type UnitedDeploymentCreateUpdateHandler struct {
	// To use the client, you need to do the following:
	// - uncomment it
	// - import sigs.k8s.io/controller-runtime/pkg/client
	// - uncomment the InjectClient method at the bottom of this file.
	// Client  client.Client

	// Decoder decodes objects
	Decoder types.Decoder
}

var _ admission.Handler = &UnitedDeploymentCreateUpdateHandler{}

// Handle handles admission requests.
func (h *UnitedDeploymentCreateUpdateHandler) Handle(ctx context.Context, req types.Request) types.Response {
	obj := &appsv1alpha1.UnitedDeployment{}

	err := h.Decoder.Decode(req, obj)
	if err != nil {
		return admission.ErrorResponse(http.StatusBadRequest, err)
	}

	switch req.AdmissionRequest.Operation {
	case admissionv1beta1.Create:
		if allErrs := validateUnitedDeployment(obj); len(allErrs) > 0 {
			return admission.ErrorResponse(http.StatusUnprocessableEntity, allErrs.ToAggregate())
		}
	case admissionv1beta1.Update:
		oldObj := &appsv1alpha1.UnitedDeployment{}
		if err := h.Decoder.Decode(types.Request{
			AdmissionRequest: &admissionv1beta1.AdmissionRequest{Object: req.AdmissionRequest.OldObject},
		}, oldObj); err != nil {
			return admission.ErrorResponse(http.StatusBadRequest, err)
		}

		validationErrorList := validateUnitedDeployment(obj)
		updateErrorList := ValidateUnitedDeploymentUpdate(obj, oldObj)
		if allErrs := append(validationErrorList, updateErrorList...); len(allErrs) > 0 {
			return admission.ErrorResponse(http.StatusUnprocessableEntity, allErrs.ToAggregate())
		}
	}

	return admission.ValidationResponse(true, "")
}

var _ inject.Decoder = &UnitedDeploymentCreateUpdateHandler{}

// InjectDecoder injects the decoder into the UnitedDeploymentCreateUpdateHandler
func (h *UnitedDeploymentCreateUpdateHandler) InjectDecoder(d types.Decoder) error {
	h.Decoder = d
	return nil
}
