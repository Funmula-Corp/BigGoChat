#!/usr/bin/env python3

import os
import glob
import subprocess
import xml.etree.ElementTree as ET
from html import escape
from dataclasses import dataclass, field
from typing import Dict, List, Optional, Tuple

RESULTS_DIR: str = "cypress/results/junit"
OUTPUT_FILE: str = "failed_tests_report.html"

@dataclass
class TestCaseFailure:
    classname: str
    test_name: str
    failure_message: str

@dataclass
class SpecSummary:
    testsuite_name: str
    total_testcases: int
    total_failed: int
    total_success: int
    last_modifier: Optional[str] = None

@dataclass
class SummaryData:
    total_specs: int = 0
    total_testcases: int = 0
    total_failed: int = 0
    total_success: int = 0
    specs: Dict[str, SpecSummary] = field(default_factory=dict)

def clear_output_file() -> None:
    with open(OUTPUT_FILE, 'w') as f:
        # Initialize HTML structure
        f.write("<!DOCTYPE html>\n<html lang='en'>\n<head>\n")
        f.write("    <meta charset='UTF-8'>\n")
        f.write("    <meta http-equiv='X-UA-Compatible' content='IE=edge'>\n")
        f.write("    <meta name='viewport' content='width=device-width, initial-scale=1.0'>\n")
        f.write("    <title>Failed Cypress Tests Report</title>\n")
        f.write("    <style>\n")
        f.write("        body { font-family: Arial, sans-serif; margin: 20px; }\n")
        f.write("        h1 { color: #333; }\n")
        f.write("        .summary { margin-bottom: 40px; }\n")
        f.write("        .testsuite { margin-bottom: 30px; }\n")
        f.write("        .testsuite-header { background-color: #f2f2f2; padding: 10px; border-radius: 5px; }\n")
        f.write("        table { width: 100%; border-collapse: collapse; margin-top: 10px; }\n")
        f.write("        th, td { border: 1px solid #ddd; padding: 8px; }\n")
        f.write("        th { background-color: #4CAF50; color: white; }\n")
        f.write("        tr:nth-child(even) { background-color: #f9f9f9; }\n")
        f.write("    </style>\n</head>\n<body>\n")
        f.write("<h1>Failed Cypress Tests Report</h1>\n")

def get_last_modifier_email(spec_file: str) -> Optional[str]:
    try:
        result = subprocess.run(
            ['git', 'log', '-1', '--format=%ae', spec_file],
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE,
            text=True,
            check=True
        )
        return result.stdout.strip()
    except subprocess.CalledProcessError:
        return None

def parse_junit_results() -> SummaryData:
    summary_data: SummaryData = SummaryData()

    for xml_file in glob.glob(os.path.join(RESULTS_DIR, "*.xml")):
        if os.path.isfile(xml_file):
            try:
                tree = ET.parse(xml_file)
                root = tree.getroot()

                spec_file_path: str = ''

                for testsuite in root.findall(".//testsuite"):
                    testsuite_name: str = testsuite.attrib.get('name', 'Unknown Test Suite')
                    if spec_file_path == '':
                        spec_file: str = testsuite.attrib.get('file', '')
                        spec_file_path: str = os.path.join("cypress", spec_file) if spec_file else ''

                    total_testcases: int = int(testsuite.attrib.get('tests', 0))
                    total_failed: int = int(testsuite.attrib.get('failures', 0))
                    total_success: int = total_testcases - total_failed

                    last_modifier_email = get_last_modifier_email(spec_file_path)
                    
                    summary_data.total_specs += 1
                    summary_data.total_testcases += total_testcases
                    summary_data.total_failed += total_failed
                    summary_data.total_success += total_success

                    spec_summary = SpecSummary(
                        testsuite_name=testsuite_name,
                        total_testcases=total_testcases,
                        total_failed=total_failed,
                        total_success=total_success,
                        last_modifier=last_modifier_email
                    )

                    # Collect failed test cases
                    failed_testcases: List[TestCaseFailure] = []
                    for testcase in testsuite.findall(".//testcase[failure]"):
                        classname: str = escape(testcase.attrib.get('classname', 'Unknown Class'))
                        test_name: str = escape(testcase.attrib.get('name', 'Unknown Test'))
                        failure: Optional[ET.Element] = testcase.find('failure')
                        failure_msg: str = escape(failure.attrib.get('message', 'No message provided')) if failure is not None else 'No message provided'
                        failed_testcases.append(TestCaseFailure(classname, test_name, failure_msg))
                    
                    spec_summary.failed_testcases = failed_testcases

                    summary_data.specs[spec_file_path] = spec_summary

            except ET.ParseError:
                continue

    return summary_data

def write_summary(summary_data: SummaryData) -> None:
    with open(OUTPUT_FILE, 'a') as report:

        report.write("<h3>Filter Specs</h3>\n")
        report.write("""
        <button onclick="filterSpecs('all')">Show All Specs</button>
        <button onclick="filterSpecs('failed')">Show Failed Specs Only</button>
        <script>
            function filterSpecs(filter) {
                var rows = document.querySelectorAll('.spec-row');
                rows.forEach(function(row) {
                    if (filter === 'all') {
                        row.style.display = '';
                    } else if (filter === 'failed') {
                        var failed = row.getAttribute('data-failed');
                        row.style.display = failed === 'true' ? '' : 'none';
                    }
                });
            }
        </script>
        """)

        report.write("<h3>Spec List</h3>\n")
        report.write("<table>\n")
        report.write("  <tr>\n")
        report.write("    <th>Spec File</th>\n")
        report.write("    <th>Total Testcases</th>\n")
        report.write("    <th>Failed</th>\n")
        report.write("    <th>Passed</th>\n")
        report.write("    <th>Last Modifier</th>\n")
        report.write("  </tr>\n")

        for spec, data in summary_data.specs.items():
            has_failed = 'true' if data.total_failed > 0 else 'false'
            report.write(f"  <tr class='spec-row' data-failed='{has_failed}'>\n")
            report.write(f"    <td><a href=#{escape(spec)}>{escape(spec)}</a></td>\n")
            report.write(f"    <td>{data.total_testcases}</td>\n")
            report.write(f"    <td>{data.total_failed}</td>\n")
            report.write(f"    <td>{data.total_success}</td>\n")
            report.write(f"    <td>{escape(data.last_modifier) if data.last_modifier else 'N/A'}</td>\n")
            report.write("  </tr>\n")

        report.write("</table>\n")
        report.write("<hr>\n</div>\n")

def parse_junit_results_and_generate_report(summary_data: SummaryData) -> None:
    with open(OUTPUT_FILE, 'a') as report:
        write_summary(summary_data)

        for spec, data in summary_data.specs.items():
            report.write(f"<div class='testsuite'>\n")
            report.write(f"  <div id='{escape(spec)}' class='testsuite-header'><strong>Spec:</strong> {escape(spec)}</div>\n")
            report.write("  <hr>\n")
            
            testsuite_info: str = f"Test Suite: {escape(data.testsuite_name)}"
            report.write(f"  <p><strong>{testsuite_info}</strong></p>\n")
            
            if data.failed_testcases:
                report.write("  <table>\n")
                report.write("    <tr>\n")
                report.write("      <th>Class Name</th>\n")
                report.write("      <th>Test Name</th>\n")
                report.write("      <th>Failure Message</th>\n")
                report.write("    </tr>\n")
                
                for testcase in data.failed_testcases:
                    report.write("    <tr>\n")
                    report.write(f"      <td>{testcase.classname}</td>\n")
                    report.write(f"      <td>{testcase.test_name}</td>\n")
                    report.write(f"      <td>{testcase.failure_message}</td>\n")
                    report.write("    </tr>\n")
                
                report.write("  </table>\n")
            else:
                report.write("  <p>No failed tests.</p>\n")
            
            report.write("</div>\n")

def print_summary(summary_data: SummaryData) -> None:
    try:
        with open(OUTPUT_FILE, 'r') as f:
            content: str = f.read()
        total_failed: int = summary_data.total_failed
    except FileNotFoundError:
        total_failed = 0

    summary_html: str = f"""
    <hr>
    <p><strong>Total Failed Tests:</strong> {total_failed}</p>
    <p>Details have been saved to <a href="{escape(OUTPUT_FILE)}">{escape(OUTPUT_FILE)}</a></p>
    <hr>
    </body>\n</html>
    """
    with open(OUTPUT_FILE, 'a') as f:
        f.write(summary_html)
    
    print("Test Results Summary:")
    print("----------------------------------------")
    print(f"Total Failed Tests: {total_failed}")
    print(f"Details have been saved to {OUTPUT_FILE}")

    if total_failed > 0:
        exit(1)
    else:
        exit(0)

def main() -> None:
    if not os.path.isdir(RESULTS_DIR):
        print(f"Error: Results directory not found at {RESULTS_DIR}")
        exit(1)
    
    clear_output_file()
    print("Analyzing Cypress test results...")
    summary_data = parse_junit_results()
    parse_junit_results_and_generate_report(summary_data)
    print_summary(summary_data)

if __name__ == "__main__":
    main()
