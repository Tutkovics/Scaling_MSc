import logging
import subprocess
import time
import requests
import json
import urllib
from halo import Halo
from datetime import date

def orchestrate_measure(config):
    logging.info("Start 'orchestrate_measure()' function")

    # Iterate through given servicegraphs
    for servicegraph in config["servicegraphs"]:
        logging.info("Servicegraph: %s", servicegraph)

        # Create arrange to use in For loop
        qps_range = range(
            int(config["qps_from"]),
            int(config["qps_to"])+1, # plus 1 to get upper bound
            int(config["qps_granularity"]))

        for qps in qps_range:
            logging.info("New QPS value: %s", qps)
            
            # install servicegraph
            arg = "-f " +  str(servicegraph)
            kubectl("apply", arg)

            # Start loading
            load(config, qps)

            # Delete current service graph
            kubectl("delete", arg)

    logging.info("End 'orchestrate_measure()' function")


def kubectl(command, arguments):
    logging.info("Start 'kubectl(%s , %s)' function" ,command, arguments)

    # Create full command to run
    cmd = "kubectl " + command + " " + arguments
    return_param = subprocess.run(cmd, shell=True, universal_newlines=True, capture_output=True)
    
    # wait all pods to be ready
    wait_to_running_pods()
    logging.info("End 'kubectl()' function") 


def load(config, qps):
    # use Fortio to load generate
    logging.info("Start 'load()' function")

    # Start spinner to be fancy
    spinner = Halo(text='Loading', spinner='dots')
    spinner.start("Fortio measurement is running.")

    if config["load_preheat"] > 0:
        # Create command to run
        # eg: fortio load -qps 10 -t 5s -a http://192.168.49.2:30000/instant
        cmd = "fortio load -qps {qps} -t {time}s -a -data-dir {location} http://{ip}:{port}/{path}?{query} ".format(
            qps = qps,
            time = config["load_preheat"],
            ip = config["load_ip"],
            port = config["load_port"],
            path = config["load_path"],
            query = config["load_query"],
            location = config["result_location"]
        )
        logging.debug("Fortio (preheat) full query: %s", cmd)
        return_param = subprocess.run(cmd, shell=True, universal_newlines=True, capture_output=True)
    
    # Actually measurement
    cmd = "fortio load -qps {qps} -t {time}s -a -data-dir {location} http://{ip}:{port}/{path}?{query}".format(
        qps = qps,
        time = config["load_time"],
        ip = config["load_ip"],
        port = config["load_port"],
        path = config["load_path"],
        query = config["load_query"],
        location = config["result_location"]
    )
    logging.debug("Fortio full query: %s", cmd)
    start_time = time.time()
    return_param = subprocess.run(cmd, shell=True, universal_newlines=True, capture_output=True)
    end_time = time.time()

    # stop spinner
    spinner.stop()

    # Fetch results from Prometheus
    prometheus_results = fetch_data(start_time, end_time, config)
    write_results_file(config["result_location"], prometheus_results)
    
    logging.info("End 'load()' function")


def wait_to_running_pods():
    logging.info("Start 'wait_to_running_pods()' function")
    #return_param = subprocess.run(cmd, shell=True, universal_newlines=True, capture_output=True)
    
    spinner = Halo(text='Loading', spinner='dots')
    spinner.start("Wait Pods to be ready.")
    
    start_time = time.time()
    warning = True

    ready = False
    while not ready:
        # wait few seconds
        current_time = time.time()

        # If more than 5 minutes stopped the execution
        if current_time - start_time >  3 * 60 and warning:
            logging.warning("Should check pods status manually")
            # Change spinner text to be 
            spinner.text = "CHECK POD STATUS MANUALLY"
            warning = False # Print this meassage only once 
        
        time.sleep(4)

        # get current states by a little shell magic
        # Todo: use kubectl wait
        cmd = "kubectl get pods | tail -n +2 | awk '!/Running/ {print}'"
        return_param = subprocess.run(cmd, shell=True, universal_newlines=True, capture_output=True)

        # Terminating, ContainerCreating, Error, 
        if return_param.stdout == "":
            ready = True

    spinner.stop()

    logging.info("End 'wait_to_running_pods()' function")
    

def fetch_data(start_time, end_time, config):
    logging.info("Start 'fetch_data()' function")
    logging.debug("Time: %s - %s", start_time, end_time)

    # Get parameters to Prometheus query
    cpu_query = {"query": "sum(rate(container_cpu_usage_seconds_total{image!='',namespace=~'default|metrics',container!='POD'}[3m])) by (container)",
                    "start": str(start_time),
                    "end": str(end_time),
                    "step": "1",
                    "timeout": "1000ms"
    }

    memory_query = {"query": "sum(rate(container_memory_usage_bytes{image!='',namespace=~'default|metrics',container!='POD'}[3m])) by (container)",
                    "start": str(start_time),
                    "end": str(end_time),
                    "step": "1",
                    "timeout": "1000ms"
    }

    pods_query = {"query": "sum(kube_pod_container_status_running{namespace=~'default|metrics',container!='POD'}) by (container)",
                    "start": str(start_time),
                    "end": str(end_time),
                    "step": "1",
                    "timeout": "1000ms"
    }

    # Assemble prometheus base query URL
    prometheus_ip = config["prometheus_ip"]
    prometheus_port = config["prometheus_port"]

    url = "http://" + str(prometheus_ip) + ":" + str(prometheus_port) + "/api/v1/query_range?"


    # Get full querys for debug reason
    cpu_full_query = url + urllib.parse.urlencode(cpu_query)
    memory_full_query = url + urllib.parse.urlencode(memory_query)
    pods_full_query = url + urllib.parse.urlencode(pods_query)

    # Print full query
    logging.debug("Full CPU query: %s", cpu_full_query)
    logging.debug("Full Memory query: %s", memory_full_query)
    logging.debug("Full Pods query: %s", pods_full_query)

    # Get results from Prometheus
    cpu_res = json.loads(requests.get(cpu_full_query).text)
    memory_res = json.loads(requests.get(memory_full_query).text)
    pods_res = json.loads(requests.get(pods_full_query).text)

     # Log and return results
    prometheus_results = {"cpu": cpu_res, "memory": memory_res, "pods": pods_res}
    logging.debug("Results from Prometheus: %s", prometheus_results)
    logging.info("End 'fetch_data()' function")

    return prometheus_results


def write_results_file(location, results):
    logging.info("Start 'write_results_file()' function")
    logging.debug("Results file location: %s", location)
    logging.debug("Results to write: %s", results)

    # today = date.today()strftime("%Y-%d-%b")
    current_time = time.strftime("%Y/%b/%d-%H:%M:%S", gmtime())

    filename =  str(current_time) + "-measurement.txt"
    logging.debug("Resultfile: %s", filename)


    with open(filename, 'w') as outfile:
        json.dump(results, outfile)

    logging.info("End 'write_results_file()' function")