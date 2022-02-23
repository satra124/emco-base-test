import json
import requests
import os
from requests.exceptions import HTTPError
from json_parser import InputParser
from utils import utils


class TestCase_postProj:
    def test_postProjects(self, cwd, orchestrator_url, emcotestlogging, api_input):
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
                    response = self.post_projects(
                        orchestrator_url,
                        params_from_inputfile.get_anchor(0),
                        params_from_inputfile.get_request_body(0, item),
                    )
                    logfile.log_debug(
                        "Received Response Code: " + str(response.status_code)
                    )
                    logfile.log_debug("Received Response Text: " + str(response.text))
                    # logfile.log_debug("Received Response Header SessionId: " + response.headers.get('SessionId'))
                    self.post_projects_assertion_check(
                        logfile,
                        response,
                        params_from_inputfile.get_response_code(0, item),
                        params_from_inputfile.get_response_text(0, item),
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
    def post_projects(self, orchestrator_url, anchor, request_body=None):
        headers = {"Content-Type": "application/json", "accept": "application/json"}
        payload = request_body
        response = requests.post(
            orchestrator_url + anchor, data=payload, headers=headers
        )
        return response

    @classmethod
    def post_projects_assertion_check(
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
