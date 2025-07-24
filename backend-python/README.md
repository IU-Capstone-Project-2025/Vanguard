### **Auth Service Endpoints**  
Handles user registration, authentication, tokens, and OAuth.

| Method | Endpoint                     | Description                                                                 | Auth Required |
|--------|------------------------------|-----------------------------------------------------------------------------|--------------|
| `GET`  | `/health`               | Health check endpoint.                                                      | No           |
| `POST` | `/register`             | Register a new user. Returns user ID.                                       | No           |
| `POST` | `/login`                | Authenticate user. Returns **access token** and **refresh token**.          | No           |
| `POST` | `/refresh`              | Refresh expired access token using a valid refresh token.                   | No¹          |
| `GET`  | `/me`                   | Get current user's profile (ID, email, username).                          | Yes (JWT)    |
| `PUT`  | `/me`                   | Update current user's profile (email, password, etc.).                     | Yes (JWT)    |
| `POST` | `/logout`               | Invalidate refresh token (server-side revocation).                          | Yes (JWT)    |
| `POST` | `/logout/all`           | Invalidate all refresh tokens for the user.                                 | Yes (JWT)    |
| `GET`  | `/oauth/{provider}`     | Initiate OAuth flow (e.g., `google`, `yandex`, `vk`). Redirects to provider. | No           |
| `GET`  | `/oauth/{provider}/callback` | OAuth callback. Exchanges code for tokens, issues app tokens.             | No           |

> **Notes**:  
> ¹ `refresh` requires a valid refresh token in the request body.  

---

### **Quiz Service Endpoints**  
Handles quizzes, tags, images, and filtering. Uses JWT for auth.

| Method | Endpoint                     | Description                                                                 | Auth Required | Parameters/Request Body |
|--------|------------------------------|-----------------------------------------------------------------------------|---------------|--------------------------|
| `GET`  | `/health`             | Health check endpoint.                                                      | No            | -                               |
| `POST` | `/`                   | Create a new quiz. Automatically creates/links tags.                        | Yes (JWT)     | **Body:** `title`, `description`, `is_public`, `questions` (JSON), `tags` (list of strings). |
| `GET`  | `/{quiz_id}`         | Get quiz by ID. Unauthenticated users see only public quizzes.              | Conditional²  | -                        |
| `PUT`  | `/{quiz_id}`         | Update quiz (title, description, questions, tags). Owner only.             | Yes (JWT)     | **Body:** Same as `POST`, partial updates allowed. |
| `DELETE`| `/{quiz_id}`         | Delete quiz. Owner only.                                                    | Yes (JWT)     | -                        |
| `GET`  | `/`                   | List/filter quizzes. Supports pagination.                                  | Optional      | **Query Params:**<br>- `public` (bool): Only public quizzes.<br>- `mine` (bool): Only current user's quizzes.<br>- `user_id` (UUID): Public quizzes by a user.<br>- `search` (str): Text search in title/description.<br>- `tag` (list): Filter by tags (e.g., `?tag=math&tag=science`).<br>- `page` (int), `size` (int): Pagination. |
| `POST` | `/images`            | Upload image for quiz/question/answer. Returns **S3 URL**.                  | Yes (JWT)     | **Body:** `image` (file upload). |
| `GET`  | `/tags`                      | List tags. Supports search and pagination.                                  | No            | **Query Params:**<br>- `name` (str): Partial tag name search.<br>- `page` (int), `size` (int). |
| `GET`  | `/tags/{tag_id}`             | Get tag details by ID.                                                      | No            | -                        |

**Notes**:  
> ² `/{quiz_id}`:  
>   - Public quizzes: No auth required.  
>   - Private quizzes: Requires JWT and user must be the owner.  

**Filtering Logic** for `GET /`:  
>   - **Unauthenticated users:** Only `public=true` and `search`/`tag` filters allowed.  
>   - **Authenticated users:**  
>     - Default: Returns public quizzes **and** the user's own quizzes.  
>     - `public=true`: Public quizzes only.  
>     - `mine=true`: User's quizzes (public + private).  
>     - `user_id=UUID`: Public quizzes by another user.  
>     - `tag=*`: AND filter (quiz must have all specified tags).

## Metrics

### Auth Service
| Metric Name                            | Type      | Labels                                        | Description                                                                           |
|----------------------------------------|-----------|-----------------------------------------------|---------------------------------------------------------------------------------------|
| `http_requests_total`                  | Counter   | `service`, `method`, `handler`, `status`      | Total count of HTTP requests received by the service                                  |
| `http_request_duration_seconds`        | Histogram | `service`, `method`, `handler`                | Request latency distribution (buckets: [0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5]) |
| `http_requests_in_flight`              | Gauge     | `service`                                     | Current number of in-flight HTTP requests                                             |
| `auth_user_registrations_total`        | Counter   | `service`, `status`                           | Count of successful user registrations                                                |
| `auth_user_logins_total`               | Counter   | `service`, `status`                           | Login attempts (success/failure)                                                      |
| `auth_token_refreshes_total`           | Counter   | `service`, `status`                           | Token refresh attempts (success/failure)                                              |
| `auth_user_logouts_total`              | Counter   | `service`, `scope`                            | Logout operations (single/all sessions)                                               |
| `auth_active_sessions`                 | Gauge     | `service`                                     | Current number of active authenticated sessions                                       |

### Quiz Service
| Metric Name                      | Type      | Labels                                | Description                                                                               |
|----------------------------------|-----------|---------------------------------------|-------------------------------------------------------------------------------------------|
| `http_requests_total`            | Counter   | `service`, `method`, `handler`, `status` | Total count of HTTP requests received by the service                                      |
| `http_request_duration_seconds`  | Histogram | `service`, `method`, `handler`        | Request latency distribution (buckets: [0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5]) |
| `http_requests_in_flight`        | Gauge     | `service`                             | Current number of in-flight HTTP requests                                                 |
| `quiz_creations_total`           | Counter   | `service`, `status`, `visibility`     | Count of new quizzes created (`visibility`: public/private)                               |
| `quiz_fetches_total`             | Counter   | `service`, `status`, `public_only`    | Quiz retrieval operations (`public_only`: true for anonymous access)                      |
| `quiz_updates_total`             | Counter   | `service`, `status`                   | Quiz modification operations                                                              |
| `quiz_deletes_total`             | Counter   | `service`, `status`                   | Quiz deletion operations                                                                  |
| `quiz_listing_requests_total`    | Counter   | `service`, `status`, `filter_type`    | Quiz listing/filter operations (`filter_type`: public/mine/search/tag/user)               |
| `quiz_image_uploads_total`       | Counter   | `service`, `status`                   | Image upload operations                                                                   |
| `quiz_image_upload_size_bytes`   | Summary   | `service`, `status`                   | Distribution of uploaded image sizes                                                      |

### Real-Time Service
| Metric Name                           | Type      | Labels                                                                 | Description                                                                 |
|---------------------------------------|-----------|------------------------------------------------------------------------|-----------------------------------------------------------------------------|
| `websocket_connection_attempts_total` | Counter   | `service`, `status` ("success"/"failure"), `reason`                   | Total WS connection attempts with failure reasons                           |
| `websocket_active_connections`        | Gauge     | `service`, `user_type` ("admin"/"participant")                        | Current active WS connections                                               |
| `websocket_connection_duration_seconds` | Histogram | `service`, `user_type`                                                | Connection duration (buckets: [60, 300, 600, 1800, 3600] = 1m,5m,10m,30m,1h) |
| `websocket_messages_sent_total`       | Counter   | `service`, `msg_type` ("welcome","question","game_start","leaderboard") | Messages sent to clients                                                    |
| `websocket_messages_received_total`   | Counter   | `service`, `msg_type` ("answer")                                      | Messages received from clients                                              |
| `websocket_message_processing_seconds` | Histogram | `service`, `msg_type`                                                 | Message processing time (buckets: [0.001, 0.01, 0.1, 0.5])                |
| `websocket_message_errors_total`      | Counter   | `service`, `msg_type`, `reason`                                       | Message processing errors                                                   |
| `quiz_answers_submitted_total`        | Counter   | `service`                                                             | Total answers submitted (business metric)                                   |
| `sessions_in_progress`                | Gauge     | `service`                                                             | Active sessions not yet ended                                               |
| `session_events_processed_total`      | Counter   | `service`, `event_type` ("question_start","session_end")              | Session lifecycle events                                                    |

### Session service
| Metric Name                           | Type      | Labels                                                                 | Description                                                                 |
|---------------------------------------|-----------|------------------------------------------------------------------------|-----------------------------------------------------------------------------|
| `http_requests_total`                 | Counter   | `service`, `method`, `handler`, `status`                               | Total HTTP requests for session management                                  |
| `http_request_duration_seconds`       | Histogram | `service`, `method`, `handler`                                         | HTTP request latency (buckets: [0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5]) |
| `websocket_connection_attempts_total` | Counter   | `service`, `status` ("success"/"failure"), `reason`                    | WS connection attempts with failure reasons                                 |
| `websocket_active_connections`        | Gauge     | `service`, `user_type` ("admin"/"participant")                         | Current active WS connections                                               |
| `websocket_messages_sent_total`       | Counter   | `service`, `msg_type`                                                  | Messages sent to clients (types: "question", "game_start", "leaderboard")   |
| `websocket_messages_received_total`   | Counter   | `service`, `msg_type`                                                  | Messages received from clients (primarily "answer")                         |
| `sessions_created_total`              | Counter   | `service`, `status`                                                    | New sessions created via `/sessions` endpoint                               |
| `sessions_started_total`              | Counter   | `service`                                                              | Sessions transitioned to active state                                      |
| `session_joins_total`                 | Counter   | `service`, `status` ("success"/"failure")                              | Participant join attempts with success/failure                              |
| `questions_advanced_total`            | Counter   | `service`                                                              | "Next question" triggers by admins                                          |
| `sessions_ended_total`                | Counter   | `service`                                                              | Completed sessions via `/session/{id}/end`                                  |
| `sessions_active`                     | Gauge     | `service`                                                              | Current active sessions (not yet ended)                                     |
| `user_removals_total`                 | Counter   | `service`                                                              | User kick operations via `/delete-user`                                     |
