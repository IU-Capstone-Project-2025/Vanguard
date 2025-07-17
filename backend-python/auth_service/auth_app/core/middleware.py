import logging
import time

from fastapi import Request
from starlette.exceptions import HTTPException as StarletteHTTPException
from starlette.middleware.base import BaseHTTPMiddleware

from auth_app.core.metrics import (
    HTTP_REQUESTS_TOTAL,
    HTTP_REQUEST_DURATION_SECONDS,
    HTTP_REQUESTS_IN_FLIGHT,
    SERVICE
)

logger = logging.getLogger("app")

EXCLUDED_PATHS = {"/health", "/metrics"}


class MetricsMiddleware(BaseHTTPMiddleware):
    async def dispatch(self, request: Request, call_next):
        path = request.url.path.rstrip('/')
        if path in EXCLUDED_PATHS:
            return await call_next(request)

        method = request.method
        status = "500"

        HTTP_REQUESTS_IN_FLIGHT.labels(service=SERVICE).inc()
        start = time.time()

        try:
            response = await call_next(request)
            status = str(response.status_code)
            return response
        except StarletteHTTPException as exc:
            status = str(exc.status_code)
            raise
        except Exception:
            raise
        finally:
            duration = time.time() - start
            HTTP_REQUESTS_TOTAL.labels(service=SERVICE, method=method, handler=path, status=status).inc()
            HTTP_REQUEST_DURATION_SECONDS.labels(service=SERVICE, method=method, handler=path).observe(duration)
            HTTP_REQUESTS_IN_FLIGHT.labels(service=SERVICE).dec()


class LoggingMiddleware(BaseHTTPMiddleware):
    async def dispatch(self, request: Request, call_next):
        path = request.url.path.rstrip('/')
        if path in EXCLUDED_PATHS:
            return await call_next(request)

        start_time = time.time()
        status_code = 500

        try:
            response = await call_next(request)
            status_code = response.status_code
            return response
        except StarletteHTTPException as exc:
            status_code = exc.status_code
            raise
        except Exception:
            raise
        finally:
            if path not in EXCLUDED_PATHS:
                process_time = round(time.time() - start_time, 4)
                client = request.client.host if request.client else "unknown"

                logger.info({
                    "method": request.method,
                    "path": path,
                    "status_code": status_code,
                    "process_time": f"{process_time}s",
                    "client": client,
                })
