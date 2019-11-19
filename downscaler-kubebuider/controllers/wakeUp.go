package controllers

import (
	"context"
	"encoding/base64"

	appsv1 "k8s.io/api/apps/v1"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	client "sigs.k8s.io/controller-runtime/pkg/client"
)

func wakeUp(ingressName string, ingressNamespace string, r *ScalingBackInfoReconciler) {
	log := r.Log
	log.Info("wakeUp", "ingressName", ingressName, "ingressNamespace", ingressNamespace)

	ctx := context.Background()

	// get ingress
	// TODO check that iongress is updated not less than a minute ago

	namespacedIngressName := client.ObjectKey{
		Namespace: ingressNamespace,
		Name:      ingressName,
	}
	ingress := &extensionsv1beta1.Ingress{}

	if err := r.Get(ctx, namespacedIngressName, ingress); err != nil {
		log.Error(err, "unable to get Ingress in wakeUp")
		return
	}

	// TODO restore ingress
	specBackup, err := base64.StdEncoding.DecodeString(ingress.ObjectMeta.Annotations["zero-scaling/backup"])
	if err != nil {
		log.Error(err, "unable to decode backup in wakeUp")
		return
	}

	ingress.Spec.Rules = []extensionsv1beta1.IngressRule{}
	err = ingress.Spec.Unmarshal(specBackup)
	log.Info("Restored rules", "rules", ingress.Spec.Rules, "specBackup", specBackup)

	if err != nil {
		log.Error(err, "unable to Unmarshal backup in wakeUp")
		return
	}

	delete(ingress.ObjectMeta.Annotations, "zero-scaling/backup")

	err = r.Update(ctx, ingress)

	if err != nil {
		log.Error(err, "unable to update ingress in wakeUp")
		return
	}

	log.Info("wakeUp complete", "ingressName", ingressName, "ingressNamespace", ingressNamespace)

	//  scale deployment back to 1

	namespacedDeploymentName := client.ObjectKey{
		Namespace: ingressNamespace,
		Name:      ingress.ObjectMeta.Annotations["zero-scaling/deploymentName"],
	}
	deployment := &appsv1.Deployment{}

	if err := r.Get(ctx, namespacedDeploymentName, deployment); err != nil {
		log.Error(err, "unable to get Deployment "+namespacedDeploymentName.String()+" in putToSleep")
		return
	}

	one := int32(1)
	deployment.Spec.Replicas = &one

	err = r.Update(ctx, deployment)

	if err != nil {
		log.Error(err, "unable to update deployment in putToSleep")
		return
	}
}
