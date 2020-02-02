# my-mutate-webhook
Learn to use mutating admission webhook
This webhook adds environment variables according to annotation,
e.g.
 com.xxx.add.env.HELLO: "WORLD"
This will generate an `HELLO=WORLD` env var to all containers

## Build webhook server binary
```
✗ make
```

## Generate ssl key
```
✗ cd ssl
✗ make cert
```

## Build docker image
```
✗ docker build --no-cache -t suexcxine/annotation-to-env:v1 .
✗ docker push suexcxine/annotation-to-env:v1
```

## Config caBundle
Replace caBundle section in your MutatingWebhookConfiguration config in webhook.yaml with output from the following command
```
✗ kubectl config view --raw --minify --flatten -o jsonpath='{.clusters[].cluster.certificate-authority-data}'
```

## Deploy
```
✗ kubectl label namespace default my-mutate=enabled
✗ kubectl get namespace -L my-mutate
✗ kubectl create -f deploy/webhook.yaml
✗ kubectl create -f deploy/sleep_for_test.yaml
```

## Verify
You can see environment variable HELLO with value WORLD has been added:
```
✗ kubectl get po
✗ kubectl get pod sleep-765c8fd8d8-fxkjp -o yaml
apiVersion: v1
kind: Pod
metadata:
  annotations:
    com.xxx.add.env.HELLO: WORLD
  creationTimestamp: "2020-02-02T04:52:19Z"
  generateName: sleep-765c8fd8d8-
  labels:
    app: sleep
    pod-template-hash: 765c8fd8d8
  name: sleep-765c8fd8d8-fxkjp
  namespace: default
  ownerReferences:
  - apiVersion: apps/v1
    blockOwnerDeletion: true
    controller: true
    kind: ReplicaSet
    name: sleep-765c8fd8d8
    uid: ca9c2ab7-4577-11ea-9375-025000000001
  resourceVersion: "1431305"
  selfLink: /api/v1/namespaces/default/pods/sleep-765c8fd8d8-fxkjp
  uid: caa7847c-4577-11ea-9375-025000000001
spec:
  containers:
  - command:
    - /bin/sleep
    - infinity
    env:
    - name: HELLO
      value: WORLD
    image: alpine
    imagePullPolicy: IfNotPresent
```

## Reference
https://medium.com/ovni/writing-a-very-basic-kubernetes-mutating-admission-webhook-398dbbcb63ec
https://github.com/morvencao/kube-mutating-webhook-tutorial

