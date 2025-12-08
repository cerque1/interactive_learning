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

          // Получаем параметр redirect из URL страницы входа
          const params = new URLSearchParams(window.location.search);
          const redirectUrl = params.get('redirect');

          // Переходим на redirectUrl, либо на главную если параметр отсутствует
          if (redirectUrl) {
            window.location.href = redirectUrl;
          } else {
            window.location.href = '/static/main.html';
          }
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

const registerLink = document.getElementById('registerLink');

registerLink.addEventListener('click', (e) => {
  e.preventDefault(); // отменяем переход по ссылке

  const params = new URLSearchParams(window.location.search);
  const redirect = params.get('redirect');

  let targetUrl = '/static/register.html';
  if (redirect) {
    targetUrl += `?redirect=${encodeURIComponent(redirect)}`;
  }

  window.location.href = targetUrl;
});
