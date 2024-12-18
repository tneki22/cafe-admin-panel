import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer } from 'recharts';
import '../styles/AnalyticsPage.css';

function AnalyticsPage() {
  const navigate = useNavigate();
  const [period, setPeriod] = useState('day');
  const [revenueData, setRevenueData] = useState([]);
  const [orderCountsData, setOrderCountsData] = useState([]);

  const goBack = () => {
    navigate('/main');
  };

  const fetchRevenue = async (selectedPeriod) => {
    try {
      const token = localStorage.getItem('token');
      const response = await fetch(`http://127.0.0.1:8885/api/revenue?period=${selectedPeriod}`, {
        headers: {
          'Authorization': `Bearer ${token}`,
        },
      });

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.message || 'Ошибка загрузки данных выручки');
      }

      const data = await response.json();
      setRevenueData(data);
    } catch (error) {
      console.error('Ошибка загрузки выручки:', error);
    }
  };

  const fetchOrderCounts = async (selectedPeriod) => {
    try {
      const token = localStorage.getItem('token');
      const response = await fetch(`http://127.0.0.1:8885/api/order_counts?period=${selectedPeriod}`, {
        headers: {
          'Authorization': `Bearer ${token}`,
        },
      });

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.message || 'Ошибка загрузки данных количества заказов');
      }

      const data = await response.json();
      setOrderCountsData(data);
    } catch (error) {
      console.error('Ошибка загрузки количества заказов:', error);
    }
  };

  useEffect(() => {
    fetchRevenue(period);
    fetchOrderCounts(period);
  }, [period]);

  const handlePeriodChange = (e) => {
    setPeriod(e.target.value);
  };

  return (
    <div className="analytics-page">
      <header>
        <button className="back-button" onClick={goBack}>← Назад</button>
        <h1>Аналитика</h1>
      </header>
      <div className="content">
        <div className="controls">
          <label htmlFor="period">Выберите период:</label>
          <select id="period" value={period} onChange={handlePeriodChange}>
            <option value="day">Последний день</option>
            <option value="week">Последняя неделя</option>
            <option value="month">Последний месяц</option>
            <option value="year">Последний год</option>
          </select>
        </div>
        {/* График выручки */}
        <div className="chart-placeholder">
          <h2>Выручка Руб</h2>
          <ResponsiveContainer width="100%" height={400}>
              <BarChart data={revenueData}>
                <CartesianGrid strokeDasharray="3 3" />
                <XAxis dataKey="time_unit" />
                <YAxis tickFormatter={(value) => value.toFixed(1)} />
                <Tooltip formatter={(value) => value.toFixed(1)} />
              <Bar dataKey="total" fill="#8884d8" />
      </BarChart>
          </ResponsiveContainer>
        </div>
        {/* График количества заказов */}
        <div className="chart-placeholder">
          <h2>Количество заказов</h2>
          <ResponsiveContainer width="100%" height={400}>
            <BarChart data={orderCountsData}>
              <CartesianGrid strokeDasharray="3 3" />
              <XAxis dataKey="time_unit" />
              <YAxis />
              <Tooltip />
              <Bar dataKey="count" fill="#82ca9d" />
            </BarChart>
          </ResponsiveContainer>
        </div>
      </div>
    </div>
  );
}

export default AnalyticsPage;