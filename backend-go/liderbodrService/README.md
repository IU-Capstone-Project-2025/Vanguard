link to diagram https://www.mermaidchart.com/app/projects/b6fc6c6e-5ef3-428b-bc09-cbd1b7b3e539/diagrams/5bbb9939-c088-4b25-941d-cb75b3b5c217/share/invite/eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJkb2N1bWVudElEIjoiNWJiYjk5MzktYzA4OC00YjI1LTk0MWQtY2I3NWIzYjVjMjE3IiwiYWNjZXNzIjoiRWRpdCIsImlhdCI6MTc1MTI4OTc1OX0.4fldTWXMVotRlmTeL1nd-XH3S5pIwNXdabroN9r1DNs

In Memory structures

type ScoreEntry = {
    userId: string;
    totalScore: number;
};

type LeaderboardSession = Map<string, ScoreEntry>;

const leaderboardMap: Map<string, LeaderboardSession> = new Map();

request example
json:
{
    "sessionCode": "ABC123",
    "answers": 
    [
        { "userId": "user42", "AnswerTimestamp": 12332132, "QuestionStartTimestamp" : 1231321321, "IsCorrect" : True},
        { "userId": "user42", "AnswerTimestamp": 12332132, "QuestionStartTimestamp" : 1231321321, "IsCorrect" : True},
        { "userId": "user42", "AnswerTimestamp": 12332132, "QuestionStartTimestamp" : 1231321321, "IsCorrect" : True}
    ]
}

response example
{
    "leaderboard": 
    [
        { "userId": "user42", "score": 3240 },
        { "userId": "user17", "score": 3120 },
        { "userId": "user88", "score": 2790 }
    ]
}

wsService collect answers from users. If user didn't give an answer, set IsCorrect as false
when time for question is end, send all results to LeaderBord and get results
In this system heavy load falls on the WsService

Second variant
diagram link https://www.mermaidchart.com/app/projects/b6fc6c6e-5ef3-428b-bc09-cbd1b7b3e539/diagrams/395f3034-1efb-45a3-aa21-deedf6a5a795/share/invite/eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJkb2N1bWVudElEIjoiMzk1ZjMwMzQtMWVmYi00NWEzLWFhMjEtZGVlZGY2YTVhNzk1IiwiYWNjZXNzIjoiRWRpdCIsImlhdCI6MTc1MTIwNDIyMH0.Vbs6ap7Sco5SLNN71sRRp04caLU33k0fLmC8BHOM7Ao

request example
json:
{
    "sessionCode": "ABC123",
    "questionId": "q1",
    "userId": "user42",
    "isCorrect": true,
    "answerTimestamp": 172398374,         // Время нажатия на ответ (на клиенте)
    "questionStartTimestamp": 172398370   // Время начала вопроса (от ведущего/сервера)
}

response example
{
    "leaderboard": 
    [
        { "userId": "user42", "score": 3240 },
        { "userId": "user17", "score": 3120 },
        { "userId": "user88", "score": 2790 }
    ]
}

in this system wsService do not need to save all answers, it just sends it to LeaderbordService


