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

    measurementDuration = obj["config"]["load_time"]

    avgResponseTime = obj["fortio"]["DurationHistogram"]["Avg"]
    
    values = {"reqQPS" : obj["fortio"]["RequestedQPS"],
              "actQPS" : obj["fortio"]["ActualQPS"],
              "runningPods" : len(runningPods),
              "runningPodsByContainer": runningPodsByContainer,
              "cpu" : cpu,
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

def draw(datas, args):
    x = []  # qps
    x_act = [] # actual QPS
    y = {}  # cpu usage
    y2 = [] # average response time

    # Init cpu dict
    for container in datas[0]["runningPodsByContainer"]:
        y[container] = []


    last_data = None
    
    end = 0 # for X axis ticks

    for data in datas:
        reqQps = data["reqQPS"]
        x.append(reqQps)  # add qps to x axis
        x_act.append(data["actQPS"])

        # update latest qps
        if end < int(reqQps):
                end = int(reqQps)

        for container in data["runningPodsByContainer"]:
            calculated_cpu = data["cpu"][container] / (int(data["runningPodsByContainer"][container]) * data["measurementDuration"])
            y[container].append(calculated_cpu)

        # add response time
        y2.append(data["avgResponseTime"])

    
    fig, axs = plt.subplots(2, 2)
    plt.title("asd") #args.title)  # Set title

    for container in y:
        axs[0, 0].plot(x,y[container], label=str(container) + " - CPU usage")
    axs[0, 0].set_title("Req QPS / mCPU")
    axs[0,0].set_xticks(np.arange(0, end, 5), minor=False)
    #axs[0,0].xticks(np.arange(0, end, 5))

    # response time
    axs[1, 0].plot(x, y2)
    axs[1, 0].set_title("Response time")
    axs[1, 0].sharex(axs[0, 0])
    
    axs[0, 1].scatter(x_act, y2)
    axs[0, 1].set_title("Actual QPS")

    axs[1, 1].plot(x, x_act)
    axs[1, 1].set_title("Actual vs requested")
    fig.tight_layout()
    print("plotting...")
    #ax.set_xlabel(args.x)
    #ax.set_ylabel(args.y)

    # Configure X axis ticks
    #ax.xaxis.set_ticks(np.arange(0, end, 5))



    # plt.ylim(-0.02, 0.5)

    #ax.plot(x,y2, label="Average Response Time")


    #ax.legend() 
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