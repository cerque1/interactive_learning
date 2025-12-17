window.addEventListener('DOMContentLoaded', () => {
  const token = localStorage.getItem('token');
  if (!token) {
    window.location.href = '/static/login.html?redirect=' + encodeURIComponent(window.location.href);
    return;
  }

  const params = new URLSearchParams(window.location.search);
  const categoryId = params.get('category_id');
  if (!categoryId) {
    document.getElementById('category-name').textContent = 'Ошибка: не указан id категории';
    return;
  }

  // Получить имя пользователя
  fetch('http://localhost:8080/api/v1/user/me?is_full=f', {
    headers: { 'Authorization': `Bearer ${token}` }
  })
  .then(res => {
    if (res.status === 401) {
      window.location.href = '/static/login.html?redirect=' + encodeURIComponent(window.location.href);
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
  .catch(() => {});

  // Получить данные категории
  fetch(`http://localhost:8080/api/v1/category/${categoryId}?is_full=t`, {
    headers: { 'Authorization': `Bearer ${token}` }
  })
  .then(res => {
    if (res.status === 401) {
      window.location.href = '/static/login.html?redirect=' + encodeURIComponent(window.location.href);
      return;
    }
    if (!res.ok) {
      throw new Error('Ошибка загрузки категории');
    }
    return res.json();
  })
  .then(categoryData => {
    categoryData = categoryData.category
    const categoryNameElem = document.getElementById('category-name');
    const modulesContainer = document.getElementById('modules-container');
    const emptyMessage = document.getElementById('empty-message');

    categoryNameElem.textContent = categoryData.name || 'Без названия';

    if (!categoryData.modules || categoryData.modules.length === 0) {
      modulesContainer.innerHTML = '';
      emptyMessage.style.display = 'block';
      return;
    }

    emptyMessage.style.display = 'none';
    modulesContainer.innerHTML = '';

    categoryData.modules.forEach(module => {
      const cardCount = (module.cards && module.cards.length) || 0;
      const cardElem = document.createElement('div');
      cardElem.className = 'card';
      cardElem.innerHTML = `
        <div class="card-title">${module.name}</div>
        <div class="card-count">Карточек: ${cardCount}</div>
      `;
      cardElem.style.cursor = 'pointer';
      cardElem.onclick = () => {
        window.location.href = `/static/module.html?module_id=${encodeURIComponent(module.id)}`;
      };
      modulesContainer.appendChild(cardElem);
    });
  })
  .catch((err) => {
    document.getElementById('category-name').textContent = 'Ошибка загрузки категории';
    document.getElementById('modules-container').innerHTML = '';
  });
});

const navToggle = document.getElementById('nav-toggle');
const navPanel = document.getElementById('nav-panel');

const navMainBut = document.getElementById('main-btn');

navToggle.addEventListener('click', function() {
  navPanel.classList.toggle('open');
  navToggle.classList.toggle('open');
});

navMainBut.addEventListener('click', function() {
  window.location.href = '/static/main.html';
});

const navModulesBut = document.getElementById('modules-btn');
navModulesBut.addEventListener('click', function() {
  window.location.href = "/static/modules.html";
});

const navCategoriesBut = document.getElementById('categories-btn');
navCategoriesBut.addEventListener('click', function() {
  window.location.href = '/static/categories.html';
});

document.getElementById('head').addEventListener('click', () => {
  window.location.href = '/static/main.html';
});