// results_script.js
let modulesCache = new Map();
let categoriesCache = new Map();

const API_BASE_URL = window.location.origin;

window.addEventListener('DOMContentLoaded', async () => {
  const token = localStorage.getItem('token');
  if (!token) {
    window.location.href = `/static/login.html?redirect=${encodeURIComponent(window.location.href)}`;
    return;
  }
  
  let userData = null;
  let modulesMap = {};
  let categoriesMap = {};
  
  const editBtn = document.getElementById('editResultsBtn');
  const resultsContainer = document.getElementById('results-container');
  const resultsEmpty = document.getElementById('results-empty');
  
  // Инициализация навигации
  initNavigation();
  
  try {
    // Загружаем данные пользователя
    const userRes = await fetch(`${API_BASE_URL}/api/v1/user/me?is_full=t`, {
      headers: { 'Authorization': `Bearer ${token}` }
    });
    
    if (userRes.status === 401) {
      window.location.href = '/static/login.html';
      return;
    }
    if (!userRes.ok) {
      throw new Error('Ошибка загрузки данных пользователя');
    }
    
    userData = await userRes.json();
    
    // Имя пользователя
    const usernameElem = document.getElementById('username');
    usernameElem.textContent = userData.user.name || 'Пользователь';
    usernameElem.onclick = () => {
      window.location.href = '/static/profile.html';
    };

    // Создаем карты модулей и категорий для поиска названий
    modulesMap = {};
    if (userData.user.modules) {
      userData.user.modules.forEach(module => {
        modulesMap[module.id] = module.name;
      });
    }

    categoriesMap = {};
    if (userData.user.categories) {
      userData.user.categories.forEach(category => {
        categoriesMap[category.id] = category.name;
      });
    }

    // Загружаем результаты
    await loadResults(token, userData.user.id, modulesMap, categoriesMap);
    
  } catch (error) {
    console.error('Ошибка загрузки данных:', error);
    document.body.innerHTML = '<p style="text-align:center; margin-top:50px; font-size:18px; color:#c75c5c;">Ошибка загрузки данных. Попробуйте перезагрузить страницу.</p>';
  }

  // Обработчик кнопки редактирования
  editBtn.addEventListener('click', function() {
    const isEditMode = this.classList.contains('edit-mode');
    
    if (isEditMode) {
      this.textContent = 'Редактировать результаты';
      this.classList.remove('edit-mode');
      resultsContainer.classList.remove('edit-mode');
      document.querySelectorAll('.card').forEach(card => {
        card.classList.remove('edit-mode');
        card.querySelector('.card-actions')?.classList.remove('show');
      });
    } else {
      this.textContent = 'Завершить редактирование';
      this.classList.add('edit-mode');
      resultsContainer.classList.add('edit-mode');
      document.querySelectorAll('.card').forEach(card => {
        card.classList.add('edit-mode');
        const actions = card.querySelector('.card-actions');
        if (actions) actions.classList.add('show');
      });
    }
  });
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

async function loadResults(token, userId, modulesMap, categoriesMap) {
  try {
    const resultsRes = await fetch(`${API_BASE_URL}/api/v1/results/to_user/${userId}`, {
      headers: { 'Authorization': `Bearer ${token}` }
    });

    if (!resultsRes.ok) {
      throw new Error('Ошибка загрузки результатов');
    }

    const resultsData = await resultsRes.json();
    const resultsContainer = document.getElementById('results-container');
    const resultsEmpty = document.getElementById('results-empty');
    
    const allResults = [];
    
    // Категории результаты
    if (resultsData.categories_results) {
      for (const catResult of resultsData.categories_results) {
        let categoryName = categoriesMap[catResult.category_id];
        if (!categoryName) {
          categoryName = await getCategoryName(catResult.category_id, token);
        }
        const firstModuleResult = catResult.modules_res[0]?.result;
        const resultType = firstModuleResult ? firstModuleResult.type : 'unknown';
        
        allResults.push({
          type: 'category',
          name: categoryName,
          time: catResult.time,
          resultType: resultType,
          id: catResult.category_result_id,
          category_id: catResult.category_id
        });
      }
    }
    
    // Модули результаты
    if (resultsData.modules_results) {
      for (const modResult of resultsData.modules_results) {
        let moduleName = modulesMap[modResult.module_id];
        if (!moduleName) {
          moduleName = await getModuleName(modResult.module_id, token);
        }
        allResults.push({
          type: 'module',
          name: moduleName,
          time: modResult.time,
          resultType: modResult.result.type,
          id: modResult.result.result_id,
          module_id: modResult.module_id
        });
      }
    }
    
    allResults.sort((a, b) => new Date(b.time) - new Date(a.time));
    
    resultsContainer.innerHTML = '';
    
    if (allResults.length === 0) {
      resultsEmpty.style.display = 'block';
    } else {
      resultsEmpty.style.display = 'none';
      allResults.forEach(result => {
        renderResultCard(result, token);
      });
    }
  } catch (error) {
    console.error('Ошибка загрузки результатов:', error);
    document.getElementById('results-empty').textContent = 'Ошибка загрузки результатов';
  }
}

function renderResultCard(result, token) {
  const resultsContainer = document.getElementById('results-container');
  
  const card = document.createElement('div');
  card.className = 'card result-card';
  
  const resultTypeClass = result.resultType === 'test' ? 'test' : 'learning';
  const resultTypeText = result.resultType === 'test' ? 'Тест' : 'Заучивание';
  const subtitle = result.type === 'category' ? 'Категория' : 'Модуль';
  
  card.innerHTML = `
    <div class="card-title">${result.name}</div>
    <div class="card-subtitle">${subtitle}</div>
    <div class="card-time">${formatDate(result.time)}</div>
    <div class="card-type ${resultTypeClass}">${resultTypeText}</div>
    <div class="card-actions">
      <button class="delete-btn" data-result-id="${result.id}" data-result-type="${result.type}">×</button>
    </div>
  `;

  card.addEventListener('click', (e) => {
    if (e.target.classList.contains('delete-btn')) return;
    const param = result.type === 'category' ? `category_res_id=${result.id}` : `result_id=${result.id}`;
    window.location.href = `/static/result.html?${param}`;
  });
  
  const deleteBtn = card.querySelector('.delete-btn');
  deleteBtn.addEventListener('click', async (e) => {
    e.stopPropagation();
    const resultId = deleteBtn.dataset.resultId;
    const resultType = deleteBtn.dataset.resultType;
    
    if (confirm(`Удалить результат "${result.name}"?`)) {
      const success = await deleteResult(resultType, resultId, token);
      if (success) {
        card.remove();
        if (resultsContainer.children.length === 0) {
          document.getElementById('results-empty').style.display = 'block';
        }
      }
    }
  });
  
  resultsContainer.appendChild(card);
}

async function deleteResult(resultType, resultId, token) {
  const endpoint = resultType === 'category' 
    ? `${API_BASE_URL}/api/v1/results/category_result/delete/${resultId}`
    : `${API_BASE_URL}/api/v1/results/module_result/delete/${resultId}`;
    
  try {
    const res = await fetch(endpoint, {
      method: 'DELETE',
      headers: { 'Authorization': `Bearer ${token}` }
    });
    
    if (!res.ok) {
      throw new Error('Ошибка удаления');
    }
    return true;
  } catch (error) {
    console.error('Ошибка удаления результата:', error);
    alert('Ошибка при удалении результата');
    return false;
  }
}

function formatDate(isoString) {
  const date = new Date(isoString);
  return date.toLocaleString('ru-RU', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    hour12: false
  });
}

function initNavigation() {
  const navToggle = document.getElementById('nav-toggle');
  const navPanel = document.getElementById('nav-panel');
  const navOverlay = document.getElementById('nav-panel-overlay');
  
  if (navToggle) {
    navToggle.addEventListener('click', function() {
      navPanel.classList.toggle('open');
      navToggle.classList.toggle('open');
      navOverlay.classList.toggle('open');
    });
  }
  
  if (navOverlay) {
    navOverlay.addEventListener('click', function() {
      navPanel.classList.remove('open');
      navToggle.classList.remove('open');
      navOverlay.classList.remove('open');
    });
  }

  const navButtons = {
    'main-btn': '/static/main.html',
    'modules-btn': '/static/modules.html',
    'categories-btn': '/static/categories.html',
    'results-btn': '/static/results.html'
  };

  Object.entries(navButtons).forEach(([id, url]) => {
    const btn = document.getElementById(id);
    if (btn) {
      btn.addEventListener('click', () => {
        window.location.href = url;
      });
    }
  });
}
