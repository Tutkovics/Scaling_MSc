$ kubectl get horizontalpodautoscalers
NAME         REFERENCE               TARGETS   MINPODS MAXPODS REPLICAS AGE
back-end-1   Deployment/back-end-1   100%/70%  1       5       5        5m34s
front-end-1  Deployment/front-end-1  80%/70%   1       5       4        5m34s

$ kubectl get all
NAME                               READY   STATUS    RESTARTS   AGE
pod/back-end-1-7cb9d49594-7btf8    1/1     Running   0          2m48s
pod/back-end-1-7cb9d49594-fpgcj    0/1     Pending   0          108s
pod/back-end-1-7cb9d49594-lhjvl    0/1     Pending   0          108s
pod/back-end-1-7cb9d49594-rvsgk    1/1     Running   0          5m34s
pod/back-end-1-7cb9d49594-vdl7m    1/1     Running   0          3m48s
pod/front-end-1-6677c56849-26nmn   1/1     Running   0          5m34s
pod/front-end-1-6677c56849-26zr2   1/1     Running   0          3m48s
pod/front-end-1-6677c56849-8s9ph   0/1     Pending   0          48s
pod/front-end-1-6677c56849-fndnw   1/1     Running   0          2m48s