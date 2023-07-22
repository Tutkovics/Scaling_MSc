# Scaling project #

This project implements a Kubernetes operator that enables users to describe the configuration of micro-services and their connections, resource requirements for serving incoming requests (memory, CPU), and introduces additional delays if defined. The operator is written in Go language and manages the deployment and HPA scaler of a generic application (`worker-node`), which can be configured to handle the incoming requests based on the micro-service specifications.

There are some articles about the background under [articles](./articles/) directory and my final [thesis](./thesis/thesis.pdf) also attached.

## How it Works
1. **Define service-mesh**: Users create custom resource definitions (CRDs) that define the micro-services they want to deploy. Each CRD contains specifications for services, their connections, resource requirements, and optional additional delays. Examples under [yaml](./yaml/) directory.

1. **Kubernetes operator**: The [operator](./operator/) continuously monitors the Kubernetes cluster for new or updated CRDs. When it detects a new configuration or changes to an existing one, it takes appropriate actions to reconcile the desired state with the current state of the cluster.

1. **Application deployment**: Based on the micro-service configurations, the operator deploys instances of the generic ([worker-node](./worker-node/)) application, connecting them as specified in the CRDs.
  - Resource Management - The operator doesn't manage resource allocation just configure the `Deployment`s and `HPA`s.
  - Additional Delay Management - If configured to add additional delays, the operator configure it at application level.

- +1. **Python script for automatic benchmark**: After the service-mesh specifications (CRD) there is a Python script to help with automatic benchmark. User can configure what service-mesh they want to use, and what benchmark parameters. More examples [here](./python_scripts/measurement_configs).

## Installation
To install the Kubernetes operator and deploy the service-mesh follow these steps:

1. Clone the repository and build the operator and generic application using the provided Go build commands.
1. Deploy the operator to your Kubernetes cluster using the generated deployment manifest.
1. Create custom resource definitions (CRDs) for the micro-services you want to configure and deploy.
1. Apply the CRDs to the Kubernetes cluster using kubectl or your preferred Kubernetes management tool.

## Contributing
Contributions to this project are welcome. If you encounter any issues or have suggestions for improvements, please open an issue or submit a pull request.

## License
This project is licensed under the [Apache-2.0](LICENSE). Feel free to use, modify, and distribute it as permitted by the license.

## ToDo ##
- [ ] Put an architecture to README
- [ ] Write a description how to use the operator
- [ ] Write operator's limitations
- [ ] Document worker-node (application) usage 
- [ ] 
- [ ] 
