# How to use secret manager with EKS

**Create EKS cluster and Create a node group with some worker nodes**



**Connect kubectl to EKS cluster**

1. Use aws cli `aws configure` and put in your AWS access key id, AWS secret access key, region code and output format(json)

2. Connect to EKS cluster by command `aws eks update-kubeconfig --region region-code --name cluster-name` and `kubectl get svc`




**Install ASCP with these 3 commands**

```
helm repo add secrets-store-csi-driver https://kubernetes-sigs.github.io/secrets-store-csi-driver/charts

helm install csi-secrets-store secrets-store-csi-driver/secrets-store-csi-driver --set syncSecret.enabled=true --namespace kube-system

kubectl apply -f https://raw.githubusercontent.com/aws/secrets-store-csi-driver-provider-aws/main/deployment/aws-provider-installer.yaml
```




**Create OIDC Provider (from https://docs.aws.amazon.com/eks/latest/userguide/enable-iam-roles-for-service-accounts.html)**

1. Open the Amazon EKS console at https://console.aws.amazon.com/eks/home#/clusters.

2. In the left pane, select Clusters, and then select the name of your cluster on the Clusters page.

3. In the Details section on the Overview tab, note the value of the OpenID Connect provider URL.

4. Open the IAM console at https://console.aws.amazon.com/iam/.

5. In the left navigation pane, choose Identity Providers under Access management. If a Provider is listed that matches the URL for your cluster, then you already have a provider for your cluster. If a provider isn't listed that matches the URL for your cluster, then you must create one.

6. To create a provider, choose Add provider.

7. For Provider type, select OpenID Connect.

8. For Provider URL, enter the OIDC provider URL for your cluster, and then choose Get thumbprint.

9. For Audience, enter sts.amazonaws.com and choose Add provider.




**Store a secret in Secret Manager**

1. Go to Secret Manager

2. Click 'Store a new secret' at right hand side

3. Choose 'Other type of secret'

4. add a pair/pairs of key-value

5. Choose Encrption key as 'aws/secretsmanager' and click 'Next'

6. Name your secret e.g. 'dev/Sumato/Refresher' and click 'Next'

7. click 'Next' again

8. Review and click 'Store'

9. Go into created secret and keep 'Secret ARN' for incoming processes




**Create a service account in EKS cluster**

-- Create policy

1. Go to IAM

2. Click 'Policies' and click 'Create policy' 

3. Click tab 'JSON'

4. Use iam-policy.json as template

5. You need to fill your account id, region code and Secret ARN into the template

6. Click 'Next: Tags' and click 'Next: Review'

6. Name your policy e.g. policy-for-pod

7. Review and click 'Create policy'


-- Create role and attach policy above to the role

1. Go to IAM

2. Click 'Roles' and click 'Create role'

3. Choose 'custom trust policy' and use 'iam-role.json' as template

4. You need to fill your account id, oidc provider, namespace and service account name

5. click 'Next'

6. In 'Add permissions' section, check the policy that you just created above

7. Click 'Next', name your role e.g. role-for-pod and click 'Create role'


-- Apply service account to the pod with created role

1. Create serviceaccount.yaml or you can name file anything you want

2. Fill service account name and ARN of role that you just created above
-----------------------------------------------------------------
apiVersion: v1
kind: ServiceAccount
metadata:
  name: <service-account-name>
  namespace: default
  annotations:
    eks.amazonaws.com/role-arn: <arn-of-role>

-----------------------------------------------------------------

3. `kubectl apply -f serviceaccount.yaml` to create a new service account to Kubernetes cluster




**Create Secret Provider Class in EKS cluster**

1. Create secretprovider.yaml or you can name file anything you want

2. Fill your data
-----------------------------------------------------------------
apiVersion: secrets-store.csi.x-k8s.io/v1
kind: SecretProviderClass
metadata:
  name: <secret-provider-class-name>
spec:
  provider: aws
  parameters:
    objects: |
        - objectName: "<arn-of-secret>"
          jmesPath: 
              - path: <secret-key>
                objectAlias: <secret-alias>

  secretObjects:
    - secretName: <secret-name>
      type: Opaque
      data:
      - objectName: <secret-alias>
        key: <secret-key>

-----------------------------------------------------------------

- <secret-provider-class-name> will be used to identify secret provider class in deployment.yaml
- <secret-name> will be used to identify secret in deployment.yaml
- <arn-of-secret> is Secret ARN
- <secret-key> is the key of key-value that you store in Secret Manager
- <secret-alias> is to link between objects and secretObjects

3. `kubectl apply -f secretprovider.yaml` to create a new secret provider to Kubernetes cluster




**Deploy a pod to EKS cluster**

1. Create deployment.yaml or you can name file anything you want

2. Fill your data

-----------------------------------------------------------------
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
              name: <secret-name-from-secretprovider.yaml>
              key: <secret-key>
      volumes:
      - name: <app-name>
        csi:
          driver: secrets-store.csi.k8s.io
          readOnly: true
          volumeAttributes:
            secretProviderClass: "<secret-provider-class-name-from-secretprovider.yaml>"
        
-----------------------------------------------------------------

- <app-name> your application name
- <service-account-name> service account name
- <env> is environment variable that you use in your application
- <secret-key> is the key of key-value that you store in Secret Manager

3. `kubectl apply -f deployment.yaml` to deploy an app to EKS cluster




**Expose your pod with service**

1. Create deployment.yaml or you can name file anything you want

2. Fill your data

-----------------------------------------------------------------
apiVersion: v1
kind: Service
metadata:
  name: <app-name>
spec:
  type: NodePort
  selector:
    app: <app-name>
  ports:
    - nodePort: <node-port>
      port: 80
      targetPort: <container-port>

-----------------------------------------------------------------

- <node-port> is the port that you want to expose and You also need to expose this port in AWS security group
