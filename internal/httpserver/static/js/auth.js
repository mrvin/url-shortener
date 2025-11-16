class AuthManager {
    constructor() {
        this.checkAuthStatus();
    }

    // Проверка статуса аутентификации
    checkAuthStatus() {
        const credentials = localStorage.getItem('credentials');
        const isLoggedIn = !!credentials;
        
        // Обновляем навигацию в зависимости от статуса
        this.updateNavigation(isLoggedIn);
        return isLoggedIn;
    }

    // Обновление навигации
    updateNavigation(isLoggedIn) {
        const authElements = document.querySelectorAll('.auth-only');
        const unauthElements = document.querySelectorAll('.unauth-only');
        
        if (isLoggedIn) {
            authElements.forEach(el => el.style.display = 'block');
            unauthElements.forEach(el => el.style.display = 'none');
        } else {
            authElements.forEach(el => el.style.display = 'none');
            unauthElements.forEach(el => el.style.display = 'block');
        }
    }

    // Вход
    async login(username, password) {
        try {
            const result = await api.login({ username, password });
            
            if (result.status === 'OK') {
                localStorage.setItem('credentials', JSON.stringify({ username, password }));
                this.checkAuthStatus();
                return { success: true };
            } else {
                return { success: false, error: 'Неверные учетные данные' };
            }
        } catch (error) {
            return { success: false, error: 'Ошибка сети' };
        }
    }

    // Регистрация
    async register(username, password) {
        try {
            const result = await api.register({ username, password });
            
            if (result.status === 'OK') {
                // Автоматически логиним после регистрации
                return await this.login(username, password);
            } else {
                return { success: false, error: 'Ошибка регистрации' };
            }
        } catch (error) {
            return { success: false, error: 'Ошибка сети' };
        }
    }

    // Выход
    logout() {
        localStorage.removeItem('credentials');
        this.checkAuthStatus();
        window.location.href = 'index.html';
    }

    // Получение текущего пользователя
    getCurrentUser() {
        const credentials = JSON.parse(localStorage.getItem('credentials') || '{}');
        return credentials.username;
    }
}

const authManager = new AuthManager();