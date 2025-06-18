import json
import logging
from datetime import datetime, UTC
from logging.config import dictConfig

from shared.core.config import settings


if settings.DEBUG:  # development
    LOG_LEVEL_APP = "DEBUG"
    LOG_LEVEL_SQLALCHEMY = "DEBUG"
    LOG_LEVEL_HANDLER = "DEBUG"
    LOG_LEVEL_ROOT = "DEBUG"
else:  # production
    LOG_LEVEL_APP = "INFO"
    LOG_LEVEL_SQLALCHEMY = "WARNING"
    LOG_LEVEL_HANDLER = "INFO"
    LOG_LEVEL_ROOT = "WARNING"


class JsonFormatter(logging.Formatter):
    def format(self, record):
        log_record = {
            "timestamp": datetime.now(UTC).isoformat(),
            "level": record.levelname,
            "logger": record.name,
            "module": record.module,
            "line": record.lineno,
            "message": record.getMessage()
        }

        if record.exc_info:
            log_record["exception"] = self.formatException(record.exc_info)

        return json.dumps(log_record)


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


def setup_logging():
    dictConfig(log_config)
