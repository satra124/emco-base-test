#!/usr/bin/env python
import logging  # To access Python's logging module.
import sys

class LogManager(object):
    """Class to implement logging feature using Python's logging module

    (public interface)
    self.object_python_logger - Logger object used by all classes/modules.
    """

    def __init__(self, log_file_name, header_string):
        """Constructor for this class will do the following tasks:

        1.) Use Python's FileHandler to specify a logging filename and location (without Output folder).
        This creates a logger object.
        2.) Use Python's Formatter to determine the message contents and add it as an event to the logger object.
        3.) Set a logging level to determine what messages show up.
        """
        self.object_python_logger = logging.getLogger(header_string)
        object_file_handler = logging.FileHandler(log_file_name)
        self.object_python_logger.addHandler(object_file_handler)
        object_console_handler = logging.StreamHandler(sys.stdout)
        self.object_python_logger.addHandler(object_console_handler)
        self.log_set_log_level(logging.DEBUG)

    def log_set_log_level(self, object_log_level_parameter):
        """Method to set the logging level for the Python logger object.

        object_log_level_parameter - (object) A level defined in the Python logging module.
        """
        self.object_python_logger.setLevel(object_log_level_parameter)
        return 0

    def log_debug(self, string_debug_message):
        """Log a DEBUG message to the log file.

        @parameter string_message - (string) Message to write to log file.
        """
        self.object_python_logger.debug(msg=string_debug_message)
        return 0

    def log_error(self, string_error_message):
        """Log an ERROR message to the log file.

        @parameter string_message - (string) Message to write to log file.
        """
        self.object_python_logger.error(msg=string_error_message)
        return 0

    def log_shut_down(self):
        """"""
        handlers = self.object_python_logger.handlers[:]
        for handler in handlers:
            handler.close()
            self.object_python_logger.removeHandler(handler)
            logging.shutdown()
        print("\nShutting down Python LOGGING")
        return 0


