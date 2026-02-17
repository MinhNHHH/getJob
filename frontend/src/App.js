import logo from './logo.svg';
import './App.css';
import HomePage from './HomePage/homepage';
import WebRTCClient from './WebRTCClient';
import { BrowserRouter as Router, Routes, Route, Link } from 'react-router-dom';
import { AppBar, Toolbar, Typography, Button, Box } from '@mui/material';

function App() {
  return (
    <Router>
      <AppBar position="static">
        <Toolbar>
          <Typography variant="h6" component="div" sx={{ flexGrow: 1 }}>
            Get Job App
          </Typography>
          <Button color="inherit" component={Link} to="/">
            Home
          </Button>
          <Button color="inherit" component={Link} to="/webrtc">
            WebRTC Client
          </Button>
        </Toolbar>
      </AppBar>
      
      <Routes>
        <Route path="/" element={<HomePage />} />
        <Route path="/webrtc" element={<WebRTCClient />} />
      </Routes>
    </Router>
  );
}

export default App;
