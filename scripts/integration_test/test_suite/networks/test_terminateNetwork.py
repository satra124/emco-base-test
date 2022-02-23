import json
import requests
import os
from requests.exceptions import HTTPError
from json_parser import InputParser
from utils import utils


class TestCase_terminateNetwork:
    def test_terminatetNetwork(self, cwd, ncm_url, emcotestlogging, api_input):
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
                    response = self.post_terminate_network(
                        ncm_url, params_from_inputfile.get_anchor(0)
                    )
                    logfile.log_debug(
                        "Received Response Code: " + str(response.status_code)
                    )
                    self.post_terminate_network_assertion_check(
                        logfile,
                        response,
                        params_from_inputfile.get_response_code(0, item),
                    )

            except HTTPError as http_err:
                logfile.log_debug(utils.printException())
                logfile.log_debug(
                    "Error - HTTP error occurred: '{0}' ".format(http_err)
                )
                raise
            except Exception as err:
                logfile.log_debug(utils.printException())
                logfile.log_debug("Error - other error occurred: '{0}' ".format(err))
                raise

    @classmethod
    def post_terminate_network(self, ncm_url, anchor):
        headers = {"Content-Type": "application/json", "accept": "application/json"}
        payload = request_body
        response = requests.post(
            ncm_url + anchor, data=json.dumps(payload), headers=headers
        )
        return response

    @classmethod
    def post_terminate_network_assertion_check(
        self, logfile, response, responseCode=None
    ):
        msgResCode = "Response Code does not match, expected: {0}, actual: {1}".format(
            responseCode, response.status_code
        )
        utils.logAssert(response.status_code == responseCode, msgResCode, logfile)
