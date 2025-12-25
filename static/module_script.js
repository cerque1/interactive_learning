window.addEventListener('DOMContentLoaded', () => {
  const token = localStorage.getItem('token');
  if (!token) {
    window.location.href = '/static/login.html?redirect=' + encodeURIComponent(window.location.href);
    return;
  }

  const params = new URLSearchParams(window.location.search);
  const moduleId = params.get('module_id');
  if (!moduleId) {
    document.getElementById('module-name').textContent = 'Ошибка: не указан id модуля';
    return;
  }

  let currentUserId = null;
  let moduleOwnerId = null;
  let isEditMode = false;

  const addCardsBtn = document.getElementById('add-cards-btn');
  const editModuleBtn = document.getElementById('edit-module-btn');

  let userLoaded = false;
  let moduleLoaded = false;

  // Функция загрузки данных пользователя
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

  // Функция загрузки данных модуля
  function fetchModule() {
    return fetch(`http://localhost:8080/api/v1/module/${moduleId}`, {
      headers: { 'Authorization': `Bearer ${token}` }
    })
    .then(res => {
      if (res.status === 401) {
        window.location.href = '/static/login.html?redirect=' + encodeURIComponent(window.location.href);
        return Promise.reject();
      }
      if (!res.ok) {
        throw new Error('Ошибка загрузки модуля');
      }
      return res.json();
    })
    .then(moduleData => {
      moduleData = moduleData.module;
      moduleOwnerId = moduleData.user_id || moduleData.owner_id;

      const moduleNameElem = document.getElementById('module-name');
      const cardsContainer = document.getElementById('cards-container');
      const emptyMessage = document.getElementById('empty-message');

      moduleNameElem.textContent = moduleData.name || 'Без названия';

      if (!moduleData.cards || moduleData.cards.length === 0) {
        cardsContainer.innerHTML = '';
        emptyMessage.style.display = 'block';
      } else {
        emptyMessage.style.display = 'none';
        cardsContainer.innerHTML = '';
        moduleData.cards.forEach((card) => {
          const cardElem = document.createElement('div');
          cardElem.className = 'card';
          cardElem.dataset.cardId = card.id;
          cardElem.innerHTML = `
            <div class="card-title">${card.term.text}</div>
            <div class="card-definition">${card.definition.text}</div>
            <div class="card-languages">
              <span class="lang-badge">${card.term.lang.toUpperCase()}</span>
              <div class="card-actions">
                <button class="action-btn delete" title="Удалить карточку">×</button>
                <button class="action-btn edit" title="Редактировать карточку">✎</button>
              </div>
              <span class="lang-badge">${card.definition.lang.toUpperCase()}</span>
            </div>
          `;
          cardsContainer.appendChild(cardElem);
        });        
      }

      moduleLoaded = true;
      checkShowButtons();
    })
    .catch(() => {
      document.getElementById('module-name').textContent = 'Ошибка загрузки модуля';
      document.getElementById('cards-container').innerHTML = '';
      moduleLoaded = true;
      checkShowButtons();
    });
  }

  // Функция проверки и показа кнопок
  function checkShowButtons() {
    console.log('Проверка кнопок:', {
      userLoaded,
      moduleLoaded,
      currentUserId,
      moduleOwnerId,
      isOwner: currentUserId && moduleOwnerId && Number(currentUserId) === Number(moduleOwnerId)
    });

    if (userLoaded && moduleLoaded && currentUserId && moduleOwnerId && 
        Number(currentUserId) === Number(moduleOwnerId)) {
      addCardsBtn.style.display = 'inline-block';
      editModuleBtn.style.display = 'inline-block';
    }
  }

  // Запускаем параллельные запросы
  fetchUser();
  fetchModule();

  // Функция для показа уведомления об успехе
  function showSuccessMessage(message) {
    const notification = document.createElement('div');
    notification.className = 'success-message';
    notification.textContent = message;
    notification.style.cssText = `
      position: fixed;
      top: 100px;
      right: 20px;
      background: #28a745;
      color: white;
      padding: 12px 20px;
      border-radius: 6px;
      box-shadow: 0 4px 12px rgba(0,0,0,0.15);
      z-index: 10000;
      font-weight: 500;
      max-width: 300px;
      opacity: 0;
      transform: translateX(100%);
      transition: all 0.3s ease;
      font-family: inherit;
    `;
    
    document.body.appendChild(notification);
    
    requestAnimationFrame(() => {
      notification.style.opacity = '1';
      notification.style.transform = 'translateX(0)';
    });
    
    setTimeout(() => {
      notification.style.opacity = '0';
      notification.style.transform = 'translateX(100%)';
      setTimeout(() => {
        document.body.removeChild(notification);
      }, 300);
    }, 3000);
  }

  // Функция обновления состояния пустого контейнера
  function updateEmptyState() {
    const cardsContainer = document.getElementById('cards-container');
    const emptyMessage = document.getElementById('empty-message');
    const remainingCards = cardsContainer.querySelectorAll('.card');
    
    if (remainingCards.length === 0) {
      cardsContainer.style.display = 'none';
      emptyMessage.style.display = 'block';
    } else {
      cardsContainer.style.display = 'flex';
      emptyMessage.style.display = 'none';
    }
  }

  // Обработчик удаления карточки
  function handleDeleteClick(e) {
    e.stopPropagation();
    const card = e.target.closest('.card');
    const cardId = parseInt(card.dataset.cardId);
    
    if (!cardId || !confirm('Вы уверены, что хотите удалить эту карточку?')) {
      return;
    }
    
    fetch(`http://localhost:8080/api/v1/card/delete/${cardId}`, {
      method: 'DELETE',
      headers: { 'Authorization': `Bearer ${token}` }
    })
    .then(res => {
      if (res.status === 401) {
        window.location.href = '/static/login.html?redirect=' + encodeURIComponent(window.location.href);
        return;
      }
      if (!res.ok) {
        if (res.status === 404) throw new Error('Карточка не найдена');
        if (res.status === 403) throw new Error('Нет прав для удаления этой карточки');
        throw new Error('Ошибка удаления карточки');
      }
      return res.json();
    })
    .then(() => {
      card.style.transition = 'opacity 0.3s ease';
      card.style.opacity = '0';
      setTimeout(() => {
        card.remove();
        updateEmptyState();
        showSuccessMessage('Карточка успешно удалена');
      }, 300);
    })
    .catch(err => {
      console.error('Ошибка удаления:', err);
      alert('Ошибка удаления карточки: ' + err.message);
    });
  }

  // Обработчик редактирования карточки
  function handleEditClick(e) {
    e.stopPropagation();
    const card = e.target.closest('.card');
    const cardId = parseInt(card.dataset.cardId);
    
    if (!cardId) {
      alert('Ошибка: ID карточки не найден');
      return;
    }
    
    const cardTitle = card.querySelector('.card-title').textContent;
    const cardDefinition = card.querySelector('.card-definition').textContent;
    const termLangBadge = card.querySelector('.lang-badge');
    const defLangBadge = card.querySelector('.lang-badge:last-child');
    
    const termLang = termLangBadge.textContent.toLowerCase();
    const defLang = defLangBadge.textContent.toLowerCase();
    
    window.currentEditingCardId = cardId;
    
    document.getElementById('edit-term-input').value = cardTitle;
    document.getElementById('edit-definition-input').value = cardDefinition;
    document.getElementById('edit-term-lang').value = termLang;
    document.getElementById('edit-definition-lang').value = defLang;
    
    document.getElementById('edit-card-error').style.display = 'none';
    document.getElementById('update-card-btn').disabled = false;
    
    document.getElementById('editCardModal').style.display = 'flex';
  }

  // Обработчик кнопки редактирования модуля
  editModuleBtn.addEventListener('click', () => {
    isEditMode = !isEditMode;
    
    const cards = document.querySelectorAll('.card');
    cards.forEach(card => {
      if (isEditMode) {
        card.classList.add('edit-mode');
        card.querySelector('.card-actions').classList.add('show');
      } else {
        card.classList.remove('edit-mode');
        card.querySelector('.card-actions').classList.remove('show');
      }
    });
    
    editModuleBtn.textContent = isEditMode ? 'Сохранить изменения' : 'Редактировать модуль';
  });

  // Обработчики кнопок действий карточек
  document.addEventListener('click', (e) => {
    if (e.target.classList.contains('delete')) {
      handleDeleteClick(e);
    } else if (e.target.classList.contains('edit')) {
      handleEditClick(e);
    }
  });

  // Навигация
  const navToggle = document.getElementById('nav-toggle');
  const navPanel = document.getElementById('nav-panel');
  const navMainBut = document.getElementById('main-btn');
  const navModulesBut = document.getElementById('modules-btn');
  const navCategoriesBut = document.getElementById('categories-btn');
  const navResultsBut = document.getElementById('results-btn');
  const head = document.getElementById('head');

  if (navToggle && navPanel) {
    navToggle.addEventListener('click', () => {
      navPanel.classList.toggle('open');
      navToggle.classList.toggle('open');
    });
  }

  if (navMainBut) navMainBut.addEventListener('click', () => window.location.href = '/static/main.html');
  if (navModulesBut) navModulesBut.addEventListener('click', () => window.location.href = '/static/modules.html');
  if (navCategoriesBut) navCategoriesBut.addEventListener('click', () => window.location.href = '/static/categories.html');
  if (navResultsBut) navResultsBut.addEventListener('click', () => window.location.href = '/static/results.html');
  if (head) head.addEventListener('click', () => window.location.href = '/static/main.html');

  // ====== ЛОГИКА МОДАЛКИ ДОБАВЛЕНИЯ КАРТОЧЕК ======
  const addCardsModal = document.getElementById('addCardsModal');
  const closeAddCardsModal = document.getElementById('closeAddCardsModal');
  const cancelAddCardsModal = document.getElementById('cancelAddCardsModal');
  const addCardRowBtn = document.getElementById('add-card-row-btn');
  const saveCardsBtn = document.getElementById('save-cards-btn');
  const cardRowsContainer = document.getElementById('card-rows-container');
  const addCardsError = document.getElementById('add-cards-error');

  let cardRowCount = 0;
  const MAX_ROWS = 100;

  function createRowElement(firstRow = false) {
    const row = document.createElement('div');
    row.className = `card-row ${firstRow ? 'first-row' : ''}`;

    row.innerHTML = `
      <div>
        <label>Термин</label>
        <textarea class="term-input" placeholder="Введите термин..."></textarea>
        <select class="term-lang">
          <option value="ru">🇷🇺 Русский</option>
          <option value="en">🇺🇸 English</option>
        </select>
      </div>
      <div>
        <label>Определение</label>
        <textarea class="definition-input" placeholder="Введите определение..."></textarea>
        <select class="definition-lang">
          <option value="ru">🇷🇺 Русский</option>
          <option value="en">🇺🇸 English</option>
        </select>
      </div>
      ${!firstRow ? '<div class="row-actions"><button type="button" class="secondary-btn remove-row-btn">Удалить</button></div>' : ''}
    `;

    const termInput = row.querySelector('.term-input');
    const defInput = row.querySelector('.definition-input');

    termInput.addEventListener('input', updateButtonsState);
    defInput.addEventListener('input', updateButtonsState);

    if (!firstRow) {
      const removeBtn = row.querySelector('.remove-row-btn');
      removeBtn.addEventListener('click', () => {
        cardRowsContainer.removeChild(row);
        cardRowCount--;
        updateButtonsState();
      });
    }

    return row;
  }

  function resetModal() {
    cardRowsContainer.innerHTML = '';
    cardRowCount = 0;
    addCardsError.style.display = 'none';
    addCardsError.textContent = '';
    addCardRowBtn.disabled = false;
    saveCardsBtn.disabled = true;
  }

  function openModal() {
    resetModal();
    const firstRow = createRowElement(true);
    cardRowsContainer.appendChild(firstRow);
    cardRowCount = 1;
    addCardsModal.style.display = 'flex';
    updateButtonsState();
  }

  function closeModal() {
    addCardsModal.style.display = 'none';
  }

  addCardsBtn.addEventListener('click', () => {
    if (!(currentUserId && moduleOwnerId && Number(currentUserId) === Number(moduleOwnerId))) {
      return;
    }
    openModal();
  });

  closeAddCardsModal.addEventListener('click', closeModal);
  cancelAddCardsModal.addEventListener('click', closeModal);

  addCardsModal.addEventListener('click', (e) => {
    if (e.target === addCardsModal) closeModal();
  });

  function addOneRow() {
    if (cardRowCount >= MAX_ROWS) return;
    const row = createRowElement(false);
    cardRowsContainer.appendChild(row);
    cardRowCount++;
  }

  function allRowsFilled() {
    const rows = cardRowsContainer.querySelectorAll('.card-row');
    if (rows.length === 0) return false;
    for (const row of rows) {
      const term = row.querySelector('.term-input').value.trim();
      const def = row.querySelector('.definition-input').value.trim();
      if (!term || !def) return false;
    }
    return true;
  }

  function updateButtonsState() {
    const canAddMore = allRowsFilled() && cardRowCount < MAX_ROWS;
    addCardRowBtn.disabled = !canAddMore;
    saveCardsBtn.disabled = !(cardRowCount > 0 && allRowsFilled());
    if (cardRowCount >= MAX_ROWS) {
      addCardRowBtn.disabled = true;
      saveCardsBtn.disabled = true;
    }
  }

  addCardRowBtn.addEventListener('click', () => {
    if (cardRowCount >= MAX_ROWS || !allRowsFilled()) return;
    addOneRow();
    updateButtonsState();
  });

  // Обработчик сохранения карточек
  saveCardsBtn.addEventListener('click', () => {
    addCardsError.style.display = 'none';
    addCardsError.textContent = '';

    const rows = cardRowsContainer.querySelectorAll('.card-row');
    if (rows.length === 0) {
      addCardsError.textContent = 'Добавьте хотя бы одну карточку';
      addCardsError.style.display = 'block';
      return;
    }

    const cardsPayload = [];
    for (const row of rows) {
      const term = row.querySelector('.term-input').value.trim();
      const def = row.querySelector('.definition-input').value.trim();
      const termLang = row.querySelector('.term-lang').value;
      const defLang = row.querySelector('.definition-lang').value;

      if (!term || !def) {
        addCardsError.textContent = 'Все поля терминов и определений должны быть заполнены';
        addCardsError.style.display = 'block';
        return;
      }

      cardsPayload.push({
        term: { text: term, lang: termLang },
        definition: { text: def, lang: defLang }
      });
    }

    fetch(`http://localhost:8080/api/v1/card/insert_to_module`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({ cards: cardsPayload, parent_module: parseInt(moduleId) })
    })
    .then(res => {
      if (res.status === 401) {
        window.location.href = '/static/login.html?redirect=' + encodeURIComponent(window.location.href);
        return;
      }
      if (res.status === 400) throw new Error('Неверные данные');
      if (res.status === 500) throw new Error('Ошибка сервера');
      if (!res.ok) throw new Error('Ошибка сохранения карточек');
      return res.json();
    })
    .then(response => {
      closeModal();
      
      const newIds = response.new_ids || [];
      if (!newIds.length) {
        showSuccessMessage('Карточки сохранены, но ID не получены');
        return;
      }
      
      const cardsContainer = document.getElementById('cards-container');
      const emptyMessage = document.getElementById('empty-message');
      
      emptyMessage.style.display = 'none';
      
      newIds.forEach((cardId, index) => {
        const newCardData = cardsPayload[index];
        const cardElem = document.createElement('div');
        cardElem.className = 'card';
        cardElem.dataset.cardId = cardId;
        cardElem.innerHTML = `
          <div class="card-title">${newCardData.term.text}</div>
          <div class="card-definition">${newCardData.definition.text}</div>
          <div class="card-languages">
            <span class="lang-badge">${newCardData.term.lang.toUpperCase()}</span>
            <div class="card-actions">
              <button class="action-btn delete" title="Удалить карточку">×</button>
              <button class="action-btn edit" title="Редактировать карточку">✎</button>
            </div>
            <span class="lang-badge">${newCardData.definition.lang.toUpperCase()}</span>
          </div>
        `;
        
        cardElem.querySelector('.delete').addEventListener('click', handleDeleteClick);
        cardElem.querySelector('.edit').addEventListener('click', handleEditClick);
        
        cardsContainer.appendChild(cardElem);
        
        cardElem.style.opacity = '0';
        cardElem.style.transform = 'translateY(20px)';
        requestAnimationFrame(() => {
          cardElem.style.transition = 'all 0.4s ease';
          cardElem.style.opacity = '1';
          cardElem.style.transform = 'translateY(0)';
        });
      });
      
      showSuccessMessage(`${newIds.length} карточек${newIds.length === 1 ? 'а' : 'ек'} успешно добавлено${newIds.length === 1 ? '' : 'ы'}`);
    })
    .catch(err => {
      addCardsError.textContent = err.message || 'Ошибка сохранения карточек';
      addCardsError.style.display = 'block';
    });
  });

  // ====== ЛОГИКА МОДАЛКИ РЕДАКТИРОВАНИЯ КАРТОЧКИ ======
  const editCardModal = document.getElementById('editCardModal');
  const closeEditCardModal = document.getElementById('closeEditCardModal');
  const cancelEditCardModal = document.getElementById('cancelEditCardModal');
  const updateCardBtn = document.getElementById('update-card-btn');
  const editCardError = document.getElementById('edit-card-error');
  const editTermInput = document.getElementById('edit-term-input');
  const editDefInput = document.getElementById('edit-definition-input');

  function checkEditFieldsFilled() {
    return editTermInput.value.trim() && editDefInput.value.trim();
  }

  editTermInput.addEventListener('input', () => {
    updateCardBtn.disabled = !checkEditFieldsFilled();
  });

  editDefInput.addEventListener('input', () => {
    updateCardBtn.disabled = !checkEditFieldsFilled();
  });

  closeEditCardModal.addEventListener('click', closeEditModal);
  cancelEditCardModal.addEventListener('click', closeEditModal);

  editCardModal.addEventListener('click', (e) => {
    if (e.target === editCardModal) closeEditModal();
  });

  function closeEditModal() {
    editCardModal.style.display = 'none';
    window.currentEditingCardId = null;
  }

  updateCardBtn.addEventListener('click', () => {
    editCardError.style.display = 'none';
    editCardError.textContent = '';

    const term = editTermInput.value.trim();
    const definition = editDefInput.value.trim();
    const termLang = document.getElementById('edit-term-lang').value;
    const defLang = document.getElementById('edit-definition-lang').value;

    if (!term || !definition) {
      editCardError.textContent = 'Все поля терминов и определений должны быть заполнены';
      editCardError.style.display = 'block';
      return;
    }

    const cardId = window.currentEditingCardId;
    if (!cardId) {
      editCardError.textContent = 'Ошибка: ID карточки не найден';
      editCardError.style.display = 'block';
      return;
    }

    fetch(`http://localhost:8080/api/v1/card/update/${cardId}`, {
      method: 'PUT',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        term: { text: term, lang: termLang },
        definition: { text: definition, lang: defLang }
      })
    })
    .then(res => {
      if (res.status === 401) {
        window.location.href = '/static/login.html?redirect=' + encodeURIComponent(window.location.href);
        return;
      }
      if (res.status === 400) throw new Error('Неверные данные');
      if (res.status === 404) throw new Error('Карточка не найдена');
      if (res.status === 403) throw new Error('Нет прав для редактирования');
      if (res.status === 500) throw new Error('Ошибка сервера');
      if (!res.ok) throw new Error('Ошибка обновления карточки');
      return res.json();
    })
    .then(() => {
      closeEditModal();
      
      const card = document.querySelector(`[data-card-id="${cardId}"]`);
      if (card) {
        card.querySelector('.card-title').textContent = term;
        card.querySelector('.card-definition').textContent = definition;
        card.querySelector('.lang-badge').textContent = termLang.toUpperCase();
        card.querySelector('.lang-badge:last-child').textContent = defLang.toUpperCase();
      }
      
      showSuccessMessage('Карточка успешно обновлена');
    })
    .catch(err => {
      editCardError.textContent = err.message || 'Ошибка обновления карточки';
      editCardError.style.display = 'block';
    });
  });
});
