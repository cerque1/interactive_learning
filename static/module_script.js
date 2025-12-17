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

  const addCardsBtn = document.getElementById('add-cards-btn');

  // Получить имя пользователя и id
  fetch('http://localhost:8080/api/v1/user/me?is_full=f', {
    headers: { 'Authorization': `Bearer ${token}` }
  })
    .then(res => {
      if (res.status === 401) {
        window.location.href = '/static/login.html?redirect=' + encodeURIComponent(window.location.href);
        return;
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
    })
    .catch(() => {
      // не критично
    });

  // Получить данные модуля
  fetch(`http://localhost:8080/api/v1/module/${moduleId}`, {
    headers: { 'Authorization': `Bearer ${token}` }
  })
    .then(res => {
      if (res.status === 401) {
        window.location.href = '/static/login.html?redirect=' + encodeURIComponent(window.location.href);
        return;
      }
      if (!res.ok) {
        throw new Error('Ошибка загрузки модуля');
      }
      return res.json();
    })
    .then(moduleData => {
      moduleData = moduleData.module;
      // подстрой это поле под свой backend (user_id / owner_id / author_id и т.п.)
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
        moduleData.cards.forEach(card => {
          const cardElem = document.createElement('div');
          cardElem.className = 'card';
          cardElem.innerHTML = `
            <div class="card-title">${card.term.text}</div>
            <div>${card.definition.text}</div>
          `;
          cardsContainer.appendChild(cardElem);
        });
      }

      // Показать кнопку "Добавить карточки", если текущий пользователь — владелец
      if (currentUserId && moduleOwnerId && Number(currentUserId) === Number(moduleOwnerId)) {
        addCardsBtn.style.display = 'inline-block';
      }
    })
    .catch(() => {
      document.getElementById('module-name').textContent = 'Ошибка загрузки модуля';
      document.getElementById('cards-container').innerHTML = '';
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
    navToggle.addEventListener('click', function () {
      navPanel.classList.toggle('open');
      navToggle.classList.toggle('open');
    });
  }

  if (navMainBut) {
    navMainBut.addEventListener('click', function () {
      window.location.href = '/static/main.html';
    });
  }

  if (navModulesBut) {
    navModulesBut.addEventListener('click', function () {
      window.location.href = '/static/modules.html';
    });
  }

  if (navCategoriesBut) {
    navCategoriesBut.addEventListener('click', function () {
      window.location.href = '/static/categories.html';
    });
  }

  if (navResultsBut) {
    navResultsBut.addEventListener('click', function () {
      window.location.href = '/static/results.html';
    });
  }

  if (head) {
    head.addEventListener('click', () => {
      window.location.href = '/static/main.html';
    });
  }

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
    addCardsModal.style.display = 'flex';
    addOneRow();
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
    if (e.target === addCardsModal) {
      closeModal();
    }
  });

  function createRowElement() {
    const row = document.createElement('div');
    row.className = 'card-row';

    row.innerHTML = `
      <div>
        <label>Термин</label>
        <textarea class="term-input"></textarea>
        <select class="term-lang">
          <option value="ru">Русский</option>
          <option value="en">Английский</option>
        </select>
      </div>
      <div>
        <label>Определение</label>
        <textarea class="definition-input"></textarea>
        <select class="definition-lang">
          <option value="ru">Русский</option>
          <option value="en">Английский</option>
        </select>
      </div>
      <div class="row-actions">
        <button type="button" class="secondary-btn remove-row-btn">Удалить</button>
      </div>
    `;

    const termInput = row.querySelector('.term-input');
    const defInput = row.querySelector('.definition-input');

    termInput.addEventListener('input', updateButtonsState);
    defInput.addEventListener('input', updateButtonsState);

    const removeBtn = row.querySelector('.remove-row-btn');
    removeBtn.addEventListener('click', () => {
      cardRowsContainer.removeChild(row);
      cardRowCount--;
      updateButtonsState();
    });

    return row;
  }

  function addOneRow() {
    if (cardRowCount >= MAX_ROWS) return;
    const row = createRowElement();
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

    // Сохранить доступно, если есть хотя бы одна строка и все заполнены
    saveCardsBtn.disabled = !(cardRowCount > 0 && allRowsFilled());

    if (cardRowCount >= MAX_ROWS) {
      addCardRowBtn.disabled = true;
      // если хочешь жёстко блокировать сохранение при 100 — раскомментируй:
      saveCardsBtn.disabled = true;
    }
  }

  addCardRowBtn.addEventListener('click', () => {
    if (cardRowCount >= MAX_ROWS) return;
    if (!allRowsFilled()) return;
    addOneRow();
    updateButtonsState();
  });

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
        if (res.status === 400) {
          throw new Error('Неверные данные');
        }
        if (res.status === 500) {
          throw new Error('Ошибка сервера');
        }
        if (!res.ok) {
          throw new Error('Ошибка сохранения карточек');
        }
        return res.json();
      })
      .then(() => {
        closeModal();
        // Можно перезагрузить страницу или переотправить запрос модуля
        window.location.reload();
      })
      .catch(err => {
        addCardsError.textContent = err.message || 'Ошибка сохранения карточек';
        addCardsError.style.display = 'block';
      });
  });
});
