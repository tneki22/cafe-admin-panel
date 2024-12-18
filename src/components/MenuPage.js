import React, { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import '../styles/MenuPage.css';

function MenuPage() {
  const [menuItems, setMenuItems] = useState([]);
  const [error, setError] = useState('');
  const [editingItem, setEditingItem] = useState(null);
  const [newItem, setNewItem] = useState({ name: '', description: '', price: '' });
  const [isAdding, setIsAdding] = useState(false);
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

  useEffect(() => {
    fetchMenu();
  }, []);

  const handleDelete = async (id) => {
    if (!window.confirm('Вы точно хотите удалить этот товар?')) return;

    try {
      const token = localStorage.getItem('token');
      const response = await fetch(`http://127.0.0.1:8885/api/menu/${id}`, {
        method: 'DELETE',
        headers: {
          'Authorization': `Bearer ${token}`,
        },
      });

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.message || 'Ошибка удаления элемента');
      }

      // После удаления заново загружаем меню
      await fetchMenu();
    } catch (error) {
      setError(error.message);
    }
  };

  const handleEditClick = (item) => {
    setEditingItem(item);
  };

  const handleCancelEdit = () => {
    setEditingItem(null);
  };

  const handleSaveEdit = async () => {
    try {
      const token = localStorage.getItem('token');
      const updatedData = {
        name: editingItem.name,
        description: editingItem.description,
        price: parseFloat(editingItem.price),
      };

      const response = await fetch(`http://127.0.0.1:8885/api/menu/${editingItem.id}`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`,
        },
        body: JSON.stringify(updatedData),
      });

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.message || 'Ошибка обновления элемента');
      }

      // После обновления заново загружаем меню
      await fetchMenu();
      setEditingItem(null);
    } catch (error) {
      setError(error.message);
    }
  };

  const handleChange = (e) => {
    const { name, value } = e.target;
    setEditingItem((prev) => ({
      ...prev,
      [name]: name === 'price' ? parseFloat(value) : value,
    }));
  };

  const handleAddChange = (e) => {
    const { name, value } = e.target;
    setNewItem((prev) => ({
      ...prev,
      [name]: name === 'price' ? parseFloat(value) : value,
    }));
  };

  const handleAddItem = async () => {
    try {
      const token = localStorage.getItem('token');
      const response = await fetch('http://127.0.0.1:8885/api/menu', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`,
        },
        body: JSON.stringify(newItem),
      });

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.message || 'Ошибка добавления элемента');
      }

      // После добавления заново загружаем меню
      await fetchMenu();
      setNewItem({ name: '', description: '', price: '' });
      setIsAdding(false);
    } catch (error) {
      setError(error.message);
    }
  };

  return (
    <div className="menu-page">
      <header>
        <button className="back-button" onClick={goBack}>← Назад</button>
        <h1>Меню</h1>
        <button className="add-button" onClick={() => setIsAdding(true)}>Добавить</button>
      </header>
      <div className="content">
        {error && <div className="error-message">{error}</div>}
        {isAdding && (
          <div className="add-form">
            <input
              name="name"
              placeholder="Название"
              value={newItem.name}
              onChange={handleAddChange}
            />
            <input
              name="description"
              placeholder="Описание"
              value={newItem.description}
              onChange={handleAddChange}
            />
            <input
              name="price"
              type="number"
              placeholder="Цена"
              value={newItem.price}
              onChange={handleAddChange}
            />
            <button onClick={handleAddItem}>Сохранить</button>
            <button onClick={() => setIsAdding(false)}>Отмена</button>
          </div>
        )}
        <div className="table">
          <table>
            <thead>
              <tr>
                <th>№</th>
                <th>Товар</th>
                <th>Описание</th>
                <th>Цена</th>
                <th>Дата создания</th>
                <th>Действия</th>
              </tr>
            </thead>
            <tbody>
              {menuItems.length > 0 ? (
                menuItems.map((item, index) => (
                  <tr key={item.id}>
                    <td>{index + 1}</td>
                    <td>
                      {editingItem && editingItem.id === item.id ? (
                        <input
                          name="name"
                          value={editingItem.name}
                          onChange={handleChange}
                        />
                      ) : (
                        item.name
                      )}
                    </td>
                    <td>
                      {editingItem && editingItem.id === item.id ? (
                        <input
                          name="description"
                          value={editingItem.description}
                          onChange={handleChange}
                        />
                      ) : (
                        item.description
                      )}
                    </td>
                    <td>
                      {editingItem && editingItem.id === item.id ? (
                        <input
                          name="price"
                          type="number"
                          value={editingItem.price}
                          onChange={handleChange}
                        />
                      ) : (
                        item.price
                      )}
                    </td>
                    <td>{item.created_at}</td>
                    <td>
                      {editingItem && editingItem.id === item.id ? (
                        <>
                          <button onClick={handleSaveEdit}>Сохранить</button>
                          <button onClick={handleCancelEdit}>Отмена</button>
                        </>
                      ) : (
                        <>
                          <button onClick={() => handleEditClick(item)}>Изменить</button>
                          <button onClick={() => handleDelete(item.id)}>Удалить</button>
                        </>
                      )}
                    </td>
                  </tr>
                ))
              ) : (
                <tr>
                  <td colSpan="6">Меню пусто</td>
                </tr>
              )}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
}

export default MenuPage;