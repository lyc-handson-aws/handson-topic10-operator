#update the generated code for that resource type:
make generate

#Once the API is defined with spec/status fields and CRD validation markers or the controller is changed, generate or update the CRD manifests
make manifests

# build and push custom operator controller image
make docker-build docker-push


# deploy controller to the cluster
make deploy