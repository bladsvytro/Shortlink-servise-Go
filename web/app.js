// Конфигурация
const API_BASE = 'http://localhost:8080/api/v1';
const TOKEN_KEY = 'pet_ssl_token';

// Элементы DOM
const urlInput = document.getElementById('url');
const titleInput = document.getElementById('title');
const customInput = document.getElementById('custom');
const shortenBtn = document.getElementById('shorten-btn');
const resultDiv = document.getElementById('result');
const shortUrlInput = document.getElementById('short-url');
const copyBtn = document.getElementById('copy-btn');
const linksList = document.getElementById('links-list');
const refreshBtn = document.getElementById('refresh-stats-btn');
const totalLinksEl = document.getElementById('total-links');
const totalClicksEl = document.getElementById('total-clicks');
const todayClicksEl = document.getElementById('today-clicks');
const notification = document.getElementById('notification');

// Токен аутентификации
let authToken = localStorage.getItem(TOKEN_KEY);

// Автоматическая аутентификация демо-пользователя
async function ensureAuth() {
    if (authToken) {
        // Проверяем, валиден ли токен
        try {
            const response = await fetch(`${API_BASE}/auth/me`, {
                headers: { 'Authorization': `Bearer ${authToken}` }
            });
            if (response.ok) {
                console.log('Токен валиден');
                return authToken;
            }
        } catch (error) {
            console.log('Токен невалиден, получаем новый');
        }
    }
    
    // Регистрируем/логиним демо-пользователя
    try {
        const demoCredentials = {
            email: `demo-${Date.now()}@example.com`,
            password: 'demo123',
            name: 'Demo User'
        };
        
        // Пытаемся зарегистрировать
        const registerResponse = await fetch(`${API_BASE}/auth/register`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(demoCredentials)
        });
        
        let token;
        if (registerResponse.ok) {
            const data = await registerResponse.json();
            token = data.token;
        } else {
            // Если регистрация не удалась, пробуем логин
            const loginResponse = await fetch(`${API_BASE}/auth/login`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({
                    email: 'demo@example.com',
                    password: 'demo123'
                })
            });
            
            if (loginResponse.ok) {
                const data = await loginResponse.json();
                token = data.token;
            } else {
                throw new Error('Не удалось аутентифицироваться');
            }
        }
        
        // Сохраняем токен
        authToken = token;
        localStorage.setItem(TOKEN_KEY, token);
        console.log('Аутентификация успешна');
        return token;
    } catch (error) {
        console.error('Ошибка аутентификации:', error);
        // Используем демо-режим
        return 'demo-token';
    }
}

// Уведомления
function showNotification(message, type = 'success') {
    notification.textContent = message;
    notification.className = 'notification';
    notification.classList.add('show');
    notification.style.background = type === 'success' ? '#10b981' : 
                                   type === 'error' ? '#ef4444' : '#3b82f6';
    
    setTimeout(() => {
        notification.classList.remove('show');
    }, 3000);
}

// Проверка URL
function isValidUrl(string) {
    try {
        new URL(string);
        return true;
    } catch (_) {
        return false;
    }
}

// Создание короткой ссылки
async function shortenUrl() {
    const url = urlInput.value.trim();
    const title = titleInput.value.trim();
    const customCode = customInput.value.trim();
    
    if (!url) {
        showNotification('Введите URL', 'error');
        return;
    }
    
    if (!isValidUrl(url)) {
        showNotification('Введите корректный URL', 'error');
        return;
    }
    
    shortenBtn.disabled = true;
    shortenBtn.innerHTML = '<i class="fas fa-spinner fa-spin"></i> Обработка...';
    
    try {
        // Убеждаемся, что у нас есть валидный токен
        const token = await ensureAuth();
        
        const payload = { url, title: title || undefined };
        if (customCode) payload.custom_code = customCode;
        
        const response = await fetch(`${API_BASE}/links`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            },
            body: JSON.stringify(payload)
        });
        
        if (!response.ok) {
            const error = await response.text();
            throw new Error(error || 'Ошибка сервера');
        }
        
        const data = await response.json();
        
        // Показываем результат
        shortUrlInput.value = data.short_url;
        resultDiv.style.display = 'block';
        
        // Прокручиваем к результату
        resultDiv.scrollIntoView({ behavior: 'smooth' });
        
        showNotification('Ссылка успешно создана!');
        
        // Обновляем список ссылок и статистику
        loadUserLinks();
        loadStats();
        
        // Очищаем поля
        urlInput.value = '';
        titleInput.value = '';
        customInput.value = '';
        
    } catch (error) {
        console.error('Error shortening URL:', error);
        showNotification(`Ошибка: ${error.message}`, 'error');
    } finally {
        shortenBtn.disabled = false;
        shortenBtn.innerHTML = '<i class="fas fa-magic"></i> Сократить';
    }
}

// Копирование ссылки в буфер обмена
function copyToClipboard() {
    shortUrlInput.select();
    shortUrlInput.setSelectionRange(0, 99999); // Для мобильных устройств
    
    try {
        navigator.clipboard.writeText(shortUrlInput.value)
            .then(() => {
                showNotification('Ссылка скопирована в буфер обмена!');
            })
            .catch(err => {
                // Fallback для старых браузеров
                document.execCommand('copy');
                showNotification('Ссылка скопирована!');
            });
    } catch (err) {
        document.execCommand('copy');
        showNotification('Ссылка скопирована!');
    }
}

// Загрузка ссылок пользователя
async function loadUserLinks() {
    try {
        const token = await ensureAuth();
        const response = await fetch(`${API_BASE}/links`, {
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });
        
        if (!response.ok) {
            // Если ошибка авторизации, используем демо-данные
            if (response.status === 401) {
                renderDemoLinks();
                return;
            }
            throw new Error('Failed to load links');
        }
        
        const links = await response.json();
        renderLinks(links);
    } catch (error) {
        console.error('Error loading links:', error);
        renderDemoLinks();
    }
}

// Рендеринг списка ссылок
function renderLinks(links) {
    if (!links || links.length === 0) {
        linksList.innerHTML = `
            <div class="empty-state">
                <i class="fas fa-history"></i>
                <p>Здесь появятся ваши сокращённые ссылки</p>
            </div>
        `;
        return;
    }
    
    const html = links.slice(0, 5).map(link => `
        <div class="link-item">
            <div class="link-info">
                <h4>${link.title || 'Без названия'}</h4>
                <p>${link.original_url.substring(0, 50)}...</p>
                <small>Создано: ${new Date(link.created_at).toLocaleDateString('ru-RU')}</small>
            </div>
            <div class="link-stats">
                <div class="clicks">${link.click_count || 0} кликов</div>
                <a href="${link.short_url}" target="_blank" class="btn-secondary" style="padding: 5px 10px; font-size: 0.8rem; margin-top: 5px;">
                    <i class="fas fa-external-link-alt"></i> Открыть
                </a>
            </div>
        </div>
    `).join('');
    
    linksList.innerHTML = html;
}

// Демо-данные для отображения
function renderDemoLinks() {
    const demoLinks = [
        {
            title: 'Пример ссылки',
            original_url: 'https://example.com/very-long-url-path',
            short_url: 'http://localhost:8080/demo1',
            created_at: new Date().toISOString(),
            click_count: 42
        },
        {
            title: 'GitHub',
            original_url: 'https://github.com',
            short_url: 'http://localhost:8080/demo2',
            created_at: new Date(Date.now() - 86400000).toISOString(),
            click_count: 18
        }
    ];
    
    renderLinks(demoLinks);
}

// Загрузка статистики
async function loadStats() {
    try {
        const token = await ensureAuth();
        const response = await fetch(`${API_BASE}/stats`, {
            headers: { 'Authorization': `Bearer ${token}` }
        });
        
        if (!response.ok) {
            throw new Error(`HTTP ${response.status}: ${response.statusText}`);
        }
        
        const data = await response.json();
        
        // Обновляем DOM с реальными данными
        totalLinksEl.textContent = data.total_links;
        totalClicksEl.textContent = data.total_clicks;
        todayClicksEl.textContent = data.today_clicks;
        
        // Добавляем анимацию обновления
        totalLinksEl.classList.add('updated');
        totalClicksEl.classList.add('updated');
        todayClicksEl.classList.add('updated');
        setTimeout(() => {
            totalLinksEl.classList.remove('updated');
            totalClicksEl.classList.remove('updated');
            todayClicksEl.classList.remove('updated');
        }, 500);
        
    } catch (error) {
        console.error('Error loading stats:', error);
        // Демо-статистика в случае ошибки
        totalLinksEl.textContent = '0';
        totalClicksEl.textContent = '0';
        todayClicksEl.textContent = '0';
        showNotification('Не удалось загрузить статистику. Используются демо-данные.', 'warning');
    }
}

// Инициализация
async function init() {
    // Показываем уведомление о загрузке
    showNotification('Загрузка...', 'info');
    
    try {
        // Аутентификация
        await ensureAuth();
        
        // Загрузка начальных данных
        await loadUserLinks();
        await loadStats();
        
        showNotification('Готово к работе!');
    } catch (error) {
        console.error('Ошибка инициализации:', error);
        showNotification('Демо-режим: данные могут быть ограничены', 'info');
    }
    
    // Обработчики событий
    shortenBtn.addEventListener('click', shortenUrl);
    copyBtn.addEventListener('click', copyToClipboard);
    refreshBtn.addEventListener('click', () => {
        loadUserLinks();
        loadStats();
        showNotification('Данные обновлены');
    });
    
    // Отправка формы по Enter
    urlInput.addEventListener('keypress', (e) => {
        if (e.key === 'Enter') shortenUrl();
    });
}

// Запуск приложения
document.addEventListener('DOMContentLoaded', init);