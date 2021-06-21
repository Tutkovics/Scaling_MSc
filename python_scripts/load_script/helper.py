import time
import yaml


def read_config_file(file_name):
    """Setup config parameters
    Read and store configuration parameters
    from given config file
    """

    with open(file_name) as config_file:
        config_values = {}
        config = yaml.safe_load(config_file)

        # get values from yaml
        config_values["benchmark_name"] = config["Name"]
        config_values["cluster_name"] = config["Cluster"]["Name"]
        config_values["cluster_ip"] = config["Cluster"]["IP"]

        config_values["servicegraphs"] = config["Servicegraphs"]

        config_values["result_location"] = config["Result_location"]
        
        config_values["qps_from"] = config["QPS"]["From"]
        config_values["qps_to"] = config["QPS"]["To"]
        config_values["qps_granularity"] = config["QPS"]["Granularity"]

        config_values["load_preheat"] = config["Load"]["Preheat"]
        config_values["load_time"] = config["Load"]["Time"]
        config_values["load_users"] = config["Load"]["Users"]
        config_values["load_ip"] = config["Load"]["ServiceIP"]
        config_values["load_timeout"] = config["Load"]["TimeOut"]
        config_values["load_port"] = config["Load"]["ServicePort"]
        config_values["load_path"] = config["Load"]["ServicePath"]
        config_values["load_query"] = config["Load"]["ServiceQuery"]


        config_values["prometheus_ip"] = config["Prometheus"]["IP"]
        config_values["prometheus_port"] = config["Prometheus"]["Port"]
        

        config_values["cluster_ip"] = config["Cluster"]["IP"]
        config_values["cluster_ip"] = config["Cluster"]["IP"]
        
    return config_values
