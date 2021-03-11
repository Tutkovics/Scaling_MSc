# Source of measurement script

For the project we need to measure the behaviour of Kubernetes cluster under loading.
To create deterministic measurements a Python script will orchestrate it.

## Prerequisite
- Installed [Fortio](https://github.com/fortio/fortio)
- Install necessary Python modules 
  ```
  $ pip freeze > requirements.txt    # to freeze 
  $ pip install -r requirements.txt  # to install
  ```
- Configured ***kubectl*** command

# Start measurement
- Create measurement config file
- Start the script
```
python main.py path/to/measurement_configs/example.yaml
```