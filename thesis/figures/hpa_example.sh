$ cat hpa/hpa.yaml
apiVersion: autoscaling/v1
kind: HorizontalPodAutoscaler
metadata:
  name: php-apache
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: php-apache
  minReplicas: 1
  maxReplicas: 10
  targetCPUUtilizationPercentage: 50

$ kubectl apply -f hpa/hpa.yaml
$ while sleep 0.01; do wget -q -O- http://localhost:30100; done

$ kubectl get horizontalpodautoscalers                                               
NAME        REFERENCE              TARGETS   MINPODS  MAXPODS  REPLICAS  AGE                        
php-apache  Deployment/php-apache  179%/50%  1        10       1         4m53s                      
...
php-apache  Deployment/php-apache  81%/50%   1        10       4         6m35s                      
...
php-apache  Deployment/php-apache  42%/50%   1        10       7         10m  
