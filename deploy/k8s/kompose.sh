#!/bin/bash

BASE_PATH=$(cd `dirname $0`; cd ..; pwd)
cd $BASE_PATH
cd depploy

# kompose使用: https://kubernetes.io/zh-cn/docs/tasks/configure-pod-container/translate-compose-kubernetes/

kompose convert

kubectl apply -f "$BASE_PATH/deploy/k8s"
