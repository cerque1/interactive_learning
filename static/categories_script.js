let isEditMode = false;
let currentCards = [];
let windowMyId = null;

const API_BASE_URL = window.location.origin;

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
    if (userData?.user) {
      const usernameElem = document.getElementById('username');
      if (usernameElem) {
        usernameElem.textContent = userData.user.name || 'Пользователь';
        usernameElem.style.cursor = 'pointer';
        usernameElem.onclick = () => window.location.href = '/static/profile.html';
      }
    }
    return userData;
  })
  .catch(() => null);
}

function createCategoryCard(category) {
  const card = document.createElement('div');
  card.className = 'card';
  card.dataset.categoryId = category.id;
  card.style.cursor = 'pointer';
  
  const privacyType = category.type === 0 ? 'public' : 'private';
  const privacyText = category.type === 0 ? 'Открытая' : 'Приватная';
  
  card.innerHTML = `
    <div class="card-header">
      <div class="category-type ${privacyType}">${privacyText}</div>
    </div>
    <div class="card-title">${category.name}</div>
    <div class="card-actions">
      <button class="category-action-btn edit" title="Редактировать категорию">✎</button>
      <button class="category-action-btn delete" title="Удалить категорию">×</button>
    </div>
  `;
  
  const deleteBtn = card.querySelector('.delete');
  deleteBtn.addEventListener('click', (e) => {
    e.stopPropagation();
    handleDeleteCategory(category.id, card);
  });

  const editBtn = card.querySelector('.edit');
  editBtn.addEventListener('click', (e) => {
    e.stopPropagation();
    if (window.openEditCategoryModal) {
      window.openEditCategoryModal(category.id, category.name);
    }
  });

  card.addEventListener('click', (e) => {
    if (!e.target.closest('.category-action-btn')) {
      window.location.href = `/static/category.html?category_id=${category.id}`;
    }
  });
  
  currentCards.push(card);
  return card;
}

function handleDeleteCategory(categoryId, cardElement) {
  if (!confirm('Вы уверены, что хотите удалить эту категорию?')) {
    return;
  }

  const token = localStorage.getItem('token');
  
  fetch(`${API_BASE_URL}/api/v1/category/delete/${categoryId}`, {
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
    alert('Ошибка при удалении категории: ' + err.message);
  });
}

function checkEmptyState() {
  const container = document.getElementById('categories-container');
  const emptyMsg = document.getElementById('categories-empty');
  
  if (currentCards.length === 0) {
    emptyMsg.style.display = 'block';
  } else {
    emptyMsg.style.display = 'none';
  }
}

async function loadCategories(token, userId, myId) {
  let userid = myId
  if (userId) {
    userid = userId
  } 
  const url = `${API_BASE_URL}/api/v1/category/to_user/${userid}`;

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
  .then(categories => {
    const container = document.getElementById('categories-container');
    const emptyMsg = document.getElementById('categories-empty');
    const pageTitle = document.getElementById('page-title');
    
    container.innerHTML = '';
    currentCards = [];
    
    if (!categories?.categories?.length) {
      emptyMsg.style.display = 'block';
    } else {
      emptyMsg.style.display = 'none';
      categories.categories.forEach(category => {
        const card = createCategoryCard(category);
        container.appendChild(card);
      });
      
      if (windowMyId) {
        setupEditMode();
      }
    }
    
    if (pageTitle) pageTitle.textContent = userId ? 'Категории пользователя' : 'Мои категории';
  })
  .catch(() => {
    document.getElementById('categories-empty').style.display = 'block';
  });
}

function setupEditMode() {
  const editBtn = document.getElementById('editCategoriesBtn');
  
  if (!windowMyId || !editBtn) {
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
    
    editBtn.textContent = isEditMode ? 'Сохранить изменения' : 'Редактировать категории';
  });
}

async function loadModules(token, userId) {
  const url = `${API_BASE_URL}/api/v1/module/to_user/${userId}?with_cards=t`;
  
  try {
    const res = await fetch(url, {
      headers: { 'Authorization': `Bearer ${token}` }
    });
    
    if (res.status === 401) {
      window.location.href = `/static/login.html?redirect=${encodeURIComponent(window.location.href)}`;
      return [];
    }
    
    if (!res.ok) throw new Error('Network response was not ok');
    
    const data = await res.json();
    return data.modules || [];
  } catch {
    return [];
  }
}

async function loadModulesForModal(token, userId) {
  const modulesList = document.getElementById('modules-list');
  const modules = await loadModules(token, userId);
  
  modulesList.innerHTML = '';
  
  if (modules.length === 0) {
    modulesList.innerHTML = '<div class="no-modules">Модули отсутствуют</div>';
    return;
  }
  
  modules.forEach((module) => {
    const div = document.createElement('div');
    div.className = 'module-item';
    
    const label = document.createElement('label');
    label.className = 'module-label';
    
    const checkbox = document.createElement('input');
    checkbox.type = 'checkbox';
    checkbox.className = 'module-checkbox';
    checkbox.dataset.moduleId = module.id;
    checkbox.id = `module-${module.id}`; // ✅ Уникальный ID без index
    
    const span = document.createElement('span');
    span.textContent = `${module.name} (${module.cards_count || module.cards?.length || 0} карточек)`;
    
    // ✅ Правильный порядок: checkbox -> span
    label.appendChild(checkbox);
    label.appendChild(span);
    label.htmlFor = checkbox.id; // ✅ Правильная связь label с checkbox
    
    // ✅ Предотвращаем всплытие для модального окна
    label.addEventListener('click', (e) => {
      e.stopPropagation();
    });
    
    div.appendChild(label);
    modulesList.appendChild(div);
  });
}

function showModalError(message) {
  const errorEl = document.getElementById('modal-error');
  errorEl.textContent = message;
  errorEl.style.display = 'block';
}

function setupEditCategoryModal(token) {
  const editModal = document.getElementById('editCategoryModal');
  const closeEditBtn = document.getElementById('closeEditCategoryModal');
  const cancelEditBtn = document.getElementById('cancelEditCategoryModal');
  const editConfirmBtn = document.getElementById('editCategoryConfirm');
  const editCategoryNameInput = document.getElementById('editCategoryName');
  const editError = document.getElementById('edit-category-error');

  let currentEditingCategoryId = null;
  let originalCategoryName = '';

  function validateEditForm() {
    const newName = editCategoryNameInput.value.trim();
    const isValid = newName && newName !== originalCategoryName;
    editConfirmBtn.disabled = !isValid;
    
    if (newName === originalCategoryName) {
      editError.textContent = 'Название должно отличаться от текущего';
      editError.style.display = 'block';
    } else if (!newName) {
      editError.style.display = 'none';
    } else {
      editError.style.display = 'none';
    }
  }

  editCategoryNameInput.addEventListener('input', validateEditForm);

  window.openEditCategoryModal = function(categoryId, currentName) {
    currentEditingCategoryId = categoryId;
    originalCategoryName = currentName;
    editCategoryNameInput.value = currentName;
    editError.style.display = 'none';
    editConfirmBtn.disabled = true;
    editModal.style.display = 'flex';
    editCategoryNameInput.focus();
  };

  function closeEditModal() {
    editModal.style.display = 'none';
    currentEditingCategoryId = null;
    originalCategoryName = '';
    editCategoryNameInput.value = '';
  }

  closeEditBtn.onclick = closeEditModal;
  cancelEditBtn.onclick = closeEditModal;
  editModal.onclick = e => e.target === editModal && closeEditModal();

  editConfirmBtn.onclick = () => {
    const newName = editCategoryNameInput.value.trim();
    
    if (!newName || newName === originalCategoryName) {
      editError.textContent = 'Введите новое название категории';
      editError.style.display = 'block';
      return;
    }

    fetch(`${API_BASE_URL}/api/v1/category/rename/${currentEditingCategoryId}`, {
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
        throw new Error('Неверное название категории');
      }
      if (res.status === 409) {
        throw new Error('Категория с таким названием уже существует');
      }
      if (!res.ok) {
        throw new Error(`Ошибка: ${res.status}`);
      }
      return res.json();
    })
    .then(() => {
      const card = document.querySelector(`[data-category-id="${currentEditingCategoryId}"]`);
      if (card) {
        const title = card.querySelector('.card-title');
        if (title) {
          title.textContent = newName;
        }
      }
      closeEditModal();
    })
    .catch(err => {
      editError.textContent = err.message || 'Ошибка переименования категории';
      editError.style.display = 'block';
    });
  };
}

function setupModal(token, userId, myId) {
  const modal = document.getElementById('createCategoryModal');
  const createBtn = document.getElementById('createCategoryBtn');
  const closeBtn = document.getElementById('closeCategoryModal');
  const cancelBtn = document.getElementById('cancelCategoryModal');
  const confirmBtn = document.getElementById('createCategoryConfirm');
  const categoryNameInput = document.getElementById('categoryName');
  const errorElement = document.getElementById('modal-error');

  function validateForm() {
    const name = categoryNameInput.value.trim();
    confirmBtn.disabled = !name;
  }

  categoryNameInput.addEventListener('input', validateForm);

  createBtn.addEventListener('click', async () => {
    modal.style.display = 'flex';
    errorElement.style.display = 'none';
    categoryNameInput.value = '';
    categoryNameInput.focus();
    confirmBtn.disabled = true;
    await loadModulesForModal(token, myId || userId);
  });

  function closeModal() {
    modal.style.display = 'none';
    document.getElementById('modules-list').innerHTML = '';
    categoryNameInput.value = '';
    errorElement.style.display = 'none';
    confirmBtn.disabled = true;
  }

  closeBtn.onclick = closeModal;
  cancelBtn.onclick = closeModal;
  modal.onclick = e => e.target === modal && closeModal();

  confirmBtn.onclick = async () => {
    const name = categoryNameInput.value.trim();
    const selectedModules = Array.from(document.querySelectorAll('.module-checkbox:checked')).map(cb => parseInt(cb.dataset.moduleId, 10));
    
    if (!name) {
      showModalError('Введите название категории');
      return;
    }

    if (selectedModules.length === 0) {
      showModalError('Выберите хотя бы один модуль');
      return;
    }

    try {
      const res = await fetch(`${API_BASE_URL}/api/v1/category/create`, {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({ name, modules_ids: selectedModules })
      });
      
      if (res.status === 401) {
        window.location.href = `/static/login.html?redirect=${encodeURIComponent(window.location.href)}`;
        return;
      }
      
      if (!res.ok) {
        showModalError('Ошибка создания категории');
        return;
      }
      
      const newCategory = await res.json();
      closeModal();
      const container = document.getElementById('categories-container');
      const newCard = createCategoryCard({ 
        id: newCategory.new_id, 
        name 
      });
      container.appendChild(newCard);
      document.getElementById('categories-empty').style.display = 'none';
    } catch (err) {
      showModalError('Ошибка создания категории');
    }
  };
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

  const head = document.getElementById('head');
  if (head) {
    head.addEventListener('click', e => {
      e.preventDefault();
      window.location.href = '/static/main.html';
    });
  }
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
  windowMyId = myUserData?.user?.id;
  
  loadCategories(token, userId, windowMyId);
  setupModal(token, userId, windowMyId);
  setupEditCategoryModal(token);
  setupNavigation();
});
