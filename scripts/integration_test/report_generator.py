import os


import uuid
from datetime import datetime
from enum import Enum
from typing import Dict, List, Optional
from pydantic import UUID4, BaseModel, Field
from pathlib import Path

from configuration.env_config import TEST_REPORT_FILE_PREFIX
from utils.models import AutomationType


def file_appender(input_data, output_path):
    with open(output_path, "a") as f:
        f.writelines(input_data)


def table_html(headers, data):
    pre_existing_template = "<table style='width:50%'>"
    pre_existing_template += "<tr>"
    for header_name in headers:
        pre_existing_template += (
            "<th style='background-color:#3DBBDB;width:85;color:white'>"
            + header_name
            + "</th>"
        )
    pre_existing_template += "</tr>"
    for i in range(len(data[0])):
        sub_template = "<tr style='text-align:center'>"
        for j in range(len(headers)):
            sub_template += "<td>" + str(data[j][i]) + "</td>"
        sub_template += "<tr/>"
        pre_existing_template += sub_template
    pre_existing_template += "</table>"
    return pre_existing_template



class TestCaseResult(BaseModel):
    section: str
    test_case: str
    result: str
    duration: str
    reason: str




class ReportGenerator:
    """Report Generator class
    """    

    def __init__(self) -> None:
        self.report_folder_path = os.path.join(os.getcwd(), "outputs", "reports")
        # handle other exceptions
        try:
            Path(self.report_folder_path).mkdir(parents=True, exist_ok=True)
        except Exception as exc:
            print(exc)  # Change afterwards

        self.report_file_path = os.path.join(
            self.report_folder_path,
            f'{TEST_REPORT_FILE_PREFIX}{datetime.now().strftime("%Y%m%d_%H%M%S")}.html',
        )

    def html_writer(self, *args):
        all_args = args
        op_str = ""
        tables = ""
        for arg in args:
            if type(arg) != dict:
                op_str += arg
            elif type(arg) == dict:
                if "table" in arg.keys():
                    headers, data = arg["table"]
                    tables += table_html(headers, data)
        op_str = "<p >" + op_str + " " + "</p>" + "\n"
        op_str += tables
        file_appender(op_str, self.report_file_path)

    def initial_report(self, test_type: str, execution_start_date: str, total_execution_time: str, total: int, executed: int, passed: int, failed: int):
        style_content = """
        table, th, td {border: 1px solid black;border-collapse: collapse;border-spacing:8px}
        body {
            font-family: 'Roboto';
        }
        """

        html_content = f"""
        <!DOCTYPE html>
        <html>
        <head>
        <title>EMCO Automation test results.</title>
        <link href='https://fonts.googleapis.com/css?family=Roboto' rel='stylesheet'>

        <style>
        {style_content}

        </style>
        </head>
        <body>

        <h1>EMCO API Automation test results.</h1>

        <h2>API Test Results summary:</h2>
        <p><strong>Test case type</strong>: {test_type} </p> 
        <p><strong>Automation start date</strong>: {execution_start_date} </p> 
        <p><strong>Total Duration</strong>: {total_execution_time} </p> 

        <table style='width:50%'>
        <tr>
            <th scope="col" style='background-color:#3DBBDB;width:85;color:white'>Planned Test Cases</th>
            <th scope="col" style='background-color:#3DBBDB;width:85;color:white'>Test Case Executed</th>
            <th scope="col" style='background-color:#3DBBDB;width:85;color:white'>Passed</th>
            <th scope="col" style='background-color:#3DBBDB;width:85;color:white'>Failed</th>
        </tr> 
        <tr>
            <td>{total}</td>
            <td>{executed}</td>
            <td>{passed}</td>
            <td>{failed}</td>
        </tr>
        </table>

        <h2>API Test Results details:</h2>


        """
        return html_content

    def generate_report(self, test_case_results: List[TestCaseResult], automation_type: AutomationType, total_duration: str, execution_start_date: str) -> str:
        """Generates a HTML Report

        Args:
            test_case_results (List[TestCaseResult]): List of test case results.
        """

        section_list = []
        test_case_list = []
        result_list = []
        duration_list = []
        reason_list = []

        results = [i.result for i in test_case_results]

        total = automation_type.total_test_cases
        executed = len(test_case_results)
        passed = results.count("PASS")
        failed = results.count("FAIL")

        test_type = f"{automation_type.type} ({automation_type.description})"

        initial_report_content = self.initial_report(test_type=test_type, execution_start_date= execution_start_date, total_execution_time= total_duration, total=total, executed=executed, passed=passed, failed=failed)

        file_appender(initial_report_content, self.report_file_path)


        for test_case_result in test_case_results:
            section_list.append(test_case_result.section)
            test_case_list.append(test_case_result.test_case)
            result_list.append(test_case_result.result)
            duration_list.append(test_case_result.duration)
            reason_list.append(test_case_result.reason)

        self.html_writer(
            {
                "table": (
                    ["Section", "Test Case", "Result", "Duration", "Reason"],
                    [
                        section_list,
                        test_case_list,
                        result_list,
                        duration_list,
                        reason_list,
                    ],
                )
            }
        )



        
        # closing tags 
        ending_content = """
        </body>
        </html>
        """
        file_appender(ending_content, self.report_file_path)

        return self.report_file_path
