window.addEventListener('DOMContentLoaded', async () => {
  const token = localStorage.getItem('token');
  if (!token) {
    window.location.href = `/static/login.html?redirect=${encodeURIComponent(window.location.href)}`;
    return;
  }

  // Получаем id пользователя из URL
  const params = new URLSearchParams(window.location.search);
  const userId = params.get('id');
  
  // Загружаем данные пользователя асинхронно
  const myUserData = await loadUserName(token);
  const myId = myUserData ? myUserData.user?.id : null;

  const headerActions = document.getElementById('header-actions');
  if (!myId) {
    headerActions.style.display = 'none';
  }
  
  loadModules(token, userId, myId);
  setupModal(token, userId);
});

// Функция возвращает данные пользователя
function loadUserName(token) {
  return fetch('http://localhost:8080/api/v1/user/me?is_full=f', {
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
      usernameElem.textContent = userData.user.name || 'Пользователь';
      usernameElem.onclick = () => {
        window.location.href = '/static/profile.html';
      };
    }
    return userData; // Возвращаем полные данные пользователя
  })
  .catch(() => {
    return null;
  });
}

// Функция для создания карточки модуля (используется и при загрузке, и при создании)
function createModuleCard(module_id, module_name, module_type) {
  const card = document.createElement('div');
  card.className = 'card';
  
  // Определяем тип модуля для отображения
  const typeText = module_type === 1 ? 'Открытый' : 'Приватный';
  
  card.innerHTML = `
    <div class="card-header">
      <span class="module-type">${typeText}</span>
    </div>
    <div class="card-title">${module_name}</div>
  `;
  card.onclick = () => {
    window.location.href = `/static/module.html?module_id=${module_id}`;
  };
  return card;
}

function loadModules(token, userId, myId) {
  let user_id = myId;
  if (userId != null) {
    user_id = userId;
  }
  let url = `http://localhost:8080/api/v1/module/to_user/${user_id}`;

  fetch(url, {
    headers: { 'Authorization': `Bearer ${token}` }
  })
  .then(res => {
    if (res.status === 401) {
      window.location.href = `/static/login.html?redirect=${encodeURIComponent(window.location.href)}`;
      return;
    }
    if (!res.ok) {
      throw new Error('Network response was not ok');
    }
    return res.json();
  })
  .then(modules => {
    const container = document.getElementById('modules-container');
    const emptyMsg = document.getElementById('modules-empty');
    const pageTitle = document.getElementById('page-title');
    
    container.innerHTML = '';
    if (!modules.modules || modules.modules.length === 0) {
      emptyMsg.style.display = 'block';
    } else {
      emptyMsg.style.display = 'none';
      modules.modules.forEach(module => {
        const card = createModuleCard(module.id, module.name, module.type);
        container.appendChild(card);
      });
    }

    // Обновляем заголовок страницы
    if (userId) {
      pageTitle.textContent = 'Модули пользователя';
    } else {
      pageTitle.textContent = 'Мои модули';
    }
  })
  .catch(() => {
    document.getElementById('modules-empty').textContent = 'Попробуйте позже';
    document.getElementById('modules-empty').style.display = 'block';
  });
}

function setupModal(token, userId) {
  const modal = document.getElementById('createModal');
  const createBtn = document.getElementById('createModuleBtn');
  const closeBtn = document.getElementById('closeModal');
  const cancelBtn = document.getElementById('cancelModal');
  const confirmBtn = document.getElementById('createModuleConfirm');
  
  // Создаем элемент для ошибок
  let errorElement = document.querySelector('.error-message');
  if (!errorElement) {
    errorElement = document.createElement('div');
    errorElement.className = 'error-message';
    errorElement.style.cssText = `
      margin-top: 12px;
      color: #c75c5c;
      font-weight: 600;
      font-size: 14px;
      display: none;
    `;
    const modalBody = document.querySelector('.modal-body');
    modalBody.appendChild(errorElement);
  }

  createBtn.onclick = () => {
    modal.style.display = 'flex';
    errorElement.style.display = 'none';
    errorElement.textContent = '';
  };
  
  function closeModal() {
    modal.style.display = 'none';
    document.getElementById('moduleName').value = '';
    document.getElementById('moduleType').value = 'private';
    errorElement.style.display = 'none';
    errorElement.textContent = '';
  }

  closeBtn.onclick = closeModal;
  cancelBtn.onclick = closeModal;

  modal.onclick = (e) => {
    if (e.target === modal) closeModal();
  };

  confirmBtn.onclick = () => {
    const name = document.getElementById('moduleName').value.trim();
    const type = document.getElementById('moduleType').value;
    
    const typeInt = getTypeAsInt(type);

    if (!name) {
      errorElement.textContent = 'Введите название модуля';
      errorElement.style.display = 'block';
      return;
    }

    fetch('http://localhost:8080/api/v1/module/create', {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({ name, type: typeInt })
    })
    .then(res => {
      if (res.status === 401) {
        window.location.href = `/static/login.html?redirect=${encodeURIComponent(window.location.href)}`;
        return;
      }
      if (res.status === 400) {
        throw new Error('BadRequest');
      }
      if (res.status === 500) {
        throw new Error('InternalServerError');
      }
      if (!res.ok) {
        throw new Error('Network response was not ok');
      }
      return res.json();
    })
    .then(newModule => {
      closeModal();
      
      // Добавляем новую карточку в конец списка
      const container = document.getElementById('modules-container');
      const emptyMsg = document.getElementById('modules-empty');
      
      emptyMsg.style.display = 'none';
      
      // Используем ID из ответа API и выбранный тип
      const newCard = createModuleCard(newModule.new_module_id, name, typeInt);
      container.appendChild(newCard);
    })
    .catch(err => {
      errorElement.style.display = 'block';
      if (err.message === 'BadRequest') {
        errorElement.textContent = 'Невалидные данные';
      } else if (err.message === 'InternalServerError') {
        errorElement.textContent = 'Попробуйте позже';
      } else {
        errorElement.textContent = 'Произошла ошибка';
      }
    });
  };
}


// Навигационная панель (если элементы существуют)
const navToggle = document.getElementById('nav-toggle');
const navPanel = document.getElementById('nav-panel');

if (navToggle && navPanel) {
  navToggle.addEventListener('click', function() {
    navPanel.classList.toggle('open');
    navToggle.classList.toggle('open');
  });
}

const navModulesBut = document.getElementById('modules-btn');
navModulesBut.addEventListener('click', function() {
  window.location.href = "/static/modules.html";
});

const navMainBut = document.getElementById('main-btn');
navMainBut.addEventListener('click', function() {
  window.location.href = '/static/main.html';
});

const navCategoriesBut = document.getElementById('categories-btn');
navCategoriesBut.addEventListener('click', function() {
  window.location.href = '/static/categories.html';
});

const head = document.getElementById('head');
head.addEventListener('click', (e) => {
  e.preventDefault();
  window.location.href = '/static/main.html';
});