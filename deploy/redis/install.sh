#!/bin/bash

helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo update
helm install dogfood-gateway bitnami/redis -n dogfood -f values.yaml
