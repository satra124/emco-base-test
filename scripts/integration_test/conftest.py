import pytest
import datetime
import os
import time
import logging.config
import shutil
from utils import log_manager

from configuration.env_config import LOG_PATH, INPUT_BASE, orchestrator, clm, dcm, dtc, hpa, ncm, HELM_CHART_FOLDER_PATH

def pytest_addoption(parser):
    parser.addoption(
        '--polling-timeout',
        action="store",
        default=55,
        help='This is polling function timeout')
    parser.addoption(
        '--sleep-before-get-app-status',
        action="store",
        default=15,
        help='this is sleep time before get app state')
    parser.addoption(
        '--sleep-before-get-app-status-for-stop',
        action="store",
        default=0,
        help='this is sleep time before get app state for stop case')
    parser.addoption(
        '--sleep-for-logical-cloud-instantiation',
        action="store",
        default=4,
        help='this is time to allow logical cloud instantiation')
    parser.addoption(
        '--is_auth_enabled',
        action="store",
        default=False,
        help='flag to be set to True if authenication flag is enabled')
    parser.addoption(
        '--is_token_to_be_generated',
        action="store",
        default=False,
        help='flag to be set to True if token is to be generated')
    parser.addoption(
        '--input-path',
        action="store",
        default="inputs/system_test_app_deployment_TC1/",
        help='path to the input json files in the format "inputs/path_to_input_json/"')
    parser.addoption(
        '--op-type',
        action="store",
        default=None,
        help='operation type: post/delete/get/put"')
    parser.addoption(
        '--input-path-api',
        action="store",
        default=None,
        help='path to the project input json files in the format "inputs/path_to_input_json/file.json"')



@pytest.fixture
def api_input(request):
    return request.config.getoption("--input-path-api")

@pytest.fixture
def sleep_for_logical_cloud_instantiation(request):
    return request.config.getoption("--sleep-for-logical-cloud-instantiation")

@pytest.fixture
def operation_type(request):
    return request.config.getoption("--op-type")

@pytest.fixture
def polling_timeout(request):
    return request.config.getoption("--polling-timeout")

@pytest.fixture
def sleep_before_get_app_status(request):
    return request.config.getoption("--sleep-before-get-app-status")

@pytest.fixture
def sleep_before_get_app_status_for_stop(request):
    return request.config.getoption("--sleep-before-get-app-status-for-stop")

@pytest.fixture
def input_file_path_systemtestcase(request):
    input = request.config.getoption("--input-path")
    return input

@pytest.fixture
def cwd():
    return os.getcwd() + os.path.sep

@pytest.fixture
def get_auth_value(request):

    if request.config.getoption("--is_auth_enabled"):
       print("Auth set to true")
       return True
    else:
       print("Auth is set to false")
       return False

@pytest.fixture
def token_to_be_generated(request):

    if request.config.getoption("--is_token_to_be_generated"):
       print("Token is to be generated")
       return True
    else:
       print("Token is not to be generated")
       return False

# @pytest.fixture
# def helm_chart_folder_path():
#     return HELM_CHART_FOLDER_PATH

@pytest.fixture
def orchestrator_url():
    EMCO_URL = f"http://{orchestrator}/v2/"
    return EMCO_URL

@pytest.fixture
def dtc_url():
    EMCO_URL = f"http://{dtc}/v2/"
    return EMCO_URL

@pytest.fixture
def clm_url():
    EMCO_URL = f"http://{clm}/v2/"
    return EMCO_URL

@pytest.fixture
def dcm_url():
    EMCO_URL = f"http://{dcm}/v2/"
    return EMCO_URL

@pytest.fixture
def hpa_url():
    EMCO_URL = f"http://{hpa}/v2/"
    return EMCO_URL

# @pytest.fixture
# def url():
#     # TODO: Check if the IP is an actual IP address using regex
#     with open('EMCO_SERVER_IP_PORT_dcm.txt', 'r') as f:
#         server_ip_port = f.read().strip("\n")
#     EMCO_URL = f"http://{server_ip_port}/v2"
#     return EMCO_URL


@pytest.fixture
def ncm_url():
    EMCO_URL = f"http://{ncm}/v2/"
    return EMCO_URL

@pytest.fixture(scope="session")
def emcotestlogging():
    cwd = os.getcwd()
    TESTRESULT_ROOT = "TestResults"
    TESTRESULT_FOLDER_NAME = "TestResults-" + str(time.strftime("%Y-%m-%d-%H%M%S"))
    test_result_path = os.path.join(cwd, TESTRESULT_ROOT, TESTRESULT_FOLDER_NAME)
    if os.path.exists(test_result_path):
        shutil.rmtree(test_result_path)

    os.makedirs(test_result_path)
    log_name = "EMCO_Test" + ".log"
    log_file_name = repr(os.path.join(test_result_path, log_name)).strip("'")
    logging.config.fileConfig('logging.conf', disable_existing_loggers=False,
                              defaults={'logfilename': log_file_name})
    return log_manager.LogManager(log_file_name, "EMCO_TEST")



