let allModules = [], allCategories = [], allResults = [];
let modulesMap = {}, categoriesMap = [];
let currentSearchQuery = '';

window.addEventListener('DOMContentLoaded', async () => {
  const token = localStorage.getItem('token');
  if (!token) {
    window.location.href = '/static/login.html';
    return;
  }
  
  initNavigation();
  initAccordions();
  initSearch();  // ← ИЗМЕНЕНО - теперь только кнопка

  try {
    const res = await fetch('http://localhost:8080/api/v1/user/me?is_full=t', {
      headers: { 'Authorization': `Bearer ${token}` }
    });
    
    if (res.status === 401) {
      window.location.href = '/static/login.html';
      return;
    }
    if (!res.ok) throw new Error('Ошибка загрузки данных');

    const userData = await res.json();
    document.getElementById('username').textContent = userData.user.name || 'Пользователь';
    document.getElementById('username').onclick = () => window.location.href = '/static/profile.html';

    modulesMap = {};
    if (userData.user.modules) {
      userData.user.modules.forEach(m => modulesMap[m.id] = m.name);
      allModules = userData.user.modules;
    }
    
    categoriesMap = {};
    if (userData.user.categories) {
      userData.user.categories.forEach(c => categoriesMap[c.id] = c.name);
      allCategories = userData.user.categories;
    }

    await loadAllResults(token, userData.user.id);
    
    // ← Изначально показываем превью при открытии аккордеонов
    updateEmptyMessages();

  } catch (error) {
    console.error('Ошибка:', error);
    document.body.innerHTML = '<p style="text-align:center; margin-top:150px; font-size:1.4em; color:#c75c5c;">Ошибка загрузки данных. Перезагрузите страницу.</p>';
  }
});

function initSearch() {
  const searchInput = document.getElementById('search-input');
  const searchButton = document.getElementById('search-button');
  
  // ← УБРАНО: input event - теперь только кнопка!
  
  searchButton.addEventListener('click', performSearch);
  
  // Enter в поле поиска = клик кнопки
  searchInput.addEventListener('keypress', (e) => {
    if (e.key === 'Enter') {
      e.preventDefault();
      performSearch();
    }
  });
}

function performSearch() {
  const query = document.getElementById('search-input').value.trim();
  currentSearchQuery = query;
  
  // Анимация поиска
  const wrapper = document.querySelector('.search-input-wrapper');
  wrapper.classList.add('searching');
  
  setTimeout(() => {
    wrapper.classList.remove('searching');
  }, 1000);
  
  // ← ТОЛЬКО ПРИ НАЖАТИИ обновляем контент
  if (query) {
    renderModulesPreview(allModules.filter(m => m.name.toLowerCase().includes(query.toLowerCase())).slice(0, 4));
    renderCategoriesPreview(allCategories.filter(c => c.name.toLowerCase().includes(query.toLowerCase())).slice(0, 4));
    renderResultsPreview(allResults.filter(r => r.name.toLowerCase().includes(query.toLowerCase())).slice(0, 4));
  } else {
    // Если пустой поиск - показываем все
    renderModulesPreview(allModules.slice(0, 4));
    renderCategoriesPreview(allCategories.slice(0, 4));
    renderResultsPreview(allResults.slice(0, 4));
  }
}

// Остальные функции БЕЗ ИЗМЕНЕНИЙ
function initAccordions() {
  document.querySelectorAll('.accordion-header').forEach(header => {
    header.addEventListener('click', function(e) {
      e.preventDefault();
      toggleAccordion(this.id);
    });
  });
}

function toggleAccordion(headerId) {
  const header = document.getElementById(headerId);
  const content = header.parentNode.querySelector('.accordion-content');
  const isActive = header.classList.contains('active');

  document.querySelectorAll('.accordion-header').forEach(h => h.classList.remove('active'));
  document.querySelectorAll('.accordion-content').forEach(c => c.classList.remove('active'));

  if (!isActive) {
    header.classList.add('active');
    content.classList.add('active');
    
    // Показываем отфильтрованные результаты поиска ИЛИ обычные
    if (currentSearchQuery) {
      if (headerId === 'modules-header') renderModulesPreview(allModules.filter(m => m.name.toLowerCase().includes(currentSearchQuery.toLowerCase())).slice(0, 4));
      if (headerId === 'categories-header') renderCategoriesPreview(allCategories.filter(c => c.name.toLowerCase().includes(currentSearchQuery.toLowerCase())).slice(0, 4));
      if (headerId === 'results-header') renderResultsPreview(allResults.filter(r => r.name.toLowerCase().includes(currentSearchQuery.toLowerCase())).slice(0, 4));
    } else {
      if (headerId === 'modules-header') renderModulesPreview(allModules.slice(0, 4));
      if (headerId === 'categories-header') renderCategoriesPreview(allCategories.slice(0, 4));
      if (headerId === 'results-header') renderResultsPreview(allResults.slice(0, 4));
    }
  }
}

async function loadAllResults(token, userId) {
  try {
    const res = await fetch(`http://localhost:8080/api/v1/results/to_user/${userId}`, {
      headers: { 'Authorization': `Bearer ${token}` }
    });
    if (!res.ok) return;

    const data = await res.json();
    allResults = [];

    if (data.categories_results) {
      for (const catResult of data.categories_results) {
        const name = categoriesMap[catResult.category_id] || `Категория ${catResult.category_id}`;
        allResults.push({
          type: 'category', name, time: catResult.time, 
          resultType: catResult.modules_res[0]?.result?.type || 'test',
          id: catResult.category_result_id
        });
      }
    }

    if (data.modules_results) {
      for (const modResult of data.modules_results) {
        const name = modulesMap[modResult.module_id] || `Модуль ${modResult.module_id}`;
        allResults.push({
          type: 'module', name, time: modResult.time, 
          resultType: modResult.result.type,
          id: modResult.result.result_id
        });
      }
    }

    allResults.sort((a, b) => new Date(b.time) - new Date(a.time));
  } catch (e) {
    console.error('Ошибка загрузки результатов:', e);
  }
}

function renderModulesPreview(modules) {
  const container = document.getElementById('modules-container');
  document.getElementById('modules-empty').style.display = modules.length ? 'none' : 'block';
  
  if (modules.length) {
    container.innerHTML = modules.map(m => `
      <div class="card" onclick="window.location.href='/static/module.html?module_id=${m.id}'">
        <div class="card-title">${m.name}</div>
        <div class="card-count">Карточек: ${m.cards?.length || 0}</div>
      </div>
    `).join('') + `
      <div class="view-all-link" onclick="window.location.href='/static/modules.html'">
        → Все модули (${allModules.length})
      </div>
    `;
  }
}

function renderCategoriesPreview(categories) {
  const container = document.getElementById('categories-container');
  document.getElementById('categories-empty').style.display = categories.length ? 'none' : 'block';
  
  if (categories.length) {
    container.innerHTML = categories.map(c => `
      <div class="card" onclick="window.location.href='/static/category.html?category_id=${c.id}'">
        <div class="card-title">${c.name}</div>
        <div class="card-count">Модулей: ${c.modules?.length || 0}</div>
      </div>
    `).join('') + `
      <div class="view-all-link" onclick="window.location.href='/static/categories.html'">
        → Все категории (${allCategories.length})
      </div>
    `;
  }
}

function renderResultsPreview(results) {
  const container = document.getElementById('results-container');
  document.getElementById('results-empty').style.display = results.length ? 'none' : 'block';
  
  if (results.length) {
    container.innerHTML = results.map(r => {
      const typeClass = r.resultType === 'test' ? 'test' : 'learning';
      const typeText = r.resultType === 'test' ? 'Тест' : 'Заучивание';
      const time = r.time.replace('T', ' ').replace('Z', '');
      const param = r.type === 'category' ? `category_res_id=${r.id}` : `result_id=${r.id}`;
      return `
        <div class="card result-card" onclick="window.location.href='/static/result.html?${param}'">
          <div class="card-title">${r.name}</div>
          <div class="card-subtitle">${r.type === 'category' ? 'Категория' : 'Модуль'}</div>
          <div class="card-time">${time}</div>
          <div class="card-type ${typeClass}">${typeText}</div>
        </div>
      `;
    }).join('') + `
      <div class="view-all-link" onclick="window.location.href='/static/results.html'">
        → Все результаты (${allResults.length})
      </div>
    `;
  }
}

function updateEmptyMessages() {
  document.getElementById('modules-empty').style.display = allModules.length ? 'none' : 'block';
  document.getElementById('categories-empty').style.display = allCategories.length ? 'none' : 'block';
  document.getElementById('results-empty').style.display = allResults.length ? 'none' : 'block';
}

function initNavigation() {
  const navToggle = document.getElementById('nav-toggle');
  const navPanel = document.getElementById('nav-panel');
  
  navToggle.onclick = () => {
    navPanel.classList.toggle('open');
    navToggle.classList.toggle('open');
  };

  document.getElementById('modules-btn').onclick = () => window.location.href = '/static/modules.html';
  document.getElementById('categories-btn').onclick = () => window.location.href = '/static/categories.html';
  document.getElementById('results-btn').onclick = () => window.location.href = '/static/results.html';
}
