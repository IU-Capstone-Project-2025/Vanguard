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
