import readline, glob, sys
import logging
#from yaml import full_load, load, FullLoader
from helper import read_config_file
from load import orchestrate_measure

def main(config_file):
    """Main function.
    This will orchestrate the measurement
    """

    # read the given config file and store values
    logging.debug("Read config file: %s", config_file)
    config = read_config_file(config_file)
    logging.debug("Config values: %s", config)
    
    orchestrate_measure(config)
    



if __name__ == '__main__':
    try:
        # Setup Logging parameters
        logging.basicConfig(format='%(asctime)s - %(levelname)s : %(message)s', level=logging.DEBUG)
        logging.info("Logger setup finished")

        # check if all ready to launch
        if len(sys.argv) != 2:
            logging.info("Wrong number of argument was given")
            # if config file is not given
            print("Please give (one and only one) configuration yaml file")

            # get config file (on runtime) with autocomplete  
            def complete(text, state):
                return (glob.glob(text+'*')+[None])[state]

            readline.set_completer_delims(' \t\n;')
            readline.parse_and_bind("tab: complete")
            readline.set_completer(complete)

            config_file = input("Config file: ")

            main(config_file)
        else:
            main(str(sys.argv[1]))
    finally:
        logging.info("Script exited")