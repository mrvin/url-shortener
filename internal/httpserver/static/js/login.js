document.addEventListener('DOMContentLoaded', function() {
    const loginForm = document.getElementById('login-form');
    const registerForm = document.getElementById('register-form-inner');
    const showRegister = document.getElementById('show-register');
    const showLogin = document.getElementById('show-login');
    
    // Переключение между формами
    showRegister.addEventListener('click', function(e) {
        e.preventDefault();
        document.getElementById('login-form').parentElement.style.display = 'none';
        document.getElementById('register-form').style.display = 'block';
    });
    
    showLogin.addEventListener('click', function(e) {
        e.preventDefault();
        document.getElementById('register-form').style.display = 'none';
        document.getElementById('login-form').parentElement.style.display = 'block';
    });
    
    // Обработка входа
    loginForm.addEventListener('submit', async function(e) {
        e.preventDefault();
        
        const username = document.getElementById('username').value;
        const password = document.getElementById('password').value;
        const button = loginForm.querySelector('button');
        
        button.disabled = true;
        button.textContent = 'Вход...';
        
        const result = await authManager.login(username, password);
        
        if (result.success) {
            window.location.href = 'dashboard.html';
        } else {
            showAuthResult(result.error, 'danger');
        }
        
        button.disabled = false;
        button.textContent = 'Войти';
    });
    
    // Обработка регистрации
    registerForm.addEventListener('submit', async function(e) {
        e.preventDefault();
        
        const username = document.getElementById('reg-username').value;
        const password = document.getElementById('reg-password').value;
        const button = registerForm.querySelector('button');
        
        button.disabled = true;
        button.textContent = 'Регистрация...';
        
        const result = await authManager.register(username, password);
        
        if (result.success) {
            window.location.href = 'dashboard.html';
        } else {
            showAuthResult(result.error, 'danger');
        }
        
        button.disabled = false;
        button.textContent = 'Зарегистрироваться';
    });
});

function showAuthResult(message, type) {
    const resultDiv = document.getElementById('auth-result');
    resultDiv.innerHTML = `
        <div class="alert alert-${type} alert-dismissible fade show">
            ${message}
            <button type="button" class="btn-close" data-bs-dismiss="alert"></button>
        </div>
    `;
}