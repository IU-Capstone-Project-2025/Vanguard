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

## 2.2 Receiving an acknowledgement that game was started (Only Participants before 1st question)

- **When**: After the admin triggers the next question for first time and receives question payload, participants receive this message.
- **Response**: Server broadcasts to participants a **`game_start`** message:
  ```json
  {
    "type": "game_start",
    "isGameStarted": true,
  }

## 3. Submitting an Answer (Participant Only)

- **When**: After receiving the `question` message by admin.
- **Request**: Send a WebSocket message with the chosen option index:

  ```json
  {
    "option": <integer zero-based index>
  }

## 4. Receiving final leaderBoard after game (All users)
### *P.s. now just an map of arrays, that show which answers were correct*

- **When**: After triggering `session end` event.
- **Response**: Server broadcasts to participants a **`leaderboard`** message:

  ```json
  {
    "type": "leaderboard",
    "payload": {
                  "<user_id_1>": [true/false, true/false, ...],
                  ...
               },
  }
  
- e.g. game has 5 questions, then the array of user's answers will contain of 5 `true/false` elements, in the order of the questions appearing in the game.
- `[true, true, false, true, false]` means that answer for the questions 1, 2, 4 was correct, and for questions 3, 5 - incorrect.
