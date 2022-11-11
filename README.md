# deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: <app-name>
  labels:
    app: <app-name>
spec:
  replicas: 1
  selector:
    matchLabels:
      app: <app-name>
  template:
    metadata:
      labels:
        app: <app-name>
    spec:
      serviceAccountName: <service-account-name>
      containers:
      - name: <app-name>
        image: <image-url>
        ports:
        - containerPort: <container-port>
        volumeMounts:
        - name: <app-name>
          mountPath: "/mnt/secrets-store"
          readOnly: true
        env:
        - name: <env>
          valueFrom:
            secretKeyRef:
              name: <secret-name-in-secretProviderClass>
              key: <secret-key>
      volumes:
      - name: <app-name>
        csi:
          driver: secrets-store.csi.k8s.io
          readOnly: true
          volumeAttributes:
            secretProviderClass: <secret-provider-class-name-in-secretProviderClass>
        
-----------------------------------------------------------------

# serviceaccount.yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: <service-account-name>
  namespace: default
  annotations:
    eks.amazonaws.com/role-arn: <arn-of-role>
        
-----------------------------------------------------------------

# secretprovider.yaml
apiVersion: secrets-store.csi.x-k8s.io/v1
kind: SecretProviderClass
metadata:
  name: <secret-provider-class-name>
spec:
  provider: aws
  parameters:
    objects: |
        - objectName: <arn-of-secret>
          jmesPath: 
              - path: <secret-key>
                objectAlias: <secret-alias>

  secretObjects:
    - secretName: my-secret-name
      type: Opaque
      data:
      - objectName: <secret-alias>
        key: <secret-key>

-----------------------------------------------------------------

# How to use secret manager with EKS

**1.Connect kubectl to Kube cluster**

**2.Install ASCP**

 ```
helm repo add secrets-store-csi-driver https://kubernetes-sigs.github.io/secrets-store-csi-driver/charts
helm install csi-secrets-store secrets-store-csi-driver/secrets-store-csi-driver --set syncSecret.enabled=true --namespace kube-system

helm repo add aws-secrets-manager https://aws.github.io/secrets-store-csi-driver-provider-aws
helm install -n kube-system secrets-provider-aws aws-secrets-manager/secrets-store-csi-driver-provider-aws
```
or

```
helm repo add secrets-store-csi-driver https://kubernetes-sigs.github.io/secrets-store-csi-driver/charts
helm install csi-secrets-store secrets-store-csi-driver/secrets-store-csi-driver --set syncSecret.enabled=true --namespace kube-system

kubectl apply -f https://raw.githubusercontent.com/aws/secrets-store-csi-driver-provider-aws/main/deployment/aws-provider-installer.yaml
```

