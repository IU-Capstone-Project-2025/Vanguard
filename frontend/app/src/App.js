import React from 'react';
import './App.css';
import {BrowserRouter as Router, Route, Routes} from 'react-router-dom';
import HomePage from './components/HomePage';
import CreateSessionPage from './components/CreateSessionPage';
import JoinGamePage from './components/JoinGamePage';
import PlayGamePage from './components/PlayGamePage';

function App() {
  return (
    <Router>
      <Routes>
        <Route path='/' element={<HomePage/>}/>
        <Route path='/create' element={<CreateSessionPage/>} />
        <Route path='/join' element={<JoinGamePage/>} />
        <Route path='/play' element={<PlayGamePage/>} />
      </Routes>
    </Router>
  );
}

export default App;
