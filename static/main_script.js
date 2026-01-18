let allModules = [], allCategories = [], allResults = [];
let modulesMap = {}, categoriesMap = [];
let currentSearchQuery = '';
let modulesCache = new Map();
let categoriesCache = new Map();

const API_BASE_URL = window.location.origin;

// Состояние поиска
let isSearchMode = false;
let searchOffsets = {
  users: 0,
  categories: 0,
  modules: 0
};
let searchResults = {
  users: [],
  categories: [],
  modules: []
};

window.addEventListener('DOMContentLoaded', async () => {
  const token = localStorage.getItem('token');
  if (!token) {
    window.location.href = '/static/login.html';
    return;
  }
  
  initNavigation();
  initAccordions();
  initSearch();

  try {
    const res = await fetch(`${API_BASE_URL}/api/v1/user/me?is_full=t`, {
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
    updateEmptyMessages();
    
  } catch (error) {
    console.error('Ошибка:', error);
    document.body.innerHTML = '<p style="text-align:center; margin-top:150px; font-size:1.4em; color:#c75c5c;">Ошибка загрузки данных. Перезагрузите страницу.</p>';
  }
});

async function fetchModuleInfo(moduleId, token) {
  const res = await fetch(`${API_BASE_URL}/api/v1/module/${moduleId}`, {
    headers: { 'Authorization': `Bearer ${token}` }
  });
  if (res.status === 401) {
    window.location.href = `/static/login.html?redirect=${encodeURIComponent(window.location.href)}`;
    throw new Error('Unauthorized');
  }
  if (res.status === 406) {
    console.warn('Модуль недоступен:', moduleId);
    return { module: { name: `Модуль ${moduleId} (недоступен)` } };
  }
  if (!res.ok) {
    console.error('Ошибка загрузки информации о модуле:', res.status);
    return null;
  }
  return await res.json();
}

async function fetchCategoryInfo(categoryId, token) {
  const res = await fetch(`${API_BASE_URL}/api/v1/category/${categoryId}`, {
    headers: { 'Authorization': `Bearer ${token}` }
  });
  if (res.status === 401) {
    window.location.href = `/static/login.html?redirect=${encodeURIComponent(window.location.href)}`;
    throw new Error('Unauthorized');
  }
  if (res.status === 406) {
    console.warn('Категория недоступна:', categoryId);
    return { category: { name: `Категория ${categoryId} (недоступна)` } };
  }
  if (!res.ok) {
    console.error('Ошибка загрузки информации о категории:', res.status);
    return null;
  }
  return await res.json();
}

async function getModuleName(moduleId, token) {
  if (modulesCache.has(moduleId)) {
    return modulesCache.get(moduleId);
  }
  
  try {
    const moduleInfo = await fetchModuleInfo(moduleId, token);
    const name = moduleInfo?.module?.name || `Модуль ${moduleId}`;
    modulesCache.set(moduleId, name);
    return name;
  } catch (error) {
    console.error('Ошибка получения названия модуля:', moduleId, error);
    const fallbackName = `Модуль ${moduleId}`;
    modulesCache.set(moduleId, fallbackName);
    return fallbackName;
  }
}

async function getCategoryName(categoryId, token) {
  if (categoriesCache.has(categoryId)) {
    return categoriesCache.get(categoryId);
  }
  
  try {
    const categoryInfo = await fetchCategoryInfo(categoryId, token);
    const name = categoryInfo?.category?.name || `Категория ${categoryId}`;
    categoriesCache.set(categoryId, name);
    return name;
  } catch (error) {
    console.error('Ошибка получения названия категории:', categoryId, error);
    const fallbackName = `Категория ${categoryId}`;
    categoriesCache.set(categoryId, fallbackName);
    return fallbackName;
  }
}

// Инициализация поиска
function initSearch() {
  const searchInput = document.getElementById('search-input');
  const searchButton = document.getElementById('search-button');
  const cancelButton = document.getElementById('cancel-search');
  
  searchButton.addEventListener('click', performSearch);
  cancelButton.addEventListener('click', cancelSearch);
  
  searchInput.addEventListener('keypress', (e) => {
    if (e.key === 'Enter') {
      e.preventDefault();
      performSearch();
    }
  });
}

// Основной поиск
async function performSearch() {
  const query = document.getElementById('search-input').value.trim();
  
  if (query.length <= 1) {
    if (isSearchMode) {
      cancelSearch();
    }
    return;
  }

  currentSearchQuery = query;
  const token = localStorage.getItem('token');
  
  // Анимация поиска
  const wrapper = document.querySelector('.search-input-wrapper');
  wrapper.classList.add('searching');
  
  try {
    // Сброс офсетов
    searchOffsets = { users: 0, categories: 0, modules: 0 };
    searchResults = { users: [], categories: [], modules: [] };
    
    // Параллельные запросы
    const [usersRes, categoriesRes, modulesRes] = await Promise.all([
      fetch(`${API_BASE_URL}/api/v1/search/users?name=${encodeURIComponent(query)}&limit=12&offset=0`, {
        headers: { 'Authorization': `Bearer ${token}` }
      }),
      fetch(`${API_BASE_URL}/api/v1/search/categories?name=${encodeURIComponent(query)}&limit=12&offset=0`, {
        headers: { 'Authorization': `Bearer ${token}` }
      }),
      fetch(`${API_BASE_URL}/api/v1/search/modules?name=${encodeURIComponent(query)}&limit=12&offset=0`, {
        headers: { 'Authorization': `Bearer ${token}` }
      })
    ]);

    const [usersData, categoriesData, modulesData] = await Promise.all([
      usersRes.ok ? usersRes.json() : { found_users: [] },
      categoriesRes.ok ? categoriesRes.json() : { found_categories: [] },
      modulesRes.ok ? modulesRes.json() : { found_modules: [] }
    ]);

    searchResults.users = usersData.found_users || [];
    searchResults.categories = categoriesData.found_categories || [];
    searchResults.modules = modulesData.found_modules || [];

    showSearchResults();
    
  } catch (error) {
    console.error('Ошибка поиска:', error);
  } finally {
    setTimeout(() => wrapper.classList.remove('searching'), 1000);
  }
}

// Показать результаты поиска
function showSearchResults() {
  isSearchMode = true;
  
  document.getElementById('main-content').classList.add('hidden');
  document.getElementById('search-results').style.display = 'block';
  document.getElementById('cancel-search').style.display = 'inline-block';
  
  renderSearchUsers();
  renderSearchCategories();
  renderSearchModules();
}

// Отмена поиска
function cancelSearch() {
  isSearchMode = false;
  currentSearchQuery = '';
  
  document.getElementById('main-content').classList.remove('hidden');
  document.getElementById('search-results').style.display = 'none';
  document.getElementById('cancel-search').style.display = 'none';
  document.getElementById('search-input').value = '';
}

// Рендер пользователей
function renderSearchUsers() {
  const container = document.getElementById('users-container');
  const loadMoreBtn = document.getElementById('users-load-more');
  
  if (searchResults.users.length === 0) {
    container.innerHTML = '<div class="empty-message">Пользователи не найдены</div>';
    loadMoreBtn.style.display = 'none';
    return;
  }
  
  container.innerHTML = searchResults.users.map(user => `
    <div class="card" onclick="window.location.href='/static/profile.html?user_id=${user.id}'">
      <div class="card-title">${user.name || 'Пользователь'}</div>
      <div class="card-subtitle">${user.email || ''}</div>
    </div>
  `).join('');
  
  loadMoreBtn.style.display = searchResults.users.length >= 12 ? 'block' : 'none';
  loadMoreBtn.onclick = () => loadMore('users');
}

// Рендер категорий поиска
function renderSearchCategories() {
  const container = document.getElementById('search-categories-container');
  const loadMoreBtn = document.getElementById('categories-load-more');
  
  if (searchResults.categories.length === 0) {
    container.innerHTML = '<div class="empty-message">Категории не найдены</div>';
    loadMoreBtn.style.display = 'none';
    return;
  }
  
  container.innerHTML = searchResults.categories.map(cat => `
    <div class="card" onclick="window.location.href='/static/category.html?category_id=${cat.id}'">
      <div class="card-title">${cat.name}</div>
      <div class="card-count">Модулей: ${cat.modules_count || 0}</div>
    </div>
  `).join('');
  
  loadMoreBtn.style.display = searchResults.categories.length >= 12 ? 'block' : 'none';
  loadMoreBtn.onclick = () => loadMore('categories');
}

// Рендер модулей поиска
function renderSearchModules() {
  const container = document.getElementById('search-modules-container');
  const loadMoreBtn = document.getElementById('modules-load-more');
  
  if (searchResults.modules.length === 0) {
    container.innerHTML = '<div class="empty-message">Модули не найдены</div>';
    loadMoreBtn.style.display = 'none';
    return;
  }
  
  container.innerHTML = searchResults.modules.map(mod => `
    <div class="card" onclick="window.location.href='/static/module.html?module_id=${mod.id}'">
      <div class="card-title">${mod.name}</div>
      <div class="card-count">Карточек: ${mod.cards_count || 0}</div>
    </div>
  `).join('');
  
  loadMoreBtn.style.display = searchResults.modules.length >= 12 ? 'block' : 'none';
  loadMoreBtn.onclick = () => loadMore('modules');
}

// Загрузка большего количества результатов
async function loadMore(type) {
  const token = localStorage.getItem('token');
  const offset = searchOffsets[type] + 12;
  
  try {
    const res = await fetch(`${API_BASE_URL}/api/v1/search/${type}s?name=${encodeURIComponent(currentSearchQuery)}&limit=12&offset=${offset}`, {
      headers: { 'Authorization': `Bearer ${token}` }
    });
    
    if (!res.ok) return;
    
    const data = await res.json();
    const results = data[`found_${type}s`] || [];
    
    searchResults[type] = searchResults[type].concat(results);
    searchOffsets[type] = offset;
    
    if (type === 'users') renderSearchUsers();
    if (type === 'categories') renderSearchCategories();
    if (type === 'modules') renderSearchModules();
    
  } catch (error) {
    console.error(`Ошибка загрузки ${type}:`, error);
  }
}

// Аккордеоны (оригинальная логика)
function initAccordions() {
  document.querySelectorAll('.accordion-header').forEach(header => {
    header.addEventListener('click', function(e) {
      if (isSearchMode) return;
      e.preventDefault();
      toggleAccordion(this.id);
    });
  });
}

async function toggleAccordion(headerId) {
  const header = document.getElementById(headerId);
  const content = header.parentNode.querySelector('.accordion-content');
  const isActive = header.classList.contains('active');

  document.querySelectorAll('.accordion-header').forEach(h => h.classList.remove('active'));
  document.querySelectorAll('.accordion-content').forEach(c => c.classList.remove('active'));

  if (!isActive) {
    header.classList.add('active');
    content.classList.add('active');
    
    if (currentSearchQuery && !isSearchMode) {
      if (headerId === 'modules-header') renderModulesPreview(allModules.filter(m => m.name.toLowerCase().includes(currentSearchQuery.toLowerCase())).slice(0, 4));
      if (headerId === 'categories-header') renderCategoriesPreview(allCategories.filter(c => c.name.toLowerCase().includes(currentSearchQuery.toLowerCase())).slice(0, 4));
      if (headerId === 'results-header') {
        const filteredResults = allResults.filter(r => r.name.toLowerCase().includes(currentSearchQuery.toLowerCase())).slice(0, 4);
        await renderResultsPreview(filteredResults);
      }
    } else {
      if (headerId === 'modules-header') renderModulesPreview(allModules.slice(0, 4));
      if (headerId === 'categories-header') renderCategoriesPreview(allCategories.slice(0, 4));
      if (headerId === 'results-header') await renderResultsPreview(allResults.slice(0, 4));
    }
  }
}

// Остальная оригинальная логика...
async function loadAllResults(token, userId) {
  try {
    const res = await fetch(`${API_BASE_URL}/api/v1/results/to_user/${userId}`, {
      headers: { 'Authorization': `Bearer ${token}` }
    });
    if (!res.ok) return;

    const data = await res.json();
    allResults = [];

    if (data.categories_results) {
      for (const catResult of data.categories_results) {
        const hasName = categoriesMap[catResult.category_id];
        const name = hasName || `Категория ${catResult.category_id}`;
        allResults.push({
          type: 'category', 
          name, 
          time: catResult.time, 
          resultType: catResult.modules_res[0]?.result?.type || 'test',
          id: catResult.category_result_id,
          category_id: catResult.category_id
        });
      }
    }

    if (data.modules_results) {
      for (const modResult of data.modules_results) {
        const hasName = modulesMap[modResult.module_id];
        const name = hasName || `Модуль ${modResult.module_id}`;
        allResults.push({
          type: 'module', 
          name, 
          time: modResult.time, 
          resultType: modResult.result.type,
          id: modResult.result.result_id,
          module_id: modResult.module_id
        });
      }
    }

    allResults.sort((a, b) => new Date(b.time) - new Date(a.time));
    await preloadVisibleResultsNames(token);
    
  } catch (e) {
    console.error('Ошибка загрузки результатов:', e);
  }
}

async function preloadVisibleResultsNames(token) {
  const firstResults = allResults.slice(0, 4);
  for (const result of firstResults) {
    if (result.type === 'category' && (!categoriesMap[result.category_id] || result.name.includes('Категория'))) {
      result.name = await getCategoryName(result.category_id, token);
    } else if (result.type === 'module' && (!modulesMap[result.module_id] || result.name.includes('Модуль'))) {
      result.name = await getModuleName(result.module_id, token);
    }
  }
}

async function renderResultsPreview(results) {
  const container = document.getElementById('results-container');
  document.getElementById('results-empty').style.display = results.length ? 'none' : 'block';
  
  if (results.length) {
    const token = localStorage.getItem('token');
    for (const result of results) {
      if (result.type === 'category' && (!result.name || result.name.includes('Категория') && !categoriesMap[result.category_id])) {
        result.name = await getCategoryName(result.category_id, token);
      } else if (result.type === 'module' && (!result.name || result.name.includes('Модуль') && !modulesMap[result.module_id])) {
        result.name = await getModuleName(result.module_id, token);
      }
    }
    
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

function updateEmptyMessages() {
  document.getElementById('modules-empty').style.display = allModules.length ? 'none' : 'block';
  document.getElementById('categories-empty').style.display = allCategories.length ? 'none' : 'block';
  document.getElementById('results-empty').style.display = allResults.length ? 'none' : 'block';
}

function initNavigation() {
  const navToggle = document.getElementById('nav-toggle');
  const navPanel = document.getElementById('nav-panel');
  
  if (navToggle) {
    navToggle.onclick = () => {
      navPanel.classList.toggle('open');
      navToggle.classList.toggle('open');
    };
  }

  const navButtons = {
    'modules-btn': '/static/modules.html',
    'categories-btn': '/static/categories.html',
    'results-btn': '/static/results.html'
  };

  Object.entries(navButtons).forEach(([id, url]) => {
    const btn = document.getElementById(id);
    if (btn) {
      btn.onclick = () => window.location.href = url;
    }
  });
}
