import React, { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import '../styles/MainPage.css';

function MainPage() {
  const [dateTime, setDateTime] = useState(new Date());
  const navigate = useNavigate();
  


  useEffect(() => {
    const timer = setInterval(() => setDateTime(new Date()), 1000);
    return () => clearInterval(timer);
  }, []);

  const goToAnalytics = () => {
    navigate('/analytics');
  };

  const goToMenu = () => {
    navigate('/menu');
  };

  const goToOrders = () => {
    navigate('/orders');
  };

  const handleLogout = () => {
    // Очистка данных при выходе
    localStorage.removeItem('token');
    localStorage.removeItem('username');
    navigate('/');
  };

  return (
    <div className="main-page">
      <header>
        <button className="logout-button" onClick={handleLogout}>Выйти</button>
        <div className="header-info">
          <h1>Добро пожаловать в панель управления кафе!</h1>
          <div className="date-time">{dateTime.toLocaleString()}</div>
        </div>
      </header>
      <div className="main-buttons">
        <button onClick={goToMenu}>Меню</button>
        <button onClick={goToOrders}>Заказы</button>
        <button onClick={goToAnalytics}>Аналитика</button>
      </div>
    </div>
  );
}

export default MainPage;