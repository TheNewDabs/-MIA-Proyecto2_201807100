import { BrowserRouter, Route, Routes } from 'react-router-dom';
import Main from './components/main';
import './App.css';

function App() {
  return (
    <div className="App">
      <BrowserRouter>
        <Routes>
          <Route path="/" element={<Main></Main>}></Route>
        </Routes>
      </BrowserRouter>
    </div>
  );
}

export default App;
