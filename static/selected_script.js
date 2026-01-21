// selected_script.js
const API_BASE_URL = window.location.origin;
let currentModules = [];
let currentCategories = [];
let currentTab = 'modules';

function loadUserName(token) {
  return fetch(`${API_BASE_URL}/api/v1/user/me?is_full=f`, {
    headers: { 'Authorization': `Bearer ${token}` }
  })
  .then(res => {
    if (res.status === 401) {
      window.location.href = `/static/login.html?redirect=${encodeURIComponent(window.location.href)}`;
      return null;
    }
    return res.json();
  })
  .then(userData => {
    if (userData && userData.user) {
      const usernameElem = document.getElementById('username');
      if (usernameElem) {
        usernameElem.textContent = userData.user.name || 'Пользователь';
        usernameElem.style.cursor = 'pointer';
        usernameElem.onclick = () => {
          window.location.href = '/static/profile.html';
        };
      }
    }
    return userData;
  })
  .catch(() => null);
}

function createFavoriteCard(item, type) {
  const card = document.createElement('div');
  card.className = 'card';
  card.dataset.itemId = item.id;
  card.dataset.itemType = type;
  
  const countText = item.card_count !== undefined ? `${item.card_count} карт` : '';
  
  card.innerHTML = `
    <button class="star-action" title="Убрать из избранного">⭐</button>
    <h3 class="card-title">${item.name}</h3>
    ${countText ? `<p class="card-count">${countText}</p>` : ''}
  `;
  
  // Кнопка удаления из избранного
  const starBtn = card.querySelector('.star-action');
  starBtn.addEventListener('click', (e) => {
    e.stopPropagation();
    handleRemoveFromFavorites(item.id, type, card);
  });
  
  // Переход по клику на карточку
  card.addEventListener('click', (e) => {
    if (!e.target.closest('.star-action')) {
      if (type === 'module') {
        window.location.href = `/static/module.html?module_id=${item.id}`;
      } else {
        window.location.href = `/static/category.html?category_id=${item.id}`;
      }
    }
  });
  
  return card;
}

function handleRemoveFromFavorites(itemId, type, cardElement) {
  if (!confirm(`Убрать "${cardElement.querySelector('.card-title').textContent}" из избранного?`)) {
    return;
  }
  
  const token = localStorage.getItem('token');
  const url = `${API_BASE_URL}/api/v1/selected/${type}s/delete?${type}_id=${itemId}`;
  
  fetch(url, {
    method: 'DELETE',
    headers: { 'Authorization': `Bearer ${token}` }
  })
  .then(res => {
    if (res.status === 401) {
      window.location.href = `/static/login.html?redirect=${encodeURIComponent(window.location.href)}`;
      return;
    }
    if (!res.ok) throw new Error(`Ошибка удаления: ${res.status}`);
    return res.json();
  })
  .then(() => {
    // Анимация удаления
    cardElement.style.transition = 'all 0.3s ease';
    cardElement.style.opacity = '0';
    cardElement.style.transform = 'scale(0.8) translateX(-20px)';
    
    setTimeout(() => {
      cardElement.remove();
      if (type === 'module') {
        currentModules = currentModules.filter(c => c !== cardElement);
      } else {
        currentCategories = currentCategories.filter(c => c !== cardElement);
      }
      checkEmptyState();
    }, 300);
  })
  .catch(err => {
    alert('Ошибка при удалении из избранного: ' + err.message);
  });
}

function checkEmptyState() {
  const modulesContainer = document.getElementById('modules-container');
  const categoriesContainer = document.getElementById('categories-container');
  const emptyMsg = document.getElementById('empty-state');
  
  // Проверяем только активную вкладку
  const activeTabBtn = document.querySelector('.tab-btn.active');
  if (!activeTabBtn) return;
  
  const activeTab = activeTabBtn.dataset.tab;
  
  if (activeTab === 'modules') {
    // Показываем модули или скрываем контейнер
    if (currentModules.length > 0) {
      modulesContainer.style.display = 'flex';
    } else {
      modulesContainer.style.display = 'none';
    }
    categoriesContainer.style.display = 'none';
    emptyMsg.style.display = currentModules.length === 0 ? 'block' : 'none';
  } else {
    // Показываем категории или скрываем контейнер  
    if (currentCategories.length > 0) {
      categoriesContainer.style.display = 'flex';
    } else {
      categoriesContainer.style.display = 'none';
    }
    modulesContainer.style.display = 'none';
    emptyMsg.style.display = currentCategories.length === 0 ? 'block' : 'none';
  }
}

function loadModules(token) {
  const container = document.getElementById('modules-container');
  container.innerHTML = '';
  currentModules = [];
  
  fetch(`${API_BASE_URL}/api/v1/selected/modules/`, {
    headers: { 'Authorization': `Bearer ${token}` }
  })
  .then(res => {
    if (res.status === 401) {
      window.location.href = `/static/login.html?redirect=${encodeURIComponent(window.location.href)}`;
      return;
    }
    if (!res.ok) throw new Error('Ошибка загрузки модулей');
    return res.json();
  })
  .then(data => {
    if (data.selected_modules?.length) {
      data.selected_modules.forEach(module => {
        const card = createFavoriteCard(module, 'module');
        currentModules.push(card);
        container.appendChild(card);
      });
    }
    checkEmptyState();
  })
  .catch(() => {
    container.innerHTML = '<div style="text-align:center;color:#999;padding:40px;">Ошибка загрузки модулей</div>';
  });
}

function loadCategories(token) {
  const container = document.getElementById('categories-container');
  container.innerHTML = '';
  currentCategories = [];
  
  fetch(`${API_BASE_URL}/api/v1/selected/categories/`, {
    headers: { 'Authorization': `Bearer ${token}` }
  })
  .then(res => {
    if (res.status === 401) {
      window.location.href = `/static/login.html?redirect=${encodeURIComponent(window.location.href)}`;
      return;
    }
    if (!res.ok) throw new Error('Ошибка загрузки категорий');
    return res.json();
  })
  .then(data => {
    if (data.selected_categories?.length) {
      data.selected_categories.forEach(category => {
        const card = createFavoriteCard(category, 'category');
        currentCategories.push(card);
        container.appendChild(card);
      });
    }
    checkEmptyState();
  })
  .catch(() => {
    container.innerHTML = '<div style="text-align:center;color:#999;padding:40px;">Ошибка загрузки категорий</div>';
  });
}

function setupTabs() {
  const tabBtns = document.querySelectorAll('.tab-btn');
  
  tabBtns.forEach(btn => {
    btn.addEventListener('click', () => {
      const tab = btn.dataset.tab;
      
      // Убрать active у всех
      tabBtns.forEach(b => b.classList.remove('active'));
      document.querySelectorAll('.tab-content').forEach(c => c.classList.remove('active'));
      
      // Добавить active к выбранным
      btn.classList.add('active');
      document.getElementById(`${tab}-container`).classList.add('active');
      
      currentTab = tab;
      
      // Обновляем видимость контейнеров
      checkEmptyState();
    });
  });
}

function setupNavigation() {
  const navToggle = document.getElementById('nav-toggle');
  const navPanel = document.getElementById('nav-panel');
  const navOverlay = document.getElementById('nav-panel-overlay');
  
  function toggleNav() {
    navPanel.classList.toggle('open');
    navToggle.classList.toggle('open');
    navOverlay.classList.toggle('open');
    
    const header = document.querySelector('header');
    if (navPanel.classList.contains('open')) {
      header.style.paddingLeft = '20%';
    } else {
      header.style.paddingLeft = '20px';
    }
  }

  if (navToggle && navPanel && navOverlay) {
    navToggle.addEventListener('click', toggleNav);
    navOverlay.addEventListener('click', toggleNav);
    
    document.addEventListener('keydown', (e) => {
      if (e.key === 'Escape' && navPanel.classList.contains('open')) {
        toggleNav();
      }
    });
  }

  ['main-btn', 'modules-btn', 'categories-btn', 'selected-btn', 'results-btn'].forEach(id => {
    const btn = document.getElementById(id);
    if (btn) {
      btn.addEventListener('click', () => {
        const page = id.replace('-btn', '.html');
        window.location.href = `/static/${page}`;
      });
    }
  });
}

window.addEventListener('DOMContentLoaded', async () => {
  const token = localStorage.getItem('token');
  if (!token) {
    window.location.href = `/static/login.html?redirect=${encodeURIComponent(window.location.href)}`;
    return;
  }

  await loadUserName(token);
  
  // Загружаем оба списка при старте
  loadModules(token);
  loadCategories(token);
  
  setupTabs();
  setupNavigation();
});
