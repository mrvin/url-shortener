const API_BASE = 'http://localhost:8080';

class URLShortenerAPI {
    constructor() {
        this.baseURL = API_BASE;
    }

    // Получение заголовков с авторизацией
    getAuthHeaders() {
        const credentials = JSON.parse(localStorage.getItem('credentials') || '{}');
        if (credentials.username && credentials.password) {
            return {
                'Authorization': 'Basic ' + btoa(credentials.username + ':' + credentials.password),
                'Content-Type': 'application/json'
            };
        }
        return { 'Content-Type': 'application/json' };
    }

    // Проверка здоровья сервиса
    async healthCheck() {
        const response = await fetch(`${this.baseURL}/api/health`);
        return await response.json();
    }

    // Регистрация пользователя
    async register(userData) {
        const response = await fetch(`${this.baseURL}/api/users`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(userData)
        });
        return await response.json();
    }

    // Вход пользователя
    async login(userData) {
        const response = await fetch(`${this.baseURL}/api/users/login`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(userData)
        });
        return await response.json();
    }

    // Создание короткой ссылки
    async shortenUrl(urlData) {
        const response = await fetch(`${this.baseURL}/api/urls`, {
            method: 'POST',
            headers: this.getAuthHeaders(),
            body: JSON.stringify(urlData)
        });
        return await response.json();
    }

    // Проверка доступности алиаса
    async checkAlias(alias) {
        const response = await fetch(`${this.baseURL}/api/urls/check/${alias}`);
        return await response.json();
    }

    // Получение списка URL с пагинацией
    async getUserUrls(limit = 10, offset = 0) {
        const response = await fetch(`${this.baseURL}/api/urls?limit=${limit}&offset=${offset}`, {
            headers: this.getAuthHeaders()
        });
        return await response.json();
    }

    // Удаление ссылки
    async deleteUrl(alias) {
        const response = await fetch(`${this.baseURL}/api/urls/${alias}`, {
            method: 'DELETE',
            headers: this.getAuthHeaders()
        });
        return await response.json();
    }
}

const api = new URLShortenerAPI();