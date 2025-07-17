from prometheus_client import Counter, Histogram, Gauge, Summary

SERVICE = "quiz"

# General HTTP metrics
HTTP_REQUESTS_TOTAL = Counter(
    'http_requests_total',
    'Total count of HTTP requests received by the service',
    ['service', 'method', 'handler', 'status']
)

HTTP_REQUEST_DURATION_SECONDS = Histogram(
    'http_request_duration_seconds',
    'Request latency distribution',
    ['service', 'method', 'handler'],
    buckets=[0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5]
)

HTTP_REQUESTS_IN_FLIGHT = Gauge(
    'http_requests_in_flight',
    'Current number of in-flight HTTP requests',
    ['service']
)

# Business metrics
QUIZ_CREATIONS_TOTAL = Counter(
    'quiz_creations_total',
    'Count of new quizzes created',
    ['service', 'status', 'visibility']
)

QUIZ_FETCHES_TOTAL = Counter(
    'quiz_fetches_total',
    'Quiz retrieval operations',
    ['service', 'status', 'public_only']
)

QUIZ_UPDATES_TOTAL = Counter(
    'quiz_updates_total',
    'Quiz modification operations',
    ['service', 'status']
)

QUIZ_DELETES_TOTAL = Counter(
    'quiz_deletes_total',
    'Quiz deletion operations',
    ['service', 'status']
)

QUIZ_LISTING_REQUESTS_TOTAL = Counter(
    'quiz_listing_requests_total',
    'Quiz listing/filter operations',
    ['service', 'status', 'filter_type']
)

QUIZ_IMAGE_UPLOADS_TOTAL = Counter(
    'quiz_image_uploads_total',
    'Image upload operations',
    ['service', 'status']
)

QUIZ_IMAGE_UPLOAD_SIZE_BYTES = Summary(
    'quiz_image_upload_size_bytes',
    'Distribution of uploaded image sizes',
    ['service', 'status']
)
