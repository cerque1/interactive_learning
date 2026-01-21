const API_BASE_URL = window.location.origin;

window.addEventListener('DOMContentLoaded', async () => {
    const token = localStorage.getItem('token');
    if (!token) {
        window.location.href = `/static/login.html?redirect=${encodeURIComponent(window.location.href)}`;
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
    let moduleType = null;
    let isEditMode = false;
    let moduleCards = [];
    let isModuleFavorite = false;
    let userLoaded = false;
    let moduleLoaded = false;
    let originalModuleName = '';

    const studyModuleBtn = document.getElementById('study-module-btn');
    const testModuleBtn = document.getElementById('test-module-btn');
    const favoriteBtn = document.getElementById('toggle-favorite-btn');
    const addCardsBtn = document.getElementById('add-cards-btn');
    const editModuleBtn = document.getElementById('edit-module-btn');
    const toggleModuleTypeBtn = document.getElementById('toggle-module-type-btn');
    const toggleTypeText = document.getElementById('toggle-type-text');
    const deleteModuleBtn = document.getElementById('delete-module-btn');
    const renameModuleBtn = document.getElementById('rename-module-btn');

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
            setTimeout(() => {
                document.body.removeChild(notification);
            }, 300);
        }, 3000);
    }

    function fetchUser() {
        return fetch(`${API_BASE_URL}/api/v1/user/me?is_full=f`, {
            headers: { 'Authorization': `Bearer ${token}` }
        })
        .then(res => {
            if (res.status === 401) {
                window.location.href = `/static/login.html?redirect=${encodeURIComponent(window.location.href)}`;
                return Promise.reject();
            }
            return res.json();
        })
        .then(userData => {
            if (userData && userData.user) {
                currentUserId = userData.user.id;
                const usernameElem = document.getElementById('username');
                usernameElem.textContent = userData.user.name || 'Пользователь';
                usernameElem.style.cursor = 'pointer';
                usernameElem.onclick = () => window.location.href = '/static/profile.html';
            }
            userLoaded = true;
            checkShowButtons();
        })
        .catch(() => {
            userLoaded = true;
            checkShowButtons();
        });
    }

    function loadModuleFavorite() {
        return fetch(`${API_BASE_URL}/api/v1/selected/modules/`, {
            headers: { 'Authorization': `Bearer ${token}` }
        })
        .then(res => {
            if (res.status === 401) {
                window.location.href = `/static/login.html?redirect=${encodeURIComponent(window.location.href)}`;
                return false;
            }
            if (!res.ok) return false;
            return res.json();
        })
        .then(data => {
            if (data.selected_modules?.length) {
                isModuleFavorite = data.selected_modules.some(m => m.id == moduleId);
            }
            updateFavoriteButton();
            return isModuleFavorite;
        })
        .catch(() => {
            isModuleFavorite = false;
            updateFavoriteButton();
            return false;
        });
    }

    function updateFavoriteButton() {
        const favoriteBtn = document.getElementById('toggle-favorite-btn');
        const favoriteText = document.getElementById('favorite-text');
        if (!favoriteBtn || !favoriteText) return;
        
        if (isModuleFavorite) {
            favoriteBtn.classList.add('filled');
            favoriteText.textContent = 'Убрать из избранного';
            favoriteBtn.title = 'Убрать из избранного';
        } else {
            favoriteBtn.classList.remove('filled');
            favoriteText.textContent = 'Добавить в избранное';
            favoriteBtn.title = 'Добавить в избранное';
        }
    }

    function toggleModuleFavorite() {
        const favoriteBtn = document.getElementById('toggle-favorite-btn');
        if (!favoriteBtn) return;
        
        if (isModuleFavorite) {
            fetch(`${API_BASE_URL}/api/v1/selected/modules/delete?module_id=${moduleId}`, {
                method: 'DELETE',
                headers: { 'Authorization': `Bearer ${token}` }
            })
            .then(res => {
                if (res.status === 401) {
                    window.location.href = `/static/login.html?redirect=${encodeURIComponent(window.location.href)}`;
                    return;
                }
                if (res.ok) {
                    isModuleFavorite = false;
                    updateFavoriteButton();
                    showSuccessMessage('Модуль убран из избранного');
                } else {
                    alert('Ошибка удаления из избранного');
                }
            })
            .catch(err => {
                console.error('Ошибка:', err);
                alert('Ошибка удаления из избранного');
            });
        } else {
            fetch(`${API_BASE_URL}/api/v1/selected/modules/insert?module_id=${moduleId}`, {
                method: 'POST',
                headers: { 
                    'Authorization': `Bearer ${token}`,
                    'Content-Type': 'application/json'
                }
            })
            .then(res => {
                if (res.status === 401) {
                    window.location.href = `/static/login.html?redirect=${encodeURIComponent(window.location.href)}`;
                    return;
                }
                if (res.ok) {
                    isModuleFavorite = true;
                    updateFavoriteButton();
                    showSuccessMessage('Модуль добавлен в избранное');
                } else {
                    alert('Ошибка добавления в избранное');
                }
            })
            .catch(err => {
                console.error('Ошибка:', err);
                alert('Ошибка добавления в избранное');
            });
        }
    }

    function fetchModule() {
        return fetch(`${API_BASE_URL}/api/v1/module/${moduleId}`, {
            headers: { 'Authorization': `Bearer ${token}` }
        })
        .then(res => {
            if (res.status === 401) {
                window.location.href = `/static/login.html?redirect=${encodeURIComponent(window.location.href)}`;
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
            moduleType = moduleData.type || 0;
            moduleCards = moduleData.cards || [];

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
            updateModuleTypeButton();
        })
        .catch(() => {
            document.getElementById('module-name').textContent = 'Ошибка загрузки модуля';
            document.getElementById('cards-container').innerHTML = '';
            moduleLoaded = true;
            checkShowButtons();
        });
    }

    function checkShowButtons() {
        if (userLoaded && moduleLoaded) {
            const isOwner = currentUserId && moduleOwnerId && currentUserId == moduleOwnerId;
            
            if (studyModuleBtn) studyModuleBtn.style.display = moduleCards.length > 0 ? 'inline-block' : 'none';
            if (testModuleBtn) testModuleBtn.style.display = moduleCards.length > 0 ? 'inline-block' : 'none';
            
            const editButtonsContainer = document.getElementById('edit-buttons-container');
            if (editButtonsContainer) {
                editButtonsContainer.style.display = isOwner ? 'flex' : 'none';
            }
        }
    }

    function updateModuleTypeButton() {  
        if (!toggleModuleTypeBtn || moduleType === null) return; 
        
        if (moduleType === 0) {
            toggleTypeText.textContent = 'Сделать приватным'; 
            toggleModuleTypeBtn.title = 'Изменить модуль на приватный (только для владельцев)'; 
        } else {
            toggleTypeText.textContent = 'Сделать публичным'; 
            toggleModuleTypeBtn.title = 'Изменить модуль на публичный (доступен всем)'; 
        } 
    }

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
        
        moduleCards = Array.from(remainingCards).map(card => ({
            id: parseInt(card.dataset.cardId)
        }));
        checkShowButtons();
    }

    function handleDeleteClick(e) {
        e.stopPropagation();
        const card = e.target.closest('.card');
        const cardId = parseInt(card.dataset.cardId);
        
        if (!cardId || !confirm('Вы уверены, что хотите удалить эту карточку?')) {
            return;
        }
        
        fetch(`${API_BASE_URL}/api/v1/card/delete/${cardId}`, {
            method: 'DELETE',
            headers: { 'Authorization': `Bearer ${token}` }
        })
        .then(res => {
            if (res.status === 401) {
                window.location.href = `/static/login.html?redirect=${encodeURIComponent(window.location.href)}`;
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

    if (editModuleBtn) {
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
    }

    if (toggleModuleTypeBtn) {
        toggleModuleTypeBtn.addEventListener('click', () => {
            if (!(currentUserId && moduleOwnerId && Number(currentUserId) === Number(moduleOwnerId))) {
                return;
            }

            const newType = moduleType === 0 ? 1 : 0;
            
            fetch(`${API_BASE_URL}/api/v1/module/change_type/${moduleId}`, {
                method: 'PUT',
                headers: {
                    'Authorization': `Bearer ${token}`,
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({ type: newType })
            })
            .then(res => {
                if (res.status === 401) {
                    window.location.href = `/static/login.html?redirect=${encodeURIComponent(window.location.href)}`;
                    return Promise.reject();
                }
                if (!res.ok) {
                    if (res.status === 403) throw new Error('Нет прав для изменения типа модуля');
                    if (res.status === 404) throw new Error('Модуль не найден');
                    throw new Error('Ошибка изменения типа модуля');
                }
            })
            .then(() => {
                moduleType = newType;
                updateModuleTypeButton();
                showSuccessMessage(  
                    newType === 0 
                        ? 'Модуль теперь публичный (доступен всем)' 
                        : 'Модуль теперь приватный (только для владельца)' 
                ); 
            })
            .catch(err => {
                console.error('Ошибка изменения типа модуля:', err);
                alert('Ошибка изменения типа модуля: ' + err.message);
            });
        });
    }

    if (deleteModuleBtn) {
        deleteModuleBtn.addEventListener('click', () => {
            if (!(currentUserId && moduleOwnerId && Number(currentUserId) === Number(moduleOwnerId))) {
                return;
            }
            if (confirm('Вы уверены, что хотите удалить этот модуль? Все карточки будут удалены навсегда.')) {
                fetch(`${API_BASE_URL}/api/v1/module/delete/${moduleId}`, {
                    method: 'DELETE',
                    headers: { 'Authorization': `Bearer ${token}` }
                })
                .then(res => {
                    if (res.status === 401) {
                        window.location.href = `/static/login.html?redirect=${encodeURIComponent(window.location.href)}`;
                        return;
                    }
                    if (!res.ok) throw new Error('Ошибка удаления модуля');
                    return res.json();
                })
                .then(() => {
                    showSuccessMessage('Модуль успешно удален');
                    setTimeout(() => {
                        window.location.href = '/static/modules.html';
                    }, 1500);
                })
                .catch(err => {
                    alert('Ошибка удаления модуля: ' + err.message);
                });
            }
        });
    }

    const renameModuleModal = document.getElementById('renameModuleModal');
    const closeRenameModal = document.getElementById('closeRenameModal');
    const cancelRenameModal = document.getElementById('cancel-rename-modal');
    const confirmRenameBtn = document.getElementById('confirm-rename-btn');
    const renameModuleInput = document.getElementById('rename-module-input');
    const renameError = document.getElementById('rename-error');

    if (renameModuleBtn) {
        renameModuleBtn.addEventListener('click', () => {
            if (!(currentUserId && moduleOwnerId && Number(currentUserId) === Number(moduleOwnerId))) {
                return;
            }
            originalModuleName = document.getElementById('module-name').textContent;
            renameModuleInput.value = originalModuleName;
            renameError.style.display = 'none';
            confirmRenameBtn.disabled = true;
            renameModuleModal.style.display = 'flex';
            renameModuleInput.focus();
        });
    }

    function validateRenameInput() {
        const newName = renameModuleInput.value.trim();
        const isValid = newName && newName !== originalModuleName;
        confirmRenameBtn.disabled = !isValid;
        
        if (newName === originalModuleName) {
            renameError.textContent = 'Название должно отличаться от текущего';
            renameError.style.display = 'block';
        } else if (!newName) {
            renameError.style.display = 'none';
        } else {
            renameError.style.display = 'none';
        }
    }

    if (renameModuleInput) {
        renameModuleInput.addEventListener('input', validateRenameInput);
    }

    if (closeRenameModal) closeRenameModal.addEventListener('click', () => {
        renameModuleModal.style.display = 'none';
    });

    if (cancelRenameModal) cancelRenameModal.addEventListener('click', () => {
        renameModuleModal.style.display = 'none';
    });

    if (renameModuleModal) {
        renameModuleModal.addEventListener('click', (e) => {
            if (e.target === renameModuleModal) renameModuleModal.style.display = 'none';
        });
    }

    if (confirmRenameBtn) {
        confirmRenameBtn.addEventListener('click', () => {
            const newName = renameModuleInput.value.trim();
            if (!newName || newName === originalModuleName) {
                validateRenameInput();
                return;
            }

            fetch(`${API_BASE_URL}/api/v1/module/rename/${moduleId}`, {
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
                if (res.status === 400) throw new Error('Неверное название модуля');
                if (res.status === 409) throw new Error('Модуль с таким названием уже существует');
                if (!res.ok) throw new Error('Ошибка переименования модуля');
                return res.json();
            })
            .then(() => {
                document.getElementById('module-name').textContent = newName;
                renameModuleModal.style.display = 'none';
                showSuccessMessage('Модуль переименован');
            })
            .catch(err => {
                renameError.textContent = err.message || 'Ошибка переименования модуля';
                renameError.style.display = 'block';
            });
        });
    }

    if (favoriteBtn) {
        favoriteBtn.addEventListener('click', () => toggleModuleFavorite());
    }

    if (studyModuleBtn && moduleId) {
        studyModuleBtn.addEventListener('click', () => {
            window.location.href = `/static/learning.html?modules_ids=${moduleId}`;
        });
    }

    if (testModuleBtn && moduleId) {
        testModuleBtn.addEventListener('click', () => {
            window.location.href = `/static/test.html?modules_ids=${moduleId}`;
        });
    }

    document.addEventListener('click', (e) => {
        if (e.target.classList.contains('delete')) {
            handleDeleteClick(e);
        } else if (e.target.classList.contains('edit')) {
            handleEditClick(e);
        }
    });

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

    if (addCardsBtn) {
        addCardsBtn.addEventListener('click', () => {
            if (!(currentUserId && moduleOwnerId && Number(currentUserId) === Number(moduleOwnerId))) {
                return;
            }
            openModal();
        });
    }

    if (closeAddCardsModal) closeAddCardsModal.addEventListener('click', closeModal);
    if (cancelAddCardsModal) cancelAddCardsModal.addEventListener('click', closeModal);

    if (addCardsModal) {
        addCardsModal.addEventListener('click', (e) => {
            if (e.target === addCardsModal) closeModal();
        });
    }

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

    if (addCardRowBtn) {
        addCardRowBtn.addEventListener('click', () => {
            if (cardRowCount >= MAX_ROWS || !allRowsFilled()) return;
            addOneRow();
            updateButtonsState();
        });
    }

    if (saveCardsBtn) {
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

            fetch(`${API_BASE_URL}/api/v1/card/insert_to_module`, {
                method: 'POST',
                headers: {
                    'Authorization': `Bearer ${token}`,
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({ cards: cardsPayload, parent_module: parseInt(moduleId) })
            })
            .then(res => {
                if (res.status === 401) {
                    window.location.href = `/static/login.html?redirect=${encodeURIComponent(window.location.href)}`;
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
                
                moduleCards = moduleCards.concat(newIds.map(id => ({ id })));
                checkShowButtons();
                
                showSuccessMessage(`${newIds.length} карточек${newIds.length === 1 ? 'а' : 'ек'} успешно добавлено${newIds.length === 1 ? '' : 'ы'}`);
            })
            .catch(err => {
                addCardsError.textContent = err.message || 'Ошибка сохранения карточек';
                addCardsError.style.display = 'block';
            });
        });
    }

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

    if (editTermInput) {
        editTermInput.addEventListener('input', () => {
            updateCardBtn.disabled = !checkEditFieldsFilled();
        });
    }

    if (editDefInput) {
        editDefInput.addEventListener('input', () => {
            updateCardBtn.disabled = !checkEditFieldsFilled();
        });
    }

    function closeEditModal() {
        if (editCardModal) editCardModal.style.display = 'none';
        window.currentEditingCardId = null;
    }

    if (closeEditCardModal) closeEditCardModal.addEventListener('click', closeEditModal);
    if (cancelEditCardModal) cancelEditCardModal.addEventListener('click', closeEditModal);

    if (editCardModal) {
        editCardModal.addEventListener('click', (e) => {
            if (e.target === editCardModal) closeEditModal();
        });
    }

    if (updateCardBtn) {
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

            fetch(`${API_BASE_URL}/api/v1/card/update/${cardId}`, {
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
                    window.location.href = `/static/login.html?redirect=${encodeURIComponent(window.location.href)}`;
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
    }

    const navToggle = document.getElementById('nav-toggle');
    const navPanel = document.getElementById('nav-panel');
    if (navToggle && navPanel) {
        navToggle.addEventListener('click', () => {
            navPanel.classList.toggle('open');
            navToggle.classList.toggle('open');
        });
    }

    ['main-btn', 'modules-btn', 'categories-btn', 'selected-btn', 'results-btn'].forEach(id => {
        const btn = document.getElementById(id);
        if (btn) btn.addEventListener('click', () => window.location.href = `/static/${id.replace('-btn', '.html')}`);
    });

    const head = document.getElementById('head');
    if (head) head.addEventListener('click', () => window.location.href = '/static/main.html');

    await fetchUser();
    await loadModuleFavorite();
    await fetchModule();
});
