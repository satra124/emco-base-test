import os
import time
import logging
from requests.exceptions import HTTPError

from json_parser import InputParser
from configuration.env_config import BASE_COMMAND
from utils.utils import emcotestlogging
from utils import utils, models

from report_generator import TestCaseResult



class CriticalFailureException(Exception):
    def __init__(self, message, status_code: int = 500):
        super().__init__(message)
        self.status_code = status_code


class TestExecutor:
    def __init__(self, test_suite_path: str, input_path: str) -> None:
        self.test_suite_path = test_suite_path
        self.input_path = input_path

    def execute_test_case(
        self, section: str, test_suite_file: str, data_file: str, logger: any
    ) -> TestCaseResult:
        """Executes Test case and collects results and other metrics.

        Args:
            section (str): API Section
            test_suite_file (str): Test script name
            data_file (str): Input and data validation file name

        Returns:
            TestCaseResult
        """

        start_time = time.time()

        script = os.path.join(self.test_suite_path, section, test_suite_file)
        rest_method = test_suite_file.split(".")[0].split("_")[1][:4]
        input_path_api = os.path.join(self.input_path, section, data_file)
        base_command = f"{BASE_COMMAND} {script}"

        command = utils.create_command(
            base_command=base_command, options={"input-path-api": input_path_api}
        )
        logger.info(f"Command: {command}")

        cmd_output = utils.run_command(command)

        logger.info(f"Command Executed")

        critical_failure = False
        if cmd_output.return_code == 0:
            if data_file == 'test_instantiateDeploymentBeforeApprove.json' or data_file == 'test_terminateDeploymentBeforeInstantiate.json':
                logger.info(f"Test case: {data_file.split('.')[0]} failed.")
                result = "FAIL"
                li = cmd_output.output.split("\n")

                reason = f"{li[6]} \n {li[8]} \n {li[9]} \n {li[11]}"
                logger.info(f"Reason: {reason}")
            else:
                logger.info(f"Test case: {data_file.split('.')[0]} passed.")
                result = "PASS"
                reason = ""

        else:
            if data_file == 'test_instantiateDeploymentBeforeApprove.json' or data_file == 'test_terminateDeploymentBeforeInstantiate.json':
                logger.info(f"Test case: {data_file.split('.')[0]} passed.")
                result = "PASS"
                reason = ""
            else:
                logger.info(f"Test case: {data_file.split('.')[0]} failed.")
                result = "FAIL"
                li = cmd_output.output.split("\n")

                reason = f"{li[6]} \n {li[8]} \n {li[9]} \n {li[11]}"
                logger.info(f"Reason: {reason}")

                if rest_method.lower() == "post" and data_file.split('/')[0] != 'negative_error_codes':
                    critical_failure = True
                if test_suite_file == 'test_postClusterProviderNeg422.py':
                    critical_failure = False

        end_time = time.time()
        duration = f"{str(round(end_time-start_time, 2))} s"
        test_case_result = TestCaseResult(
            section=section,
            test_case=data_file.split(".")[0],
            result=result,
            reason=reason,
            duration=duration,
        )

        return test_case_result, critical_failure
