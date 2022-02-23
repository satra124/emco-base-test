import os
import datetime
import time
import shutil
import sys
import subprocess
from subprocess import CalledProcessError
import linecache
import logging
from utils import log_manager
from pathlib import Path

from typing import Dict, Optional

from configuration.env_config import LOG_PATH, COMMAND_EXEC_TIMEOUT, HELM_CHART_FOLDER_PATH, LOG_LEVEL
from utils.models import CommandOutput

helm_chart_folder_path = HELM_CHART_FOLDER_PATH


def setup_framework_logging():

    logger = logging.getLogger('emco_test_automation')
    logger.setLevel(logging.DEBUG)    


    # create console handler and set level to info
    stream_handler = logging.StreamHandler()
    stream_handler.setLevel(logging.DEBUG)

    # create file handler and set level to info
    try:
        Path(LOG_PATH).mkdir(parents=True, exist_ok=True)
    except Exception as exc:
        print(exc)
        print("File logs will not be generated")
    
    file_handler = logging.FileHandler(filename=os.path.join(LOG_PATH, f'emco_test_logs_{datetime.datetime.now()}.log'))
    file_handler.setLevel(logging.DEBUG)

    # create formatter
    formatter = logging.Formatter("[%(asctime)s] - {%(name)s.%(funcName)s:%(lineno)d} - %(levelname)s ""- %(message)s",)

    # add formatters to our handlers
    stream_handler.setFormatter(formatter)
    file_handler.setFormatter(formatter)

    # add Handlers to our logger
    logger.addHandler(stream_handler)
    logger.addHandler(file_handler)

    return logger

def emcotestlogging():
    log_path = LOG_PATH
    TESTRESULT_ROOT = "TestResults"
    TESTRESULT_FOLDER_NAME = "TestResults-" + str(time.strftime("%Y-%m-%d-%H%M%S"))
    test_result_path = os.path.join(log_path, TESTRESULT_ROOT, TESTRESULT_FOLDER_NAME)
    if os.path.exists(test_result_path):
        shutil.rmtree(test_result_path)

    os.makedirs(test_result_path)
    log_name = "EMCO_Test" + ".log"
    log_file_name = repr(os.path.join(test_result_path, log_name)).strip("'")
    logging.config.fileConfig('logging.conf', disable_existing_loggers=False,
                              defaults={'logfilename': log_file_name})
    return log_manager.LogManager(log_file_name, "EMCO_TEST")


def logAssert(test, msg, logfile):
    if not test:
        logfile.log_debug(msg)
        assert test, msg
        return False
    return True


def printException():
    exc_type, exc_obj, tb = sys.exc_info()
    f = tb.tb_frame
    lineno = tb.tb_lineno
    filename = f.f_code.co_filename
    linecache.checkcache(filename)
    line = linecache.getline(filename, lineno, f.f_globals)
    print("Exception in " + str(filename), " ," + "line " + str(lineno) + " " + str(line.strip()) + ": " + str(exc_obj))


def add_command_options(options: Optional[Dict[str, Optional[str]]] = None) -> str:
    """
    Creates the command line args required for helm command
    Parameters
    ----------
        options: A dictionary of command line arguments to pass to helm
    Returns
    -------
        Command line args required for helm command
    """
    option_string = ""
    if options is None:
        return option_string
    for option in options:
        option_string += f" --{option}"
        # Add value after flag if one is given
        value = options[option]
        if value:
            option_string += f" {value}"
    return option_string

def create_command(base_command: str, options: Optional[Dict[str, Optional[str]]] = None) -> str:
    """
    Creates the command
    Parameters
    ----------
        base_command: Base command
        options: A dictionary of command line arguments to pass to helm
    Returns
    -------
        Command with added namespace and command line args required for helm command
    """

    command = base_command
    command = command + add_command_options(options)
    return command


def run_command(command: str) -> CommandOutput:
    """
        Runs the command using subprocess
        Parameters
        ----------
            command: Command which need to be run
        Returns
        -------
            Output of the given command
    """
    return_code = ""
    output = ""
    
    try:
        process = subprocess.run(command, shell= True, capture_output=True, timeout=COMMAND_EXEC_TIMEOUT)
        output = process.stdout.decode('utf-8')
        return_code = process.returncode
    except Exception as exc:
        logger.error(f"ERROR IN RUN COMMAND: {exc}")
        return_code = 1
        output = str(exc)
    

    # except CalledProcessError as err:
    #     # message = err.output.decode("utf-8")
    #     raise KubectlUtilsException(err.output.decode("utf-8"))
    return CommandOutput(return_code = return_code, 
                            output = output)