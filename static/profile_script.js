window.addEventListener('DOMContentLoaded', async () => {
  const token = localStorage.getItem('token');
  if (!token) {
    window.location.href = '/static/login.html';
    return;
  }

  initNavigation();

  let userData = null;
  let modulesMap = {};
  let categoriesMap = {};

  try {
    const res = await fetch('http://localhost:8080/api/v1/user/me?is_full=t', {
      headers: { 'Authorization': `Bearer ${token}` }
    });

    if (res.status === 401) {
      window.location.href = '/static/login.html';
      return;
    }
    if (!res.ok) {
      throw new Error('Ошибка загрузки данных пользователя');
    }

    userData = await res.json();

    // ПРОФИЛЬ
    document.getElementById('username').textContent = userData.user.name || 'Пользователь';
    document.getElementById('profile-name').textContent = userData.user.name || 'Пользователь';
    document.getElementById('profile-login').textContent = `Логин: ${userData.user.login || 'Не указан'}`;

    // Карты модулей/категорий
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

    // Отображаем контент
    await displayModules(userData.user.modules);
    await displayCategories(userData.user.categories);
    await loadAndDisplayResults(token, userData.user.id, modulesMap, categoriesMap);

    // КНОПКА СМЕНИТЬ ПАРОЛЬ
    document.getElementById('change-password-btn').onclick = () => {
      const newPassword = prompt('Введите новый пароль:');
      if (newPassword && newPassword.length >= 6) {
        changePassword(token, newPassword);
      } else if (newPassword) {
        alert('Пароль должен содержать минимум 6 символов');
      }
    };

  } catch (error) {
    console.error('Ошибка загрузки данных:', error);
    document.body.innerHTML = '<p style="text-align:center; margin-top:50px; font-size:18px; color:#c75c5c;">Ошибка загрузки данных. Попробуйте перезагрузить страницу.</p>';
  }
});

async function changePassword(token, newPassword) {
  try {
    const res = await fetch('http://localhost:8080/api/v1/user/password', {
      method: 'PUT',
      headers: { 
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({ password: newPassword })
    });

    if (res.ok) {
      alert('Пароль успешно изменен!');
    } else {
      const error = await res.json();
      alert('Ошибка смены пароля: ' + (error.message || 'Попробуйте еще раз'));
    }
  } catch (error) {
    console.error('Ошибка:', error);
    alert('Ошибка соединения. Попробуйте еще раз.');
  }
}

// Функции отображения (копия из main_script.js)
async function displayModules(modules) {
  const modulesContainer = document.getElementById('modules-container');
  const modulesEmpty = document.getElementById('modules-empty');
  modulesContainer.innerHTML = '';

  if (!modules || modules.length === 0) {
    modulesEmpty.style.display = 'block';
  } else {
    modulesEmpty.style.display = 'none';
    modules.forEach(module => {
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
}

async function displayCategories(categories) {
  const categoriesContainer = document.getElementById('categories-container');
  const categoriesEmpty = document.getElementById('categories-empty');
  categoriesContainer.innerHTML = '';

  if (!categories || categories.length === 0) {
    categoriesEmpty.style.display = 'block';
  } else {
    categoriesEmpty.style.display = 'none';
    categories.forEach(category => {
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
}

async function loadAndDisplayResults(token, userId, modulesMap, categoriesMap) {
  try {
    const resultsRes = await fetch(`http://localhost:8080/api/v1/results/to_user/${userId}`, {
      headers: { 'Authorization': `Bearer ${token}` }
    });

    if (!resultsRes.ok) {
      throw new Error('Ошибка загрузки результатов');
    }

    const resultsData = await resultsRes.json();
    const resultsContainer = document.getElementById('results-container');
    const resultsEmpty = document.getElementById('results-empty');

    const allResults = [];

    if (resultsData.categories_results) {
      for (const catResult of resultsData.categories_results) {
        const categoryName = categoriesMap[catResult.category_id] || await fetchEntityName('category', catResult.category_id, token);
        const firstModuleResult = catResult.modules_res[0]?.result;
        const resultType = firstModuleResult ? firstModuleResult.type : 'unknown';

        allResults.push({
          type: 'category',
          name: categoryName,
          time: catResult.time,
          resultType: resultType,
          id: catResult.category_result_id
        });
      }
    }

    if (resultsData.modules_results) {
      for (const modResult of resultsData.modules_results) {
        const moduleName = modulesMap[modResult.module_id] || await fetchEntityName('module', modResult.module_id, token);
        allResults.push({
          type: 'module',
          name: moduleName,
          time: modResult.time,
          resultType: modResult.result.type,
          id: modResult.result.result_id
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
        const card = document.createElement('div');
        card.className = 'card result-card';

        const resultTypeClass = result.resultType === 'test' ? 'test' : 'learning';
        const resultTypeText = result.resultType === 'test' ? 'Тест' : 'Заучивание';
        const subtitle = result.type === 'category' ? 'Категория' : 'Модуль';

        result.time = result.time.replace('T', ' ').replace('Z', '');

        card.innerHTML = `
          <div class="card-title">${result.name}</div>
          <div class="card-subtitle">${subtitle}</div>
          <div class="card-time">${result.time}</div>
          <div class="card-type ${resultTypeClass}">${resultTypeText}</div>
        `;

        card.addEventListener('click', () => {
          const param = result.type === 'category' ? `category_res_id=${result.id}` : `result_id=${result.id}`;
          window.location.href = `/static/result.html?${param}`;
        });

        resultsContainer.appendChild(card);
      });
    }
  } catch (error) {
    console.error('Ошибка загрузки результатов:', error);
    document.getElementById('results-empty').textContent = 'Ошибка загрузки результатов';
  }
}

async function fetchEntityName(entityType, entityId, token) {
  try {
    const url = `http://localhost:8080/api/v1/${entityType}/${entityId}`;
    const res = await fetch(url, {
      headers: { 'Authorization': `Bearer ${token}` }
    });

    if (res.ok) {
      const data = await res.json();
      return data.name || `ID ${entityId}`;
    }
  } catch (error) {
    console.error(`Ошибка загрузки ${entityType}:`, error);
  }
  return `ID ${entityId}`;
}

function initNavigation() {
  const navToggle = document.getElementById('nav-toggle');
  const navPanel = document.getElementById('nav-panel');

  navToggle.addEventListener('click', () => {
    navPanel.classList.toggle('open');
    navToggle.classList.toggle('open');
  });

  const navButtons = {
    'main-btn': '/static/main.html',
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
