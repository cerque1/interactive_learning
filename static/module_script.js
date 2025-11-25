window.addEventListener('DOMContentLoaded', () => {
    const token = localStorage.getItem('token');
    if (!token) {
      window.location.href = '/static/login.html';
      return;
    }
  
    const params = new URLSearchParams(window.location.search);
    const moduleId = params.get('module_id');
    if (!moduleId) {
      document.getElementById('module-name').textContent = 'Ошибка: не указан id модуля';
      return;
    }
  
    // Получить имя пользователя
    fetch('http://localhost:8080/api/v1/user/me?is_full=f', {
      headers: {
        'Authorization': `Bearer ${token}`
      }
    })
    .then(res => {
      if (res.status === 401) {
        window.location.href = '/static/login.html';
        return;
      }
      return res.json();
    })
    .then(userData => {
      if (userData) {
        const usernameElem = document.getElementById('username');
        usernameElem.textContent = userData.user.name || 'Пользователь';
        usernameElem.onclick = () => {
          window.location.href = '/static/profile.html';
        };
      }
    })
    .catch(() => {
      // Ошибка не критична для показа модуля
    });
  
    // Получить данные модуля
    fetch(`http://localhost:8080/api/v1/module/${moduleId}`, {
      headers: {
        'Authorization': `Bearer ${token}`
      }
    })
    .then(res => {
      if (res.status === 401) {
        window.location.href = '/static/login.html';
        return;
      }
      if (!res.ok) {
        throw new Error('Ошибка загрузки модуля');
      }
      return res.json();
    })
    .then(moduleData => {
      moduleData = moduleData.module
      const moduleNameElem = document.getElementById('module-name');
      const cardsContainer = document.getElementById('cards-container');
      const emptyMessage = document.getElementById('empty-message');
  
      moduleNameElem.textContent = moduleData.name || 'Без названия';
  
      if (!moduleData.cards || moduleData.cards.length === 0) {
        cardsContainer.innerHTML = '';
        emptyMessage.style.display = 'block';
        return;
      }
  
      emptyMessage.style.display = 'none';
      cardsContainer.innerHTML = '';
  
      moduleData.cards.forEach(card => {
        const cardElem = document.createElement('div');
        cardElem.className = 'card';
        cardElem.innerHTML = `
          <div class="card-title">${card.term.text}</div>
          <div>${card.definition.text}</div>
        `;
        cardsContainer.appendChild(cardElem);
      });
    })
    .catch(err => {
      document.getElementById('module-name').textContent = 'Ошибка загрузки модуля';
      document.getElementById('cards-container').innerHTML = '';
    });
  });
  