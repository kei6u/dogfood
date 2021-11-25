#!/bin/bash

kubectl create namespace datadog

kubectl -n datadog create secret generic datadog-credentials \
--from-literal=api-key=${API_KEY} \
--from-literal=app-key=${APP_KEY}

helm repo update

helm repo add datadog https://helm.datadoghq.com
helm install datadog-operator datadog/datadog-operator -n datadog
kubectl apply -f datadogagent.yaml

helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm install kube-state-metrics prometheus-community/kube-state-metrics -n datadog
