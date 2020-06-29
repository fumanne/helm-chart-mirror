### 背景
之前国内使用阿里云的helm chart 镜像库, 但由于https://kubernetes.oss-cn-hangzhou.aliyuncs.com/charts 一直不更新, 所以决定自己同步
参考了网上的同步方案, 决定利用golang 的goroutine 来快速同步

### 使用
    $ helm repo add mirror https://fumanne.github.io/helm-chart-mirror
    $ helm repo list
    $ NAME     	URL
      stable   	https://kubernetes-charts.storage.googleapis.com
      incubator	https://kubernetes-charts-incubator.storage.googleapis.com
      mirror  	https://fumanne.github.io/helm-chart-mirror
    
    # search chart
    $ helm search repo grafana
    NAME             	CHART VERSION	APP VERSION	DESCRIPTION
    mirror/grafana  	5.3.0        	7.0.3      	The leading tool for querying and visualizing t...
    incubator/grafana	0.1.3        	           	A Helm chart for Kubernetes
    stable/grafana   	5.3.0        	7.0.3      	The leading tool for querying and visualizing t...
    
    
    # Install Chart   
    $  helm install grafana --namespace=kube-system mirror/grafana
    
### 附
会不定期同步官方repo...需要同步请联系