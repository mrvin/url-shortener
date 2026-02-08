document.addEventListener('DOMContentLoaded', function() {
    const shortenForm = document.getElementById('shorten-form');
    const aliasInput = document.getElementById('alias-input');
    const aliasStatus = document.getElementById('alias-status');
    
    // Проверка доступности алиаса в реальном времени
    aliasInput.addEventListener('input', async function() {
        const alias = this.value.trim();
        
        if (alias.length === 0) {
            aliasStatus.textContent = '';
            return;
        }
        
        // Валидация паттерна
        const pattern = /^[a-zA-Z0-9_-]+$/;
        if (!pattern.test(alias)) {
            aliasStatus.innerHTML = '<span class="text-danger">❌ Используйте только буквы, цифры, - и _</span>';
            return;
        }
        
        try {
            const result = await api.checkAlias(alias);
            if (result.exists) {
                aliasStatus.innerHTML = '<span class="text-danger">❌ Этот алиас уже занят</span>';
            } else {
                aliasStatus.innerHTML = '<span class="text-success">✅ Алиас доступен</span>';
            }
        } catch (error) {
            aliasStatus.innerHTML = '<span class="text-warning">⚠️ Ошибка проверки</span>';
        }
    });
    
    // Обработка формы сокращения
    shortenForm.addEventListener('submit', async function(e) {
        e.preventDefault();
        
        const url = document.getElementById('url-input').value;
        const alias = aliasInput.value.trim() || undefined;
        const shortenBtn = document.getElementById('shorten-btn');
        
        // Проверка авторизации
        if (!authManager.checkAuthStatus()) {
            showResult('Для создания ссылки необходимо <a href="login.html">войти в систему</a>.', 'warning');
            return;
        }
        
        shortenBtn.disabled = true;
        shortenBtn.textContent = 'Создание...';
        
        try {
            const result = await api.shortenUrl({ url, alias });
            
            if (result.status === 'OK') {
                const shortUrl = `${API_BASE}/${alias}`;
                showResult(`
                    <h5>✅ Ссылка создана!</h5>
                    <div class="mt-2">
                        <strong>Короткая ссылка:</strong><br>
                        <a href="${shortUrl}" target="_blank">${shortUrl}</a>
                    </div>
                    <button class="btn btn-sm btn-outline-secondary mt-2" onclick="copyToClipboard('${shortUrl}')">
                        📋 Скопировать
                    </button>
                `, 'success');
                
                shortenForm.reset();
                aliasStatus.textContent = '';
            } else {
                showResult('❌ Ошибка при создании ссылки', 'danger');
            }
        } catch (error) {
            showResult('❌ Ошибка сети', 'danger');
        } finally {
            shortenBtn.disabled = false;
            shortenBtn.textContent = 'Сократить ссылку';
        }
    });
});

// Функция показа результата
function showResult(message, type = 'info') {
    const resultDiv = document.getElementById('result');
    resultDiv.innerHTML = `
        <div class="alert alert-${type} alert-dismissible fade show">
            ${message}
            <button type="button" class="btn-close" data-bs-dismiss="alert"></button>
        </div>
    `;
}

// Функция копирования в буфер обмена
function copyToClipboard(text) {
    navigator.clipboard.writeText(text).then(() => {
        const alert = document.createElement('div');
        alert.className = 'alert alert-success alert-dismissible fade show';
        alert.innerHTML = '✅ Скопировано в буфер обмена!';
        document.getElementById('result').appendChild(alert);
        
        setTimeout(() => alert.remove(), 3000);
    });
}

class AliasGenerator {
    static generate(length = 6) {
        // Используем Base62 подход как самый популярный
        const chars = '0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ';
        let result = '';
        
        for (let i = 0; i < length; i++) {
            result += chars.charAt(Math.floor(Math.random() * chars.length));
        }
        
        return result;
    }
    
    // Альтернативный метод с проверкой на лету
    static async generateUnique(length = 6, maxAttempts = 5) {
        for (let attempt = 0; attempt < maxAttempts; attempt++) {
            const alias = this.generate(length);
            
            try {
                const result = await api.checkAlias(alias);
                if (!result.exists) {
                    return alias; // Нашли уникальный алиас
                }
            } catch (error) {
                console.warn('Ошибка проверки алиаса:', error);
            }
        }
        
        // Если не нашли за maxAttempts, генерируем с увеличенной длиной
        return this.generate(length + 2);
    }
}

// Функция для генерации случайного алиаса
async function generateRandomAlias() {
    const aliasInput = document.getElementById('alias-input');
    const aliasStatus = document.getElementById('alias-status');
    
    aliasInput.disabled = true;
    aliasStatus.innerHTML = '<span class="text-info">⏳ Генерация...</span>';
    
    try {
        // Генерируем уникальный алиас
        const alias = await AliasGenerator.generateUnique(6);
        aliasInput.value = alias;
        aliasStatus.innerHTML = '<span class="text-success">✅ Алиас сгенерирован и доступен</span>';
    } catch (error) {
        // Если API недоступно, генерируем локально
        const alias = AliasGenerator.generate(6);
        aliasInput.value = alias;
        aliasStatus.innerHTML = '<span class="text-warning">⚠️ Алиас сгенерирован (требуется проверка)</span>';
    } finally {
        aliasInput.disabled = false;
    }
}

// Автогенерация при загрузке, если поле пустое
document.addEventListener('DOMContentLoaded', function() {
    const aliasInput = document.getElementById('alias-input');
    
    // Если пользователь начинает вводить - не генерируем автоматически
    aliasInput.addEventListener('focus', function() {
        if (!this.value.trim()) {
            generateRandomAlias();
        }
    });
});