import logging
import subprocess
# import os
# from spinners import Spinners
import time
from halo import Halo

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
            load(config["load_preheat"], config["load_time"])

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


def load(preheat_time, load_time):
    logging.info("Start 'load()' function")

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
        cmd = "kubectl get pods | tail -n +2 | awk '!/Running/ {print}'"
        return_param = subprocess.run(cmd, shell=True, universal_newlines=True, capture_output=True)

        # Terminating, ContainerCreating, Error, 
        if return_param.stdout == "":
            ready = True

    spinner.stop()

    logging.info("End 'wait_to_running_pods()' function")
    

def fetch_data(prometheus_ip, prometheus_port):
    pass
    # logging.info("Start 'load()' function")

    # logging.info("End 'load()' function")

def write_results_file(location, results):
    pass
    # logging.info("Start 'load()' function")

    # logging.info("End 'load()' function")