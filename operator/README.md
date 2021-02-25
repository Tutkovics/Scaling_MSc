# How to write an operator 
The following commands will generate your operator. You only need to set the variables to your preference. 

```
# if you're not in GOPATH:
go mod init example.com/m/v2
operator-sdk init --owner "Tutkovics András" --project-version 3 --project-name service-mesh-operator
operator-sdk create api --help
```