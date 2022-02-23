import json
import requests
import os
from requests.exceptions import HTTPError
from json_parser import InputParser
from utils import utils


class TestCase_deleteLogicalClouds:
    def test_deleteLogicalClouds(self, cwd, dcm_url, emcotestlogging, api_input):
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
                    # print(params_from_inputfile.get_request_body(0, item))
                    # print(params_from_inputfile.get_response_code(0, item))
                    # print(params_from_inputfile.get_response_text(0, item))
                    # print(params_from_inputfile.get_param_projectName(0, item))
                    response = self.delete_logical_clouds(
                        dcm_url, params_from_inputfile.get_anchor(0)
                    )
                    logfile.log_debug(
                        "Received Response Code: " + str(response.status_code)
                    )
                    self.delete_logical_clouds_assertion_check(
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
    def delete_logical_clouds(self, dcm_url, anchor):
        headers = {"accept": "*/*"}
        response = requests.delete(dcm_url + anchor, headers=headers)
        return response

    @classmethod
    def delete_logical_clouds_assertion_check(
        self, logfile, response, responseCode=None
    ):
        msgResCode = "Response Code does not match, expected: {0}, actual: {1}".format(
            responseCode, response.status_code
        )
        utils.logAssert(response.status_code == responseCode, msgResCode, logfile)
