import React from 'react';
import './App.css';
import {BrowserRouter as Router, Route, Routes} from 'react-router-dom';
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
import ProtectedRoute from './components/ProtectedRoute';

function App() {
  return (
    <Router>
      <Routes>
        <Route path="/register" element={<RegisterPage />} />
        <Route path="/login" element={<AuthPage/>} />
        {/* Catch-all route for 404 Not Found */}
        <Route path="*" element={<NotFoundPage />} /> 
        <Route path='/' element={<HomePage/>}/>
        
        <Route path='/create' element={
            <ProtectedRoute>
              <CreateSessionPage/>
            </ProtectedRoute>
          } 
        />
        <Route path='/enter-nickname' element={
            <ProtectedRoute>
              <EnterNicknamePage/>
            </ProtectedRoute>
          }
        />
        <Route path='/join' element={
            <ProtectedRoute>
              <JoinGamePage/>
            </ProtectedRoute>
          } 
        />
        <Route path="/ask-to-join/:sessionCode" element={
            <ProtectedRoute>
              <AskToJoinSessionPage/>
            </ProtectedRoute>
          }
        /> 
        <Route path="/sessionAdmin/:sessionCode" element={
            <ProtectedRoute> 
              <WaitGameStartAdmin/>
            </ProtectedRoute>
          } 
        />
        <Route path="/wait/:sessionCode" element={
            <ProtectedRoute>
              <WaitGameStartPlayer/>
            </ProtectedRoute>
          } 
        />
        <Route path="/game-process/:sessionCode" element={
            <ProtectedRoute>
              <GameController/>
            </ProtectedRoute>
          } 
        />
        <Route path="/game-controller/:sessionCode" element={
            <ProtectedRoute>
              <GameProcessAdmin/>
            </ProtectedRoute>
          } 
        />
      </Routes>
    </Router>
  );
}

export default App;
