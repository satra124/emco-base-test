import json
import requests
import os
from requests.exceptions import HTTPError
from json_parser import InputParser
from utils import utils


class TestCase_statusDeployment:
    def test_statusDeployment(self, cwd, orchestrator_url, emcotestlogging, api_input):
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
                    response = self.get_statusDeployment(
                        orchestrator_url, params_from_inputfile.get_anchor(0)
                    )
                    logfile.log_debug(
                        "Received Response Code: " + str(response.status_code)
                    )
                    logfile.log_debug("Received Response Text: " + str(response.text))
                    self.get_status_deployment_assertion_check(
                        logfile,
                        response,
                        params_from_inputfile.get_response_code(0, item),
                        params_from_inputfile.get_response_text(0, item),
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
    def get_statusDeployment(self, orchestrator_url, anchor):
        headers = {"Content-Type": "application/json", "accept": "application/json"}
        response = requests.get(orchestrator_url + anchor, headers=headers)
        return response

    @classmethod
    def get_status_deployment_assertion_check(
        self, logfile, response, responseCode=None, responseText=None
    ):
        msgResCode = "Response Code does not match, expected: {0}, actual: {1}".format(
            responseCode, response.status_code
        )
        utils.logAssert(response.status_code == responseCode, msgResCode, logfile)
        if responseText:
            if type(responseText) == str:
                msgText = "Response Text does not contain expected string, expected: {0}, actual: {1}".format(
                    responseText, response.text
                )
                if responseCode == 200:
                    if "status" in response.text:
                        utils.logAssert(
                            responseText == json.loads(response.text)["status"],
                            msgText,
                            logfile,
                        )
                    else:
                        raise Exception("status key not found")
                else:
                    if responseText not in response.text:
                        utils.logAssert(True == False, msgText, logfile)
            else:
                msgText = "Response Text does not contain expected string, expected: {0}, actual: {1}".format(
                    json.dumps(responseText), response.text
                )
                # comparing dictionaries
                utils.logAssert(
                    json.loads(json.dumps(responseText)) == json.loads(response.text),
                    msgText,
                    logfile,
                )
