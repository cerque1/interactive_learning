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

  let currentUserId = null;
  let categoryOwnerId = null;
  let isEditMode = false;

  const addModuleBtn = document.getElementById('add-module-btn');
  const editCategoryBtn = document.getElementById('edit-category-btn');

  let userLoaded = false;
  let categoryLoaded = false;

  // Переменные для модального окна
  let availableModules = [];
  let selectedModuleIds = new Set();
  const addModuleModal = document.getElementById('add-module-modal');
  const closeModalBtn = document.getElementById('close-modal');
  const confirmAddModulesBtn = document.getElementById('confirm-add-modules');
  const availableModulesContainer = document.getElementById('available-modules-container');
  const noModulesMessage = document.getElementById('no-modules-message');

  // Глобальный объект для хранения данных модулей из модального окна
  const moduleDataMap = new Map();

  // Функция загрузки пользователя
  function fetchUser() {
    return fetch('http://localhost:8080/api/v1/user/me?is_full=f', {
      headers: { 'Authorization': `Bearer ${token}` }
    })
    .then(res => {
      if (res.status === 401) {
        window.location.href = '/static/login.html?redirect=' + encodeURIComponent(window.location.href);
        return Promise.reject();
      }
      return res.json();
    })
    .then(userData => {
      if (userData && userData.user) {
        currentUserId = userData.user.id;
        const usernameElem = document.getElementById('username');
        usernameElem.textContent = userData.user.name || 'Пользователь';
        usernameElem.onclick = () => {
          window.location.href = '/static/profile.html';
        };
      }
      userLoaded = true;
      checkShowButtons();
    })
    .catch(() => {
      userLoaded = true;
      checkShowButtons();
    });
  }

  // Функция загрузки категории
  function fetchCategory() {
    return fetch(`http://localhost:8080/api/v1/category/${categoryId}?is_full=t`, {
      headers: { 'Authorization': `Bearer ${token}` }
    })
    .then(res => {
      if (res.status === 401) {
        window.location.href = '/static/login.html?redirect=' + encodeURIComponent(window.location.href);
        return Promise.reject();
      }
      if (!res.ok) {
        throw new Error('Ошибка загрузки категории');
      }
      return res.json();
    })
    .then(categoryData => {
      categoryData = categoryData.category;
      categoryOwnerId = categoryData.user_id || categoryData.owner_id;

      const categoryNameElem = document.getElementById('category-name');
      const modulesContainer = document.getElementById('modules-container');
      const emptyMessage = document.getElementById('empty-message');

      categoryNameElem.textContent = categoryData.name || 'Без названия';

      if (!categoryData.modules || categoryData.modules.length === 0) {
        modulesContainer.innerHTML = '';
        emptyMessage.style.display = 'block';
      } else {
        emptyMessage.style.display = 'none';
        modulesContainer.innerHTML = '';
        categoryData.modules.forEach(module => createModuleCard(module, modulesContainer));
      }

      categoryLoaded = true;
      checkShowButtons();
    })
    .catch(err => {
      document.getElementById('category-name').textContent = 'Ошибка загрузки категории';
      document.getElementById('modules-container').innerHTML = '';
      categoryLoaded = true;
      checkShowButtons();
    });
  }

  function createModuleCard(module, container) {
    const cardCount = (module.cards && module.cards.length) || 0;
    const moduleElem = document.createElement('div');
    moduleElem.className = 'card';
    moduleElem.dataset.moduleId = module.id;
    moduleElem.innerHTML = `
      <div class="card-title">${module.name}</div>
      <div class="card-count">Карточек: ${cardCount}</div>
      <div class="module-actions">
        <button class="action-btn delete" title="Удалить из категории">×</button>
      </div>
    `;
    moduleElem.style.cursor = 'pointer';

    // Клик по карточке (переход в модуль)
    moduleElem.addEventListener('click', (e) => {
      if (!e.target.classList.contains('action-btn')) {
        window.location.href = `/static/module.html?module_id=${encodeURIComponent(module.id)}`;
      }
    });

    // Удаление модуля
    moduleElem.querySelector('.delete').addEventListener('click', (e) => {
      e.stopPropagation();
      const moduleId = parseInt(module.id);
      if (!moduleId || !confirm('Удалить модуль из категории?')) return;

      fetch(`http://localhost:8080/api/v1/category/${categoryId}/${moduleId}/delete`, {
        method: 'DELETE',
        headers: { 'Authorization': `Bearer ${token}` }
      })
      .then(res => {
        if (!res.ok) throw new Error('Ошибка удаления');
        moduleElem.style.opacity = '0';
        setTimeout(() => {
          moduleElem.remove();
          const remaining = container.querySelectorAll('.card');
          if (remaining.length === 0) {
            document.getElementById('empty-message').style.display = 'block';
          }
        }, 300);
      })
      .catch(err => {
        console.error(err);
        alert('Ошибка удаления модуля');
      });
    });

    container.appendChild(moduleElem);
    return moduleElem;
  }

  // Проверка показа кнопок
  function checkShowButtons() {
    if (userLoaded && categoryLoaded && currentUserId && categoryOwnerId && 
        Number(currentUserId) === Number(categoryOwnerId)) {
      addModuleBtn.style.display = 'inline-block';
      editCategoryBtn.style.display = 'inline-block';
    }
  }

  // Запуск запросов
  fetchUser();
  fetchCategory();

  // Функция загрузки доступных модулей
  function fetchAvailableModules() {
    const modulesContainer = document.getElementById('modules-container');
    const currentCategoryModules = Array.from(modulesContainer.querySelectorAll('.card')).map(card => 
      parseInt(card.dataset.moduleId)
    );

    return fetch(`http://localhost:8080/api/v1/module/to_user/${currentUserId}?with_cards=t`, {
      headers: { 'Authorization': `Bearer ${token}` }
    })
    .then(res => {
      if (res.status === 401) {
        window.location.href = '/static/login.html?redirect=' + encodeURIComponent(window.location.href);
        return Promise.reject();
      }
      if (!res.ok) {
        throw new Error('Ошибка загрузки модулей');
      }
      return res.json();
    })
    .then(data => {
      moduleDataMap.clear();
      availableModules = (data.modules || []).filter(module => {
        const moduleId = parseInt(module.id);
        if (!currentCategoryModules.includes(moduleId)) {
          moduleDataMap.set(moduleId.toString(), module); // Сохраняем полные данные модуля
          return true;
        }
        return false;
      });
      renderAvailableModules();
    })
    .catch(err => {
      console.error('Ошибка загрузки доступных модулей:', err);
      noModulesMessage.style.display = 'block';
      availableModulesContainer.innerHTML = '';
    });
  }

  // Отрисовка доступных модулей
  function renderAvailableModules() {
    availableModulesContainer.innerHTML = '';
    noModulesMessage.style.display = 'none';

    if (availableModules.length === 0) {
      noModulesMessage.style.display = 'block';
      confirmAddModulesBtn.disabled = true;
      return;
    }

    availableModules.forEach(module => {
      const cardCount = (module.cards && module.cards.length) || 0;
      const moduleElem = document.createElement('div');
      moduleElem.className = 'module-checkbox';
      moduleElem.innerHTML = `
        <input type="checkbox" id="module_${module.id}" value="${module.id}">
        <label for="module_${module.id}" class="module-checkbox-label">${module.name}</label>
        <span class="module-checkbox-count">Карточек: ${cardCount}</span>
      `;
      
      moduleElem.addEventListener('click', (e) => {
        if (e.target.tagName !== 'INPUT') {
          const checkbox = moduleElem.querySelector('input[type="checkbox"]');
          checkbox.checked = !checkbox.checked;
          updateSelection(checkbox);
        }
      });

      availableModulesContainer.appendChild(moduleElem);
    });

    availableModulesContainer.querySelectorAll('input[type="checkbox"]').forEach(checkbox => {
      checkbox.addEventListener('change', () => updateSelection(checkbox));
    });

    confirmAddModulesBtn.disabled = true;
  }

  // Обновление выбора
  function updateSelection(checkbox) {
    if (checkbox.checked) {
      selectedModuleIds.add(checkbox.value);
    } else {
      selectedModuleIds.delete(checkbox.value);
    }
    
    confirmAddModulesBtn.disabled = selectedModuleIds.size === 0;
  }

  // Обработчик кнопки добавления модуля
  addModuleBtn.addEventListener('click', () => {
    selectedModuleIds.clear();
    fetchAvailableModules();
    addModuleModal.style.display = 'flex';
  });

  // Закрытие модального окна
  function closeModal() {
    addModuleModal.style.display = 'none';
    selectedModuleIds.clear();
  }

  closeModalBtn.addEventListener('click', closeModal);
  addModuleModal.addEventListener('click', (e) => {
    if (e.target === addModuleModal) closeModal();
  });

  confirmAddModulesBtn.addEventListener('click', () => {
    if (selectedModuleIds.size === 0) return;

    const modulesIdsArray = Array.from(selectedModuleIds).map(id => parseInt(id));
    const modulesContainer = document.getElementById('modules-container');
    const emptyMessage = document.getElementById('empty-message');

    fetch(`http://localhost:8080/api/v1/category/${categoryId}/add_modules`, {
      method: 'POST',
      headers: { 
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({ 
        modules_ids: modulesIdsArray
      })
    })
    .then(res => {
      if (!res.ok) throw new Error('Ошибка добавления модулей');
      return res.json();
    })
    .then(() => {
      closeModal();
      
      if (emptyMessage.style.display !== 'none') {
        emptyMessage.style.display = 'none';
      }
      
      modulesIdsArray.forEach(moduleId => {
        const moduleData = moduleDataMap.get(moduleId.toString());
        if (moduleData) {
          createModuleCard(moduleData, modulesContainer);
        }
      });
    })
    .catch(err => {
      console.error('Ошибка добавления модулей:', err);
      alert('Ошибка при добавлении модулей');
    });
  });

  // Обработчик редактирования категории
  editCategoryBtn.addEventListener('click', () => {
    isEditMode = !isEditMode;
    
    const modules = document.querySelectorAll('.card');
    modules.forEach(module => {
      if (isEditMode) {
        module.classList.add('edit-mode');
        module.querySelector('.module-actions').classList.add('show');
      } else {
        module.classList.remove('edit-mode');
        module.querySelector('.module-actions').classList.remove('show');
      }
    });
    
    editCategoryBtn.textContent = isEditMode ? 'Сохранить изменения' : 'Редактировать категорию';
  });

  // Навигация
  const navToggle = document.getElementById('nav-toggle');
  const navPanel = document.getElementById('nav-panel');
  const navMainBtn = document.getElementById('main-btn');
  const navModulesBtn = document.getElementById('modules-btn');
  const navCategoriesBtn = document.getElementById('categories-btn');
  const navResultsBtn = document.getElementById('results-btn');
  const head = document.getElementById('head');

  if (navToggle && navPanel) {
    navToggle.addEventListener('click', () => {
      navPanel.classList.toggle('open');
      navToggle.classList.toggle('open');
    });
  }

  if (navMainBtn) navMainBtn.addEventListener('click', () => window.location.href = '/static/main.html');
  if (navModulesBtn) navModulesBtn.addEventListener('click', () => window.location.href = '/static/modules.html');
  if (navCategoriesBtn) navCategoriesBtn.addEventListener('click', () => window.location.href = '/static/categories.html');
  if (navResultsBtn) navResultsBtn.addEventListener('click', () => window.location.href = '/static/results.html');
  if (head) head.addEventListener('click', () => window.location.href = '/static/main.html');
});
