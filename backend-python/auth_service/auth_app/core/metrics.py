from prometheus_client import Counter, Histogram, Gauge

SERVICE = "auth"

# General HTTP metrics
HTTP_REQUESTS_TOTAL = Counter(
    "http_requests_total",
    "Total HTTP requests",
    ["service", "method", "handler", "status"]
)

HTTP_REQUEST_DURATION_SECONDS = Histogram(
    "http_request_duration_seconds",
    "HTTP request latency (seconds)",
    ["service", "method", "handler"],
    buckets=[0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5]
)

HTTP_REQUESTS_IN_FLIGHT = Gauge(
    "http_requests_in_flight",
    "Number of in-flight requests",
    ["service"]
)

# Business metrics
AUTH_USER_REGISTRATIONS_TOTAL = Counter(
    "auth_user_registrations_total",
    "Count of user registrations",
    ["service", "status"]
)

AUTH_USER_LOGINS_TOTAL = Counter(
    "auth_user_logins_total",
    "Login attempts",
    ["service", "status"]
)

AUTH_TOKEN_REFRESHES_TOTAL = Counter(
    "auth_token_refreshes_total",
    "Token refresh attempts",
    ["service", "status"]
)

AUTH_USER_LOGOUTS_TOTAL = Counter(
    "auth_user_logouts_total",
    "Logout operations",
    ["service", "scope"]
)

AUTH_ACTIVE_SESSIONS = Gauge(
    "auth_active_sessions",
    "Active authenticated sessions",
    ["service"]
)
