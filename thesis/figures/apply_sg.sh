$ kubectl apply -f path/to/servicegraph.yaml                          
servicegraph.dipterv.my.domain/servicegraph created 

$ kubectl get all
NAME                             READY   STATUS    RESTARTS   AGE
pod/back-end-59698d4b95-2tmhm    1/1     Running   0          24s
pod/back-end-59698d4b95-ftgtb    1/1     Running   0          24s
pod/back-end-59698d4b95-pnf2p    1/1     Running   0          24s
pod/db-65bb88cdcc-bclv6          1/1     Running   0          24s
pod/front-end-54b6ffc64c-9rjg2   1/1     Running   0          24s
pod/front-end-54b6ffc64c-cd8kl   1/1     Running   0          24s

NAME                 TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)        AGE
service/back-end     ClusterIP   10.108.141.216   <none>        80/TCP         24s
service/db           ClusterIP   10.96.131.117    <none>        36/TCP         24s
service/front-end    NodePort    10.102.37.115    <none>        80:30000/TCP   24s
service/kubernetes   ClusterIP   10.96.0.1        <none>        443/TCP        25d

NAME                        READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/back-end    3/3     3            3           24s
deployment.apps/db          1/1     1            1           24s
deployment.apps/front-end   2/2     2            2           24s

NAME                                   DESIRED   CURRENT   READY   AGE
replicaset.apps/back-end-59698d4b95    3         3         3       24s
replicaset.apps/db-65bb88cdcc          1         1         1       24s
replicaset.apps/front-end-54b6ffc64c   2         2         2       24s

NAME                                            REFERENCE              TARGETS         MINPODS   MAXPODS   REPLICAS   AGE                                                                                           
horizontalpodautoscaler.autoscaling/front-end   Deployment/front-end   <unknown>/10%   2         5         2          24s             
