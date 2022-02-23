import json
import requests
import os
from requests.exceptions import HTTPError
from json_parser import InputParser
from utils import utils


class TestCase_approveDeployment:
    def test_approveDeployment(self, cwd, orchestrator_url, emcotestlogging, api_input):
        try:
            logfile = emcotestlogging
            logfile.log_debug("=" * 15)
            INPUTFILE_PATH = cwd + api_input
            params_from_inputfile = InputParser(INPUTFILE_PATH)
            no_of_api_calls = params_from_inputfile.get_type_metadata_length(0)
        except Exception as err:
            logfile.log_debug(utils.printException())
            logfile.log_debug(
                "Error - other error occurred before API operation: '{0}' ".format(err)
            )
            raise
        else:
            try:
                for item in range(0, no_of_api_calls):
                    response = self.post_approveDeployment(
                        orchestrator_url, params_from_inputfile.get_anchor(0)
                    )
                    logfile.log_debug(
                        "Received Response Code: " + str(response.status_code)
                    )
                    logfile.log_debug("Received Response Text: " + str(response.text))
                    self.post_approve_deployment_assertion_check(
                        logfile,
                        response,
                        params_from_inputfile.get_response_code(0, item),
                    )
            except HTTPError as http_err:
                logfile.log_debug(utils.printException())
                logfile.log_debug(
                    "Error - HTTP error occurred during api call item '{0}': '{1}' ".format(
                        item, http_err
                    )
                )
                raise
            except Exception as err:
                logfile.log_debug(utils.printException())
                logfile.log_debug(
                    "Error - other error occurred during api call item '{0}': '{1}' ".format(
                        item, err
                    )
                )
                raise

    @classmethod
    def post_approveDeployment(self, orchestrator_url, anchor):
        headers = {"Content-Type": "application/json", "accept": "application/json"}
        response = requests.post(orchestrator_url + anchor, headers=headers)
        return response

    @classmethod
    def post_approve_deployment_assertion_check(
        self, logfile, response, responseCode=None
    ):
        msgResCode = "Response Code does not match, expected: {0}, actual: {1}".format(
            responseCode, response.status_code
        )
        accepted_response_codes = [int(x) for x in responseCode.split(",")]
        if response.status_code in accepted_response_codes:
            utils.logAssert(1 == 1, msgResCode, logfile)
        else:
            utils.logAssert(1 == 0, msgResCode, logfile)
