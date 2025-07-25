import React from 'react';
import { BrowserRouter as Router, Route, Routes } from 'react-router-dom';
import './App.css';

import HomePage from './components/HomePage';
import CreateSessionPage from './components/CreateSessionPage';
import AskToJoinSessionPage from './components/AskToJoinSessionPage';
import WaitGameStartAdmin from './components/WaitGameStartAdmin';
import WaitGameStartPlayer from './components/WaitGameStartPage';
import JoinGamePage from './components/JoinGamePage';
import EnterNicknamePage from './components/EnterNicknamePage';
import GameController from './components/GameController';
import GameProcessAdmin from './components/GameProcessAdmin';
import AuthPage from './components/AuthPage';
import NotFoundPage from './components/NotFoundPage';
import RegisterPage from './components/RegisterPage';
import QuizStorePage from './components/QuizStorePage';
import FinalAdminPage from './components/FinalAdminPage';

import ProtectedRoute from './components/ProtectedRoute';

import { SessionWebSocketProvider } from './contexts/SessionWebSocketContext';
import { RealtimeWebSocketProvider } from './contexts/RealtimeWebSocketContext';
import ConstructorPage from './components/ConstructorPage';
import FloatingBackground from './components/background/FloatingBackground';

import WithMusicLayout from './components/WithMusicLayout'

function App() {
  return (
    <div className='app'>
    <FloatingBackground />
    <Router>
      <Routes>
        {/* Страницы без WebSocket */}
        <Route path="/" element={<HomePage />} />
        <Route path="/register" element={<RegisterPage />} />
        <Route path="/login" element={<AuthPage />} />
        <Route path="/create" element={<CreateSessionPage />} />
        <Route path="/enter-nickname" element={<EnterNicknamePage />} />
        <Route path="/join" element={<JoinGamePage />} />
        <Route path="/ask-to-join/:sessionCode" element={<AskToJoinSessionPage />} />
        <Route path="/final" element={<FinalAdminPage />} />
        <Route path="*" element={<NotFoundPage />} />
        <Route
          path="/store"
          element={
              <QuizStorePage/>
          }
        />

        {/* Protected страницы */}
        <Route
          path="/constructor/:quizID/"
          element={
            <ProtectedRoute>
              <ConstructorPage />
            </ProtectedRoute>
          } />

        {/* Страницы с WebSocket */}
        
        <Route
          path="/game-process/:sessionCode"
          element={
            <RealtimeWebSocketProvider>
              <SessionWebSocketProvider>
                <GameController />
              </SessionWebSocketProvider>
            </RealtimeWebSocketProvider>
          }
           />
        <Route
          path="/wait/:sessionCode" 
          element={
            <RealtimeWebSocketProvider>
              <SessionWebSocketProvider>
                <WaitGameStartPlayer />
              </SessionWebSocketProvider>
            </RealtimeWebSocketProvider>
          }
        />
        <Route element={<WithMusicLayout/>}>
          <Route
            path="/sessionAdmin/:sessionCode"
            element={
              <SessionWebSocketProvider>
                <RealtimeWebSocketProvider>
                  <WaitGameStartAdmin />
                </RealtimeWebSocketProvider>
              </SessionWebSocketProvider>
            }
          />
          <Route
            path="/game-controller/:sessionCode"
            element={
              <SessionWebSocketProvider>
                <RealtimeWebSocketProvider>
                  <GameProcessAdmin />
                </RealtimeWebSocketProvider>
              </SessionWebSocketProvider>
            }
          />
        </Route>
      </Routes>
    </Router>
  </div>
  );
}

export default App;
