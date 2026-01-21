let detailedResultsData = null;
let moduleCardDetails = new Map();

const API_BASE_URL = window.location.origin;

window.addEventListener('DOMContentLoaded', async () => {
  const token = localStorage.getItem('token');
  if (!token) {
    window.location.href = `/static/login.html?redirect=${encodeURIComponent(window.location.href)}`;
    return;
  }

  initNavigation();
  
  const urlParams = new URLSearchParams(window.location.search);
  const categoryResId = urlParams.get('category_res_id');
  const resultId = urlParams.get('result_id');
  
  if (!categoryResId && !resultId) {
    document.body.innerHTML = '<p style="text-align:center; margin-top:50px;">Неверные параметры</p>';
    return;
  }

  try {
    let resultData;
    
    if (categoryResId) {
      resultData = await fetchCategoryResult(categoryResId, token);
      detailedResultsData = resultData.category_result;
    } else {
      resultData = await fetchModuleResult(resultId, token);
      detailedResultsData = resultData.module_res;
    }
    
    await loadResultDetails(detailedResultsData, token);
  } catch (error) {
    console.error('Ошибка загрузки результата:', error);
    document.body.innerHTML = '<p style="text-align:center; margin-top:50px; color:#c75c5c;">Ошибка загрузки результата</p>';
  }
});

async function fetchCategoryResult(categoryResId, token) {
  const res = await fetch(`${API_BASE_URL}/api/v1/results/category_result/${categoryResId}`, {
    headers: { 'Authorization': `Bearer ${token}` }
  });
  if (res.status === 401) {
    window.location.href = `/static/login.html?redirect=${encodeURIComponent(window.location.href)}`;
    throw new Error('Unauthorized');
  }
  if (!res.ok) throw new Error('Ошибка загрузки результата категории');
  return await res.json();
}

async function fetchModuleResult(resultId, token) {
  const res = await fetch(`${API_BASE_URL}/api/v1/results/module_result/${resultId}`, {
    headers: { 'Authorization': `Bearer ${token}` }
  });
  if (res.status === 401) {
    window.location.href = `/static/login.html?redirect=${encodeURIComponent(window.location.href)}`;
    throw new Error('Unauthorized');
  }
  if (!res.ok) throw new Error('Ошибка загрузки результата модуля');
  return await res.json();
}

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


async function loadResultDetails(resultData, token) {
  const userRes = await fetch(`${API_BASE_URL}/api/v1/user/me?is_full=f`, {
    headers: { 'Authorization': `Bearer ${token}` }
  });
  if (userRes.status === 401) {
    window.location.href = `/static/login.html?redirect=${encodeURIComponent(window.location.href)}`;
    return;
  }
  const userData = await userRes.json();

  document.getElementById('username').textContent = userData.user.name || 'Пользователь';
  document.getElementById('resultTime').textContent = formatDate(resultData.time);
  
  let title, type, categoryName, moduleName, resultIds = [];
  const isCategory = !!resultData.category_result_id;
  
  if (isCategory) {
    // Загружаем информацию о категории отдельным запросом
    const categoryInfo = await fetchCategoryInfo(resultData.category_id, token);
    title = categoryInfo?.category?.name || `Категория ${resultData.category_id}`;
    categoryName = title;
    type = resultData.modules_res[0]?.result.type || 'test';
    resultIds = resultData.modules_res.map(m => m.result.result_id);
    document.getElementById('resultModule').style.display = 'none';
  } else {
    // Загружаем информацию о модуле отдельным запросом
    const moduleInfo = await fetchModuleInfo(resultData.module_id, token);
    title = moduleInfo?.module?.name || `Модуль ${resultData.module_id}`;
    moduleName = title;
    type = resultData.result.type;
    resultIds = [resultData.result.result_id];
    document.getElementById('resultCategory').style.display = 'none';
  }
  
  document.getElementById('resultTitle').textContent = title;
  document.getElementById('resultType').textContent = type === 'test' ? 'Тест' : 'Заучивание';
  document.getElementById('resultType').className = `info-value type-${type}`;
  document.getElementById('categoryName').textContent = categoryName || '';
  document.getElementById('moduleName').textContent = moduleName || '';

  const stats = await loadCardStats(resultIds, token);
  renderChart(stats);
  
  await renderDetailedResults(resultData, token, isCategory, type);
}

// Обновленная функция renderDetailedResults - теперь получает названия модулей отдельно
async function renderDetailedResults(resultData, token, isCategory, resultType) {
  const container = document.getElementById('detailed-results-container');
  if (!container) return;

  container.innerHTML = `
    <div class="detailed-results-header">
      <h3 class="detailed-results-title">Подробные результаты</h3>
    </div>
  `;

  if (isCategory) {
    for (const moduleRes of resultData.modules_res) {
      const moduleId = moduleRes.module_id;
      // Загружаем информацию о модуле отдельным запросом
      const moduleInfo = await fetchModuleInfo(moduleId, token);
      const moduleName = moduleInfo?.module?.name || `Модуль ${moduleId}`;
      const moduleResultType = moduleRes.result.type || resultType;
      const resultId = moduleRes.result.result_id;
      
      const cardsData = await fetchCardsResult(resultId, token);
      renderModuleSummary(container, moduleName, cardsData, moduleResultType, resultId);
    }
  } else {
    const resultId = resultData.result.result_id;
    // Название модуля уже получено в loadResultDetails
    const moduleName = document.getElementById('moduleName').textContent || 'Результаты модуля';
    const cardsData = await fetchCardsResult(resultId, token);
    renderModuleSummary(container, moduleName, cardsData, resultType, resultId);
  }
}

// Остальные функции остаются без изменений
function renderModuleSummary(container, moduleName, cardsData, resultType, resultId) {
  const cardsResult = cardsData.cards_results || [];
  const correct = cardsResult.filter(c => c.result === 'correct').length;
  const incorrect = cardsResult.filter(c => c.result !== 'correct').length;
  const total = correct + incorrect;
  const percent = total > 0 ? Math.round((correct / total) * 100) : 0;

  const moduleDiv = document.createElement('div');
  moduleDiv.className = 'module-result-card';
  moduleDiv.dataset.resultId = resultId;
  moduleDiv.dataset.resultType = resultType;
  moduleDiv.innerHTML = `
    <div class="module-header">
      <div style="display: flex; flex-direction: column; flex: 1;">
        <h4 class="module-title">${moduleName}</h4>
        <div class="module-stats-row" style="margin-top: 8px;">
         <div class="stat-badge correct">${correct} правильных</div>
         <div class="stat-badge incorrect">${incorrect} неправильных</div>
        </div>
      </div>
      <div class="module-stats-badge">
        <span class="correct-count">${correct}</span> / ${total}
        <span class="percent">${percent}%</span>
      </div>
    </div>
    <div class="module-details collapsed" style="display: none;">
      <div class="cards-section">
        <h5 class="cards-section-title">Загрузка...</h5>
      </div>
    </div>
    <button class="details-toggle" data-result-id="${resultId}">
      Подробнее <span class="toggle-icon">▼</span>
    </button>
  `;
  container.appendChild(moduleDiv);

  const toggleBtn = moduleDiv.querySelector('.details-toggle');
  toggleBtn.addEventListener('click', async (e) => {
    e.stopPropagation();
    const resultId = e.currentTarget.dataset.resultId;
    const moduleDetails = moduleDiv.querySelector('.module-details');
    const isCollapsed = moduleDetails.classList.contains('collapsed');
    
    if (isCollapsed) {
      await toggleModuleDetails(moduleDiv, resultId);
    } else {
      moduleDetails.style.display = 'none';
      moduleDetails.classList.add('collapsed');
      toggleBtn.innerHTML = 'Подробнее <span class="toggle-icon">▼</span>';
    }
  });
}

// Остальные функции (toggleModuleDetails, fetchCardDetails, renderCardDetail, fetchCardsResult, loadCardStats, renderChart, formatDate, initNavigation) остаются без изменений
async function toggleModuleDetails(moduleDiv, resultId) {
  const token = localStorage.getItem('token');
  const moduleDetails = moduleDiv.querySelector('.module-details');
  const toggleBtn = moduleDiv.querySelector('.details-toggle');
  const resultType = moduleDiv.dataset.resultType;
  
  moduleDetails.classList.remove('collapsed');
  moduleDetails.innerHTML = `
    <div class="cards-section">
      <h5 class="cards-section-title">Загрузка...</h5>
    </div>
  `;
  toggleBtn.innerHTML = 'Скрыть <span class="toggle-icon">▲</span>';

  try {
    const cardsData = await fetchCardsResult(resultId, token);
    const cardsResult = cardsData.cards_results || [];
    
    const cardsSection = moduleDetails.querySelector('.cards-section');
    
    cardsSection.innerHTML = `
      <h5 class="cards-section-title">${resultType === 'test' ? 'Вопросы теста:' : 'Слова для заучивания:'}</h5>
    `;

    for (const [index, card] of cardsResult.entries()) {
      const cardData = await fetchCardDetails(card.card_id, token);
      renderCardDetail(cardsSection, card, cardData, resultType, index);
    }
    
    moduleDetails.style.display = 'block';
  } catch (error) {
    console.error('Ошибка загрузки деталей карточек:', error);
    moduleDetails.innerHTML = `
      <div class="cards-section">
        <h5 class="cards-section-title">Ошибка загрузки деталей</h5>
      </div>
    `;
  }
}

async function fetchCardDetails(cardId, token) {
  const res = await fetch(`${API_BASE_URL}/api/v1/card/${cardId}`, {
    headers: { 'Authorization': `Bearer ${token}` }
  });
  
  if (res.status === 401) {
    window.location.href = `/static/login.html?redirect=${encodeURIComponent(window.location.href)}`;
    throw new Error('Unauthorized');
  }
  
  if (!res.ok) {
    console.error('Ошибка загрузки карточки:', res.status);
    return null;
  }
  
  return await res.json();
}

function renderCardDetail(container, cardResult, cardData, resultType, index) {
  if (!cardData?.card) return;

  const card = cardData.card;
  const cardDiv = document.createElement('div');
  cardDiv.className = `card-result ${resultType} ${cardResult.result}`;
  
  cardDiv.innerHTML = `
    <div class="card-header">
      <span class="card-number">${index + 1}</span>
      <span class="card-status ${cardResult.result}">
        ${resultType === 'test' ? 
          (cardResult.result === 'correct' ? 'Правильно' : 'Неправильно') : 
          (cardResult.result === 'correct' ? 'Знаю' : 'Не знаю')}
      </span>
    </div>
    <div class="card-content">
      <div class="card-fields-row">
        <div class="field-column term-column">
         <div class="field-label">Термин (${card.term?.lang?.toUpperCase() || 'EN'})</div>
         <div class="card-field compact">${card.term?.text || '—'}</div>
        </div>
        <div class="field-column definition-column">
         <div class="field-label">Определение (${card.definition?.lang?.toUpperCase() || 'RU'})</div>
         <div class="card-field compact">${card.definition?.text || '—'}</div>
        </div>
      </div>
    </div>
  `;
  
  container.appendChild(cardDiv);
}

async function fetchCardsResult(resultId, token) {
  const res = await fetch(`${API_BASE_URL}/api/v1/results/cards_result/${resultId}`, {
    headers: { 'Authorization': `Bearer ${token}` }
  });
  
  if (res.status === 401) {
    window.location.href = `/static/login.html?redirect=${encodeURIComponent(window.location.href)}`;
    throw new Error('Unauthorized');
  }
  
  if (!res.ok) {
    console.error('Ошибка загрузки результатов карт:', res.status);
    return { cards_results: [] };
  }
  
  return await res.json();
}

async function loadCardStats(resultIds, token) {
  let totalCorrect = 0;
  let totalIncorrect = 0;
  
  for (const resultId of resultIds) {
    const res = await fetch(`${API_BASE_URL}/api/v1/results/cards_result/${resultId}`, {
      headers: { 'Authorization': `Bearer ${token}` }
    });
    
    if (res.status === 401) {
      window.location.href = `/static/login.html?redirect=${encodeURIComponent(window.location.href)}`;
      return { correct: 0, incorrect: 0, total: 0 };
    }
    
    if (res.ok) {
      const data = await res.json();
      const correct = data.cards_results?.filter(c => c.result === 'correct').length || 0;
      const incorrect = data.cards_results?.filter(c => c.result !== 'correct').length || 0;
      totalCorrect += correct;
      totalIncorrect += incorrect;
    }
  }
  
  return { correct: totalCorrect, incorrect: totalIncorrect, total: totalCorrect + totalIncorrect };
}

function renderChart(stats) {
  const total = stats.total;
  if (total === 0) return;
  
  const correctPercent = (stats.correct / total) * 360;
  const incorrectPercent = 360 - correctPercent;
  
  const correctSlice = document.getElementById('correctSlice');
  const incorrectSlice = document.getElementById('incorrectSlice');
  
  correctSlice.style.background = '';
  incorrectSlice.style.background = '';
  
  if (correctPercent > 0) {
    correctSlice.style.background = 
      `conic-gradient(#51cf66 0deg, #51cf66 ${correctPercent}deg, transparent ${correctPercent}deg)`;
  }
  
  if (incorrectPercent > 0) {
    incorrectSlice.style.background = 
      `conic-gradient(#ff6b6b ${correctPercent}deg, #ff6b6b 360deg, transparent 360deg)`;
  }
  
  document.getElementById('statsTotal').textContent = total;
  document.getElementById('correctCount').textContent = stats.correct;
  document.getElementById('incorrectCount').textContent = stats.incorrect;
}

function formatDate(isoString) {
  const date = new Date(isoString);
  const day = date.getDate().toString().padStart(2, '0');
  const month = (date.getMonth() + 1).toString().padStart(2, '0');
  const year = date.getFullYear();
  const hours = date.getHours().toString().padStart(2, '0');
  const minutes = date.getMinutes().toString().padStart(2, '0');
  const seconds = date.getSeconds().toString().padStart(2, '0');
  return `${day}.${month}.${year}, ${hours}:${minutes}:${seconds}`;
}

function initNavigation() {
  const navToggle = document.getElementById('nav-toggle');
  const navPanel = document.getElementById('nav-panel');
  const navOverlay = document.getElementById('nav-panel-overlay');
  
  if (navToggle) {
    navToggle.addEventListener('click', () => {
      navPanel?.classList.toggle('open');
      navToggle.classList.toggle('open');
      navOverlay?.classList.toggle('open');
    });
  }
  
  if (navOverlay) {
    navOverlay.addEventListener('click', () => {
      navPanel?.classList.remove('open');
      navToggle?.classList.remove('open');
      navOverlay.classList.remove('open');
    });
  }

  const navButtons = {
    'main-btn': '/static/main.html',
    'modules-btn': '/static/modules.html',
    'categories-btn': '/static/categories.html',
    'results-btn': '/static/results.html',
    'selected-btn': '/static/selected.html'
  };

  Object.entries(navButtons).forEach(([id, url]) => {
    const btn = document.getElementById(id);
    if (btn) {
      btn.onclick = () => window.location.href = url;
    }
  });
}
