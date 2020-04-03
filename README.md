# k8s-scaler

This project is a kubernetes resource scaler. The main objective of this project is to create 10's/100's/100's of kubernetes resources with ease. Infact by just using a single command.

### Usage

Clone this repository, you can find ```k8s-scaler``` executable under ```bin``` directory. Running ```./build.sh``` generates latest **k8s-scaler** executable.

```yaml
vineeth@vineeth-Latitude-7490 /bin  $ ./k8s-scaler --help

k8s-scaler helps to create 1000's kubernetes resources uing a single command.
It is built to ease the k8s resources simulation easy in large kubernetes clusters.,

Usage:
  k8s-scaler [command]

Available Commands:
  create      To create deployments/daemonsets/pods/namespaces
  delete      To delete deployments/daemonsets/pods/namespaces
  help        Help about any command
  list        To list namespaces, deployments, pods, daemonsets.

Flags:
  -h, --help                help for k8s-scaler
  -k, --kubeconfig string   Path to the KUBECONFIG file.

Use "k8s-scaler [command] --help" for more information about a command.

```

**Note:**

You can configure the cluster details by passing the **KUBECONFIG** file path to ```--kubeconfig``` global flag.

```yaml
./k8s-scaler create pods --scale 500 --containers 10 --kubeconfig /home/vineeth/gke.yaml
```
If ```--kubeconfig``` flag is not provided. k8s-scaler tries to read **KUBECONFIG** environment variable.

If **KUBECONFIG** env variable is not set. k8s-scaler tries to find ```InClusterConfig``` using k8s **client-go** library.

#### To create deployments in a random namespace

```yaml
./k8s-scaler create deployments --scale 250 --replicas 25 --containers 10
```

#### To create deployments in a random namespace but exclude couple of namespaces

```yaml
./k8s-scaler create deployments --scale 250 --replicas 25 --containers 10 --exclude-namespaces namespace01,namespace02
```

#### To create deployments in a specific namespace

```yaml
./k8s-scaler create deployments --scale 250 --replicas 25 --containers 10 --namespace namepsace01
```

#### To create daemonsets in a random namespace

```yaml
./k8s-scaler create daemonsets --scale 50 --containers 10
```

#### To create daemonsets in a random namespace but to exclude couple of namespaces

```yaml
./k8s-scaler create daemonsets --scale 50 --containers 10 --exclude-namespaces namespace01,namespace02
```

#### To create daemonsets in specific namespace

```yaml
./k8s-scaler create daemonsets --scale 50 --containers 10 --namespace namespace01
```

#### To create pods in a random namespace

```yaml
./k8s-scaler create pods --scale 500  --containers 10
```

#### To create pods in a random namespace but exclude couple of namespaces

```yaml
./k8s-scaler create pods --scale 50 --containers 10 --exclude-namespaces namespace01,namespace02
```

#### To create pods in specific namespace

```yaml
./k8s-scaler create pods --scale 50 --containers 10 --namespace namespace01 
```

#### To create statefulsets in a random namespace

```yaml
./k8s-scaler create statefulsets --scale 500 --replicas 3 --containers 10
```

#### To create statefulsets in a random namespace but exclude couple of namespaces

```yaml
./k8s-scaler create statefulsets --scale 500 --replicas 3 --containers 10 --exclude-namespaces namespace01,namespace02
```

#### To create statefulsets in specific namespace

```yaml
./k8s-scaler create statefulsets --scale 500 --replicas 3 --containers 10 --namespace namespace01 
```

#### To create jobs in a random namespace

```yaml
./k8s-scaler create jobs --scale 500 --containers 10
```

#### To create jobs in a random namespace but exclude couple of namespaces

```yaml
./k8s-scaler create jobs --scale 500 --containers 10 --exclude-namespaces namespace01,namespace02
```

#### To create jobs in specific namespace

```yaml
./k8s-scaler create jobs --scale 500 --namespace namespace01 --containers 10
```

#### To create namespaces

```yaml
./k8s-scaler create namespaces --scale 100
```

#### To delete deployments across multiple namespaces

```yaml
./k8s-scaler delete deployments --scale 500
```

#### To delete deployments across multiple namespaces but to exclude couple of namespaces

```yaml
./k8s-scaler delete deployments --scale 500 --exclude-namespaces namespace01,namespace02
```

#### To delete deployments in a specific namespace

```yaml
./k8s-scaler delete deployments --scale 500 --namespace namespace01
```

**Note:**
Deletion of resources can be performed same as above provided example for pods/daemonsets/statefulsets/jobs.

#### To list namespaces, deployments, pods, daemonsets

```yaml
NAMESPACE         DAEMONSETS      DEPLOYMENTS     STATEFULSETS    PODS        JOBS        
aqua              5               0               0               10          0           
default           0               0               0               0           0           
kube-node-lease   0               0               0               0           0           
kube-public       0               0               0               0           0           
kube-system       9               4               0               31          0           
monitoring        1               0               3               10          4   
```
#### TODO:

1. Support custom resources creation at scale. As CRD already exists in cluster. CR creation can be done using k8s-scaler by passing ```--kind``` value but custom resources spec formation needs to be taken care in k8s-scaler.
2. Support services and other Kubernetes resources.

