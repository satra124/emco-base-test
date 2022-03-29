import os
import time
import datetime


from deps import TestExecutor
from report_generator import ReportGenerator
from configuration.env_config import INPUT_BASE, TEST_SUITE_BASE
from utils.models import AutomationType
from utils.utils import setup_framework_logging

from test_list import full_api_test_list, cleanup_test_list


full_api_automation = AutomationType(type='Complete API Automation', description='In this type of Automation all the available EMCO APIs are tested.', total_test_cases=len(full_api_test_list))
cleanup_api_automation = AutomationType(type='Cleanup API Automation', description='In this type of Automation all the data is removed and terminated', total_test_cases=len(cleanup_test_list))


# sets up logging configuration
logger = setup_framework_logging()


class AutomationFacadeException(Exception):
    """Automation Facade Exception class"""

    def __init__(self, message):
        super().__init__(message)


class AutomationFacade:
    """Class to import and use the backend services like Test Executor and Report Generator
    """

    def full_api_automation(self):
        """Executes all the API test cases.
        """
        logger.info("Application started")

        execution_start_date = datetime.datetime.now().strftime("%d/%m/%Y, %H:%M:%S")
        complete_execution_start_time = time.time()

        test_type = "full_api_test"

        INPUT_PATH = os.path.join(INPUT_BASE, test_type)

        try:
            test_executor = TestExecutor(
                test_suite_path=TEST_SUITE_BASE, input_path=INPUT_PATH
            )
            report_generator = ReportGenerator()
        except Exception as exc:
            logger.error(f"Error in creating test_executor and report_generator: {exc}")
        test_case_results = []
        logger.info("Test Executor and Report Generator initialized")

        for test_case in full_api_test_list:
            try:
                logger.info(f"Tuple INFO: {test_case}")

                test_case_result, critical_failure = test_executor.execute_test_case(
                    section=test_case[0],
                    test_suite_file=test_case[1],
                    data_file=test_case[2],
                    logger = logger
                )
                test_case_results.append(test_case_result)
                if critical_failure:
                    logger.error(
                        f"Critical Failure occured in : {test_case[2]}, Stopping Automation."
                    )
                    break

            except Exception as exc:
                logger.error(f"Exception occured: {exc}")

        complete_execution_end_time = time.time()

        total_duration = f"{str(round(complete_execution_end_time-complete_execution_start_time, 2))} s"

        # report_generator
        logger.info(f"Report Generation started.")
        try:
            report_path = report_generator.generate_report(test_case_results=test_case_results, automation_type=full_api_automation, total_duration = total_duration, execution_start_date = execution_start_date)
            logger.info(f"Report Generation Completed. Report saved at: {report_path}")
        except Exception as exc:
            logger.error(f"Report Generation failed: {exc}")

    def delete_api_automation(self):

        test_type = "full_api_test"

        execution_start_date = datetime.datetime.now().strftime("%d/%m/%Y, %H:%M:%S")
        complete_execution_start_time = time.time()

        INPUT_PATH = os.path.join(INPUT_BASE, test_type)

        test_executor = TestExecutor(
            test_suite_path=TEST_SUITE_BASE, input_path=INPUT_PATH
        )
        report_generator = ReportGenerator()
        test_case_results = []
        logger.info("Test Executor and Report Generator initialized")

        for test_case in cleanup_test_list:
            try:
                logger.info(f"Tuple INFO: {test_case}")

                test_case_result, critical_failure = test_executor.execute_test_case(
                    test_case[0], test_case[1], test_case[2], logger
                )
                test_case_results.append(test_case_result)

            except Exception as exc:
                logger.error(f"Exception occured: {exc}")

        complete_execution_end_time = time.time()

        total_duration = f"{str(round(complete_execution_end_time-complete_execution_start_time, 2))} s"

        # report_generator

        logger.info(f"Report Generation started.")
        try:
            report_path = report_generator.generate_report(test_case_results=test_case_results, automation_type=cleanup_api_automation, total_duration = total_duration, execution_start_date = execution_start_date)

            logger.info(f"Report Generation Completed. Report saved at: {report_path}")
        except Exception as exc:
            logger.error(f"Report Generation failed: {exc}")

