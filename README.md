# my-mutate-webhook
This webhook adds environment variables according to annotation

e.g.
 com.xxx.add.env.HELLO: "WORLD"
will inject an `HELLO=WORLD` env var to all containers

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
Warning: This will push ssl credentials to docker.io, you could use a private registry,
or place ssl credentials to a volume instead of docker image, involving some code change.
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
You can see environment variable HELLO with value WORLD has been injected:
```
✗ kubectl get po
NAME                         READY   STATUS    RESTARTS   AGE
my-mutate-5c99645667-95664   1/1     Running   0          4m6s
sleep-5678f9bb47-57psc       2/2     Running   0          2m59s
✗ kubectl get pod sleep-765c8fd8d8-fxkjp -o yaml
apiVersion: v1
kind: Pod
metadata:
  annotations:
    com.xxx.add.env.HELLO: WORLD
  creationTimestamp: "2020-02-02T07:39:19Z"
  generateName: sleep-5678f9bb47-
  labels:
    app: sleep
    pod-template-hash: 5678f9bb47
  name: sleep-5678f9bb47-57psc
  namespace: default
  ownerReferences:
  - apiVersion: apps/v1
    blockOwnerDeletion: true
    controller: true
    kind: ReplicaSet
    name: sleep-5678f9bb47
    uid: 1f6442bd-458f-11ea-9375-025000000001
  resourceVersion: "1443716"
  selfLink: /api/v1/namespaces/default/pods/sleep-5678f9bb47-57psc
  uid: 1f6b130d-458f-11ea-9375-025000000001
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
    name: sleep
    resources: {}
    terminationMessagePath: /dev/termination-log
    terminationMessagePolicy: File
    volumeMounts:
    - mountPath: /var/run/secrets/kubernetes.io/serviceaccount
      name: default-token-htfc2
      readOnly: true
  - command:
    - /bin/sleep
    - infinity
    env:
    - name: xxx
      value: yyy
    - name: HELLO
      value: WORLD
    image: alpine
    imagePullPolicy: IfNotPresent
    ...
```

## Reference
https://medium.com/ovni/writing-a-very-basic-kubernetes-mutating-admission-webhook-398dbbcb63ec
https://github.com/morvencao/kube-mutating-webhook-tutorial

