const registerBtn = document.getElementById('registerBtn');
const errorMsg = document.getElementById('errorMsg');
const loginLink = document.getElementById('loginLink');

loginLink.addEventListener('click', () => {
  window.location.href = '/login';
});

registerBtn.addEventListener('click', () => {
  errorMsg.textContent = '';
  const name = document.getElementById('name').value.trim();
  const login = document.getElementById('login').value.trim();
  const password = document.getElementById('password').value.trim();

  if (!name || !login || !password) {
    errorMsg.textContent = 'Заполнены не все поля.';
    return;
  }

  const url = new URL('http://localhost:8080/api/auth/register');
  url.searchParams.append('name', name);
  url.searchParams.append('login', login);
  url.searchParams.append('password', password);

  fetch(url, { method: 'POST' })
    .then(async (response) => {
      if (response.ok) {
        const data = await response.json();
        if (data.token) {
          sessionStorage.setItem('token', data.token);
          localStorage.setItem('token', data.token);
          window.location.href = 'http://localhost:8080/static/main.html';
        } else {
          errorMsg.textContent = 'Ошибка: отсутствует токен в ответе.';
        }
      } else if (response.status === 400) {
        const errData = await response.json();
        if (errData.message === 'wrong data') {
          errorMsg.textContent = 'Неверные данные.';
        } else if (
          errData.message &&
          errData.message.toLowerCase().includes('login already exists')
        ) {
          errorMsg.textContent = 'Логин уже существует.';
        } else {
          errorMsg.textContent = 'Неверные данные.';
        }
      } else if (response.status === 500) {
        errorMsg.textContent = 'Попробуйте позже.';
      } else {
        errorMsg.textContent = 'Произошла ошибка.';
      }
    })
    .catch(() => {
      errorMsg.textContent = 'Ошибка соединения. Попробуйте позже.';
    });
});

const registerLink = document.getElementById('loginLink');

registerLink.addEventListener('click', (e) => {
  e.preventDefault(); // отменяем переход по ссылке

  const params = new URLSearchParams(window.location.search);
  const redirect = params.get('redirect');

  let targetUrl = '/static/login.html';
  if (redirect) {
    targetUrl += `?redirect=${encodeURIComponent(redirect)}`;
  }

  window.location.href = targetUrl;
});
