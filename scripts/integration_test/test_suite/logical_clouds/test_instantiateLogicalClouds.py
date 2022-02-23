import json
import requests
import os
import time
from requests.exceptions import HTTPError
from json_parser import InputParser
from utils import utils


class TestCase_instantiateLogicalClouds:
    def test_instantiateLogicalClouds(
        self,
        cwd,
        dcm_url,
        emcotestlogging,
        api_input,
        sleep_for_logical_cloud_instantiation,
    ):
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
                    response = self.instantiate_logical_clouds(
                        dcm_url, params_from_inputfile.get_anchor(0)
                    )
                    logfile.log_debug(
                        "Received Response Code: " + str(response.status_code)
                    )
                    logfile.log_debug("Received Response text: " + str(response.text))
                    self.instantiate_logical_clouds_assertion_check(
                        logfile,
                        response,
                        params_from_inputfile.get_response_code(0, item),
                    )
                    # Adding a sleep time for logical cloud instantiation to enable Automation runs. The sleep time is same as that in system-test.
                    logfile.log_debug(
                        "Sleeping for {0}s to allow logical cloud instantiation...".format(
                            sleep_for_logical_cloud_instantiation
                        )
                    )
                    time.sleep(int(sleep_for_logical_cloud_instantiation))
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
    def instantiate_logical_clouds(self, dcm_url, anchor):
        headers = {"Content-Type": "application/json", "accept": "application/json"}
        response = requests.post(dcm_url + anchor, headers=headers)
        return response

    @classmethod
    def instantiate_logical_clouds_assertion_check(
        self, logfile, response, responseCode=None, responseText=None
    ):
        msgResCode = "Response Code does not match, expected: {0}, actual: {1}".format(
            responseCode, response.status_code
        )
        utils.logAssert(response.status_code == responseCode, msgResCode, logfile)
