from argparse import ArgumentParser
from collections import defaultdict
from matplotlib import pyplot as plt
from datetime import datetime, timezone
import numpy as np
# import os
import glob
import json

def readCommandLineParameters():
    parser = ArgumentParser()
    parser.add_argument("-d", "--directory", dest="directory",
        help="Directory where json files stored.", metavar="DIR",
        required=True)
    parser.add_argument("-t", "--title", dest="title",
        help="Title of the graph.", metavar="TITLE",
        default="Scaling Thesis Plot")
    parser.add_argument("-x", "--x-axis", dest="x",
        help="X axis label.", metavar="QPS", default="")
    parser.add_argument("-y", "--y-axis", dest="y",
        help="Y axis label.", metavar="mCPU", default="")
    
    parser.add_argument("-s", "--save", dest="save",
        help="Save the plotting.", nargs='?', const=True, default=False)
          
    return parser.parse_args()

def getFilesToRead(directory):
    extension = 'json'
    files = glob.glob(directory + '/*.' + extension)

    # sort by name
    files = sorted(files)

    return files

def readFileAndGetValues(fileName):
    # read file
    with open(fileName, 'r') as jsonFile:
        data=jsonFile.read()

    # parse file
    obj = json.loads(data)

    # Values should be useful: requested QPS, actual QPS, sum fe & sum be usage

    runningPods = getRunningPodsFromJson(obj)
    runningPodsByContainer = getRunningPodsByContainer(obj)

    cpu = sumCpuByRunningContainer(obj, runningPods)
    memory = sumMemoryByRunningContainer(obj, runningPods)

    measurementDuration = obj["config"]["load_time"]

    avgResponseTime = obj["fortio"]["DurationHistogram"]["Avg"]
    
    values = {"reqQPS" : obj["fortio"]["RequestedQPS"],
              "actQPS" : obj["fortio"]["ActualQPS"],
              "runningPods" : len(runningPods),
              "runningPodsByContainer": runningPodsByContainer,
              "cpu" : cpu,
              "memory" : memory,
              "measurementDuration": measurementDuration,
              "avgResponseTime": avgResponseTime,
              }

    print(values)
    return values

def getRunningPodsFromJson(obj):
    pods = []

    for result in obj["prometheus"]["running_pods"]["data"]["result"]:
        pods.append(result["metric"]["pod"])
    
    return pods

def getRunningPodsByContainer(obj):
    pods = {}

    for result in obj["prometheus"]["number_of_pods"]["data"]["result"]:
        pods[result["metric"]["container"]] = result["values"][-1][1]

    return pods

def sumCpuByRunningContainer(obj, runningPods):
    cpus = defaultdict(float)

    for result in obj["prometheus"]["cpu"]["data"]["result"]:
        # if the current pod name was in running pods
        if result["metric"]["pod"] in runningPods:
            # Check if the cpu usage raised during reported period
            diff = float(result["values"][-1][1]) - float(result["values"][0][1])
            if diff != 0:
                cpus[result["metric"]["container"]] += diff
    # print(dict(cpus))
    return dict(cpus)

def sumMemoryByRunningContainer(obj, runningPods):
    memory = defaultdict(float)

    for result in obj["prometheus"]["memory"]["data"]["result"]:
        # if the current pod name was in running pods
        if result["metric"]["pod"] in runningPods:
            # Check if the memory usage raised during reported period
            diff = float(result["values"][-1][1]) - float(result["values"][0][1])
            if diff != 0:
                memory[result["metric"]["container"]] += diff

    return dict(memory)

def draw(datas, args):
    x = []      # requested QPS
    x_act = []  # actual QPS
    y = {}      # cpu usage (per container)
    cpu_sum = [] # Summary of cpu usage
    memory = {} # Memory usage (per container)
    memory_sum = [] # Summary of memory usage 
    y2 = []     # average response time

    # Init cpu and memory dict
    for container in datas[0]["runningPodsByContainer"]:
        y[container] = []
        memory[container] = []


    last_data = None
    
    end = 0 # for X axis ticks

    for data in datas:
        reqQps = data["reqQPS"]
        x.append(reqQps)  # add qps to x axis
        x_act.append(data["actQPS"])

        # update latest qps
        if end < int(reqQps):
                end = int(reqQps)

        # Counter of cpu and memory usage
        tmp_sum_cpu = 0
        tmp_sum_memory = 0

        for container in data["runningPodsByContainer"]:
            # Calculate resource usage per container
            container_count = int(data["runningPodsByContainer"][container])

            # "Solve" zero division count
            if container_count == 0:
                container_count = 1
            
            try:
                calculated_cpu = data["cpu"][container] / (container_count * data["measurementDuration"])
                calculated_memory = data["memory"][container] / (container_count * data["measurementDuration"])
            except KeyError:
                # If something went wrong in Prometheus collecting, then use the latest data
                calculated_cpu = 1    # y[container][-1]
                calculated_memory = 1 # memory[container][-1]
                print("KeyError: " + str(container) + str(reqQps))
                

            
            # Add container resource usage
            y[container].append(calculated_cpu)
            memory[container].append(calculated_memory)

            # Increment counters
            tmp_sum_cpu += calculated_cpu
            tmp_sum_memory += calculated_memory

        cpu_sum.append(tmp_sum_cpu)
        memory_sum.append(tmp_sum_memory)

        # add response time
        y2.append(data["avgResponseTime"])

    
    # subplot(nrows, ncols, index, **kwargs)
    fig, axs = plt.subplots(4, 2)
    plt.title("asd") #args.title)  # Set title
    
    # Left upper
    axs[0, 0].plot(x, cpu_sum)
    min_x, max_x = axs[0, 0].get_xlim()
    axs[0, 0].set_xticks(np.arange(min_x, max_x+1, 5), minor=False)
    # axs[0, 0].set_xticks(np.arange(min(x), max(x)), minor=False)
    axs[0, 0].set_title("Sum CPU (mCPU) / requested QPS")
    

    # Right upper
    axs[0, 1].plot(x, memory_sum)
    min_x, max_x = axs[0, 1].get_xlim()
    axs[0, 1].set_xticks(np.arange(min_x, max_x+1, 5), minor=False)
    axs[0, 1].set_title("Sum Memory / requested QPS")

    # CPU usage per container
    for container in y:
        axs[1, 0].plot(x, y[container], label=str(container) + " - CPU usage")
    axs[1, 0].set_title("Container CPU usage / requested QPS")
    min_x, max_x = axs[1, 0].get_xlim()
    axs[1, 0].set_xticks(np.arange(min_x, max_x+1, 5), minor=False)
    #axs[1, 0].set_xticks(np.arange(0, end, 5), minor=False)

    for container in memory:
        axs[1, 1].plot(x, memory[container], label=str(container) + " - Memory usage")
    axs[1, 1].set_title("Container Memory usage / requested QPS")
    min_x, max_x = axs[1, 1].get_xlim()
    axs[1, 1].set_xticks(np.arange(min_x, max_x+1, 5), minor=False)
    #axs[1, 1].set_xticks(np.arange(0, end, 5), minor=False)

    # CPU usage per container / X axis: actual QPS
    for container in y:
        axs[2, 0].scatter(x_act, y[container], label=str(container) + " - CPU usage")
    axs[2, 0].set_title("Container CPU usage / actual QPS")
    #axs[1, 0].set_xticks(np.arange(0, end, 5), minor=False)

    for container in memory:
        axs[2, 1].plot(x_act, memory[container], label=str(container) + " - Memory usage")
    axs[2, 1].set_title("Container Memory usage / actual QPS")
    #axs[1, 1].set_xticks(np.arange(0, end, 5), minor=False)

    # Bottom left
    axs[3, 0].scatter(x_act, y2)
    axs[3, 0].set_title("Response time / actual QPS")
    
    # Bottom right
    axs[3, 1].plot(x, x_act)
    axs[3, 1].set_title("Actual vs requested")
    min_x, max_x = axs[3, 1].get_xlim()
    axs[3, 1].set_xticks(np.arange(min_x, max_x+1, 5), minor=False)

    # response time
    # axs[1, 0].plot(x, y2)
    # axs[1, 0].set_title("Response time")
    # axs[1, 0].sharex(axs[0, 0])

    fig.tight_layout()
    print("plotting...")
    #ax.set_xlabel(args.x)
    #ax.set_ylabel(args.y)

    # Configure X axis ticks
    #ax.xaxis.set_ticks(np.arange(0, end, 5))

    plt.show()

    # Save the figure if command line flag was given
    if args.save:
        plot_location = "../results/plots/"
        plot_name = str(args.title) + "_" + str(datetime.now(timezone.utc).isoformat(timespec='minutes')) + ".jpg"
        fig.savefig(plot_location + plot_name)
    

def main():
    args = readCommandLineParameters()
    
    files = getFilesToRead(args.directory)

    clean_data = []
    
    # Foreach file in directory
    for file in files:
        try:
            values = readFileAndGetValues(file)
            clean_data.append(values)
        except KeyError: 
            print("Can't read file:", file) 

    draw(clean_data, args)

main()