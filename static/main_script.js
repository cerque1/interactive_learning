window.addEventListener('DOMContentLoaded', () => {
  const token = localStorage.getItem('token');
  if (!token) {
    window.location.href = '/static/login.html';
    return;
  }
  
  fetch('http://localhost:8080/api/v1/user/me?is_full=t', {
    headers: {
      'Authorization': `Bearer ${token}`
    }
  })
  .then(async res => {
    if (res.status === 401) {
      window.location.href = '/static/login.html';
      return;
    }
    if (!res.ok) {
      throw new Error('Ошибка загрузки данных');
    }
    const data = await res.json();
    // Имя пользователя
    const usernameElem = document.getElementById('username');
    usernameElem.textContent = data.user.name || 'Пользователь';
    usernameElem.onclick = () => {
      window.location.href = '/static/profile.html';
    };

    // Модули пользователя
    const modulesContainer = document.getElementById('modules-container');
    const modulesEmpty = document.getElementById('modules-empty');
    modulesContainer.innerHTML = '';
    if (!data.user.modules || data.user.modules.length === 0) {
      modulesEmpty.style.display = 'block';
    } else {
      modulesEmpty.style.display = 'none';
      data.user.modules.forEach(module => {
        const cardCount = (module.cards && module.cards.length) || 0;
        const card = document.createElement('div');
        card.className = 'card';
        card.innerHTML = `<div class="card-title">${module.name}</div><div class="card-count">Карточек: ${cardCount}</div>`;
        card.style.cursor = 'pointer';
        card.onclick = () => {
          window.location.href = `/static/module.html?module_id=${encodeURIComponent(module.id)}`;
        };
        modulesContainer.appendChild(card);
      });
    }

    // Категории
    const categoriesContainer = document.getElementById('categories-container');
    const categoriesEmpty = document.getElementById('categories-empty');
    categoriesContainer.innerHTML = '';
    if (!data.user.categories || data.user.categories.length === 0) {
      categoriesEmpty.style.display = 'block';
    } else {
      categoriesEmpty.style.display = 'none';
      data.user.categories.forEach(category => {
        const modulesCount = (category.modules && category.modules.length) || 0;
        const card = document.createElement('div');
        card.className = 'card';
        card.innerHTML = `<div class="card-title">${category.name}</div><div class="card-count">Модулей: ${modulesCount}</div>`;
        card.style.cursor = 'pointer';
        card.onclick = () => {
          window.location.href = `/static/category.html?category_id=${encodeURIComponent(category.id)}`;
        };
        categoriesContainer.appendChild(card);
      });
    }

    // Результаты - статично
  })
  .catch(() => {
    document.body.innerHTML = '<p style="text-align:center; margin-top:50px; font-size:18px; color:#c75c5c;">Ошибка загрузки данных. Попробуйте перезагрузить страницу.</p>';
  });
});

const modulesHeader = document.getElementById('modules-header');
modulesHeader.addEventListener('click', function() {
  window.location.href = "/static/modules.html";
});

const navToggle = document.getElementById('nav-toggle');
const navPanel = document.getElementById('nav-panel');

navToggle.addEventListener('click', function() {
  navPanel.classList.toggle('open');
  navToggle.classList.toggle('open');
});

const navModulesBut = document.getElementById('modules-btn');
navModulesBut.addEventListener('click', function() {
  window.location.href = "/static/modules.html";
});