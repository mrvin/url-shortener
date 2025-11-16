document.addEventListener('DOMContentLoaded', async function() {
    // Проверка авторизации
    if (!authManager.checkAuthStatus()) {
        window.location.href = 'login.html';
        return;
    }
    
    // Показываем имя пользователя
    document.getElementById('username-display').textContent = authManager.getCurrentUser();
    
    await loadUserUrls();
});

async function loadUserUrls() {
    const loading = document.getElementById('loading');
    const container = document.getElementById('urls-container');
    const emptyMessage = document.getElementById('empty-message');
    const urlsList = document.getElementById('urls-list');
    
    try {
        const result = await api.getUserUrls();
        
        loading.style.display = 'none';
        container.style.display = 'block';
        
        if (result.urls && result.urls.length > 0) {
            emptyMessage.style.display = 'none';
            renderUrlsList(result.urls);
        } else {
            emptyMessage.style.display = 'block';
        }
    } catch (error) {
        loading.style.display = 'none';
        container.style.display = 'block';
        emptyMessage.innerHTML = '❌ Ошибка загрузки данных';
        emptyMessage.style.display = 'block';
    }
}

function renderUrlsList(urls) {
    const urlsList = document.getElementById('urls-list');
    
    const html = urls.map(url => `
        <div class="card mb-3">
            <div class="card-body">
                <div class="row">
                    <div class="col-md-8">
                        <h6 class="card-title">
                            <a href="${url.url}" target="_blank">${url.url}</a>
                        </h6>
                        <p class="card-text">
                            <strong>Короткая ссылка:</strong> 
                            <a href="${API_BASE}/${url.alias}" target="_blank">${API_BASE}/${url.alias}</a>
                        </p>
                        <small class="text-muted">
                            Создано: ${new Date(url.created_at).toLocaleDateString('ru-RU')} | 
                            Переходов: ${url.count}
                        </small>
                    </div>
                    <div class="col-md-4 text-end">
                        <button class="btn btn-sm btn-outline-primary me-2" 
                                onclick="copyUrl('${API_BASE}/${url.alias}')">
                            📋 Копировать
                        </button>
                        <button class="btn btn-sm btn-outline-danger" 
                                onclick="deleteUrl('${url.alias}')">
                            🗑️ Удалить
                        </button>
                    </div>
                </div>
            </div>
        </div>
    `).join('');
    
    urlsList.innerHTML = html;
}

async function deleteUrl(alias) {
    if (!confirm('Вы уверены, что хотите удалить эту ссылку?')) {
        return;
    }
    
    try {
        const result = await api.deleteUrl(alias);
        
        if (result.status === 'OK') {
            await loadUserUrls(); // Перезагружаем список
        } else {
            alert('Ошибка при удалении ссылки');
        }
    } catch (error) {
        alert('Ошибка сети');
    }
}

function copyUrl(url) {
    navigator.clipboard.writeText(url);
    // Можно добавить уведомление о копировании
}