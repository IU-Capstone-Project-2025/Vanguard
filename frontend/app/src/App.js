import React from 'react';
import './App.css';
import {BrowserRouter as Router, Route, Routes} from 'react-router-dom';
import HomePage from './components/HomePage';
import CreateSessionPage from './components/CreateSessionPage';
import AskToJoinSessionPage from './components/AskToJoinSessionPage';

function App() {
  return (
    <Router>
      <Routes>
        <Route path='/' element={<HomePage/>}/>
        <Route path='/create' element={<CreateSessionPage/>} />
        <Route path="/ask-to-join/:sessionCode" element={<AskToJoinSessionPage/>} /> 
      </Routes>
    </Router>
  );
}

export default App;
