import os

from dotenv import load_dotenv
from pathlib import Path

dotenv_path = Path('configuration/.env')

load_dotenv(dotenv_path=dotenv_path)


HOST = os.getenv('HOST')

orchestrator_port= os.getenv('ORCHESTRATOR_PORT')
clm_port= os.getenv('CLM_PORT')
dcm_port= os.getenv('DCM_PORT')
dtc_port= os.getenv('DTC_PORT')
hpa_port= os.getenv('HPA_PORT')
ncm_port = os.getenv('NCM_PORT')

orchestrator=f"{HOST}:{orchestrator_port}"
clm=f"{HOST}:{clm_port}"
dcm=f"{HOST}:{dcm_port}"
dtc=f"{HOST}:{dtc_port}"
hpa=f"{HOST}:{hpa_port}"
ncm=f"{HOST}:{ncm_port}"

INPUT_BASE = os.path.join('inputs', 'data')
TEST_SUITE_BASE = os.path.join('test_suite')
LOG_PATH = os.path.join('outputs', 'logs')
LOG_LEVEL = 'DEBUG'
HELM_CHART_FOLDER_PATH = os.getcwd() + os.path.sep + 'tgz_files'

BASE_COMMAND = 'python -m pytest --disable-pytest-warnings -s'
COMMAND_EXEC_TIMEOUT = 20
TEST_REPORT_FILE_PREFIX = 'test_report_'



