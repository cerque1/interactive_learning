let cards = [];
let currentQuestionIndex = 0;
let correctAnswers = 0;
let incorrectAnswers = 0;
let answeredQuestions = new Set();
let categoryId = null;
let modulesIds = [];
let isAnswerLocked = false;
let moduleData = null;

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

function setupBackButton() {
    const backBtn = document.getElementById('backBtn');
    if (!backBtn) return false;
    
    const params = new URLSearchParams(window.location.search);
    
    const categoryIdParam = params.get('category_id');
    if (categoryIdParam) {
        categoryId = parseInt(categoryIdParam);
        backBtn.textContent = '← Назад к категории';
        backBtn.title = 'Вернуться к категории';
        backBtn.style.display = 'inline-flex';
        backBtn.onclick = () => {
            window.location.href = `/static/category.html?category_id=${categoryId}`;
        };
        return true;
    }

    const moduleIdParam = params.get('modules_ids');
    if (moduleIdParam) {
        backBtn.textContent = '← Назад к модулю';
        backBtn.title = 'Вернуться к модулю';
        backBtn.style.display = 'inline-flex';
        backBtn.onclick = () => {
            window.location.href = `/static/module.html?module_id=${moduleIdParam}`;
        };
        return true;
    }

    backBtn.style.display = 'none';
    return false;
}

function shuffleArray(array) {
    const shuffled = [...array];
    for (let i = shuffled.length - 1; i > 0; i--) {
        const j = Math.floor(Math.random() * (i + 1));
        [shuffled[i], shuffled[j]] = [shuffled[j], shuffled[i]];
    }
    return shuffled;
}

function generateQuestion(index) {
    const card = cards[index];
    if (!card || !card.term || !card.definition || !card.term.text || !card.definition.text) {
        console.error('Некорректная карточка:', card);
        return null;
    }
    
    const term = card.term.text;
    const correctDefinition = card.definition.text;
    
    const allDefinitions = cards
        .filter(c => c && c.definition && c.definition.text)
        .map(c => c.definition.text);
    
    let uniqueWrongDefinitions = [];
    
    allDefinitions.forEach(def => {
        if (def !== correctDefinition && !uniqueWrongDefinitions.includes(def)) {
            uniqueWrongDefinitions.push(def);
        }
    });
    
    const totalOptionsNeeded = Math.min(4, allDefinitions.length);
    const wrongOptionsNeeded = totalOptionsNeeded - 1;
    
    let wrongDefinitions = [];
    
    if (uniqueWrongDefinitions.length >= wrongOptionsNeeded) {
        wrongDefinitions = uniqueWrongDefinitions.slice(0, wrongOptionsNeeded);
    } else {
        const availableUnique = uniqueWrongDefinitions.length;
        wrongDefinitions = [...uniqueWrongDefinitions];
        
        const additionalNeeded = wrongOptionsNeeded - availableUnique;
        for (let i = 0; i < additionalNeeded; i++) {
            wrongDefinitions.push(uniqueWrongDefinitions[i % availableUnique] || "Неизвестно");
        }
    }
    
    const answerOptions = [...wrongDefinitions, correctDefinition];
    const shuffled = shuffleArray(answerOptions);
    
    return { term, correctDefinition, shuffledAnswers: shuffled, totalOptions: totalOptionsNeeded };
}

function showQuestion(index) {
    if (index >= cards.length || index < 0 || cards.length === 0) {
        showResults();
        return;
    }

    isAnswerLocked = false;

    const questionData = generateQuestion(index);
    if (!questionData) {
        alert('Ошибка генерации вопроса');
        showResults();
        return;
    }

    const questionTermElem = document.getElementById('question-term');
    const currentQuestionElem = document.getElementById('current-question');
    const totalQuestionsCountElem = document.getElementById('total-questions-count');
    
    if (questionTermElem) questionTermElem.textContent = questionData.term;
    if (currentQuestionElem) currentQuestionElem.textContent = index + 1;
    if (totalQuestionsCountElem) totalQuestionsCountElem.textContent = cards.length;
    
    const answersContainer = document.getElementById('answers-container');
    if (!answersContainer) {
        console.error('answers-container не найден');
        return;
    }
    
    answersContainer.innerHTML = '';
    
    if (!questionData.shuffledAnswers || !Array.isArray(questionData.shuffledAnswers)) {
        console.error('shuffledAnswers некорректен:', questionData.shuffledAnswers);
        return;
    }
    
    questionData.shuffledAnswers.forEach((answer, i) => {
        if (!answer || typeof answer !== 'string') return;
        
        const option = document.createElement('div');
        option.className = 'answer-option';
        option.dataset.index = i;
        option.dataset.answerText = answer;
        option.textContent = answer;
        
        option.addEventListener('click', function clickHandler(e) {
            e.stopPropagation();
            if (isAnswerLocked) return;
            
            option.removeEventListener('click', clickHandler);
            
            document.querySelectorAll('.answer-option').forEach(opt => {
                opt.style.pointerEvents = 'none';
            });
            
            selectAnswer(i, answer, questionData.correctDefinition);
        });
        
        answersContainer.appendChild(option);
    });
    
    const nextBtn = document.getElementById('next-btn');
    if (nextBtn) nextBtn.style.display = 'none';
    
    answeredQuestions.delete(index);
}

function selectAnswer(selectedIndex, selectedAnswer, correctDefinition) {
    if (isAnswerLocked) return;
    
    isAnswerLocked = true;
    
    const options = document.querySelectorAll('.answer-option');
    let isCorrect = false;
    
    options.forEach((option, i) => {
        const optionText = option.dataset.answerText || option.textContent;
        
        option.classList.remove('selected', 'correct', 'incorrect');
        
        if (optionText === correctDefinition) {
            option.classList.add('correct');
            if (i === selectedIndex) {
                isCorrect = true;
            }
        } else if (i === selectedIndex) {
            option.classList.add('incorrect');
            isCorrect = false;
        }
    });
    
    if (isCorrect) {
        correctAnswers++;
    } else {
        incorrectAnswers++;
    }
    
    answeredQuestions.add(currentQuestionIndex);
    const nextBtn = document.getElementById('next-btn');
    if (nextBtn) {
        if (currentQuestionIndex + 1 >= cards.length){
            nextBtn.style.display = 'none';
        }
        else {
            nextBtn.style.display = 'inline-flex';
        }
    }
}

function nextQuestion() {
    currentQuestionIndex++;
    showQuestion(currentQuestionIndex);
}

function endTestPrematurely() {
    const remainingQuestions = cards.length - answeredQuestions.size;
    incorrectAnswers += remainingQuestions;
    showResults();
}

function showResults() {
    const testArea = document.getElementById('test-area');
    if (testArea) testArea.style.display = 'none';
    
    const resultsScreen = document.getElementById('results-screen');
    if (!resultsScreen) return;
    
    const percent = cards.length > 0 ? Math.round((correctAnswers / cards.length) * 100) : 0;
    
    const resultsCorrect = document.getElementById('results-correct');
    const resultsIncorrect = document.getElementById('results-incorrect');
    const resultsPercent = document.getElementById('results-percent');
    
    if (resultsCorrect) resultsCorrect.textContent = correctAnswers;
    if (resultsIncorrect) resultsIncorrect.textContent = incorrectAnswers;
    if (resultsPercent) resultsPercent.textContent = `${percent}%`;
    
    resultsScreen.classList.remove('hidden');
    
    const repeatBtn = document.getElementById('repeat-test-btn');
    if (repeatBtn) repeatBtn.onclick = restartTest;
    
    const backBtn = document.getElementById('back-from-results-btn');
    if (backBtn) {
        if (categoryId) {
            backBtn.textContent = 'К категории';
            backBtn.onclick = () => window.location.href = `/static/category.html?category_id=${categoryId}`;
        } else {
            const params = new URLSearchParams(window.location.search);
            const moduleIdParam = params.get('modules_ids');
            backBtn.textContent = 'К модулю';
            backBtn.onclick = () => {
                window.location.href = moduleIdParam ? 
                    `/static/module.html?module_id=${moduleIdParam}` : '/static/main.html';
            };
        }
    }
}

function restartTest() {
    currentQuestionIndex = 0;
    correctAnswers = 0;
    incorrectAnswers = 0;
    answeredQuestions.clear();
    isAnswerLocked = false;
    
    const testArea = document.getElementById('test-area');
    const resultsScreen = document.getElementById('results-screen');
    
    if (testArea) testArea.style.display = 'flex';
    if (resultsScreen) resultsScreen.classList.add('hidden');
    
    showQuestion(0);
}

async function loadModulesByIds(token, modulesIdsParam, hasCategoryId) {
    try {
        const response = await fetch('http://localhost:8080/api/v1/module/by_ids?with_cards=t', {
            method: 'POST',
            headers: { 
                'Authorization': `Bearer ${token}`,
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ modules_ids: modulesIdsParam })
        });

        if (response.status === 401) {
            window.location.href = `/static/login.html?redirect=${encodeURIComponent(window.location.href)}`;
            return;
        }

        if (!response.ok) {
            throw new Error('Не удалось загрузить модули');
        }

        const data = await response.json();
        
        if (!data) {
            throw new Error('Получен пустой ответ от сервера');
        }
        
        moduleData = data;
        cards = [];
        let moduleTitle = '';
        
        if (data.modules && Array.isArray(data.modules)) {
            data.modules.forEach(module => {
                if (module && module.cards && Array.isArray(module.cards)) {
                    cards = cards.concat(module.cards);
                }
                if (!moduleTitle && module && module.name) {
                    moduleTitle = module.name;
                }
            });
        }
        
        if (cards.length === 0) {
            const backUrl = hasCategoryId ? `/static/category.html?category_id=${categoryId}` : '/static/main.html';
            alert('В выбранных модулях нет карточек');
            window.location.href = backUrl;
            return;
        }
        
        const testTitle = document.getElementById('test-title');
        const totalQuestions = document.getElementById('total-questions');
        if (testTitle) {
            testTitle.textContent = modulesIdsParam.length === 1 ? moduleTitle : `${modulesIdsParam.length} модулей`;
        }
        if (totalQuestions) {
            totalQuestions.textContent = cards.length;
        }
        
        showQuestion(0);
        
    } catch (error) {
        console.error('Ошибка загрузки модулей:', error);
        const backUrl = hasCategoryId ? `/static/category.html?category_id=${categoryId}` : '/static/main.html';
        alert('Ошибка загрузки модулей: ' + (error.message || 'Неизвестная ошибка'));
        window.location.href = backUrl;
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
    
    const params = new URLSearchParams(window.location.search);
    const categoryIdParam = params.get('category_id');
    const modulesIdsParam = params.get('modules_ids');
    const moduleIdParam = params.get('module_id');
    
    if (categoryIdParam) {
        categoryId = parseInt(categoryIdParam);
    }
    
    setupBackButton();
    
    let modulesIdsArray = [];
    
    if (modulesIdsParam) {
        modulesIdsArray = modulesIdsParam.split(',').map(id => parseInt(id.trim())).filter(id => !isNaN(id));
    } else if (moduleIdParam) {
        modulesIdsArray = [parseInt(moduleIdParam)];
    }
    
    if (modulesIdsArray.length === 0) {
        const backUrl = categoryId ? `/static/category.html?category_id=${categoryId}` : '/static/main.html';
        alert('Не указаны корректные ID модулей');
        window.location.href = backUrl;
        return;
    }
    
    modulesIds = modulesIdsArray;
    const hasCategoryId = !!categoryId;
    const testTitle = document.getElementById('test-title');
    if (testTitle) testTitle.textContent = 'Загрузка теста...';
    
    await loadModulesByIds(token, modulesIdsArray, hasCategoryId);
    
    const nextBtn = document.getElementById('next-btn');
    const endTestBtn = document.getElementById('end-test-btn');
    if (nextBtn) nextBtn.onclick = nextQuestion;
    if (endTestBtn) endTestBtn.onclick = endTestPrematurely;
});
