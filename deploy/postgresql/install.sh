#!/bin/bash

helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo update
helm install dogfood-backend bitnami/postgresql -n dogfood -f values.yaml
