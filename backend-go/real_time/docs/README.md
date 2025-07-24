# Real‑Time Service WebSocket Guide

This document explains **when** to invoke the `/ws` endpoint, **what** data to send, and **what** messages to expect in response—without any client‑side code samples.

---

## 1. Establishing the WebSocket Connection

- **When**: As soon as the user (admin or participant) has a valid JWT from the Session Service.
- **Request**:
    - Method: `GET`
    - URL: `/ws?token=<JWT>`
        - `token` query parameter must contain the signed JWT with claims:
          ```yaml
          userId: string
          sessionId: string
          userType: "admin" / "participant"
          exp: integer
          ```  
- **Response**:
    - On success, the server upgrades to WebSocket and immediately sends a **`welcome`** message (just ignore it, it is an acknowledgement):
      ```json
      {
        "type": "welcome",
        "message": "Welcome to the quiz session!",
        "sessionId": "<sessionId>"
      }
      ```  
    - On failure (missing/invalid token), the connection is closed with an appropriate close code.

---

## 2.1 Receiving a New Question (Only Admin)

- **When**: After the admin triggers the next question (or when the session starts).
- **Response**: Server sends to admin a **`question`** message:
  ```json
  {
    "type": "question",
    "questionIdx": <one-based index of the question>,
    "questionsAmount": <total number of the questions in the quiz>,
    "text": "<question text>",
    "options": [
      { "text": "<option 1>", "is_correct": true/false },
      { "text": "<option 2>", "is_correct": true/false },
      …
    ]
  }

## 2.2 Receiving an acknowledgement next_question (Only Participants before 1st question)

- **When**: After the admin triggers the next question, participants receive this message.
- **Response**: Server broadcasts to participants a **`next_question`** message:
  ```json
  {
    "type": "next_question"
  }
---
### Attention: next question triggered at this moment.
### Therefore, at each new question starting from 2nd users firstly receive leaderboard / statistics, and then question payload / ack

---

## 2.4 (SKIP FOR QUESTION 1) Sending request to notify users (Only Admin)

- **When**: When Admin displayed leaderboard, and ready to show next question
    - **Request**: Admin sends to server a **`next_question`** message:
      ```json
      {
        "type": "next_question"
      }

---

## 2.4 Receiving a Leader Board (Only Admin)

- **When**: When the next question triggers, admin receives leaderboard.
  - **Response**: Server sends to admin a **`leaderboard`** message:
    ```json
    {
      "type": "leaderboard",
      "payload": {
          "session_code": "ABC123",
          "users": [
            {
              "user_id": "alice",
              "place": 1, // integer represents the current place in the rating
              "previous_place": 2, // integer represents the previous place in the rating
              "progress": "Up"/"Down"/"Same", // shows whether user's rank(place) increased/decreased/remained the same
              "total_score": 7
            },
            ...
          ]
        }
    }

## 2.5 Receiving a Question Statistics (Only Participants)

- **When**: When the next question triggers, participants receive following statistics.
    - **Response**: Server sends to admin a **`question_stat`** message:
      ```json
      {
        "type": "question_stat",
        "correct": true/false,
        "payload": {
            "session_code": "ABC123",
            "answers": {
                  "0": 8, // 8 people chose 0-th option
                  "1": 6, // 6 people chose 1-th option
                  "2": 4  // ...
           }
        },
        "options": [
          { "text": "<option 1>", "is_correct": true/false },
          { "text": "<option 2>", "is_correct": true/false },
          …
        ]  
      }
      
<div style="background-color: transparent; border-top: 4px solid red; padding: 0;">
</div>

- **Attention:** after these steps the websocket cycle goes to step [2.1](#21-receiving-a-new-question-only-admin
) after `next_question` trigger.

<div style="background-color: transparent; border-bottom: 4px solid red; padding: 0;">
</div>

## 3. Submitting an Answer (Participant Only)

- **When**: After receiving the `question` message by admin.
- **Request**: Send a WebSocket message with the chosen option index:

  ```json
  {
    "type": "answer",
    "option": <integer zero-based index>,
    "timestamp": <timestamp (in UTC) of user answer moment> "2025-07-17T12:34:56.789Z"
  }

## 3.1 Notification that one more user answered (Admin Only)

- **When**: After answer from the user.
  - **Response**: Send a WebSocket message with the chosen option index:

    ```json
    {
      "type": "user_answered",
      "payload": {
        "user_id": id,
      }
    }

## 4. Game End (Only Participants)

- **When**: After receiving triggering the `end_session` by admin.
- **Response**: Server sends to admin a **`game_end`** message:

  ```json
  {
    "type": "game_end"
  }
