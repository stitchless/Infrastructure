## Table of Contents
- [**Initial Setup**](#initial-setup)
    - [Installing MicroK8s](#installing-microk8s)
    - [MicroK8s Addons](#microk8s-addons)
    - [Kubernetes Dashboard](#kubernetes-dashboard)
- [**Automation Setup**](#automation-setup)
    - [Cert-Manager](#cert-manager)
    - [Cloudflare Origin CA Issuer](#cloudflare-origin-ca-issuer)
    - [Enteral-DNS](#enteral-dns)
    - [*Kubed*](#kubed)
- [**Deploy your Infrastructure**](#deploy-your-infrastructure)
    - [Deploy your first app](#deploy-your-first-app)
        - [All unique namespaces require the following](#all-unique-namespaces-require-the-following)
        - [Required for each unique domain name:](#required-for-each-unique-domain-name)


<br>
<br>

## **Requirements**
Throughout this project I am using the latest stable tools I could find

1. Microk8s (Kubernetes for single node deployment)
2. RBAC enabled on your cluster (explained later)
3. SSH access to your server (running Linux, assumed Ubuntu / Debian)
4. Cloudflare as your CA and DNS provider
5. Helm installed on your local computer
6. Access to cloudflare API and origin CA API
7. Assume you are splitting up your deployments across multiple namespaces
8. A place to manage secrets (This example uses [secrethub](https://secrethub.io/)) (WIP)
    * This is used because it's a simple way deploy using secrets either from your IDE, or Github actions deployments
<!--
#####################
##### MICROK8S ######
#####################
-->
<br>
<br>
<br>

---
# **Initial Setup**
<br>

## [Installing MicroK8s](https://microk8s.io/docs)
[Microk8s](https://microk8s.io/) is a slim version of kubernetes.  By nature MicroK8s is designed to run and operrate a single server, not a cluster.\
Snap will be used for all bare metal installs this will help keep things clean on your bare metal.

```bash
# Install MicroK8s
sudo snap install microk8s --classic

# Set alias to kubectl < microk8s.kubectl
# From this point on microk8s.kubectl will be referred to as kubectl
sudo snap alias microk8s.kubectl kubectl

# Join the group
sudo usermod -a -G microk8s $USER
sudo chown -f -R $USER ~/.kube

# Reload your active user to apply changes
su - $USER

# Confirm a successful installation
microk8s status --wait-ready

# jq (Optional) used to confirm cloudflare origin issuer
# `kubectl describe` can be used later on to also verify as well
sudo snap install jq

# If rbac.authorization.k8s.io/v1, or  rbac.authorization.k8s.io/v1beta1 is shown you have RBAC enabled
kubectl api-versions
```
<br>

## MicroK8s Addons
Most of these add-ons are general use to make life easier
```bash
microk8s enable \
  rbac \
  dns \
  ingress \ # or your choice of ingress controller
  metallb \ # optional
  prometheus \ # optional
  storage # WIP
```

<br>

## [Kubernetes Dashboard](https://github.com/kubernetes/dashboard)
```bash
# It's recommended to use the install script from the repo in case of any updates
# Click the title link to go to the repo
kubectl apply -f https://raw.githubusercontent.com/kubernetes/dashboard/v2.1.0/aio/deploy/recommended.yaml
```
```yaml
# Run The following ...
# kubectl config view --raw
apiVersion: v1
clusters:
- cluster:
    certificate-authority-data: <SuperSecretCert>
    server: https://127.0.0.1:16443 # Replace IP with your server IP:16443
  name: microk8s-cluster
contexts:
- context:
    cluster: microk8s-cluster
    user: admin
  name: microk8s
current-context: microk8s
kind: Config
preferences: {}
users:
- name: admin
  user:
    token: <SuperSecretToken>

# Copy the output to ~/.kube/config or your OS equivelant 
```
> In order to manage our cluster remotely we will need to get the config and install kubectl on our local computer.
>
> Follow the instructions found [here](https://kubernetes.io/docs/tasks/tools/install-kubectl/) to install kubectl on your local computer.

```bash
# From here you can access your cluser by passing the following command
kubectlk proxy

# Then navigate to the following website
http://localhost:8001/api/v1/namespaces/kubernetes-dashboard/services/https:kubernetes-dashboard:/proxy/
```
From here you have a working kubernetes setup that can be managed from a web UI.  You can stop here or continue to enable website automations

<br>
<br>

---

# **Automation Setup**
<br>

## [Cert-Manager](https://cert-manager.io/docs/installation/kubernetes/)
```bash
# Create your cert-manager namespace
kubectl create namespace cert-manager

# Install the CustomResourceDefinition resources using kubectl
kubectl apply -f https://github.com/jetstack/cert-manager/releases/download/v1.2.0/cert-manager.crds.yaml

# Install the Cert Manager
kubectl apply -f https://github.com/jetstack/cert-manager/releases/download/v1.2.0/cert-manager.yaml

# Verify Install
kubectl get pods --namespace cert-manager

# NAME                                       READY   STATUS    RESTARTS   AGE
# cert-manager-5c6866597-zw7kh               1/1     Running   0          2m
# cert-manager-cainjector-577f6d9fd7-tr77l   1/1     Running   0          2m
# cert-manager-webhook-787858fcdb-nlzsq      1/1     Running   0          2m
```

<br>

## [Cloudflare Origin CA Issuer](https://github.com/cloudflare/origin-ca-issuer)
```bash
# Download the offical cloudflare origin-ca-issuer repo
git clone https://github.com/cloudflare/origin-ca-issuer.git

# Enter the repo directory
cd ./origin-ca-issuer

# Install the Custom Resource Definitions
kubectl apply -f deploy/crds

# Install the RBAC rules
kubectl apply -f deploy/rbac

# Install the controller
kubectl apply -f deploy/manifests

# Confirm the deployment
kubectl get -n origin-ca-issuer pod

# NAME                                READY   STATUS      RESTARTS    AGE
# pod/origin-ca-issuer-1234568-abcdw  1/1     Running     0           1m
```

<br>

## [Exteral-DNS](https://github.com/kubernetes-sigs/external-dns/blob/master/docs/tutorials/cloudflare.md)
Create a yaml file as seen below and deploy it.

```yaml
# External-DNS.yaml.old
apiVersion: v1
kind: ServiceAccount
metadata:
  name: external-dns
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: external-dns
rules:
- apiGroups: [""]
  resources: ["services","endpoints","pods"]
  verbs: ["get","watch","list"]
- apiGroups: ["extensions","networking.k8s.io"]
  resources: ["ingresses"] 
  verbs: ["get","watch","list"]
- apiGroups: [""]
  resources: ["nodes"]
  verbs: ["list", "watch"]
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: external-dns-viewer
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: external-dns
subjects:
- kind: ServiceAccount
  name: external-dns
  namespace: default
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: external-dns
spec:
  strategy:
    type: Recreate
  selector:
    matchLabels:
      app: external-dns
  template:
    metadata:
      labels:
        app: external-dns
    spec:
      serviceAccountName: external-dns
      containers:
      - name: external-dns
        image: k8s.gcr.io/external-dns/external-dns:v0.7.6
        args:
        - --source=service # ingress is also possible
        - --domain-filter=example.com # (optional) limit to only example.com domains; change to match the zone created above.
        - --zone-id-filter=02d9c2d73f48310a31b8aecef3e0c353 # (optional) limit to a specific zone.
        - --provider=cloudflare
        - --cloudflare-proxied # (optional) enable the proxy feature of Cloudflare (DDOS protection, CDN...)
        env:
        - name: CF_API_KEY # CF_API_TOKEN is preferred; if used, CF_API_KEY and CF_API_EMAIL can be removed.
          value: "YOUR_CLOUDFLARE_API_KEY"
        - name: CF_API_EMAIL
          value: "YOUR_CLOUDFLARE_EMAIL"
        - name: CF_API_KEY # Can be removed if using CF_API_KEY + CF_API_EMAIL combo.
          value: "YOUR_CF_API_TOKEN"
```
> In order to make an API token you will need to
> 1. Navigate to your [cloudflate](https://dash.cloudflare.com/profile/api-tokens) profile
> 2. Copy your Origin CA key
> 3. Create new API token using the following permissions and copy the API key:
     >    * i. Zone > DNS > Edit
>    * ii. Zone > Zone > Read
>    * iii. Include > Specific Zone > "Your TLD"

<br>

## *[Kubed](https://appscode.com/products/kubed/v0.12.0/setup/install/)*
We will need to sync a file across all the namespaces in order to properly sign our certificates
```bash
# Run locally
helm repo add appscode https://charts.appscode.com/stable/
helm repo update
helm search repo appscode/kubed --version v0.12.0
# NAME            CHART VERSION APP VERSION DESCRIPTION
# appscode/kubed  v0.12.0    v0.12.0  Kubed by AppsCode - Kubernetes daemon

helm template kubed appscode/kubed \
  --version v0.12.0 \
  --namespace kube-system \
  --no-hooks | kubectl apply -f -
```

<br>

# **Deploy your Infrastructure**
Create an issue service secret that hold your Cloudflare Origin API key and deploy it
```yaml
# service-key.yaml
apiVersion: v1
kind: Secret
metadata:
  name: service-key
  namespace: cert-manager
  annotations:
    kubed.appscode.com/sync: "cert-manager-tls=example-com" # replce with a lable of your choice used to link synged secrets with labeled namespaces
type: Opaque
# Required to be in base64 format
# echo -n "yourKey" | base64 -w 0
data:
  key: |
    YOUR_BASE64_ENCODED_API_KEY
```

Create an issuer yaml file and deploy it
```yaml
# issuer.yaml
# required for each unique namespace adn cannot be synced automatically.
apiVersion: cert-manager.k8s.cloudflare.com/v1
kind: OriginIssuer
metadata:
  name: prod-issuer
  namespace: default
spec:
  requestType: OriginECC
  auth:
    serviceKeyRef:
      name: service-key
      key: key
```
```bash
# Verify successful setup and connection
kubectl get originissuer.cert-manager.k8s.cloudflare.com prod-issuer -n cert-manager -o json | jq .status.conditions
```

<br>

## Deploy your first app
### All unique namespaces require the following
```yaml
# namespace.yaml
apiVersion: v1
kind: Namespace
metadata:
  name: matchingNamespace
  labels: # required for tls automation
    cert-manager-tls: example-com # Required for tls automation and must match the label from service-key.yaml
```


```yaml
# issuer.yaml
apiVersion: cert-manager.k8s.cloudflare.com/v1
kind: OriginIssuer
metadata:
  name: prod-issuer
  namespace: matchingNamespace
spec:
  requestType: OriginECC
  auth:
    serviceKeyRef:
      name: service-key
      key: key
```


```yaml
# cert.yaml
apiVersion: cert-manager.io/v1
kind: Certificate
metadata:
  name: example-crt
  namespace: matchingNamespace
spec:
  # The secret name where cert-manager should store the signed certificate
  secretName: example-tls
  dnsNames:
    - example.com
    - sub.example.com
  # Duation of the certificate
  duration: 168h
  # Renew a day before the certificate expiration
  renewBefore: 24h
  # Reference the Origin CA Issuer you created above, which must be in the same namespace.
  issuerRef:
    group: cert-manager.k8s.cloudflare.com
    kind: OriginIssuer
    name: prod-issuer
```

<br>

### Required for each unique domain name:
```yaml
apiVersion: v1
kind: Service
metadata:
  name: service-lb
  namespace: matchingNamespace
  annotations:
    external-dns.alpha.kubernetes.io/hostname: sub.example.com # Required | or example.com
    external-dns.alpha.kubernetes.io/ttl: "120" #optional
spec:
  type: LoadBalancer # required for service discovery
  ...
```


```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: app-ingress
  namespace: matchingNamespace
  annotations:
    cert-manager.io/issuer: prod-issuer # required
    cert-manager.io/issuer-kind: OriginIssuer # Required
    cert-manager.io/issuer-group: cert-manager.k8s.cloudflare.com # Required
spec:
  tls:
    # specifying a host in the TLS section will tell cert-manager what
    # DNS SANs should be on the created certificate.
    - hosts:
        - sub.example.com
      # cert-manager will create this secret (Referenced from the cert that was previously created)
      secretName: example-tls
  rules:
    ...
```

<br>

Example Deployment:
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx
spec:
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
        - image: nginx
          name: nginx
          ports:
            - containerPort: 80
          resources:
            limits:
              memory: "300Mi"
              cpu: "600m"
---
apiVersion: v1
kind: Service
metadata:
  name: nginx
  annotations:
    external-dns.alpha.kubernetes.io/hostname: test.example.com
    external-dns.alpha.kubernetes.io/ttl: "120" #optional | default is 300
spec:
  selector:
    app: nginx
  type: LoadBalancer
  ports:
    - protocol: TCP
      port: 80
      targetPort: 80
---
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  annotations:
    cert-manager.io/issuer: prod-issuer
    cert-manager.io/issuer-kind: OriginIssuer
    cert-manager.io/issuer-group: cert-manager.k8s.cloudflare.com
  name: example
  namespace: default
spec:
  rules:
    - host: test.example.com
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: nginx
                port:
                  number: 80
  tls:
    # specifying a host in the TLS section will tell cert-manager what
    # DNS SANs should be on the created certificate.
    - hosts:
        - test.example.com
      # cert-manager will create this secret
      secretName: example-tls
```


____


### adguard + tls
Add a config map with cloudflare ca root cert
https://developers.cloudflare.com/ssl/origin-configuration/origin-ca#4-required-for-some-add-cloudflare-origin-ca-root-certificates
Name the key in the config map and write down the exact key name

Create a mount path similar to
```yaml
apiVersion: v1
kind: Pod
metadata:
  name: cacheconnectsample
spec:
      containers:
      - name: cacheconnectsample
        image: cacheconnectsample:v1
        volumeMounts:
        - name: ca-pemstore
          mountPath: /etc/ssl/certs/my-cert.pem
          subPath: my-cert.pem
          readOnly: false
      volumes:
      - name: ca-pemstore
        configMap:
          name: ca-pemstore
```

Once finished edit the deployment yaml to add:
```yaml
spec:
  containers:
    - name: auth
      image: {{my-service-image}}
      env:
        - name: NODE_ENV
          value: "docker-dev"
      resources:
        requests:
          cpu: 100m
          memory: 100Mi
      ports:
        - containerPort: 3000
// ADD THIS SECTION
    lifecycle:
        postStart:
          exec:
            command: ["/bin/sh", "-c", "cd /etc/ssl/certs && cat youraddedrootca.pem >> ca-certificates.crt"]
```