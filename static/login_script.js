const loginBtn = document.getElementById('loginBtn');
const errorMsg = document.getElementById('errorMsg');
const registerLink = document.getElementById('registerLink');

registerLink.addEventListener('click', () => {
  window.location.href = '/register';
});

loginBtn.addEventListener('click', () => {
  errorMsg.textContent = '';
  const login = document.getElementById('login').value.trim();
  const password = document.getElementById('password').value.trim();

  if (!login || !password) {
    errorMsg.textContent = 'Заполнены не все поля.';
    return;
  }

  const url = new URL('http://localhost:8080/api/auth/login');
  url.searchParams.append('login', login);
  url.searchParams.append('password', password);

  fetch(url, { method: 'POST' })
    .then(async (response) => {
      if (response.status === 200) {
        const data = await response.json();
        if (data.token) {
          sessionStorage.setItem('token', data.token);
          localStorage.setItem('token', data.token);
          window.location.href = 'http://localhost:8080/static/main.html';
        } else {
          errorMsg.textContent = 'Ошибка: отсутствует токен в ответе.';
        }
      } else if (response.status === 401) {
        errorMsg.textContent = 'Неверно введён логин или пароль.';
      } else if (response.status === 404) {
        errorMsg.textContent = 'Заполнены не все поля.';
      } else {
        errorMsg.textContent = 'Произошла ошибка.';
      }
    })
    .catch(() => {
      errorMsg.textContent = 'Ошибка соединения. Попробуйте позже.';
    });
});
