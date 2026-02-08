// Глобальные переменные для пагинации
let currentPage = 1;
const itemsPerPage = 10; // Фиксированное значение
let totalItems = 0;

document.addEventListener('DOMContentLoaded', async function() {
    // Проверка авторизации
    if (!authManager.checkAuthStatus()) {
        window.location.href = 'login.html';
        return;
    }
    
    // Показываем имя пользователя
    document.getElementById('username-display').textContent = authManager.getCurrentUser();
    
    // Настройка обработчиков
    setupEventListeners();
    
    // Загрузка данных
    await loadUserUrls();
});

function setupEventListeners() {
    // Кнопка "Назад"
    document.getElementById('prev-page').addEventListener('click', function(e) {
        e.preventDefault();
        if (currentPage > 1) {
            currentPage--;
            loadUserUrls();
        }
    });
    
    // Кнопка "Вперед"
    document.getElementById('next-page').addEventListener('click', function(e) {
        e.preventDefault();
        const totalPages = Math.ceil(totalItems / itemsPerPage);
        if (currentPage < totalPages) {
            currentPage++;
            loadUserUrls();
        }
    });
}

async function loadUserUrls() {
    const loading = document.getElementById('loading');
    const container = document.getElementById('urls-container');
    const emptyMessage = document.getElementById('empty-message');
    const paginationContainer = document.getElementById('pagination-container');
    
    // Показываем загрузку
    loading.style.display = 'block';
    container.style.display = 'none';
    paginationContainer.style.display = 'none';
    
    try {
        const offset = (currentPage - 1) * itemsPerPage;
        const result = await api.getUserUrls(itemsPerPage, offset);
        
        loading.style.display = 'none';
        container.style.display = 'block';
        
        // Сохраняем общее количество
        totalItems = result.total || 0;
        
        // Обновляем статистику
        updateStats(result.urls ? result.urls.length : 0);
        
        if (result.urls && result.urls.length > 0) {
            emptyMessage.style.display = 'none';
            renderUrlsList(result.urls);
            renderPagination();
        } else {
            emptyMessage.style.display = 'block';
            document.getElementById('urls-list').innerHTML = '';
            paginationContainer.style.display = 'none';
        }
    } catch (error) {
        loading.style.display = 'none';
        container.style.display = 'block';
        emptyMessage.innerHTML = '❌ Ошибка загрузки данных';
        emptyMessage.style.display = 'block';
        console.error('Error loading URLs:', error);
    }
}

function updateStats(showingCount) {
    document.getElementById('total-count').textContent = totalItems;
    document.getElementById('total-count-2').textContent = totalItems;
    document.getElementById('showing-count').textContent = showingCount;
}

function renderUrlsList(urls) {
    const urlsList = document.getElementById('urls-list');
    
    const html = urls.map((url, index) => {
        const globalIndex = (currentPage - 1) * itemsPerPage + index + 1;
        return `
        <div class="card mb-3">
            <div class="card-body">
                <div class="row">
                    <div class="col-md-8">
                        <div class="d-flex align-items-center mb-2">
                            <span class="badge bg-secondary me-2">${globalIndex}</span>
                            <h6 class="card-title mb-0">
                                <a href="${url.url}" target="_blank" class="text-truncate d-inline-block" style="max-width: 400px;">
                                    ${url.url}
                                </a>
                            </h6>
                        </div>
                        <p class="card-text mb-1">
                            <strong>Короткая ссылка:</strong> 
                            <a href="${API_BASE}/${url.alias}" target="_blank">${API_BASE}/${url.alias}</a>
                        </p>
                        <small class="text-muted">
                            Создано: ${new Date(url.created_at).toLocaleDateString('ru-RU')} | 
                            Переходов: <span class="badge bg-info">${url.count}</span>
                        </small>
                    </div>
                    <div class="col-md-4 text-end">
                        <button class="btn btn-sm btn-outline-primary me-2" 
                                onclick="copyUrl('${API_BASE}/${url.alias}')"
                                title="Копировать ссылку">
                            📋 Копировать
                        </button>
                        <button class="btn btn-sm btn-outline-danger" 
                                onclick="deleteUrl('${url.alias}')"
                                title="Удалить ссылку">
                            🗑️ Удалить
                        </button>
                    </div>
                </div>
            </div>
        </div>
        `;
    }).join('');
    
    urlsList.innerHTML = html;
}

function renderPagination() {
    const paginationContainer = document.getElementById('pagination-container');
    const totalPages = Math.ceil(totalItems / itemsPerPage);
    
    if (totalPages <= 1) {
        paginationContainer.style.display = 'none';
        return;
    }
    
    paginationContainer.style.display = 'block';
    
    // Обновляем информацию о страницах
    document.getElementById('page-info').textContent = currentPage;
    document.getElementById('total-pages').textContent = totalPages;
    document.getElementById('current-page').textContent = currentPage;
    
    // Обновляем состояние кнопок
    const prevButton = document.getElementById('prev-page');
    const nextButton = document.getElementById('next-page');
    
    if (currentPage === 1) {
        prevButton.classList.add('disabled');
    } else {
        prevButton.classList.remove('disabled');
    }
    
    if (currentPage === totalPages) {
        nextButton.classList.add('disabled');
    } else {
        nextButton.classList.remove('disabled');
    }
}

async function deleteUrl(alias) {
    if (!confirm('Вы уверены, что хотите удалить эту ссылку?')) {
        return;
    }
    
    try {
        const result = await api.deleteUrl(alias);
        
        if (result.status === 'OK') {
            // Если на текущей странице осталась только одна ссылка и это не первая страница
            const currentItems = document.querySelectorAll('#urls-list .card').length;
            if (currentItems === 1 && currentPage > 1) {
                currentPage--; // Переходим на предыдущую страницу
            }
            
            await loadUserUrls(); // Перезагружаем список
        } else {
            alert('Ошибка при удалении ссылки');
        }
    } catch (error) {
        alert('Ошибка сети');
    }
}

function copyUrl(url) {
    navigator.clipboard.writeText(url).then(() => {
        // Показываем временное уведомление
        showTempAlert('✅ Ссылка скопирована в буфер обмена!', 'success');
    }).catch(() => {
        showTempAlert('❌ Ошибка копирования', 'danger');
    });
}

function showTempAlert(message, type) {
    const alertDiv = document.createElement('div');
    alertDiv.className = `alert alert-${type} alert-dismissible fade show position-fixed`;
    alertDiv.style.top = '20px';
    alertDiv.style.right = '20px';
    alertDiv.style.zIndex = '1050';
    alertDiv.innerHTML = `
        ${message}
        <button type="button" class="btn-close" data-bs-dismiss="alert"></button>
    `;
    
    document.body.appendChild(alertDiv);
    
    // Автоматически скрываем через 3 секунды
    setTimeout(() => {
        if (alertDiv.parentNode) {
            alertDiv.remove();
        }
    }, 3000);
}
