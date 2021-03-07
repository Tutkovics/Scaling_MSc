import logging

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
            
            # todo: start servicegraph
            load(config["load_preheat"], config["load_time"])


    logging.info("End 'orchestrate_measure()' function")

def kubectl(command, argument):
    # eg: kubectl delete, apply,  
    pass

def load(preheat_time, load_time):
    logging.info("Start 'load()' function")

    logging.info("End 'load()' function")

def wait_to_running_pods():
    pass

def fetch_data(prometheus_ip, prometheus_port):
    pass

def write_results_file(location, results):
    pass