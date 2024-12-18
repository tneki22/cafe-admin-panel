import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import '../styles/OrdersPage.css';

function OrdersPage() {
  const [menuItems, setMenuItems] = useState([]);
  const [orders, setOrders] = useState([]);
  const [newOrder, setNewOrder] = useState([{ menuItemId: '', quantity: '1' }]);
  const [isAdding, setIsAdding] = useState(false);
  const [error, setError] = useState('');
  const navigate = useNavigate();

  const goBack = () => {
    navigate('/main');
  };

  const fetchMenu = async () => {
    try {
      const token = localStorage.getItem('token');
      const response = await fetch('http://127.0.0.1:8885/api/menu', {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`,
        },
      });

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.message || 'Ошибка загрузки меню');
      }

      const data = await response.json();
      setMenuItems(data);
    } catch (error) {
      setError(error.message);
    }
  };

  const fetchOrders = async () => {
    try {
      const token = localStorage.getItem('token');
      const response = await fetch('http://127.0.0.1:8885/api/orders', {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`,
        },
      });

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.message || 'Ошибка загрузки заказов');
      }

      const data = await response.json();
      // Сортировка заказов по дате создания (новые сначала)
      const sortedData = [...data].sort((a, b) => new Date(b.created_at) - new Date(a.created_at));
      setOrders(sortedData);
    } catch (error) {
      setError(error.message);
    }
  };

  useEffect(() => {
    fetchMenu();
    fetchOrders();
  }, []);

  const handleAddOrderChange = (index, field, value) => {
    const updatedOrder = [...newOrder];
    if (field === 'quantity') {
      updatedOrder[index][field] = value;
    } else {
      updatedOrder[index][field] = value;
    }
    setNewOrder(updatedOrder);
  };

  const handleAddOrderItem = () => {
    setNewOrder([...newOrder, { menuItemId: '', quantity: '1' }]);
  };

  const handleRemoveOrderItem = (index) => {
    const updatedOrder = newOrder.filter((_, i) => i !== index);
    setNewOrder(updatedOrder);
  };

  const handleAddOrder = async () => {
    try {
      // Проверяем, что все товары выбраны и количество является целым числом >=1
      for (let item of newOrder) {
        if (!item.menuItemId) {
          throw new Error('Пожалуйста, выберите товар для всех позиций заказа.');
        }
        const qty = parseInt(item.quantity, 10);
        if (isNaN(qty) || qty < 1) {
          throw new Error('Количество товара должно быть целым числом и не менее 1.');
        }
      }

      const token = localStorage.getItem('token');
      const orderItems = newOrder.map(item => ({
        menuItemId: parseInt(item.menuItemId, 10),
        quantity: parseInt(item.quantity, 10),
      }));

      const response = await fetch('http://127.0.0.1:8885/api/orders', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`,
        },
        body: JSON.stringify({ items: orderItems }),
      });

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.message || 'Ошибка добавления заказа');
      }

      // После добавления заново загружаем заказы
      await fetchOrders();
      setNewOrder([{ menuItemId: '', quantity: '1' }]);
      setIsAdding(false);
    } catch (error) {
      setError(error.message);
    }
  };

  const updateStatus = async (orderID, status) => {
    try {
      const token = localStorage.getItem('token');
      const response = await fetch(`http://127.0.0.1:8885/api/orders/${orderID}/status`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`,
        },
        body: JSON.stringify({ status }),
      });
      if (response.ok) {
        await fetchOrders();
      } else {
        const errorData = await response.json();
        throw new Error(errorData.message || 'Ошибка обновления статуса');
      }
    } catch (error) {
      console.error('Ошибка обновления статуса:', error);
    }
  };

  const formatDate = (dateString) => {
    const options = {
      year: 'numeric', month: '2-digit', day: '2-digit',
      hour: '2-digit', minute: '2-digit',
      hour12: false,
    };
    const date = new Date(dateString);
    return date.toLocaleString('ru-RU', options).replace(',', '');
  };

  return (
    <div className="orders-page">
      <header>
        <button className="back-button" onClick={goBack}>← Назад</button>
        <h1>Заказы</h1>
        <button className="add-button" onClick={() => setIsAdding(true)}>Добавить</button>
      </header>
      <div className="content">
        {error && <div className="error-message">{error}</div>}
        {isAdding && (
          <div className="add-order-form">
            {newOrder.map((item, index) => (
              <div key={index} className="order-item">
                <select
                  value={item.menuItemId}
                  onChange={(e) => handleAddOrderChange(index, 'menuItemId', e.target.value)}
                >
                  <option value="">Выберите товар</option>
                  {menuItems.map(menuItem => (
                    <option key={menuItem.id} value={menuItem.id}>
                      {menuItem.name}
                    </option>
                  ))}
                </select>
                <input
                  type="number"
                  value={item.quantity}
                  onChange={(e) => handleAddOrderChange(index, 'quantity', e.target.value)}
                />
                <button onClick={() => handleRemoveOrderItem(index)}>Удалить</button>
              </div>
            ))}
            <button onClick={handleAddOrderItem}>Добавить товар</button>
            <button onClick={handleAddOrder}>Сохранить</button>
            <button onClick={() => setIsAdding(false)}>Отмена</button>
          </div>
        )}
        <table className="table">
          <thead>
            <tr>
              <th>Номер заказа</th>
              <th>Сумма</th>
              <th>Статус</th>
              <th>Дата создания</th>
              <th>Действия</th>
            </tr>
          </thead>
          <tbody>
            {orders.map((order, index) => (
              <tr key={order.id}>
                <td>{String(index + 1).padStart(2, '0')}</td>
                <td>{order.total.toFixed(1)} ₽</td>
                <td>{order.status}</td>
                <td>{formatDate(order.created_at)}</td>
                <td>
                  {order.status !== "Выполнен" && (
                    <button onClick={() => updateStatus(order.id, "Выполнен")}>Выполнен</button>
                  )}
                  {order.status !== "Отменен" && (
                    <button onClick={() => updateStatus(order.id, "Отменен")}>Отменен</button>
                  )}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}

export default OrdersPage;