let cards = [];
let currentCardIndex = 0;
let isFlipped = false;
let knownCount = 0;
let unknownCount = 0;
let moduleId = null;
let moduleData = null;
let windowMyId = null;
let isDragging = false;
let isDown = false;
let dragThreshold = 30;

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

function showCard(index) {
    if (index >= cards.length) {
        showResults();
        return;
    }

    const container = document.getElementById('cards-container');
    container.innerHTML = '';

    const card = document.createElement('div');
    card.className = 'flashcard';
    card.dataset.cardIndex = index;

    const cardData = cards[index];
    const termLang = cardData.term.lang;
    const definitionLang = cardData.definition.lang;

    card.innerHTML = `
        <div class="flashcard-container">
            <div class="flashcard-side front">
                <div class="flashcard-lang">${termLang}</div>
                <div class="flashcard-content">${cardData.term.text}</div>
            </div>
            <div class="flashcard-side back">
                <div class="flashcard-lang">${definitionLang}</div>
                <div class="flashcard-content">${cardData.definition.text}</div>
            </div>
        </div>
    `;

    container.appendChild(card);

    // Принудительное центрирование
    requestAnimationFrame(() => {
        card.style.transform = 'translate(-50%, -50%)';
        setupCardEvents(card);
    });

    updateProgress();
}

function setupCardEvents(card) {
    let startX = 0;
    let startY = 0;
    let currentX = 0;
    let currentY = 0;
    let isTouch = false;

    const resetTransform = () => {
        card.style.transform = 'translate(-50%, -50%) translateX(0px) rotate(0deg) scale(1)';
    };

    card.addEventListener('click', (e) => {
        if (!isDragging && Math.abs(currentX - startX) < dragThreshold && Math.abs(currentY - startY) < dragThreshold) {
            toggleFlip();
        }
    });

    const handleStart = (e) => {
        isDragging = false;
        isDown = true;
        isTouch = e.type === 'touchstart';
        
        startX = isTouch ? e.touches[0].clientX : e.clientX;
        startY = isTouch ? e.touches[0].clientY : e.clientY;
        currentX = startX;
        currentY = startY;
        
        card.style.transition = 'none';
        resetTransform();
    };

    const handleMove = (e) => {
        if (!isDown) return;
        
        currentX = isTouch ? e.touches[0].clientX : e.clientX;
        currentY = isTouch ? e.touches[0].clientY : e.clientY;
        
        const deltaX = currentX - startX;
        const deltaY = currentY - startY;
        
        if (Math.abs(deltaX) > dragThreshold || Math.abs(deltaY) > dragThreshold) {
            isDragging = true;
            
            if (Math.abs(deltaY) > Math.abs(deltaX)) return;
            
            const rotate = deltaX / 15;
            const scale = Math.max(0.95, 1 - Math.abs(deltaX) / 500);
            
            card.style.transform = `translate(-50%, -50%) translateX(${deltaX}px) rotate(${rotate}deg) scale(${scale})`;
        }
    };

    const handleEnd = (e) => {
        if (!isDown) return;
        
        isDown = false;
        card.style.transition = 'transform 0.6s cubic-bezier(0.23, 1, 0.32, 1)';
        
        const deltaX = currentX - startX;
        const threshold = 100;
        
        if (isDragging && Math.abs(deltaX) > threshold) {
            if (deltaX > 0) {
                handleSwipe('known');
            } else {
                handleSwipe('unknown');
            }
        } else {
            resetTransform();
        }
    };

    card.addEventListener('mousedown', handleStart);
    card.addEventListener('mousemove', handleMove);
    card.addEventListener('mouseup', handleEnd);
    card.addEventListener('mouseleave', handleEnd);

    card.addEventListener('touchstart', handleStart, { passive: true });
    card.addEventListener('touchmove', handleMove, { passive: false });
    card.addEventListener('touchend', handleEnd);
    card.addEventListener('touchcancel', handleEnd);

    card.addEventListener('transitionend', resetTransform);
}

function toggleFlip() {
    const card = document.querySelector('.flashcard');
    if (!card) return;
    
    card.classList.toggle('flipped');
    isFlipped = !isFlipped;
}

function handleSwipe(result) {
    const card = document.querySelector('.flashcard');
    if (!card) return;

    const deltaX = result === 'known' ? 400 : -400;
    card.style.transform = `translate(-50%, -50%) translateX(${deltaX}px) rotate(${deltaX / 10}deg) scale(0.8)`;

    if (result === 'known') {
        knownCount++;
        document.getElementById('known-count').textContent = knownCount;
    } else {
        unknownCount++;
        document.getElementById('unknown-count').textContent = unknownCount;
    }

    setTimeout(() => {
        currentCardIndex++;
        showCard(currentCardIndex);
    }, 400);
}

function updateProgress() {
    const progress = ((currentCardIndex / cards.length) * 100);
    document.getElementById('progress-fill').style.width = `${progress}%`;
    document.getElementById('total-cards').textContent = cards.length;
}

function showResults() {
    document.getElementById('study-area').style.display = 'none';
    const resultsScreen = document.getElementById('results-screen');
    
    const percent = cards.length > 0 ? Math.round((knownCount / cards.length) * 100) : 0;
    
    document.getElementById('results-known').textContent = knownCount;
    document.getElementById('results-unknown').textContent = unknownCount;
    document.getElementById('results-percent').textContent = `${percent}%`;
    
    resultsScreen.classList.remove('hidden');
    
    document.getElementById('repeat-btn').onclick = restartStudy;
    document.getElementById('back-to-module-btn').onclick = () => {
        window.location.href = `/static/module.html?module_id=${moduleId}`;
    };
}

function restartStudy() {
    currentCardIndex = 0;
    knownCount = 0;
    unknownCount = 0;
    isFlipped = false;
    
    document.getElementById('known-count').textContent = '0';
    document.getElementById('unknown-count').textContent = '0';
    
    document.getElementById('study-area').style.display = 'flex';
    document.getElementById('results-screen').classList.add('hidden');
    
    showCard(0);
}

async function loadModule(token) {
    const params = new URLSearchParams(window.location.search);
    moduleId = params.get('module_id');
    
    if (!moduleId) {
        alert('Не указан ID модуля');
        window.location.href = '/static/modules.html';
        return;
    }

    // ✅ ПОКАЗЫВАЕМ КНОПКУ НАЗАД К МОДУЛЮ
    const backBtn = document.getElementById('backToModuleBtn');
    if (backBtn) {
        backBtn.style.display = 'inline-flex';
        backBtn.addEventListener('click', () => {
            window.location.href = `/static/module.html?module_id=${moduleId}`;
        });
    }

    try {
        const response = await fetch(`http://localhost:8080/api/v1/module/${moduleId}`, {
            headers: { 'Authorization': `Bearer ${token}` }
        });

        if (response.status === 401) {
            window.location.href = `/static/login.html?redirect=${encodeURIComponent(window.location.href)}`;
            return;
        }

        if (!response.ok) {
            throw new Error('Не удалось загрузить модуль');
        }

        moduleData = await response.json();
        document.getElementById('module-title').textContent = moduleData.module.name;
        
        cards = moduleData.module.cards || [];
        
        if (cards.length === 0) {
            alert('В модуле нет карточек');
            window.location.href = `/static/module.html?module_id=${moduleId}`;
            return;
        }

        showCard(0);
    } catch (error) {
        console.error('Ошибка загрузки модуля:', error);
        alert('Ошибка загрузки модуля');
    }
}

window.addEventListener('DOMContentLoaded', async () => {
    const token = localStorage.getItem('token');
    if (!token) {
        window.location.href = `/static/login.html?redirect=${encodeURIComponent(window.location.href)}`;
        return;
    }

    await loadUserName(token);
    setupNavigation();
    await loadModule(token);
});
