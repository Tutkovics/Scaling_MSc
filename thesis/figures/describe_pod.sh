$ kubectl describe pod front-end-54b6ffc64c-cd8kl
...
    Command:
      /app/main
      -name=front-end
      -port=80
      -cpu=100
      -memory=100
      -endpoint-url=/instant
      -endpoint-delay=1
      -endpoint-call=''
      -endpoint-cpu=10
      -endpoint-url=/chain
      -endpoint-delay=2
      -endpoint-call='back-end:80/profile_db:36/get'
      -endpoint-cpu=20
...

$ kubectl run shell --rm -it --image nicolaka/netshoot -- sh
~ # curl front-end:80/instant
{"service":"front-end","host":"front-end","config":true,"endpoint":"/instant","cpu":10,"delay":1,"calloutparameter":"''","callouts":[""],"actualDelay":16967552,"time":"2021-05-13T16:55:02.951824587Z","requestMethod":"GET","requestURL":{"Scheme":"","Opaque":"","User":null,"Host":"","Path":"/instant","RawPath":"","ForceQuery":false,"RawQuery":"","Fragment":"","RawFragment":""},"requestAddr":"10.0.1.49:43678"}   