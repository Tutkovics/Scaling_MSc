(python_venv) $ python main.py --help
usage: main.py [-h] -d DIR [-t TITLE] [-x QPS] [-y mCPU] [-s [SAVE]]

optional arguments:
  -h, --help            show this help message and exit
  -d DIR, --directory DIR
                        Directory where json files stored.
  -t TITLE, --title TITLE
                        Title of the graph.
  -x QPS, --x-axis QPS  X axis label.
  -y mCPU, --y-axis mCPU
                        Y axis label.
  -s [SAVE], --save [SAVE]
                        Save the plotting.
(python_venv) $ python main.py --directory ./path/to/results/directory/ --title "Plot title" -x QPS -y mCPU --save
