# How to write an operator 
The following commands will generate your operator. You only need to set the variables to your preference. 

```
# if you're not in GOPATH:
go mod init example.com/m/v2
operator-sdk init --owner "Tutkovics András" --project-version 3 --project-name service-mesh-operator
operator-sdk create api --group dipterv --version v1beta1 --kind Servicegraph
```

The code under comes from operator-sdk help screen. Contains: how to create an API and where to edit Scheme & Controller.

```
Examples:
  # Create a frigates API with Group: ship, Version: v1beta1 and Kind: Frigate
  operator-sdk create api --group ship --version v1beta1 --kind Frigate

  # Edit the API Scheme
  nano api/v1beta1/frigate_types.go

  # Edit the Controller
  nano controllers/frigate/frigate_controller.go

  # Edit the Controller Test
  nano controllers/frigate/frigate_controller_test.go

  # Install CRDs into the Kubernetes cluster using kubectl apply
  make install

  # Regenerate code and run against the Kubernetes cluster configured by ~/.kube/config
  make run
```
