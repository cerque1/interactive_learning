const API_BASE_URL = window.location.origin;

function getTypeAsInt(type) {
  return type === 'public' ? 0 : 1;
}

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

let isEditMode = false;
let currentCards = [];

function createModuleCard(module) { 
  const card = document.createElement('div');
  card.className = 'card';
  card.dataset.moduleId = module.id;
  card.style.cursor = 'pointer';
  
  const typeText = module.type === 0 ? 'Открытый' : 'Приватный';
  
  card.innerHTML = `
    <div class="card-header">
      <span class="module-type">${typeText}</span>
    </div>
    <div class="card-title">${module.name}</div>
    <div class="card-actions">
      <button class="module-action-btn edit" title="Редактировать модуль">✎</button>
      <button class="module-action-btn delete" title="Удалить модуль">×</button>
    </div>
  `;
  
  const deleteBtn = card.querySelector('.delete');
  deleteBtn.addEventListener('click', (e) => {
    e.stopPropagation();
    handleDeleteModule(module.id, card);
  });

  const editBtn = card.querySelector('.edit');
  editBtn.addEventListener('click', (e) => {
    e.stopPropagation();
    // Открываем модалку редактирования
    if (window.openEditModuleModal) {
      window.openEditModuleModal(module.id, module.name);
    }
  });

  card.addEventListener('click', (e) => {
    if (!e.target.closest('.module-action-btn')) {
      window.location.href = `/static/module.html?module_id=${module.id}`;
    }
  });
  
  currentCards.push(card);
  return card;
}

function handleDeleteModule(moduleId, cardElement) {
  if (!confirm('Вы уверены, что хотите удалить этот модуль? Все карточки в модуле тоже будут удалены.')) {
    return;
  }

  const token = localStorage.getItem('token');
  
  fetch(`${API_BASE_URL}/api/v1/module/delete/${moduleId}`, {
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
    cardElement.style.transition = 'opacity 0.3s ease, transform 0.3s ease';
    cardElement.style.opacity = '0';
    cardElement.style.transform = 'translateX(-20px)';
    setTimeout(() => {
      cardElement.remove();
      currentCards = currentCards.filter(c => c !== cardElement);
      checkEmptyState();
    }, 300);
  })
  .catch(err => {
    alert('Ошибка при удалении модуля: ' + err.message);
  });
}

function checkEmptyState() {
  const container = document.getElementById('modules-container');
  const emptyMsg = document.getElementById('modules-empty');
  
  if (currentCards.length === 0) {
    emptyMsg.style.display = 'block';
  } else {
    emptyMsg.style.display = 'none';
  }
}

function loadModules(token, userId, myId) {
  let user_id = myId || userId;
  const url = `${API_BASE_URL}/api/v1/module/to_user/${user_id}`;

  currentCards = [];

  fetch(url, {
    headers: { 'Authorization': `Bearer ${token}` }
  })
  .then(res => {
    if (res.status === 401) {
      window.location.href = `/static/login.html?redirect=${encodeURIComponent(window.location.href)}`;
      return;
    }
    if (!res.ok) throw new Error('Network error');
    return res.json();
  })
  .then(modules => {
    const container = document.getElementById('modules-container');
    const emptyMsg = document.getElementById('modules-empty');
    const pageTitle = document.getElementById('page-title');
    
    container.innerHTML = '';
    currentCards = [];
    
    if (!modules.modules?.length) {
      emptyMsg.style.display = 'block';
    } else {
      emptyMsg.style.display = 'none';
      modules.modules.forEach(module => {
        const card = createModuleCard(module);
        container.appendChild(card);
      });
      
      if (window.myId) {
        setupEditMode();
      }
    }
    
    if (pageTitle) pageTitle.textContent = userId ? 'Модули пользователя' : 'Мои модули';
  })
  .catch(err => {
    document.getElementById('modules-empty').style.display = 'block';
  });
}

function setupEditMode() {
  const editBtn = document.getElementById('editModulesBtn');
  
  if (!window.myId || !editBtn) {
    if (editBtn) editBtn.style.display = 'none';
    return;
  }
  
  editBtn.style.display = 'inline-block';
  
  editBtn.addEventListener('click', () => {
    isEditMode = !isEditMode;
    
    currentCards.forEach(card => {
      if (isEditMode) {
        card.classList.add('edit-mode');
        const actions = card.querySelector('.card-actions');
        if (actions) actions.classList.add('show');
      } else {
        card.classList.remove('edit-mode');
        const actions = card.querySelector('.card-actions');
        if (actions) actions.classList.remove('show');
      }
    });
    
    editBtn.textContent = isEditMode ? 'Сохранить изменения' : 'Редактировать модули';
  });
}

function setupModal(token) {
  const modal = document.getElementById('createModal');
  const createBtn = document.getElementById('createModuleBtn');
  const closeBtn = document.getElementById('closeModal');
  const cancelBtn = document.getElementById('cancelModal');
  const confirmBtn = document.getElementById('createModuleConfirm');
  const moduleNameInput = document.getElementById('moduleName');
  const errorElement = document.getElementById('create-error');

  function validateForm() {
    const name = moduleNameInput.value.trim();
    confirmBtn.disabled = !name;
  }

  moduleNameInput.addEventListener('input', validateForm);

  createBtn.addEventListener('click', () => {
    modal.style.display = 'flex';
    errorElement.style.display = 'none';
    moduleNameInput.value = '';
    confirmBtn.disabled = true;
  });

  function closeModal() {
    modal.style.display = 'none';
  }

  closeBtn.onclick = closeModal;
  cancelBtn.onclick = closeModal;
  modal.onclick = e => e.target === modal && closeModal();

  confirmBtn.onclick = () => {
    const name = moduleNameInput.value.trim();
    const type = document.getElementById('moduleType').value;
    
    if (!name) {
      errorElement.textContent = 'Введите название модуля';
      errorElement.style.display = 'block';
      return;
    }

    fetch(`${API_BASE_URL}/api/v1/module/create`, {
      method: 'POST',
      headers: { 
        'Authorization': `Bearer ${token}`, 
        'Content-Type': 'application/json' 
      },
      body: JSON.stringify({ name, type: getTypeAsInt(type) })
    })
    .then(res => {
      if (res.status === 401) throw new Error('401');
      if (!res.ok) throw new Error(res.status);
      return res.json();
    })
    .then(data => {
      closeModal();
      const container = document.getElementById('modules-container');
      const newCard = createModuleCard({ 
        id: data.new_module_id, 
        name, 
        type: getTypeAsInt(type) 
      });
      container.appendChild(newCard);
      document.getElementById('modules-empty').style.display = 'none';
    })
    .catch(err => {
      errorElement.textContent = 'Ошибка создания модуля';
      errorElement.style.display = 'block';
    });
  };
}

function setupEditModuleModal(token) {
  const editModal = document.getElementById('editModal');
  const closeEditBtn = document.getElementById('closeEditModal');
  const cancelEditBtn = document.getElementById('cancelEditModal');
  const editConfirmBtn = document.getElementById('editModuleConfirm');
  const editModuleNameInput = document.getElementById('editModuleName');
  const editError = document.getElementById('edit-error');

  let currentEditingModuleId = null;
  let originalModuleName = '';

  function validateEditForm() {
    const newName = editModuleNameInput.value.trim();
    const isValid = newName && newName !== originalModuleName;
    editConfirmBtn.disabled = !isValid;
    
    if (newName === originalModuleName) {
      editError.textContent = 'Название должно отличаться от текущего';
      editError.style.display = 'block';
    } else if (!newName) {
      editError.style.display = 'none';
    } else {
      editError.style.display = 'none';
    }
  }

  editModuleNameInput.addEventListener('input', validateEditForm);

  function openEditModal(moduleId, currentName) {
    currentEditingModuleId = moduleId;
    originalModuleName = currentName;
    editModuleNameInput.value = currentName;
    editError.style.display = 'none';
    editConfirmBtn.disabled = true;
    editModal.style.display = 'flex';
  }

  function closeEditModal() {
    editModal.style.display = 'none';
    currentEditingModuleId = null;
    originalModuleName = '';
  }

  closeEditBtn.onclick = closeEditModal;
  cancelEditBtn.onclick = closeEditModal;
  editModal.onclick = e => e.target === editModal && closeEditModal();

  editConfirmBtn.onclick = () => {
    const newName = editModuleNameInput.value.trim();
    
    if (!newName || newName === originalModuleName) {
      editError.textContent = 'Введите новое название модуля';
      editError.style.display = 'block';
      return;
    }

    fetch(`${API_BASE_URL}/api/v1/module/rename/${currentEditingModuleId}`, {
      method: 'PUT',
      headers: { 
        'Authorization': `Bearer ${token}`, 
        'Content-Type': 'application/json' 
      },
      body: JSON.stringify({ new_name: newName })
    })
    .then(res => {
      if (res.status === 401) {
        window.location.href = `/static/login.html?redirect=${encodeURIComponent(window.location.href)}`;
        return;
      }
      if (res.status === 400) {
        throw new Error('Неверное название модуля');
      }
      if (res.status === 409) {
        throw new Error('Модуль с таким названием уже существует');
      }
      if (!res.ok) {
        throw new Error(`Ошибка: ${res.status}`);
      }
      return res.json();
    })
    .then(() => {
      // Найти и обновить карточку
      const card = document.querySelector(`[data-module-id="${currentEditingModuleId}"]`);
      if (card) {
        const title = card.querySelector('.card-title');
        if (title) {
          title.textContent = newName;
        }
      }
      closeEditModal();
    })
    .catch(err => {
      editError.textContent = err.message || 'Ошибка переименования модуля';
      editError.style.display = 'block';
    });
  };

  // Функция для открытия модалки из createModuleCard
  window.openEditModuleModal = openEditModal;
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

  ['main-btn', 'modules-btn', 'categories-btn', 'results-btn'].forEach(id => {
    const btn = document.getElementById(id);
    if (btn) btn.addEventListener('click', () => window.location.href = `/static/${id.replace('-btn', '.html')}`);
  });
}

window.addEventListener('DOMContentLoaded', async () => {
  const token = localStorage.getItem('token');
  if (!token) {
    window.location.href = `/static/login.html?redirect=${encodeURIComponent(window.location.href)}`;
    return;
  }

  const params = new URLSearchParams(window.location.search);
  const userId = params.get('id');
  
  const myUserData = await loadUserName(token);
  window.myId = myUserData?.user?.id;
  
  loadModules(token, userId, window.myId);
  setupModal(token);
  setupEditModuleModal(token);
  setupNavigation();
});
