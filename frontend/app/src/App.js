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

import NotFoundPage from './components/NotFoundPage';

function App() {
  return (
    <Router>
      <Routes>
        <Route path='/' element={<HomePage/>}/>
        <Route path='/create' element={<CreateSessionPage/>} />
        <Route path='/enter-nickname' element={<EnterNicknamePage/>} />
        <Route path='/join' element={<JoinGamePage/>} />
        <Route path="/ask-to-join/:sessionCode" element={<AskToJoinSessionPage/>} /> 
        <Route path="/sessionAdmin/:sessionCode" element={<WaitGameStartAdmin/>} />
        <Route path="/wait/:sessionCode" element={<WaitGameStartPlayer/>} />
        <Route path="/game-process/:sessionCode" element={<GameController/>} />
        
        {/* Catch-all route for 404 Not Found */}
        <Route path="*" element={<NotFoundPage />} /> 
        
      </Routes>
    </Router>
  );
}

export default App;
