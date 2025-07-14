from logging.config import dictConfig

from shared.utils.logging_formatter import JsonFormatter


def get_log_config(debug: bool = False):
    if debug:  # development
        config = {
            "app": "DEBUG",
            "sqlalchemy_engine": "WARNING",
            "sqlalchemy_pool": "WARNING",
            "sqlalchemy_orm": "WARNING",
            "root": "WARNING",
            "handler": "DEBUG"
        }
    else:  # production
        config = {
            "app": "INFO",
            "sqlalchemy_engine": "WARNING",
            "sqlalchemy_pool": "ERROR",
            "sqlalchemy_orm": "ERROR",
            "root": "WARNING",
            "handler": "INFO"
        }

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
                "level": config.get("handler"),
                "formatter": "json",
                "stream": "ext://sys.stdout",
            }
        },
        "loggers": {
            "app": {
                "handlers": ["console"],
                "level": config.get("app"),
                "propagate": False
            },
            "sqlalchemy.engine": {
                "handlers": ["console"],
                "level": config.get("sqlalchemy_engine"),
                "propagate": False
            },
            "sqlalchemy.pool": {
                "handlers": ["console"],
                "level": config.get("sqlalchemy_pool"),
                "propagate": False
            },
            "sqlalchemy.orm": {
                "handlers": ["console"],
                "level": config.get("sqlalchemy_orm"),
                "propagate": False
            }
        },
        "root": {
            "handlers": ["console"],
            "level": config.get("root")
        }
    }

    return log_config


def setup_loggers(debug: bool = False):
    dictConfig(get_log_config(debug))
