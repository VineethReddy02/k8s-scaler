# k8s-scaler

This project is a kubernetes resource scaler. The main objective of this project is to create 10's/100's/1000's of kubernetes resources with ease. Infact by just using a single command.

### Usage

Clone this repository, you can find ```k8s-scaler``` executable under ```bin``` directory. Running ```./build.sh``` generates latest **k8s-scaler** executable.

```yaml
vineeth@vineeth-Latitude-7490 /bin  $ ./k8s-scaler --help

k8s-scaler helps to create 1000's kubernetes resources uing a single command.
It is built to ease the k8s resources simulation easy in large kubernetes clusters.,

Usage:
  k8s-scaler [command]

Available Commands:
  create      To create deployments/daemonsets/pods/namespaces/statefulsets/jobs/cronjobs
  delete      To delete deployments/daemonsets/pods/namespaces/statefulsets/jobs/cronjobs
  help        Help about any command
  list        To list namespaces, deployments, pods, daemonsets, statefulsets, jobs, cronjobs.

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
(or)
./k8s-scaler create d --scale 250 --replicas 25 --containers 10
```

#### To create deployments in a random namespace but exclude couple of namespaces

```yaml
./k8s-scaler create deployments --scale 250 --replicas 25 --containers 10 --exclude-namespaces namespace01,namespace02
```

#### To create deployments in a specific namespace and to schedule resources on desired node

```yaml
./k8s-scaler create deployments --scale 250 --replicas 25 --containers 10 --namespace namepsace01 --node-selector kubernetes.io/hostname=mock-kubelet --toleration mocklet.io/provider=mock
```

#### To create daemonsets in a random namespace

```yaml
./k8s-scaler create daemonsets --scale 50 --containers 10
(or)
./k8s-scaler create ds --scale 50 --containers 10
```

#### To create daemonsets in a random namespace but to exclude couple of namespaces

```yaml
./k8s-scaler create daemonsets --scale 50 --containers 10 --exclude-namespaces namespace01,namespace02
```

#### To create daemonsets in specific namespace and to schedule resources on desired node

```yaml
./k8s-scaler create daemonsets --scale 50 --containers 10 --namespace namespace01 --node-selector kubernetes.io/hostname=mock-kubelet --toleration mocklet.io/provider=mock
```

#### To create pods in a random namespace

```yaml
./k8s-scaler create pods --scale 500  --containers 10
(or)
./k8s-scaler create p --scale 500  --containers 10
```

#### To create pods in a random namespace but exclude couple of namespaces

```yaml
./k8s-scaler create pods --scale 50 --containers 10 --exclude-namespaces namespace01,namespace02
```

#### To create pods in specific namespace and to schedule resources on desired node

```yaml
./k8s-scaler create pods --scale 50 --containers 10 --namespace namespace01 --node-selector kubernetes.io/hostname=mock-kubelet --toleration mocklet.io/provider=mock
```

#### To create statefulsets in a random namespace

```yaml
./k8s-scaler create statefulsets --scale 500 --replicas 3 --containers 10
(or)
./k8s-scaler create s --scale 500 --replicas 3 --containers 10
```

#### To create statefulsets in a random namespace but exclude couple of namespaces

```yaml
./k8s-scaler create statefulsets --scale 500 --replicas 3 --containers 10 --exclude-namespaces namespace01,namespace02
```

#### To create statefulsets in specific namespace and to schedule resources on desired node

```yaml
./k8s-scaler create statefulsets --scale 500 --replicas 3 --containers 10 --namespace namespace01 --node-selector kubernetes.io/hostname=mock-kubelet --toleration mocklet.io/provider=mock
```

**Note:** All the jobs created are by default configured to sleep for 1 minute and move to completed state.

#### To create jobs in a random namespace

```yaml
./k8s-scaler create jobs --scale 500 --containers 10
(or)
./k8s-scaler create j --scale 500 --containers 10
```

#### To create jobs in a random namespace but exclude couple of namespaces

```yaml
./k8s-scaler create jobs --scale 500 --containers 10 --exclude-namespaces namespace01,namespace02
```

#### To create jobs in specific namespace and to schedule resources on desired node

```yaml
./k8s-scaler create jobs --scale 500 --namespace namespace01 --containers 10 --node-selector kubernetes.io/hostname=mock-kubelet --toleration mocklet.io/provider=mock
```

**Note:** All the cron jobs created are by default configured to sleep for 1 minute and to run for every 30 minutes.

#### To create cron jobs in a random namespace

```yaml
./k8s-scaler create cronjobs --scale 500 --containers 10
(or)
./k8s-scaler create cj --scale 500 --containers 10
```

#### To create cron jobs in a random namespace but exclude couple of namespaces

```yaml
./k8s-scaler create cronjobs --scale 500 --containers 10 --exclude-namespaces namespace01,namespace02
```

#### To create cron jobs in specific namespace and to schedule resources on desired node

```yaml
./k8s-scaler create cronjobs --scale 500 --namespace namespace01 --containers 10 --node-selector kubernetes.io/hostname=mock-kubelet --toleration mocklet.io/provider=mock
```

#### To create replicationcontrollers in a random namespace.
```
./k8s-scaler create replicationcontrollers --scale 10 --replicas 5 --containers 15
(or)
./k8s-scaler create rc --scale 10 --replicas 5 --containers 15
```

#### To create replicationcontrollers in a random namespace but exclude couple of namespaces.
```
./k8s-scaler create replicationcontrollers --scale 10 --replicas 5 --containers 15 --exclude-namespaces namespace01,namespace02
(or)
./k8s-scaler create rc --scale 10 --replicas 5 --containers 15 --exclude-namespaces namespace01,namespace02
```

#### To create replicasets in a random namespace.
```
./k8s-scaler create replicasets --scale 10 --replicas 5 --containers 15 
(or)
./k8s-scaler create rs --scale 10 --replicas 5 --containers 15
```

#### To create replicasets in a random namespace but exclude couple of namespaces.
```
./k8s-scaler create replicasets --scale 10 --replicas 5 --containers 15 --exclude-namespaces namespace01,namespace02
(or)
./k8s-scaler create rs --scale 10 --replicas 5 --containers 15 --exclude-namespaces namespace01,namespace02
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
(or)
./k8s-scaler delete d --scale 500 --namespace namespace01
```

**Note:**
Deletion of resources can be performed same as above provided example for pods/daemonsets/statefulsets/replicationcontrollers/replicasets/jobs/cronjobs.

You also add pass node-selector & toleration while resource creation to schedule the workloads on the desired node. For now, we only accept one pair of node-selector and toleration. Also, toleration operator is by default set to "Equal" and effect is set to "NoSchedule" key & value from toleration is configurable using flag.   

#### To list namespaces, deployments, pods, daemonsets, statefulsets, jobs, cronjobs

```yaml
vineeth@vineeth-Latitude-7490 /bin (master) $ ./k8s-scaler list
NAMESPACE         DEPLOYMENTS     REPLICASETS     DAEMONSETS      STATEFULSETS    PODS        JOBS        CRONJOBS    REPLICATION-CONTROLLERS
test              3000            3000            1000            500             7486        30          10          30               
default           1300            1300            456             250             5642        10          5           5                            
kube-system       8               11              4               0               15          0           0           0               
mock-kubelet      3500            4000            1200            400             9348        50          30          35     
```


#### To create & delete periodically the deployments with provided time interval (i.e in seconds)

The below cmd would create & delete 100 deployments with replicas as 50 for every 10 seconds in scale namespace.
 
```
./k8s-scaler chaos d --scale 100 -r 50 --time 10 -n scale
```

#### TODO:

1. Support custom resources creation at scale. As CRD already exists in cluster. CR creation can be done using k8s-scaler by passing ```--kind``` value but custom resources spec formation needs to be taken care in k8s-scaler.
2. Support services and other Kubernetes resources.

