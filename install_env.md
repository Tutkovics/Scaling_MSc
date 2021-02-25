# Prerequisites
First need to install some software. Some code below only works on debian based systems.

# Install VirtualBox
```
sudo apt upgrade  
wget -q https://www.virtualbox.org/download/oracle_vbox_2016.asc -O- | sudo apt-key add -
echo "deb [arch=amd64] http://download.virtualbox.org/virtualbox/debian bionic contrib" | sudo tee /etc/apt/sources.list.d/virtualbox.list 
sudo apt install virtualbox-6.1
```

# Install Minikube to have local K8s cluster
```
curl -LO https://storage.googleapis.com/minikube/releases/latest/minikube_latest_amd64.deb
sudo dpkg -i minikube_latest_amd64.deb
minikube start -p dipterv
```

# Install kubectl command
```
curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"
sudo install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl
```
## Or use package manager
```
sudo apt-get update && sudo apt-get install -y apt-transport-https gnupg2 curl
curl -s https://packages.cloud.google.com/apt/doc/apt-key.gpg | sudo apt-key add -
echo "deb https://apt.kubernetes.io/ kubernetes-xenial main" | sudo tee -a /etc/apt/sources.list.d/kubernetes.list
sudo apt-get update
sudo apt-get install -y kubectl
```
## Check the install/config
```
kubectl version --client
kubectl cluster-info
```
## Add to ~/.zshrc
```
source <(kubectl completion zsh)
```

# Install Go
```
wget https://golang.org/dl/go1.16.linux-amd64.tar.gz 
tar -C /usr/local -xzf go1.16.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
go version
```

# Install Operator-sdk
```
git clone https://github.com/operator-framework/operator-sdk
cd operator-sdk
git checkout master
make install
```
Perhaps you need to add ~/go/bin to your PATH.

# Install Docker (if need)
```
sudo apt-get install apt-transport-https ca-certificates \
    curl gnupg-agent software-properties-common
curl -fsSL https://download.docker.com/linux/debian/gpg | sudo apt-key add -
sudo add-apt-repository \
   "deb [arch=amd64] https://download.docker.com/linux/debian \
   $(lsb_release -cs) \
   stable"
sudo apt-get update
sudo apt-get install docker-ce docker-ce-cli containerd.io
```
## Add current user do docker group
So don't need sudo for docker commands.
```
sudo groupadd docker
sudo usermod -aG docker $USER 
```