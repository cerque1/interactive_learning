const API_BASE_URL = window.location.origin;

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
  let categoryType = null; // 0 - публичная, 1 - приватная
  let isEditMode = false;
  let categoryModules = [];

  document.getElementById('username').textContent = 'Загрузка...';
  document.getElementById('category-name').textContent = 'Загрузка категории...';

  // Все элементы DOM
  const addModuleBtn = document.getElementById('add-module-btn');
  const studyModulesBtn = document.getElementById('study-modules-btn');
  const testModulesBtn = document.getElementById('test-modules-btn');
  const editCategoryBtn = document.getElementById('edit-category-btn');
  const toggleCategoryTypeBtn = document.getElementById('toggle-category-type-btn');
  const toggleTypeText = document.getElementById('toggle-type-text');

  let userLoaded = false;
  let categoryLoaded = false;

  let availableModules = [];
  let selectedModuleIds = new Set();
  
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

  // Утилиты
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

  function fetchWithTimeout(url, options = {}, timeout = 5000) {
    return Promise.race([
      fetch(url, options),
      new Promise((_, reject) => 
        setTimeout(() => reject(new Error('Таймаут')), timeout)
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
      checkShowButtons();
    })
    .catch(() => {
      document.getElementById('username').textContent = 'Гость';
      currentUserId = 1;
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
      
      if (addModuleBtn) addModuleBtn.style.display = isOwner ? 'inline-block' : 'none';
      if (editCategoryBtn) editCategoryBtn.style.display = isOwner ? 'inline-block' : 'none';
      if (toggleCategoryTypeBtn) toggleCategoryTypeBtn.style.display = isOwner ? 'inline-block' : 'none';
      if (studyModulesBtn) studyModulesBtn.style.display = categoryModules.length > 0 ? 'inline-block' : 'none';
      if (testModulesBtn) testModulesBtn.style.display = categoryModules.length > 0 ? 'inline-block' : 'none';
    }
  }

  function updateCategoryTypeButton() {
    if (!toggleCategoryTypeBtn || categoryType === null) return;
    
    if (categoryType === 0) { // публичная
      toggleTypeText.textContent = 'Сделать приватной';
      toggleCategoryTypeBtn.title = 'Изменить категорию на приватную (только для владельца)';
    } else { // приватная
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
        .catch(() => alert('Ошибка удаления модуля'));
      });
    }

    container.appendChild(moduleElem);
    return moduleElem;
  }

  // Модальное окно добавления модулей
  function fetchAvailableModules() {
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
      renderAvailableModules();
    })
    .catch(() => {
      if (noModulesMessage) noModulesMessage.style.display = 'block';
      if (availableModulesContainer) availableModulesContainer.innerHTML = '';
    });
  }

  function renderAvailableModules() {
    if (!availableModulesContainer) return;
    availableModulesContainer.innerHTML = '';
    if (noModulesMessage) noModulesMessage.style.display = 'none';
    
    if (availableModules.length === 0) {
      if (noModulesMessage) noModulesMessage.style.display = 'block';
      if (confirmAddModulesBtn) safeSetDisabled(confirmAddModulesBtn, true);
      return;
    }

    availableModules.forEach(module => {
      const cardCount = (module.cards?.length || 0);
      const moduleElem = document.createElement('div');
      moduleElem.className = 'module-checkbox';
      moduleElem.innerHTML = `
        <input type="checkbox" id="module_${module.id}" value="${module.id}">
        <label for="module_${module.id}" class="module-checkbox-label">${module.name}</label>
        <span class="module-checkbox-count">Карточек: ${cardCount}</span>
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
    if (confirmAddModulesBtn) safeSetDisabled(confirmAddModulesBtn, true);
  }

  function updateSelection(checkbox) {
    if (checkbox.checked) selectedModuleIds.add(checkbox.value);
    else selectedModuleIds.delete(checkbox.value);
    if (confirmAddModulesBtn) safeSetDisabled(confirmAddModulesBtn, selectedModuleIds.size === 0);
  }

  // Модальное окно заучивания
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

  // Модальное окно тестирования
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

  // Обработчики событий
  if (addModuleBtn) {
    safeAddEventListener(addModuleBtn, 'click', () => {
      selectedModuleIds.clear();
      fetchAvailableModules();
      if (addModuleModal) addModuleModal.style.display = 'flex';
    });
  }

  if (studyModulesBtn) {
    safeAddEventListener(studyModulesBtn, 'click', () => {
      renderStudyModules();
      if (studyModal) studyModal.style.display = 'flex';
    });
  }

  if (testModulesBtn) {
    safeAddEventListener(testModulesBtn, 'click', () => {
      renderTestModules();
      if (testModal) testModal.style.display = 'flex';
    });
  }

  // Закрытие модальных окон
  if (closeModalBtn) safeAddEventListener(closeModalBtn, 'click', () => {
    if (addModuleModal) addModuleModal.style.display = 'none';
  });
  if (addModuleModal) {
    safeAddEventListener(addModuleModal, 'click', (e) => {
      if (e.target === addModuleModal) addModuleModal.style.display = 'none';
    });
  }

  if (closeStudyModalBtn) safeAddEventListener(closeStudyModalBtn, 'click', () => {
    if (studyModal) studyModal.style.display = 'none';
  });
  if (studyModal) {
    safeAddEventListener(studyModal, 'click', (e) => {
      if (e.target === studyModal) studyModal.style.display = 'none';
    });
  }

  if (closeTestModalBtn) safeAddEventListener(closeTestModalBtn, 'click', () => {
    if (testModal) testModal.style.display = 'none';
  });
  if (testModal) {
    safeAddEventListener(testModal, 'click', (e) => {
      if (e.target === testModal) testModal.style.display = 'none';
    });
  }

  // Выбор всех/снятие выбора
  if (selectAllBtn) {
    safeAddEventListener(selectAllBtn, 'click', () => {
      if (studyModulesContainer) {
        studyModulesContainer.querySelectorAll('input[type="checkbox"]').forEach(checkbox => {
          checkbox.checked = true;
          selectedModuleIds.add(checkbox.value);
        });
        if (startStudyingBtn) safeSetDisabled(startStudyingBtn, false);
      }
    });
  }

  if (deselectAllBtn) {
    safeAddEventListener(deselectAllBtn, 'click', () => {
      if (studyModulesContainer) {
        studyModulesContainer.querySelectorAll('input[type="checkbox"]').forEach(checkbox => {
          checkbox.checked = false;
          selectedModuleIds.delete(checkbox.value);
        });
        if (startStudyingBtn) safeSetDisabled(startStudyingBtn, true);
      }
    });
  }

  if (testSelectAllBtn) {
    safeAddEventListener(testSelectAllBtn, 'click', () => {
      if (testModulesContainer) {
        testModulesContainer.querySelectorAll('input[type="checkbox"]').forEach(checkbox => {
          checkbox.checked = true;
          selectedModuleIds.add(checkbox.value);
        });
        if (startTestingBtn) safeSetDisabled(startTestingBtn, false);
      }
    });
  }

  if (testDeselectAllBtn) {
    safeAddEventListener(testDeselectAllBtn, 'click', () => {
      if (testModulesContainer) {
        testModulesContainer.querySelectorAll('input[type="checkbox"]').forEach(checkbox => {
          checkbox.checked = false;
          selectedModuleIds.delete(checkbox.value);
        });
        if (startTestingBtn) safeSetDisabled(startTestingBtn, true);
      }
    });
  }

  // Подтверждение добавления модулей
  if (confirmAddModulesBtn) {
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
        if (addModuleModal) addModuleModal.style.display = 'none';
        if (emptyMessage && emptyMessage.style.display !== 'none') emptyMessage.style.display = 'none';
        if (modulesContainer) {
          modulesIdsArray.forEach(moduleId => {
            const moduleData = moduleDataMap.get(moduleId.toString());
            if (moduleData) {
              createModuleCard(moduleData, modulesContainer);
              categoryModules.push(moduleData);
            }
          });
        }
        checkShowButtons();
        showSuccessMessage('Модули успешно добавлены');
      })
      .catch(() => alert('Ошибка добавления модулей'));
    });
  }

  // Запуск заучивания/тестирования
  if (startStudyingBtn) {
    safeAddEventListener(startStudyingBtn, 'click', () => {
      if (selectedModuleIds.size === 0) return;
      const modulesIdsString = Array.from(selectedModuleIds).join(',');
      window.location.href = `/static/learning.html?category_id=${categoryId}&modules_ids=${modulesIdsString}`;
    });
  }

  if (startTestingBtn) {
    safeAddEventListener(startTestingBtn, 'click', () => {
      if (selectedModuleIds.size === 0) return;
      const modulesIdsString = Array.from(selectedModuleIds).join(',');
      window.location.href = `/static/test.html?category_id=${categoryId}&modules_ids=${modulesIdsString}`;
    });
  }

  // Редактирование категории
  if (editCategoryBtn) {
    safeAddEventListener(editCategoryBtn, 'click', () => {
      isEditMode = !isEditMode;
      const modules = document.querySelectorAll('.card');
      modules.forEach(module => {
        const actions = module.querySelector('.module-actions');
        if (isEditMode) {
          module.classList.add('edit-mode');
          if (actions) actions.classList.add('show');
        } else {
          module.classList.remove('edit-mode');
          if (actions) actions.classList.remove('show');
        }
      });
      editCategoryBtn.textContent = isEditMode ? 'Сохранить изменения' : 'Редактировать категорию';
    });
  }

  // НОВАЯ КНОПКА: Изменение типа категории
  if (toggleCategoryTypeBtn) {
    safeAddEventListener(toggleCategoryTypeBtn, 'click', () => {
      if (!(currentUserId && categoryOwnerId && currentUserId == categoryOwnerId)) {
        return;
      }

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
      .catch(err => {
        console.error('Ошибка изменения типа категории:', err);
        alert('Ошибка изменения типа категории: ' + err.message);
      });
    });
  }

  // Навигация
  const navToggle = document.getElementById('nav-toggle');
  const navPanel = document.getElementById('nav-panel');
  const navMainBtn = document.getElementById('main-btn');
  const navModulesBtn = document.getElementById('modules-btn');
  const navCategoriesBtn = document.getElementById('categories-btn');
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
  safeAddEventListener(head, 'click', () => window.location.href = '/static/main.html');

  // Таймаут для fallback
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
  }, 3000);

  // Запуск загрузки данных
  fetchUser();
  fetchCategory();
});
