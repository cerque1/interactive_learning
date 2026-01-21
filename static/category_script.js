const API_BASE_URL = window.location.origin;
let isCategoryFavorited = false;
let favoriteCategories = new Set();

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
    document.getElementById('username').textContent = 'Ошибка загрузки';
    return;
  }

  let currentUserId = null;
  let categoryOwnerId = null;
  let categoryType = null;
  let isEditMode = false;
  let categoryModules = [];

  document.getElementById('username').textContent = 'Загрузка...';
  document.getElementById('category-name').textContent = 'Загрузка категории...';

  // DOM элементы
  const addModuleBtn = document.getElementById('add-module-btn');
  const studyModulesBtn = document.getElementById('study-modules-btn');
  const testModulesBtn = document.getElementById('test-modules-btn');
  const editCategoryBtn = document.getElementById('edit-category-btn');
  const deleteCategoryBtn = document.getElementById('delete-category-btn');
  const renameCategoryBtn = document.getElementById('rename-category-btn');
  const toggleCategoryTypeBtn = document.getElementById('toggle-category-type-btn');
  const toggleTypeText = document.getElementById('toggle-type-text');

  let userLoaded = false;
  let categoryLoaded = false;

  // Поиск модулей
  let isSearchMode = false;
  let searchOffset = 0;
  let searchQuery = '';
  let searchResults = [];
  let totalSearchResults = 0;
  let selectedModuleIds = new Set();

  // Модальные окна
  const addModuleModal = document.getElementById('add-module-modal');
  const closeModalBtn = document.getElementById('close-modal');
  const confirmAddModulesBtn = document.getElementById('confirm-add-modules');
  const availableModulesContainer = document.getElementById('available-modules-container');
  const noModulesMessage = document.getElementById('no-modules-message');

  const studyModal = document.getElementById('study-modal');
  const closeStudyModalBtn = document.getElementById('close-study-modal');
  const studyModulesContainer = document.getElementById('study-modules-container');
  const noStudyModulesMessage = document.getElementById('no-study-modules-message');
  const selectAllBtn = document.getElementById('select-all-btn');
  const deselectAllBtn = document.getElementById('deselect-all-btn');
  const startStudyingBtn = document.getElementById('start-studying-btn');

  const testModal = document.getElementById('test-modal');
  const closeTestModalBtn = document.getElementById('close-test-modal');
  const testModulesContainer = document.getElementById('test-modules-container');
  const noTestModulesMessage = document.getElementById('no-test-modules-message');
  const testSelectAllBtn = document.getElementById('test-select-all-btn');
  const testDeselectAllBtn = document.getElementById('test-deselect-all-btn');
  const startTestingBtn = document.getElementById('start-testing-btn');

  const moduleDataMap = new Map();
  let availableModules = [];
  let searchToggleBtn, searchInput, prevBtn, nextBtn;
  let searchTimeout = null;

  function safeAddEventListener(element, event, handler) {
    if (element) {
      element.addEventListener(event, handler);
    }
  }

  function safeSetDisabled(element, disabled) {
    if (element) {
      element.disabled = disabled;
    }
  }

  function showSuccessMessage(message) {
    const notification = document.createElement('div');
    notification.className = 'success-message';
    notification.textContent = message;
    document.body.appendChild(notification);
    
    requestAnimationFrame(() => {
      notification.style.opacity = '1';
      notification.style.transform = 'translateX(0)';
    });
    
    setTimeout(() => {
      notification.style.opacity = '0';
      notification.style.transform = 'translateX(100%)';
      setTimeout(() => document.body.removeChild(notification), 300);
    }, 3000);
  }

  function showErrorMessage(message) {
    const notification = document.createElement('div');
    notification.className = 'success-message error-message';
    notification.textContent = message;
    notification.style.background = '#dc3545';
    document.body.appendChild(notification);
    
    requestAnimationFrame(() => {
      notification.style.opacity = '1';
      notification.style.transform = 'translateX(0)';
    });
    
    setTimeout(() => {
      notification.style.opacity = '0';
      notification.style.transform = 'translateX(100%)';
      setTimeout(() => document.body.removeChild(notification), 300);
    }, 4000);
  }

  function fetchWithTimeout(url, options = {}, timeout = 8000) {
    return Promise.race([
      fetch(url, options),
      new Promise((_, reject) => 
        setTimeout(() => reject(new Error('Таймаут запроса')), timeout)
      )
    ]);
  }

  // Загрузка пользователя
  function fetchUser() {
    fetchWithTimeout(`${API_BASE_URL}/api/v1/user/me?is_full=f`, {
      headers: { 'Authorization': `Bearer ${token}` }
    })
    .then(res => {
      if (res.status === 401) {
        window.location.href = '/static/login.html?redirect=' + encodeURIComponent(window.location.href);
        return;
      }
      if (!res.ok) throw new Error(`HTTP ${res.status}`);
      return res.json();
    })
    .then(userData => {
      if (userData?.user?.id) {
        currentUserId = userData.user.id;
        document.getElementById('username').textContent = userData.user.name || 'Пользователь';
        safeAddEventListener(document.getElementById('username'), 'click', () => {
          window.location.href = '/static/profile.html';
        });
      }
      userLoaded = true;
      
      // Загружаем избранные категории после авторизации
      if (currentUserId && categoryLoaded) {
        fetchFavoriteCategories().then(() => {
          updateFavoriteButton();
          checkShowButtons();
        });
      } else {
        checkShowButtons();
      }
    })
    .catch(() => {
      document.getElementById('username').textContent = 'Гость';
      currentUserId = null;
      userLoaded = true;
      checkShowButtons();
    });
  } 

  // Загрузка категории
  function fetchCategory() {
    fetchWithTimeout(`${API_BASE_URL}/api/v1/category/${categoryId}?is_full=t`, {
      headers: { 'Authorization': `Bearer ${token}` }
    })
    .then(res => {
      if (res.status === 401) {
        window.location.href = '/static/login.html?redirect=' + encodeURIComponent(window.location.href);
        return;
      }

      if (res.status === 406) {
        document.getElementById('category-name').textContent = 'Категория недоступна';
        document.querySelector('.buttons-container')?.style.setProperty('display', 'none');
        
        const modulesContainer = document.getElementById('modules-container');
        if (modulesContainer) {
          modulesContainer.style.display = 'flex';
          modulesContainer.style.justifyContent = 'center';
          modulesContainer.style.alignItems = 'center';
          modulesContainer.style.minHeight = '400px';
          modulesContainer.innerHTML = `
            <div style="text-align: center; padding: 60px 20px; color: #94a3b8; font-size: 1.1em; max-width: 500px;">
              <p style="margin-bottom: 24px;">Эта категория недоступна для просмотра</p>
              <button onclick="window.location.href='/static/main.html'"
                      style="background: linear-gradient(135deg, #007bfb, #5ab9ea); color: white; border: none; padding: 14px 28px; border-radius: 12px; font-size: 16px; font-weight: 600; cursor: pointer; box-shadow: 0 4px 15px rgba(0,123,251,0.3); transition: all 0.3s ease;">
                Вернуться на главную
              </button>
            </div>
          `;
        }
        
        categoryModules = [];
        categoryOwnerId = currentUserId || 1;
        categoryLoaded = true;
        checkShowButtons();
        return undefined;
      }

      if (!res.ok) throw new Error(`HTTP ${res.status}`);
      return res.json();
    })
    .then(data => {
      const categoryData = data.category || data;
      categoryOwnerId = categoryData.user_id || categoryData.owner_id || currentUserId || 1;
      categoryType = categoryData.type || 0;

      if (currentUserId && userLoaded) {
        fetchFavoriteCategories().then(() => {
          updateFavoriteButton();
        });
      }

      isCategoryFavorited = categoryData.is_favorited || false;
      updateFavoriteButton();
      
      document.getElementById('category-name').textContent = categoryData.name || 'Без названия';
      categoryModules = categoryData.modules || [];
      
      const modulesContainer = document.getElementById('modules-container');
      const emptyMessage = document.getElementById('empty-message');
      
      if (modulesContainer && emptyMessage) {
        if (categoryModules.length === 0) {
          modulesContainer.innerHTML = '';
          emptyMessage.style.display = 'block';
        } else {
          emptyMessage.style.display = 'none';
          modulesContainer.innerHTML = '';
          categoryModules.forEach(module => createModuleCard(module, modulesContainer));
        }
      }
      
      categoryLoaded = true;
      checkShowButtons();
      if (toggleCategoryTypeBtn && toggleTypeText) updateCategoryTypeButton();
    })
    .catch(() => {
      document.getElementById('category-name').textContent = 'Категория не найдена';
      const modulesContainer = document.getElementById('modules-container');
      if (modulesContainer) modulesContainer.innerHTML = '';
      categoryModules = [];
      categoryOwnerId = currentUserId || 1;
      categoryLoaded = true;
      checkShowButtons();
    });
  }

  function checkShowButtons() {
    if (userLoaded && categoryLoaded) {
      const isOwner = currentUserId && categoryOwnerId && currentUserId == categoryOwnerId;
      
      if (studyModulesBtn) studyModulesBtn.style.display = categoryModules.length > 0 ? 'inline-block' : 'none';
      if (testModulesBtn) testModulesBtn.style.display = categoryModules.length > 0 ? 'inline-block' : 'none';
      
      const editButtonsContainer = document.getElementById('edit-buttons-container');
      if (editButtonsContainer) {
        editButtonsContainer.style.display = isOwner ? 'flex' : 'none';
      }
    }
  }

  function updateCategoryTypeButton() {
    if (!toggleCategoryTypeBtn || categoryType === null) return;
    
    if (categoryType === 0) {
      toggleTypeText.textContent = 'Сделать приватной';
      toggleCategoryTypeBtn.title = 'Изменить категорию на приватную (только для владельца)';
    } else {
      toggleTypeText.textContent = 'Сделать публичной';
      toggleCategoryTypeBtn.title = 'Изменить категорию на публичную (доступна всем)';
    }
  }

  // Создание карточки модуля
  function createModuleCard(module, container) {
    const cardCount = (module.cards?.length || 0);
    const moduleElem = document.createElement('div');
    moduleElem.className = 'card';
    moduleElem.dataset.moduleId = module.id;
    moduleElem.innerHTML = `
      <div class="card-title">${module.name || 'Без названия'}</div>
      <div class="card-count">Карточек: ${cardCount}</div>
      <div class="module-actions">
        <button class="action-btn delete" title="Удалить из категории">×</button>
      </div>
    `;
    moduleElem.style.cursor = 'pointer';

    safeAddEventListener(moduleElem, 'click', (e) => {
      if (!e.target.classList.contains('action-btn')) {
        window.location.href = `/static/module.html?module_id=${module.id}`;
      }
    });

    const deleteBtn = moduleElem.querySelector('.delete');
    if (deleteBtn) {
      safeAddEventListener(deleteBtn, 'click', (e) => {
        e.stopPropagation();
        if (!confirm('Удалить модуль из категории?')) return;
        
        fetch(`${API_BASE_URL}/api/v1/category/${categoryId}/${module.id}/delete`, {
          method: 'DELETE',
          headers: { 'Authorization': `Bearer ${token}` }
        })
        .then(res => {
          if (!res.ok) throw new Error('Ошибка удаления');
          moduleElem.style.opacity = '0';
          moduleElem.style.transform = 'scale(0.8)';
          setTimeout(() => {
            moduleElem.remove();
            const remaining = container.querySelectorAll('.card');
            const emptyMessage = document.getElementById('empty-message');
            if (remaining.length === 0 && emptyMessage) {
              emptyMessage.style.display = 'block';
            }
            categoryModules = categoryModules.filter(m => m.id != module.id);
            checkShowButtons();
          }, 300);
        })
        .catch(() => showErrorMessage('Ошибка удаления модуля'));
      });
    }

    container.appendChild(moduleElem);
    return moduleElem;
  }

  // Загрузка избранных категорий пользователя
  function fetchFavoriteCategories() {
    if (!currentUserId) return Promise.resolve();
    
    return fetchWithTimeout(`${API_BASE_URL}/api/v1/selected/categories/`, {
      headers: { 'Authorization': `Bearer ${token}` }
    })
    .then(res => {
      if (res.status === 401) {
        window.location.href = '/static/login.html?redirect=' + encodeURIComponent(window.location.href);
        return Promise.reject();
      }
      if (!res.ok) throw new Error('Ошибка загрузки избранных категорий');
      return res.json();
    })
    .then(data => {
      favoriteCategories.clear();
      const categories = data.selected_categories || data || [];
      categories.forEach(cat => {
        favoriteCategories.add(cat.id.toString());
      });
      isCategoryFavorited = favoriteCategories.has(categoryId);
    })
    .catch(err => {
      console.error('Ошибка загрузки избранных:', err);
      // Не прерываем загрузку страницы при ошибке избранного
    });
  }


  // Обновление кнопки избранного
  function updateFavoriteButton() {
    const btn = document.getElementById('toggle-favorite-btn');
    if (!btn) return;
    
    if (isCategoryFavorited) {
      btn.classList.add('favorited');
      btn.innerHTML = '<span class="favorite-text">Убрать из избранного</span>';
      btn.title = 'Убрать из избранного';
    } else {
      btn.classList.remove('favorited');
      btn.innerHTML = '<span class="favorite-text">В избранное</span>';
      btn.title = 'Добавить в избранное';
    }
  }

  // Переключение избранного
  function toggleCategoryFavorite() {
    if (!currentUserId) {
      showErrorMessage('Необходимо авторизоваться');
      return;
    }

    const isFavorited = isCategoryFavorited;
    const url = isFavorited 
      ? `${API_BASE_URL}/api/v1/selected/categories/delete?category_id=${categoryId}`
      : `${API_BASE_URL}/api/v1/selected/categories/insert?category_id=${categoryId}`;

    const method = isFavorited
      ? "DELETE"
      : "POST";

    fetch(url, {
      method: method,
      headers: { 
        'Authorization': `Bearer ${token}`
      }
    })
    .then(res => {
      if (res.status === 401) {
        window.location.href = '/static/login.html?redirect=' + encodeURIComponent(window.location.href);
        return Promise.reject();
      }
      if (!res.ok) {
        if (res.status === 403) throw new Error('Нет прав для изменения избранного');
        throw new Error(isFavorited ? 'Ошибка удаления из избранного' : 'Ошибка добавления в избранное');
      }
      return;
    })
    .then(() => {
      // Обновляем локальное состояние
      if (isFavorited) {
        favoriteCategories.delete(categoryId);
        isCategoryFavorited = false;
      } else {
        favoriteCategories.add(categoryId);
        isCategoryFavorited = true;
      }
      updateFavoriteButton();
      
      showSuccessMessage(
        isCategoryFavorited 
          ? 'Категория добавлена в избранное' 
          : 'Категория убрана из избранного'
      );
    })
    .catch(err => showErrorMessage(err.message));
  }

  // ===== ПОИСК МОДУЛЕЙ =====
  function createSearchUI() {
    const searchSection = document.createElement('div');
    searchSection.className = 'search-section';
    searchSection.innerHTML = `
      <div class="tab-buttons" style="display: flex; gap: 8px; margin-bottom: 12px;">
        <button id="my-modules-tab" class="tab-btn active">Мои модули</button>
        <button id="search-tab" class="tab-btn">Поиск</button>
      </div>
      <div id="search-input-container" style="display: none;">
        <div style="display: flex; gap: 8px; align-items: end;">
          <input type="text" id="search-input" class="search-input" placeholder="Введите название модуля (минимум 2 символа)" maxlength="100">
          <button id="search-btn" class="secondary-btn">Поиск</button>
        </div>
        <div class="pagination-buttons">
          <button id="search-prev-btn" class="pagination-btn" disabled>Назад</button>
          <button id="search-next-btn" class="pagination-btn" disabled>Далее</button>
        </div>
      </div>
    `;
    
    const modalBody = addModuleModal.querySelector('.modal-body');
    const checkboxContainer = document.getElementById('available-modules-container');
    modalBody.insertBefore(searchSection, checkboxContainer);
    
    searchToggleBtn = document.getElementById('my-modules-tab');
    const searchTabBtn = document.getElementById('search-tab');
    searchInput = document.getElementById('search-input');
    prevBtn = document.getElementById('search-prev-btn');
    nextBtn = document.getElementById('search-next-btn');
    const searchBtn = document.getElementById('search-btn');
    
    safeAddEventListener(searchToggleBtn, 'click', () => switchTab('my-modules'));
    safeAddEventListener(searchTabBtn, 'click', () => switchTab('search'));
    safeAddEventListener(searchBtn, 'click', performSearch);
    safeAddEventListener(prevBtn, 'click', () => loadSearchResults(false));
    safeAddEventListener(nextBtn, 'click', () => loadSearchResults(true));
    safeAddEventListener(searchInput, 'keypress', (e) => {
      if (e.key === 'Enter') performSearch();
    });
  }

  function switchTab(activeTab) {
    const myModulesTab = document.getElementById('my-modules-tab');
    const searchTab = document.getElementById('search-tab');
    const inputContainer = document.getElementById('search-input-container');
    
    // Обновляем активные вкладки
    if (activeTab === 'my-modules') {
      myModulesTab.classList.add('active');
      searchTab.classList.remove('active');
      inputContainer.style.display = 'none';
      selectedModuleIds.clear();
      loadUserModules();
      isSearchMode = false;
    } else {
      myModulesTab.classList.remove('active');
      searchTab.classList.add('active');
      inputContainer.style.display = 'block';
      searchInput.value = '';
      searchOffset = 0;
      searchInput.focus();
      loadPopularModules();
      isSearchMode = true;
    }
  }

  function performSearch() {
    const query = searchInput.value.trim();
    if (query.length < 2) {
      showErrorMessage('Введите минимум 2 символа для поиска');
      return;
    }
    
    searchQuery = query;
    searchOffset = 0;
    loadSearchResults(false);
  }

  function loadPopularModules() {
    fetchWithTimeout(`${API_BASE_URL}/api/v1/module/popular?limit=10&offset=0`, {
      headers: { 'Authorization': `Bearer ${token}` }
    })
    .then(res => {
      if (!res.ok) throw new Error('Ошибка загрузки популярных модулей');
      return res.json();
    })
    .then(data => {
      searchResults = data.popular_modules || [];
      totalSearchResults = searchResults.length;
      renderAvailableModules(searchResults.filter(filterCurrentCategoryModules));
      updatePaginationButtons();
    })
    .catch(() => {
      noModulesMessage.style.display = 'block';
      availableModulesContainer.innerHTML = '';
    });
  }

  function loadSearchResults(nextPage = true) {
    if (nextPage) searchOffset += 10;
    else searchOffset = Math.max(0, searchOffset - 10);

    const params = new URLSearchParams({
      name: searchQuery,
      limit: '10',
      offset: searchOffset.toString()
    });

    fetchWithTimeout(`${API_BASE_URL}/api/v1/search/modules?${params}`, {
      headers: { 'Authorization': `Bearer ${token}` }
    })
    .then(res => {
      if (!res.ok) throw new Error('Ошибка поиска');
      return res.json();
    })
    .then(data => {
      searchResults = data.found_modules || [];
      totalSearchResults = data.total || 0;
      renderAvailableModules(searchResults.filter(filterCurrentCategoryModules));
      updatePaginationButtons();
    })
    .catch(() => {
      noModulesMessage.style.display = 'block';
      availableModulesContainer.innerHTML = '';
    });
  }

  function filterCurrentCategoryModules(module) {
    const modulesContainer = document.getElementById('modules-container');
    if (!modulesContainer) return true;
    const currentCategoryModules = Array.from(modulesContainer.querySelectorAll('.card')).map(card => parseInt(card.dataset.moduleId));
    const module_id = module.module?.id ? parseInt(module.module.id) : parseInt(module.id);
    return !currentCategoryModules.includes(module_id);
  }

  function updatePaginationButtons() {
    if (!prevBtn || !nextBtn) return;
    
    // ИСПРАВЛЕНИЕ: отключаем "Назад" при пустом результате или offset=0
    prevBtn.disabled = searchOffset === 0 || searchResults.length === 0;
    nextBtn.disabled = searchResults.length < 10 || totalSearchResults <= searchOffset + 10;
  }

  function loadUserModules() {
    if (!currentUserId) return;
    const modulesContainer = document.getElementById('modules-container');
    if (!modulesContainer) return;
    
    const currentCategoryModules = Array.from(modulesContainer.querySelectorAll('.card')).map(card => parseInt(card.dataset.moduleId));

    fetchWithTimeout(`${API_BASE_URL}/api/v1/module/to_user/${currentUserId}?with_cards=t`, {
      headers: { 'Authorization': `Bearer ${token}` }
    })
    .then(res => {
      if (!res.ok) throw new Error('Ошибка загрузки');
      return res.json();
    })
    .then(data => {
      moduleDataMap.clear();
      availableModules = (data.modules || []).filter(module => {
        const moduleId = parseInt(module.id);
        if (!currentCategoryModules.includes(moduleId)) {
          moduleDataMap.set(moduleId.toString(), module);
          return true;
        }
        return false;
      });
      renderAvailableModules(availableModules);
    })
    .catch(() => {
      noModulesMessage.style.display = 'block';
      availableModulesContainer.innerHTML = '';
    });
  }

  function renderAvailableModules(modules) {
    availableModulesContainer.innerHTML = '';
    noModulesMessage.style.display = 'none';
    
    if (modules.length === 0) {
      noModulesMessage.style.display = 'block';
      safeSetDisabled(confirmAddModulesBtn, true);
      return;
    }

    modules.forEach(module => {
      const moduleId = module.module?.id || module.id;
      const moduleName = module.module?.name || module.name;
      const countLabel = module.count ? `Пользователей за неделю: ${module.count}` : '';

      const moduleElem = document.createElement('div');
      moduleElem.className = 'module-checkbox';
      moduleElem.innerHTML = `
        <input type="checkbox" id="module_${moduleId}" value="${moduleId}">
        <label for="module_${moduleId}" class="module-checkbox-label">${moduleName}</label>
        <span class="module-checkbox-count">${countLabel}</span>
      `;

      safeAddEventListener(moduleElem, 'click', (e) => {
        if (e.target.tagName !== 'INPUT') {
          const checkbox = moduleElem.querySelector('input[type="checkbox"]');
          if (checkbox) {
            checkbox.checked = !checkbox.checked;
            updateSelection(checkbox);
          }
        }
      });
      availableModulesContainer.appendChild(moduleElem);
    });

    availableModulesContainer.querySelectorAll('input[type="checkbox"]').forEach(checkbox => {
      safeAddEventListener(checkbox, 'change', () => updateSelection(checkbox));
    });
    safeSetDisabled(confirmAddModulesBtn, selectedModuleIds.size === 0);
  }

  function updateSelection(checkbox) {
    if (checkbox.checked) {
      selectedModuleIds.add(checkbox.value);
    } else {
      selectedModuleIds.delete(checkbox.value);
    }
    safeSetDisabled(confirmAddModulesBtn, selectedModuleIds.size === 0);
  }

  // ===== Модальное окно заучивания =====
  function renderStudyModules() {
    if (!studyModulesContainer) return;
    studyModulesContainer.innerHTML = '';
    if (noStudyModulesMessage) noStudyModulesMessage.style.display = 'none';
    selectedModuleIds.clear();
    
    if (categoryModules.length === 0) {
      if (noStudyModulesMessage) noStudyModulesMessage.style.display = 'block';
      if (startStudyingBtn) safeSetDisabled(startStudyingBtn, true);
      return;
    }

    categoryModules.forEach(module => {
      const cardCount = (module.cards?.length || 0);
      const moduleElem = document.createElement('div');
      moduleElem.className = 'module-checkbox';
      moduleElem.innerHTML = `
        <input type="checkbox" id="study_module_${module.id}" value="${module.id}">
        <label for="study_module_${module.id}" class="module-checkbox-label">${module.name}</label>
        <span class="module-checkbox-count">Карточек: ${cardCount}</span>
      `;
      safeAddEventListener(moduleElem, 'click', (e) => {
        if (e.target.tagName !== 'INPUT') {
          const checkbox = moduleElem.querySelector('input[type="checkbox"]');
          if (checkbox) {
            checkbox.checked = !checkbox.checked;
            updateStudySelection(checkbox);
          }
        }
      });
      studyModulesContainer.appendChild(moduleElem);
    });

    studyModulesContainer.querySelectorAll('input[type="checkbox"]').forEach(checkbox => {
      safeAddEventListener(checkbox, 'change', () => updateStudySelection(checkbox));
    });
    if (startStudyingBtn) safeSetDisabled(startStudyingBtn, true);
  }

  function updateStudySelection(checkbox) {
    if (checkbox.checked) selectedModuleIds.add(checkbox.value);
    else selectedModuleIds.delete(checkbox.value);
    if (startStudyingBtn) safeSetDisabled(startStudyingBtn, selectedModuleIds.size === 0);
  }

  // ===== Модальное окно тестирования =====
  function renderTestModules() {
    if (!testModulesContainer) return;
    testModulesContainer.innerHTML = '';
    if (noTestModulesMessage) noTestModulesMessage.style.display = 'none';
    selectedModuleIds.clear();
    
    if (categoryModules.length === 0) {
      if (noTestModulesMessage) noTestModulesMessage.style.display = 'block';
      if (startTestingBtn) safeSetDisabled(startTestingBtn, true);
      return;
    }

    categoryModules.forEach(module => {
      const cardCount = (module.cards?.length || 0);
      const moduleElem = document.createElement('div');
      moduleElem.className = 'module-checkbox';
      moduleElem.innerHTML = `
        <input type="checkbox" id="test_module_${module.id}" value="${module.id}">
        <label for="test_module_${module.id}" class="module-checkbox-label">${module.name}</label>
        <span class="module-checkbox-count">Карточек: ${cardCount}</span>
      `;
      safeAddEventListener(moduleElem, 'click', (e) => {
        if (e.target.tagName !== 'INPUT') {
          const checkbox = moduleElem.querySelector('input[type="checkbox"]');
          if (checkbox) {
            checkbox.checked = !checkbox.checked;
            updateTestSelection(checkbox);
          }
        }
      });
      testModulesContainer.appendChild(moduleElem);
    });

    testModulesContainer.querySelectorAll('input[type="checkbox"]').forEach(checkbox => {
      safeAddEventListener(checkbox, 'change', () => updateTestSelection(checkbox));
    });
    if (startTestingBtn) safeSetDisabled(startTestingBtn, true);
  }

  function updateTestSelection(checkbox) {
    if (checkbox.checked) selectedModuleIds.add(checkbox.value);
    else selectedModuleIds.delete(checkbox.value);
    if (startTestingBtn) safeSetDisabled(startTestingBtn, selectedModuleIds.size === 0);
  }

  // ===== УПРАВЛЕНИЕ КАТЕГОРИЕЙ =====
  function deleteCategory() {
    if (!confirm('Вы уверены, что хотите удалить эту категорию? Все модули из неё будут удалены.')) return;
    
    fetch(`${API_BASE_URL}/api/v1/category/delete/${categoryId}`, {
      method: 'DELETE',
      headers: { 'Authorization': `Bearer ${token}` }
    })
    .then(res => {
      if (!res.ok) {
        if (res.status === 403) throw new Error('Нет прав для удаления категории');
        throw new Error('Ошибка удаления категории');
      }
      showSuccessMessage('Категория успешно удалена');
      setTimeout(() => {
        window.location.href = '/static/categories.html';
      }, 1500);
    })
    .catch(err => showErrorMessage('Ошибка удаления: ' + err.message));
  }

  // Модальное окно переименования категории
  function showRenameCategoryModal() {
    const modal = document.getElementById('renameCategoryModal');
    const input = document.getElementById('rename-category-input');
    const confirmBtn = document.getElementById('confirm-rename-category-btn');
    const errorDiv = document.getElementById('rename-category-error');
    const currentName = document.getElementById('category-name').textContent;
    
    // Сбрасываем состояние
    input.value = currentName;
    errorDiv.style.display = 'none';
    input.disabled = false;
    confirmBtn.disabled = true;
    
    modal.style.display = 'flex';
    input.focus();
    input.select();
    
    // Валидация ввода
    function validateInput() {
      const newName = input.value.trim();
      confirmBtn.disabled = newName.length === 0 || newName === currentName;
      errorDiv.style.display = 'none';
    }
    
    input.addEventListener('input', validateInput);
    
    // Подтверждение переименования
    const confirmHandler = () => {
      const newName = input.value.trim();
      if (newName.length === 0 || newName === currentName) return;
      
      input.disabled = true;
      confirmBtn.disabled = true;
      
      fetch(`${API_BASE_URL}/api/v1/category/rename/${categoryId}`, {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({ name: newName })
      })
      .then(res => {
        if (!res.ok) {
          if (res.status === 403) throw new Error('Нет прав для переименования');
          if (res.status === 409) throw new Error('Категория с таким названием уже существует');
          throw new Error('Ошибка переименования категории');
        }
        return res.json();
      })
      .then(() => {
        document.getElementById('category-name').textContent = newName;
        modal.style.display = 'none';
        showSuccessMessage('Категория успешно переименована');
      })
      .catch(err => {
        input.disabled = false;
        confirmBtn.disabled = false;
        errorDiv.textContent = err.message;
        errorDiv.style.display = 'block';
        input.focus();
      });
    };
    
    // Удаляем старые обработчики
    const oldConfirmBtn = document.getElementById('confirm-rename-category-btn');
    const oldCancelBtn = document.getElementById('cancel-rename-category-modal');
    const oldCloseBtn = document.getElementById('closeRenameCategoryModal');
    
    oldConfirmBtn.replaceWith(oldConfirmBtn.cloneNode(true));
    oldCancelBtn.replaceWith(oldCancelBtn.cloneNode(true));
    oldCloseBtn.replaceWith(oldCloseBtn.cloneNode(true));
    
    // Новые обработчики
    safeAddEventListener(document.getElementById('confirm-rename-category-btn'), 'click', confirmHandler);
    safeAddEventListener(document.getElementById('cancel-rename-category-modal'), 'click', () => {
      modal.style.display = 'none';
    });
    safeAddEventListener(document.getElementById('closeRenameCategoryModal'), 'click', () => {
      modal.style.display = 'none';
    });
    
    // Закрытие по клику вне модалки
    safeAddEventListener(modal, 'click', (e) => {
      if (e.target === modal) modal.style.display = 'none';
    });
  }


  // ===== ОБРАБОТЧИКИ СОБЫТИЙ =====
  safeAddEventListener(addModuleBtn, 'click', () => {
    selectedModuleIds.clear();
    if (!document.getElementById('my-modules-tab')) {
      createSearchUI();
    }
    // По умолчанию открываем "Мои модули"
    switchTab('my-modules');
    addModuleModal.style.display = 'flex';
  });

  safeAddEventListener(studyModulesBtn, 'click', () => {
    renderStudyModules();
    studyModal.style.display = 'flex';
  });

  safeAddEventListener(testModulesBtn, 'click', () => {
    renderTestModules();
    testModal.style.display = 'flex';
  });

  // Кнопки управления категорией
  safeAddEventListener(deleteCategoryBtn, 'click', deleteCategory);
  safeAddEventListener(renameCategoryBtn, 'click', showRenameCategoryModal);

  // После обработчиков других кнопок
  safeAddEventListener(document.getElementById('toggle-favorite-btn'), 'click', toggleCategoryFavorite);


  // Редактирование категории
  safeAddEventListener(editCategoryBtn, 'click', () => {
    isEditMode = !isEditMode;
    const modules = document.querySelectorAll('.card');
    modules.forEach(moduleElem => {
      const actions = moduleElem.querySelector('.module-actions');
      if (isEditMode) {
        moduleElem.classList.add('edit-mode');
        if (actions) actions.classList.add('show');
      } else {
        moduleElem.classList.remove('edit-mode');
        if (actions) actions.classList.remove('show');
      }
    });
    editCategoryBtn.textContent = isEditMode ? 'Сохранить изменения' : 'Редактировать';
  });

  // Переключение типа категории
  safeAddEventListener(toggleCategoryTypeBtn, 'click', () => {
    if (!(currentUserId && categoryOwnerId && currentUserId == categoryOwnerId)) return;

    const newType = categoryType === 0 ? 1 : 0;
    
    fetch(`${API_BASE_URL}/api/v1/category/change_type/${categoryId}`, {
      method: 'PUT',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({ type: newType })
    })
    .then(res => {
      if (res.status === 401) {
        window.location.href = '/static/login.html?redirect=' + encodeURIComponent(window.location.href);
        return Promise.reject();
      }
      if (!res.ok) {
        if (res.status === 403) throw new Error('Нет прав для изменения типа категории');
        if (res.status === 404) throw new Error('Категория не найдена');
        throw new Error('Ошибка изменения типа категории');
      }
    })
    .then(() => {
      categoryType = newType;
      updateCategoryTypeButton();
      showSuccessMessage(
        newType === 0 
          ? 'Категория теперь публичная (доступна всем)' 
          : 'Категория теперь приватная (только для владельца)'
      );
    })
    .catch(err => showErrorMessage('Ошибка изменения типа категории: ' + err.message));
  });

  // Подтверждение добавления модулей
  safeAddEventListener(confirmAddModulesBtn, 'click', () => {
    if (selectedModuleIds.size === 0) return;
    const modulesIdsArray = Array.from(selectedModuleIds).map(id => parseInt(id));
    const modulesContainer = document.getElementById('modules-container');
    const emptyMessage = document.getElementById('empty-message');

    fetch(`${API_BASE_URL}/api/v1/category/${categoryId}/add_modules`, {
      method: 'POST',
      headers: { 
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({ modules_ids: modulesIdsArray })
    })
    .then(res => {
      if (!res.ok) throw new Error('Ошибка добавления');
      return res.json();
    })
    .then(() => {
      addModuleModal.style.display = 'none';
      if (emptyMessage && emptyMessage.style.display !== 'none') emptyMessage.style.display = 'none';
      if (modulesContainer) {
        modulesIdsArray.forEach(moduleId => {
          const moduleData = moduleDataMap.get(moduleId.toString()) || searchResults.find(m => m.id == moduleId);
          if (moduleData) {
            createModuleCard(moduleData, modulesContainer);
            categoryModules.push(moduleData);
          }
        });
      }
      checkShowButtons();
      showSuccessMessage('Модули успешно добавлены');
    })
    .catch(() => showErrorMessage('Ошибка добавления модулей'));
  });

  // Запуск заучивания/тестирования
  safeAddEventListener(startStudyingBtn, 'click', () => {
    if (selectedModuleIds.size === 0) return;
    const modulesIdsString = Array.from(selectedModuleIds).join(',');
    window.location.href = `/static/learning.html?category_id=${categoryId}&modules_ids=${modulesIdsString}`;
  });

  safeAddEventListener(startTestingBtn, 'click', () => {
    if (selectedModuleIds.size === 0) return;
    const modulesIdsString = Array.from(selectedModuleIds).join(',');
    window.location.href = `/static/test.html?category_id=${categoryId}&modules_ids=${modulesIdsString}`;
  });

  // Выбор всех/снятие выбора
  safeAddEventListener(selectAllBtn, 'click', () => {
    if (studyModulesContainer) {
      studyModulesContainer.querySelectorAll('input[type="checkbox"]').forEach(checkbox => {
        checkbox.checked = true;
        selectedModuleIds.add(checkbox.value);
      });
      safeSetDisabled(startStudyingBtn, false);
    }
  });

  safeAddEventListener(deselectAllBtn, 'click', () => {
    if (studyModulesContainer) {
      studyModulesContainer.querySelectorAll('input[type="checkbox"]').forEach(checkbox => {
        checkbox.checked = false;
        selectedModuleIds.delete(checkbox.value);
      });
      safeSetDisabled(startStudyingBtn, true);
    }
  });

  safeAddEventListener(testSelectAllBtn, 'click', () => {
    if (testModulesContainer) {
      testModulesContainer.querySelectorAll('input[type="checkbox"]').forEach(checkbox => {
        checkbox.checked = true;
        selectedModuleIds.add(checkbox.value);
      });
      safeSetDisabled(startTestingBtn, false);
    }
  });

  safeAddEventListener(testDeselectAllBtn, 'click', () => {
    if (testModulesContainer) {
      testModulesContainer.querySelectorAll('input[type="checkbox"]').forEach(checkbox => {
        checkbox.checked = false;
        selectedModuleIds.delete(checkbox.value);
      });
      safeSetDisabled(startTestingBtn, true);
    }
  });

  // Закрытие модальных окон
  safeAddEventListener(closeModalBtn, 'click', () => {
    addModuleModal.style.display = 'none';
  });
  safeAddEventListener(addModuleModal, 'click', (e) => {
    if (e.target === addModuleModal) addModuleModal.style.display = 'none';
  });

  safeAddEventListener(closeStudyModalBtn, 'click', () => {
    studyModal.style.display = 'none';
  });
  safeAddEventListener(studyModal, 'click', (e) => {
    if (e.target === studyModal) studyModal.style.display = 'none';
  });

  safeAddEventListener(closeTestModalBtn, 'click', () => {
    testModal.style.display = 'none';
  });
  safeAddEventListener(testModal, 'click', (e) => {
    if (e.target === testModal) testModal.style.display = 'none';
  });

  // Навигация
  const navToggle = document.getElementById('nav-toggle');
  const navPanel = document.getElementById('nav-panel');
  const navMainBtn = document.getElementById('main-btn');
  const navModulesBtn = document.getElementById('modules-btn');
  const navCategoriesBtn = document.getElementById('categories-btn');
  const navSelectedBtn = document.getElementById('selected-btn');
  const navResultsBtn = document.getElementById('results-btn');
  const head = document.getElementById('head');

  safeAddEventListener(navToggle, 'click', () => {
    if (navPanel) navPanel.classList.toggle('open');
    if (navToggle) navToggle.classList.toggle('open');
  });

  safeAddEventListener(navMainBtn, 'click', () => window.location.href = '/static/main.html');
  safeAddEventListener(navModulesBtn, 'click', () => window.location.href = '/static/modules.html');
  safeAddEventListener(navCategoriesBtn, 'click', () => window.location.href = '/static/categories.html');
  safeAddEventListener(navResultsBtn, 'click', () => window.location.href = '/static/results.html');
  safeAddEventListener(navSelectedBtn, 'click', () => window.location.href = '/static/selected.html');
  safeAddEventListener(head, 'click', () => window.location.href = '/static/main.html');

  // Закрытие по Escape
  document.addEventListener('keydown', (e) => {
    if (e.key === 'Escape') {
      addModuleModal.style.display = 'none';
      studyModal.style.display = 'none';
      testModal.style.display = 'none';
    }
  });

  // Таймаут для загрузки
  setTimeout(() => {
    if (!userLoaded) {
      userLoaded = true;
      currentUserId = 1;
      document.getElementById('username').textContent = 'Гость';
    }
    if (!categoryLoaded) {
      categoryLoaded = true;
      categoryOwnerId = 1;
      document.getElementById('category-name').textContent = 'Категория не найдена';
    }
    checkShowButtons();
  }, 5000);

  fetchUser();
  fetchCategory();
});
