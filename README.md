apiVersion: apps/v1
kind: Deployment
metadata:
  name: secret
  labels:
    app: secret
spec:
  replicas: 1
  selector:
    matchLabels:
      app: secret
  template:
    metadata:
      labels:
        app: secret
    spec:
      containers:
      - name: secret
        image: image-url
        ports:
        - containerPort: 8080
        env:
        - name: SECRET
          value: "This is secret from deployment.yaml"
