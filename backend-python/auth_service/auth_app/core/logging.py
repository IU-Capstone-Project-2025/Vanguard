from logging.config import dictConfig

from shared.utils.logging_formatter import JsonFormatter


def get_log_config(debug: bool = False):
    if debug:  # development
        LOG_LEVEL_APP = "DEBUG"
        LOG_LEVEL_SQLALCHEMY = "DEBUG"
        LOG_LEVEL_HANDLER = "DEBUG"
        LOG_LEVEL_ROOT = "DEBUG"
    else:  # production
        LOG_LEVEL_APP = "INFO"
        LOG_LEVEL_SQLALCHEMY = "WARNING"
        LOG_LEVEL_HANDLER = "INFO"
        LOG_LEVEL_ROOT = "WARNING"

    log_config = {
        "version": 1,
        "disable_existing_loggers": False,
        "formatters": {
            "json": {
                "()": JsonFormatter
            }
        },
        "handlers": {
            "console": {
                "class": "logging.StreamHandler",
                "level": LOG_LEVEL_HANDLER,
                "formatter": "json",
                "stream": "ext://sys.stdout",
            }
        },
        "loggers": {
            "app": {
                "handlers": ["console"],
                "level": LOG_LEVEL_APP,
                "propagate": False
            },
            "sqlalchemy.engine": {
                "handlers": ["console"],
                "level": LOG_LEVEL_SQLALCHEMY,
                "propagate": False
            }
        },
        "root": {
            "handlers": ["console"],
            "level": LOG_LEVEL_ROOT
        }
    }

    return log_config


def setup_logging(debug: bool = False):
    dictConfig(get_log_config(debug))
