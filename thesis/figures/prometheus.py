# Prometheus query parameters
cpu_query = {"query": "container_cpu_usage_seconds_total{image!='', \
                       namespace=~'default|metrics', \
                       container!='POD'}",
    "start": str(start_time),
    "end": str(end_time),
    "step": "2",
    "timeout": "1000ms"
}

# Get parameters from config
prometheus_ip = config["prometheus_ip"]
prometheus_port = config["prometheus_port"]

# Assemble prometheus base query URL
url = "http://" + str(prometheus_ip) + ":" + str(prometheus_port) + "/api/v1/query_range?"

# Assemble full query
cpu_full_query = url + urllib.parse.urlencode(cpu_query)

# Get results from Prometheus
cpu_res = json.loads(requests.get(cpu_full_query).text)
