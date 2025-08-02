import json
import logging
from datetime import datetime, UTC


class JsonFormatter(logging.Formatter):
    def __init__(self, service_name: str, *args, **kwargs):
        super().__init__(*args, **kwargs)
        self.service_name = service_name

    def format(self, record: logging.LogRecord) -> str:
        log_record = {
            "timestamp": datetime.now(UTC).isoformat(),
            "level": record.levelname,
            "service": self.service_name,
            "file": record.filename,
            "line": record.lineno,
            "message": record.getMessage(),
            "metadata": {}
        }

        if hasattr(record, "metadata"):
            log_record["metadata"] = record.metadata

        if record.exc_info:
            log_record["exception"] = self.formatException(record.exc_info)

        return json.dumps(log_record, default=self._json_fallback)

    @staticmethod
    def _json_fallback(obj):
        if isinstance(obj, datetime):
            return obj.isoformat()
        return str(obj)
