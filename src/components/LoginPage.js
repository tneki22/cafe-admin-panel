import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import '../styles/LoginPage.css';

function LoginPage() {
  const [name, setName] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  
  const navigate = useNavigate();

  const handleLogin = async () => {
    try {
      const response = await fetch('http://127.0.0.1:8885/api/login', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ email, password }),
      });

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.message || 'Ошибка входа');
      }

      const data = await response.json();
      localStorage.setItem('token', data.token);
      
      // Сохранение имени пользователя из ответа сервера, если оно присутствует
      if (data.username) {
        localStorage.setItem('username', data.username);
      } else {
        // Если имя не возвращается, используем введённое имя
        localStorage.setItem('username', name);
      }

      navigate('/main');
    } catch (error) {
      setError(error.message);
    }
  };

  const handleRegister = async () => {
    try {
      const response = await fetch('http://127.0.0.1:8885/api/register', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ name, email, password }),
      });

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.message || 'Ошибка регистрации');
      }

      const data = await response.json();
      localStorage.setItem('token', data.token);
      localStorage.setItem('username', name); // Сохраняем имя пользователя

      navigate('/main');
    } catch (error) {
      setError(error.message);
    }
  };

  return (
    <div className="login-page">
      <h1>Вход/Регистрация</h1>
      <div className="login-form">
        <input 
          type="text" 
          placeholder="Имя" 
          value={name}
          onChange={(e) => setName(e.target.value)}
        />
        <input 
          type="email" 
          placeholder="Email" 
          value={email}
          onChange={(e) => setEmail(e.target.value)}
        />
        <input 
          type="password" 
          placeholder="Пароль"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
        />
        {error && <div className="error-message">{error}</div>}
        <div className="login-buttons">
          <button onClick={handleRegister}>Регистрация</button>
          <button onClick={handleLogin}>Вход</button>
        </div>
      </div>
    </div>
  );
}

export default LoginPage;