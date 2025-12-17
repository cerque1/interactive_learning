window.addEventListener('DOMContentLoaded', async () => {
    const token = localStorage.getItem('token');
    if (!token) {
        window.location.href = `/static/login.html?redirect=${encodeURIComponent(window.location.href)}`;
        return;
    }

    const url = new URL(window.location.href);
    const params = new URLSearchParams(window.location.search);
    const userId = params.get('id');
  
    const myUserData = await loadUserName(token);
    const myId = myUserData ? myUserData.user?.id : null;
  
    loadCategories(token, userId, myId);
    setupModal(token, userId, myId);
});

function loadUserName(token) {
    return fetch('http://localhost:8080/api/v1/user/me?isfull=f', {
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
            usernameElem.textContent = userData.user.name;
            usernameElem.onclick = () => window.location.href = '/static/profile.html';
        }
        return userData;
    })
    .catch(() => null);
}

async function loadCategories(token, userId, myId) {
    let userid = myId;
    if (userId !== null) userid = userId;
  
    let url = `http://localhost:8080/api/v1/category/to_user/${userid}`;
  
    fetch(url, {
        headers: { 'Authorization': `Bearer ${token}` }
    })
    .then(res => {
        if (res.status === 401) {
            window.location.href = `/static/login.html?redirect=${encodeURIComponent(window.location.href)}`;
            return;
        }
        if (!res.ok) throw new Error('Network response was not ok');
        return res.json();
    })
    .then(categories => {
        const container = document.getElementById('categories-container');
        const emptyMsg = document.getElementById('categories-empty');
        const pageTitle = document.getElementById('page-title');
      
        container.innerHTML = '';
      
        if (!categories || !categories.categories || categories.categories.length === 0) {
            emptyMsg.style.display = 'block';
        } else {
            emptyMsg.style.display = 'none';
            categories.categories.forEach(category => {
                const card = document.createElement('div');
                card.className = 'card';
                card.innerHTML = `<div class="card-title">${category.name}</div>`;
                card.onclick = () => window.location.href = `/static/category.html?category_id=${category.id}`;
                container.appendChild(card);
            });
        }
      
        if (userId) {
            pageTitle.textContent = `Категории пользователя`;
        } else {
            pageTitle.textContent = 'Мои категории';
        }
    })
    .catch(() => {
        document.getElementById('categories-empty').textContent = 'Ошибка загрузки категорий';
        document.getElementById('categories-empty').style.display = 'block';
    });
}

async function loadModules(token, userId) {
    const url = `http://localhost:8080/api/v1/module/to_user/${userId}?with_cards=t`;
    
    try {
        const res = await fetch(url, {
            headers: { 'Authorization': `Bearer ${token}` }
        });
        
        if (res.status === 401) {
            window.location.href = `/static/login.html?redirect=${encodeURIComponent(window.location.href)}`;
            return;
        }
        
        if (!res.ok) throw new Error('Network response was not ok');
        
        const data = await res.json();
        return data.modules || [];
    } catch {
        return [];
    }
}

function setupModal(token, userId, myId) {
    const modal = document.getElementById('createCategoryModal');
    const createBtn = document.getElementById('createCategoryBtn');
    const headerActions = document.getElementById('header-actions');
    const closeBtn = document.getElementById('closeCategoryModal');
    const cancelBtn = document.getElementById('cancelCategoryModal');
    const confirmBtn = document.getElementById('createCategoryConfirm');
  
    if (userId) {
        headerActions.style.display = 'none';
    }
  
    createBtn.onclick = async () => {
        modal.style.display = 'flex';
        await loadModulesForModal(token, myId || userId);
    };

    function closeModal() {
        modal.style.display = 'none';
        document.getElementById('categoryName').value = '';
        document.getElementById('modules-list').innerHTML = '';
        document.getElementById('modal-error').style.display = 'none';
    }

    closeBtn.onclick = closeModal;
    cancelBtn.onclick = closeModal;
  
    modal.onclick = (e) => {
        if (e.target === modal) closeModal();
    };

    confirmBtn.onclick = async () => {
        const name = document.getElementById('categoryName').value.trim();
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
            const res = await fetch('http://localhost:8080/api/v1/category/create', {
                method: 'POST',
                headers: {
                    'Authorization': `Bearer ${token}`,
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({ name, modules_ids: selectedModules })
            }).then(res => {
                if (res.status === 401) {
                    window.location.href = `/static/login.html?redirect=${encodeURIComponent(window.location.href)}`;
                    return;
                }
                
                if (res.status === 400) {
                    showModalError('Ошибка: неверные данные');
                    return;
                }
                
                if (res.status === 500) {
                    showModalError('Ошибка сервера');
                    return;
                }
                
                if (!res.ok) {
                    showModalError('Ошибка создания категории');
                    return;
                }
                return res.json();
            }).then(newCategory => {
                closeModal();
                console.log(newCategory.new_id);
                addCategoryCard(newCategory.new_id, name);
            });
        } catch (err) {
            showModalError('Ошибка создания категории');
        }
    };
}

async function loadModulesForModal(token, userId) {
    const modulesList = document.getElementById('modules-list');
    const modules = await loadModules(token, userId);
    
    modulesList.innerHTML = '';
    
    if (modules.length === 0) {
        modulesList.innerHTML = '<div class="no-modules">Модули отсутствуют</div>';
        return;
    }
    
    modules.forEach(module => {
        const div = document.createElement('div');
        div.className = 'module-item';
        div.innerHTML = `
            <label class="module-label">
                <input type="checkbox" class="module-checkbox" data-module-id="${module.id}">
                <span>${module.name} (${module.cards_count || 0} карточек)</span>
            </label>
        `;
        modulesList.appendChild(div);
    });
}

function showModalError(message) {
    const errorEl = document.getElementById('modal-error');
    errorEl.textContent = message;
    errorEl.style.display = 'block';
}

function addCategoryCard(categoryId, categoryName) {
    const container = document.getElementById('categories-container');
    const emptyMsg = document.getElementById('categories-empty');
    
    const card = document.createElement('div');
    card.className = 'card';
    card.innerHTML = `<div class="card-title">${categoryName}</div>`;
    card.onclick = () => window.location.href = `/static/category.html?category_id=${categoryId}`;
    
    container.appendChild(card);
    emptyMsg.style.display = 'none';
}

// Навигация
const navToggle = document.getElementById('nav-toggle');
const navPanel = document.getElementById('nav-panel');

if (navToggle && navPanel) {
    navToggle.addEventListener('click', function() {
        navPanel.classList.toggle('open');
        navToggle.classList.toggle('open');
    });
}

const navModulesBtn = document.getElementById('modules-btn');
if (navModulesBtn) {
    navModulesBtn.addEventListener('click', function() {
        window.location.href = '/static/modules.html';
    });
}

const navCategoriesBtn = document.getElementById('categories-btn');
if (navCategoriesBtn) {
    navCategoriesBtn.addEventListener('click', function() {
        window.location.href = '/static/categories.html';
    });
}

const navMainBtn = document.getElementById('main-btn');
if (navMainBtn) {
    navMainBtn.addEventListener('click', function() {
        window.location.href = '/static/main.html';
    });
}

const head = document.getElementById('head');
if (head) {
    head.addEventListener('click', e => {
        e.preventDefault();
        window.location.href = '/static/main.html';
    });
}
