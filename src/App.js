import React from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import LoginPage from './components/LoginPage';
import MainPage from './components/MainPage';
import AnalyticsPage from './components/AnalyticsPage';
import MenuPage from './components/MenuPage';
import OrdersPage from './components/OrdersPage';

function App() {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<LoginPage />} />
        <Route path="/main" element={<MainPage />} />
        <Route path="/analytics" element={<AnalyticsPage />} />
        <Route path="/menu" element={<MenuPage />} />
        <Route path="/orders" element={<OrdersPage />} />
      </Routes>
    </Router>
  );
}

export default App;