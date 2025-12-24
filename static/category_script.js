// Полный category_script.js
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
        categoryData.modules.forEach(module => {
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
          modulesContainer.appendChild(moduleElem);
        });
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

  // Проверка показа кнопок
  function checkShowButtons() {
    console.log('Проверка кнопок категории:', {
      userLoaded,
      categoryLoaded,
      currentUserId,
      categoryOwnerId,
      isOwner: currentUserId && categoryOwnerId && Number(currentUserId) === Number(categoryOwnerId)
    });

    if (userLoaded && categoryLoaded && currentUserId && categoryOwnerId && 
        Number(currentUserId) === Number(categoryOwnerId)) {
      addModuleBtn.style.display = 'inline-block';
      editCategoryBtn.style.display = 'inline-block';
      console.log('✅ Кнопки категории показаны');
    }
  }

  // Запуск запросов
  fetchUser();
  fetchCategory();

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

  // Обработчик удаления модуля из категории
  document.addEventListener('click', (e) => {
    if (e.target.classList.contains('delete')) {
      e.stopPropagation();
      const moduleCard = e.target.closest('.card');
      const moduleId = parseInt(moduleCard.dataset.moduleId);
      
      if (!moduleId || !confirm('Удалить модуль из категории?')) return;
      
      // TODO: API запрос на удаление связи модуль-категория
      fetch(`http://localhost:8080/api/v1/category/${categoryId}/remove_module/${moduleId}`, {
        method: 'DELETE',
        headers: { 'Authorization': `Bearer ${token}` }
      })
      .then(res => {
        if (!res.ok) throw new Error('Ошибка удаления');
        moduleCard.style.opacity = '0';
        setTimeout(() => {
          moduleCard.remove();
          // Проверка пустоты
          const remaining = document.querySelectorAll('.card');
          if (remaining.length === 0) {
            document.getElementById('empty-message').style.display = 'block';
          }
        }, 300);
      })
      .catch(err => {
        console.error(err);
        alert('Ошибка удаления модуля');
      });
    }
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

  // Клик по карточке модуля (переход в модуль)
  document.addEventListener('click', (e) => {
    const moduleCard = e.target.closest('.card');
    if (moduleCard && !e.target.classList.contains('action-btn')) {
      const moduleId = moduleCard.dataset.moduleId;
      if (moduleId) {
        window.location.href = `/static/module.html?module_id=${encodeURIComponent(moduleId)}`;
      }
    }
  });
});
