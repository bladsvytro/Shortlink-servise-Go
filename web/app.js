// Конфигурация
const API_BASE = 'http://localhost:8080/api/v1';
const TOKEN_KEY = 'pet_ssl_token';

// Локализация
const translations = {
    en: {
        // Заголовки
        'title': 'URL Shortener Service',
        'tagline': 'Fast and secure URL shortening',
        'shorten_header': 'Shorten URL',
        'stats_header': 'Statistics',
        'admin_panel': 'Admin Panel',
        'sign_in': 'Sign In',
        'sign_up': 'Sign Up',
        'logout': 'Logout',
        // Форма
        'enter_url': 'Enter URL',
        'title_label': 'Title (optional)',
        'custom_label': 'Custom code (optional)',
        'preview': 'Preview',
        'shorten_button': 'Shorten',
        'link_created': 'Link created!',
        'copy': 'Copy',
        'link_ready': 'Link is ready to use',
        'total_links': 'Total Links',
        'total_clicks': 'Total Clicks',
        'today_clicks': 'Today\'s Clicks',
        'stats_update': 'Statistics update in real time',
        'refresh_stats': 'Refresh Statistics',
        // Уведомления
        'loading': 'Loading...',
        'ready': 'Ready to work!',
        'demo_mode': 'Demo mode: data may be limited',
        'data_updated': 'Data updated',
        'enter_url_error': 'Enter URL',
        'valid_url_error': 'Enter a valid URL',
        'processing': 'Processing...',
        'link_created_success': 'Link created successfully!',
        'copy_success': 'Link copied to clipboard!',
        'sign_in_success': 'Sign in successful',
        'registration_success': 'Registration successful',
        'logout_success': 'You have logged out',
        'auth_required': 'Authorization required',
        // Админ-панель
        'users_tab': 'Users',
        'links_tab': 'Links',
        'stats_tab': 'Statistics',
        'user_management': 'User Management',
        'all_links': 'All Links',
        'system_stats': 'System Statistics',
        'refresh': 'Refresh',
        'export_csv': 'Export CSV',
        'export_report': 'Export Report',
        'loading_text': 'Loading...',
        'no_users': 'No users found',
        'no_links': 'No links found',
        'view': 'View',
        'delete': 'Delete',
        'admin': 'Admin',
        'user': 'User',
        'page': 'Page',
        'of': 'of',
        'previous': 'Previous',
        'next': 'Next',
        'total_users': 'Total Users',
        // Футер
        'footer_text': 'Pet SSL © 2026 | Shorten URLs fast and free',
        'made_by': 'Made by',
        // Формы аутентификации
        'name': 'Name',
        'email_label': 'Email',
        'username_label': 'Username',
        'username_or_email_label': 'Username or Email',
        'password': 'Password',
        'username_optional': 'Username (optional)',
    },
    ru: {
        'title': 'Сервис Сокращения Ссылок',
        'tagline': 'Быстрое и безопасное сокращение ссылок',
        'shorten_header': 'Сократить ссылку',
        'stats_header': 'Статистика',
        'admin_panel': 'Админ-панель',
        'sign_in': 'Вход',
        'sign_up': 'Регистрация',
        'logout': 'Выйти',
        'enter_url': 'Введите URL',
        'title_label': 'Название (необязательно)',
        'custom_label': 'Пользовательский код (необязательно)',
        'preview': 'Превью',
        'shorten_button': 'Сократить',
        'link_created': 'Ссылка создана!',
        'copy': 'Копировать',
        'link_ready': 'Ссылка готова к использованию',
        'total_links': 'Всего ссылок',
        'total_clicks': 'Всего кликов',
        'today_clicks': 'Кликов сегодня',
        'stats_update': 'Статистика обновляется в реальном времени',
        'refresh_stats': 'Обновить статистику',
        'loading': 'Загрузка...',
        'ready': 'Готово к работе!',
        'demo_mode': 'Демо-режим: данные могут быть ограничены',
        'data_updated': 'Данные обновлены',
        'enter_url_error': 'Введите URL',
        'valid_url_error': 'Введите корректный URL',
        'processing': 'Обработка...',
        'link_created_success': 'Ссылка успешно создана!',
        'copy_success': 'Ссылка скопирована в буфер обмена!',
        'sign_in_success': 'Вход выполнен успешно',
        'registration_success': 'Регистрация выполнена успешно',
        'logout_success': 'Вы вышли из системы',
        'auth_required': 'Требуется авторизация',
        'users_tab': 'Пользователи',
        'links_tab': 'Ссылки',
        'stats_tab': 'Статистика',
        'user_management': 'Управление пользователями',
        'all_links': 'Все ссылки',
        'system_stats': 'Системная статистика',
        'refresh': 'Обновить',
        'export_csv': 'Экспорт CSV',
        'export_report': 'Экспорт отчёта',
        'loading_text': 'Загрузка...',
        'no_users': 'Пользователи не найдены',
        'no_links': 'Ссылки не найдены',
        'view': 'Просмотр',
        'delete': 'Удалить',
        'admin': 'Админ',
        'user': 'Пользователь',
        'page': 'Страница',
        'of': 'из',
        'previous': 'Назад',
        'next': 'Вперёд',
        'total_users': 'Всего пользователей',
        'footer_text': 'Pet SSL © 2026 | Сокращайте ссылки быстро и бесплатно',
        'made_by': 'Сделано',
        // Формы аутентификации
        'name': 'Имя',
        'email_label': 'Email',
        'username_label': 'Логин',
        'username_or_email_label': 'Логин или Email',
        'password': 'Пароль',
        'username_optional': 'Логин (необязательно)',
    }
};

let currentLang = localStorage.getItem('lang') || 'en';

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
const previewUrl = document.getElementById('preview-url');
const previewText = document.getElementById('preview-text');
const previewCode = document.getElementById('preview-code');

// Новые элементы для аутентификации и админки
const userInfo = document.getElementById('user-info');
const userEmail = document.getElementById('user-email');
const logoutBtn = document.getElementById('logout-btn');
const adminPanelBtn = document.getElementById('admin-panel-btn');
const authButtons = document.getElementById('auth-buttons');
const loginBtn = document.getElementById('login-btn');
const registerBtn = document.getElementById('register-btn');
const loginModal = document.getElementById('login-modal');
const registerModal = document.getElementById('register-modal');
const loginForm = document.getElementById('login-form');
const registerForm = document.getElementById('register-form');
const closeButtons = document.querySelectorAll('.close');

// Элементы админ-панели
const adminModal = document.getElementById('admin-modal');
const adminTabs = document.querySelectorAll('.admin-tab');
const adminTabPanes = document.querySelectorAll('.admin-tab-pane');
const adminUsersTable = document.getElementById('admin-users-table');
const adminLinksTable = document.getElementById('admin-links-table');
const adminTotalUsersEl = document.getElementById('admin-total-users');
const adminTotalLinksEl = document.getElementById('admin-total-links');
const adminTotalClicksEl = document.getElementById('admin-total-clicks');
const adminTodayClicksEl = document.getElementById('admin-today-clicks');
const refreshUsersBtn = document.getElementById('refresh-users-btn');
const refreshLinksBtn = document.getElementById('refresh-links-btn');
const refreshAdminStatsBtn = document.getElementById('refresh-stats-btn');
const exportUsersBtn = document.getElementById('export-users-btn');
const exportLinksBtn = document.getElementById('export-links-btn');
const exportStatsBtn = document.getElementById('export-stats-btn');
const prevUsersBtn = document.getElementById('prev-users-btn');
const nextUsersBtn = document.getElementById('next-users-btn');
const prevLinksBtn = document.getElementById('prev-links-btn');
const nextLinksBtn = document.getElementById('next-links-btn');
const usersPageInfo = document.getElementById('users-page-info');
const linksPageInfo = document.getElementById('links-page-info');

// Токен аутентификации
let authToken = localStorage.getItem(TOKEN_KEY);

// Пагинация админ-панели
let adminUsersPage = 1;
let adminLinksPage = 1;
const itemsPerPage = 10;

// Локализация
function applyLanguage(lang) {
    currentLang = lang;
    localStorage.setItem('lang', lang);
    
    const t = translations[lang];
    
    // Обновляем заголовок страницы
    document.title = t.title;
    
    // Обновляем элементы с data-i18n (пока не используем)
    // Вместо этого обновляем конкретные элементы по ID
    
    // Заголовок и теглайн
    const h1 = document.querySelector('h1');
    if (h1) h1.textContent = t.title;
    const tagline = document.querySelector('.tagline');
    if (tagline) tagline.textContent = t.tagline;
    
    // Кнопки навигации
    const loginBtn = document.getElementById('login-btn');
    if (loginBtn) loginBtn.textContent = t.sign_in;
    const registerBtn = document.getElementById('register-btn');
    if (registerBtn) registerBtn.textContent = t.sign_up;
    const logoutBtn = document.getElementById('logout-btn');
    if (logoutBtn) logoutBtn.textContent = t.logout;
    const adminPanelBtn = document.getElementById('admin-panel-btn');
    if (adminPanelBtn) adminPanelBtn.textContent = t.admin_panel;
    
    // Форма сокращения
    const shortenHeader = document.querySelector('.main-form h2');
    if (shortenHeader) shortenHeader.innerHTML = `<i class="fas fa-compress-alt"></i> ${t.shorten_header}`;
    const urlLabel = document.querySelector('label[for="url"]');
    if (urlLabel) urlLabel.textContent = t.enter_url;
    const titleLabel = document.querySelector('label[for="title"]');
    if (titleLabel) titleLabel.textContent = t.title_label;
    const customLabel = document.querySelector('label[for="custom"]');
    if (customLabel) customLabel.textContent = t.custom_label;
    const previewStrong = document.querySelector('.preview-url strong');
    if (previewStrong) previewStrong.textContent = t.preview + ':';
    const shortenButton = document.getElementById('shorten-btn');
    if (shortenButton) shortenButton.innerHTML = `<i class="fas fa-magic"></i> ${t.shorten_button}`;
    const linkCreated = document.querySelector('.result h3');
    if (linkCreated) linkCreated.innerHTML = `<i class="fas fa-check-circle"></i> ${t.link_created}`;
    const copyButton = document.getElementById('copy-btn');
    if (copyButton) copyButton.innerHTML = `<i class="far fa-copy"></i> ${t.copy}`;
    const linkReady = document.querySelector('.stats');
    if (linkReady) linkReady.textContent = t.link_ready;
    
    // Статистика
    const statsHeader = document.querySelector('.stats-card h2');
    if (statsHeader) statsHeader.innerHTML = `<i class="fas fa-chart-line"></i> ${t.stats_header}`;
    const totalLinksLabel = document.querySelector('.stat-label:nth-of-type(1)');
    if (totalLinksLabel) totalLinksLabel.textContent = t.total_links;
    const totalClicksLabel = document.querySelector('.stat-label:nth-of-type(2)');
    if (totalClicksLabel) totalClicksLabel.textContent = t.total_clicks;
    const todayClicksLabel = document.querySelector('.stat-label:nth-of-type(3)');
    if (todayClicksLabel) todayClicksLabel.textContent = t.today_clicks;
    const statsFooter = document.querySelector('.stats-footer p');
    if (statsFooter) statsFooter.textContent = t.stats_update;
    const refreshStatsBtn = document.getElementById('refresh-stats-btn');
    if (refreshStatsBtn) refreshStatsBtn.innerHTML = `<i class="fas fa-sync-alt"></i> ${t.refresh_stats}`;
    
    // Футер
    const footerText = document.querySelector('footer p:first-of-type');
    if (footerText) footerText.textContent = t.footer_text;
    const madeBy = document.querySelector('.author');
    if (madeBy) madeBy.innerHTML = `${t.made_by} <a href="https://github.com/bladsvytro" target="_blank">Vladislav Naumov</a>`;
    
    // Модальные окна (вход/регистрация)
    const loginModalTitle = document.querySelector('#login-modal h2');
    if (loginModalTitle) loginModalTitle.innerHTML = `<i class="fas fa-sign-in-alt"></i> ${t.sign_in}`;
    const loginLabel = document.querySelector('label[for="login-email"]');
    if (loginLabel) loginLabel.textContent = t.username_or_email_label + ':';
    const passwordLabel = document.querySelector('label[for="login-password"]');
    if (passwordLabel) passwordLabel.textContent = t.password + ':';
    const loginSubmit = document.querySelector('#login-form button');
    if (loginSubmit) loginSubmit.textContent = t.sign_in;

    const registerModalTitle = document.querySelector('#register-modal h2');
    if (registerModalTitle) registerModalTitle.innerHTML = `<i class="fas fa-user-plus"></i> ${t.sign_up}`;
    const nameLabel = document.querySelector('label[for="register-name"]');
    if (nameLabel) nameLabel.textContent = t.name + ':';
    const registerEmailLabel = document.querySelector('label[for="register-email"]');
    if (registerEmailLabel) registerEmailLabel.textContent = t.email_label + ':';
    const registerUsernameLabel = document.querySelector('label[for="register-username"]');
    if (registerUsernameLabel) registerUsernameLabel.textContent = t.username_optional + ':';
    const registerPasswordLabel = document.querySelector('label[for="register-password"]');
    if (registerPasswordLabel) registerPasswordLabel.textContent = t.password + ':';
    const registerSubmit = document.querySelector('#register-form button');
    if (registerSubmit) registerSubmit.textContent = t.sign_up;
    
    // Админ-панель (частично)
    const adminModalTitle = document.querySelector('#admin-modal h2');
    if (adminModalTitle) adminModalTitle.innerHTML = `<i class="fas fa-cogs"></i> ${t.admin_panel}`;
    const userTab = document.querySelector('.admin-tab[data-tab="users"]');
    if (userTab) userTab.textContent = t.users_tab;
    const linksTab = document.querySelector('.admin-tab[data-tab="links"]');
    if (linksTab) linksTab.textContent = t.links_tab;
    const statsTab = document.querySelector('.admin-tab[data-tab="stats"]');
    if (statsTab) statsTab.textContent = t.stats_tab;
    
    // Обновляем активность кнопок переключателя языка
    document.querySelectorAll('.lang-btn').forEach(btn => {
        btn.classList.toggle('active', btn.dataset.lang === lang);
    });
    
    // Обновляем атрибут lang у html
    document.documentElement.lang = lang;
}

function initLanguageSwitcher() {
    const ruBtn = document.getElementById('lang-ru');
    const enBtn = document.getElementById('lang-en');
    
    if (ruBtn && enBtn) {
        ruBtn.addEventListener('click', () => {
            applyLanguage('ru');
        });
        enBtn.addEventListener('click', () => {
            applyLanguage('en');
        });
    }
}

// Автоматическая аутентификация демо-пользователя
async function ensureAuth() {
    if (authToken) {
        // Проверяем, валиден ли токен
        try {
            const response = await fetch(`${API_BASE}/auth/me`, {
                headers: { 'Authorization': `Bearer ${authToken}` }
            });
            if (response.ok) {
                console.log('Token is valid');
                return authToken;
            }
        } catch (error) {
            console.log('Token invalid, obtaining new');
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
                throw new Error('Failed to authenticate');
            }
        }
        
        // Сохраняем токен
        authToken = token;
        localStorage.setItem(TOKEN_KEY, token);
        console.log('Authentication successful');
        return token;
    } catch (error) {
        console.error('Authentication error:', error);
        // Используем демо-режим
        return 'demo-token';
    }
}

// Функции для управления аутентификацией и UI
async function updateUserInfo() {
    try {
        const response = await fetch(`${API_BASE}/auth/me`, {
            headers: { 'Authorization': `Bearer ${authToken}` }
        });
        if (response.ok) {
            const user = await response.json();
            // Проверяем, является ли пользователь демо-пользователем
            const isDemoUser = user.email.includes('demo-') || user.email.endsWith('@example.com');
            if (isDemoUser) {
                // Скрываем информацию о демо-пользователе, показываем кнопки аутентификации
                userInfo.style.display = 'none';
                authButtons.style.display = 'block';
                adminPanelBtn.style.display = 'none';
                return;
            }
            userEmail.textContent = user.email;
            userInfo.style.display = 'block';
            authButtons.style.display = 'none';
            // Проверяем, является ли пользователь админом
            if (user.is_admin) {
                adminPanelBtn.style.display = 'inline-block';
            } else {
                adminPanelBtn.style.display = 'none';
            }
        } else {
            throw new Error('Failed to get user info');
        }
    } catch (error) {
        console.error('Error getting user info:', error);
        // Скрываем информацию о пользователе, показываем кнопки аутентификации
        userInfo.style.display = 'none';
        authButtons.style.display = 'block';
    }
}

function showModal(modal) {
    modal.style.display = 'block';
}

function hideModal(modal) {
    modal.style.display = 'none';
}

// Обработка входа
async function handleLogin(event) {
    event.preventDefault();
    const email = document.getElementById('login-email').value;
    const password = document.getElementById('login-password').value;

    try {
        const response = await fetch(`${API_BASE}/auth/login`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ email, password })
        });
        if (!response.ok) {
            const error = await response.text();
            throw new Error(error || 'Login error');
        }
        const data = await response.json();
        authToken = data.token;
        localStorage.setItem(TOKEN_KEY, data.token);
        hideModal(loginModal);
        showNotification('Sign in successful');
        await updateUserInfo();
        // Обновляем статистику и ссылки
        loadStats();
        loadUserLinks();
    } catch (error) {
        console.error('Ошибка входа:', error);
        showNotification(`Sign in error: ${error.message}`, 'error');
    }
}

// Обработка регистрации
async function handleRegister(event) {
    event.preventDefault();
    const name = document.getElementById('register-name').value;
    const email = document.getElementById('register-email').value;
    const username = document.getElementById('register-username').value;
    const password = document.getElementById('register-password').value;

    try {
        const response = await fetch(`${API_BASE}/auth/register`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ name, email, username, password })
        });
        if (!response.ok) {
            const error = await response.text();
            throw new Error(error || 'Registration error');
        }
        const data = await response.json();
        authToken = data.token;
        localStorage.setItem(TOKEN_KEY, data.token);
        hideModal(registerModal);
        showNotification('Registration successful');
        await updateUserInfo();
        loadStats();
        loadUserLinks();
    } catch (error) {
        console.error('Ошибка регистрации:', error);
        showNotification(`Registration error: ${error.message}`, 'error');
    }
}

// Выход
function handleLogout() {
    authToken = null;
    localStorage.removeItem(TOKEN_KEY);
    userInfo.style.display = 'none';
    authButtons.style.display = 'block';
    showNotification('You have logged out');
    // Очищаем статистику и ссылки
    totalLinksEl.textContent = '0';
    totalClicksEl.textContent = '0';
    todayClicksEl.textContent = '0';
    if (linksList) linksList.innerHTML = '';
}

// Переход в админ-панель
function handleAdminPanel() {
    // Проверяем, является ли пользователь админом
    if (!authToken) {
        showNotification('Authorization required', 'error');
        return;
    }
    // Открываем модальное окно
    showModal(adminModal);
    // Загружаем данные для активной вкладки
    loadAdminUsers();
    loadAdminStats();
}

// Загрузка пользователей для админ-панели
async function loadAdminUsers() {
    if (!authToken) return;
    try {
        adminUsersTable.innerHTML = '<tr><td colspan="6">Loading...</td></tr>';
        const response = await fetch(`${API_BASE}/admin/users?page=${adminUsersPage}&limit=${itemsPerPage}`, {
            headers: { 'Authorization': `Bearer ${authToken}` }
        });
        if (!response.ok) throw new Error('Error loading users');
        const data = await response.json();
        renderAdminUsers(data.users || []);
        updateUsersPagination(data.total || 0);
    } catch (error) {
        console.error('Ошибка загрузки пользователей:', error);
        adminUsersTable.innerHTML = '<tr><td colspan="6">Load error</td></tr>';
        showNotification('Failed to load users', 'error');
    }
}

// Отрисовка таблицы пользователей
function renderAdminUsers(users) {
    if (!users.length) {
        adminUsersTable.innerHTML = '<tr><td colspan="7">No users found</td></tr>';
        return;
    }
    let html = '';
    users.forEach(user => {
        const date = new Date(user.created_at).toLocaleDateString('en-US');
        const badge = user.is_admin ? '<span class="badge admin">Admin</span>' : '<span class="badge user">User</span>';
        html += `
            <tr>
                <td>${user.id}</td>
                <td>${user.email}</td>
                <td>${user.username || '-'}</td>
                <td>${user.name || '-'}</td>
                <td>${date}</td>
                <td>${badge}</td>
                <td>
                    <button class="action-btn view" onclick="viewUser(${user.id})">View</button>
                    <button class="action-btn delete" onclick="deleteUser(${user.id})">Delete</button>
                </td>
            </tr>
        `;
    });
    adminUsersTable.innerHTML = html;
}

// Загрузка ссылок для админ-панели
async function loadAdminLinks() {
    if (!authToken) return;
    try {
        adminLinksTable.innerHTML = '<tr><td colspan="7">Loading...</td></tr>';
        const response = await fetch(`${API_BASE}/admin/links?page=${adminLinksPage}&limit=${itemsPerPage}`, {
            headers: { 'Authorization': `Bearer ${authToken}` }
        });
        if (!response.ok) throw new Error('Error loading links');
        const data = await response.json();
        renderAdminLinks(data.links || []);
        updateLinksPagination(data.total || 0);
    } catch (error) {
        console.error('Ошибка загрузки ссылок:', error);
        adminLinksTable.innerHTML = '<tr><td colspan="7">Load error</td></tr>';
        showNotification('Failed to load links', 'error');
    }
}

// Отрисовка таблицы ссылок
function renderAdminLinks(links) {
    if (!links.length) {
        adminLinksTable.innerHTML = '<tr><td colspan="7">No links found</td></tr>';
        return;
    }
    let html = '';
    links.forEach(link => {
        const date = new Date(link.created_at).toLocaleDateString('ru-RU');
        const shortUrl = `${window.location.origin}/${link.short_code}`;
        html += `
            <tr>
                <td>${link.id}</td>
                <td><a href="${shortUrl}" target="_blank">${link.short_code}</a></td>
                <td title="${link.original_url}">${link.original_url.substring(0, 50)}...</td>
                <td>${link.user_email || link.user_id}</td>
                <td>${link.clicks || 0}</td>
                <td>${date}</td>
                <td>
                    <button class="action-btn view" onclick="viewLink('${link.short_code}')">Просмотр</button>
                    <button class="action-btn delete" onclick="deleteLink(${link.id})">Удалить</button>
                </td>
            </tr>
        `;
    });
    adminLinksTable.innerHTML = html;
}

// Загрузка статистики для админ-панели
async function loadAdminStats() {
    if (!authToken) return;
    try {
        const response = await fetch(`${API_BASE}/admin/stats`, {
            headers: { 'Authorization': `Bearer ${authToken}` }
        });
        if (!response.ok) throw new Error('Error loading statistics');
        const data = await response.json();
        adminTotalUsersEl.textContent = data.total_users || 0;
        adminTotalLinksEl.textContent = data.total_links || 0;
        adminTotalClicksEl.textContent = data.total_clicks || 0;
        adminTodayClicksEl.textContent = data.today_clicks || 0;
    } catch (error) {
        console.error('Ошибка загрузки статистики:', error);
        showNotification('Failed to load statistics', 'error');
    }
}

// Обновление пагинации пользователей
function updateUsersPagination(total) {
    const totalPages = Math.ceil(total / itemsPerPage);
    usersPageInfo.textContent = `Page ${adminUsersPage} of ${totalPages}`;
    prevUsersBtn.disabled = adminUsersPage <= 1;
    nextUsersBtn.disabled = adminUsersPage >= totalPages;
}

// Обновление пагинации ссылок
function updateLinksPagination(total) {
    const totalPages = Math.ceil(total / itemsPerPage);
    linksPageInfo.textContent = `Page ${adminLinksPage} of ${totalPages}`;
    prevLinksBtn.disabled = adminLinksPage <= 1;
    nextLinksBtn.disabled = adminLinksPage >= totalPages;
}

// Переключение вкладок админ-панели
function initAdminTabs() {
    adminTabs.forEach(tab => {
        tab.addEventListener('click', () => {
            const tabName = tab.dataset.tab;
            // Убираем активный класс у всех вкладок и панелей
            adminTabs.forEach(t => t.classList.remove('active'));
            adminTabPanes.forEach(pane => pane.classList.remove('active'));
            // Активируем выбранную вкладку и панель
            tab.classList.add('active');
            document.getElementById(`admin-${tabName}-tab`).classList.add('active');
            // Загружаем данные для этой вкладки, если нужно
            if (tabName === 'users') loadAdminUsers();
            if (tabName === 'links') loadAdminLinks();
            if (tabName === 'stats') loadAdminStats();
        });
    });
}

// Вспомогательные функции (заглушки)
function viewUser(id) {
    showNotification(`View user ${id}`, 'info');
}

function deleteUser(id) {
    if (confirm(`Delete user ${id}?`)) {
        showNotification(`User ${id} deleted`, 'success');
        loadAdminUsers();
    }
}

function viewLink(shortCode) {
    window.open(`/${shortCode}`, '_blank');
}

function deleteLink(id) {
    if (confirm(`Delete link ${id}?`)) {
        showNotification(`Link ${id} deleted`, 'success');
        loadAdminLinks();
    }
}

// Инициализация обработчиков событий
function initAuth() {
    // Кнопки открытия модальных окон
    if (loginBtn) loginBtn.addEventListener('click', () => showModal(loginModal));
    if (registerBtn) registerBtn.addEventListener('click', () => showModal(registerModal));
    
    // Закрытие модальных окон
    closeButtons.forEach(btn => {
        btn.addEventListener('click', function() {
            const modal = this.closest('.modal');
            hideModal(modal);
        });
    });
    
    // Закрытие модального окна при клике вне его
    window.addEventListener('click', (event) => {
        if (event.target.classList.contains('modal')) {
            hideModal(event.target);
        }
    });
    
    // Обработка форм
    if (loginForm) loginForm.addEventListener('submit', handleLogin);
    if (registerForm) registerForm.addEventListener('submit', handleRegister);
    
    // Выход
    if (logoutBtn) logoutBtn.addEventListener('click', handleLogout);
    
    // Админ-панель
    if (adminPanelBtn) adminPanelBtn.addEventListener('click', handleAdminPanel);
    
    // Инициализация админ-панели
    if (adminModal) {
        initAdminTabs();
        // Кнопки обновления
        if (refreshUsersBtn) refreshUsersBtn.addEventListener('click', loadAdminUsers);
        if (refreshLinksBtn) refreshLinksBtn.addEventListener('click', loadAdminLinks);
        if (refreshAdminStatsBtn) refreshAdminStatsBtn.addEventListener('click', loadAdminStats);
        // Кнопки экспорта (заглушки)
        if (exportUsersBtn) exportUsersBtn.addEventListener('click', () => showNotification('Export users in development', 'info'));
        if (exportLinksBtn) exportLinksBtn.addEventListener('click', () => showNotification('Export links in development', 'info'));
        if (exportStatsBtn) exportStatsBtn.addEventListener('click', () => showNotification('Export statistics in development', 'info'));
        // Пагинация
        if (prevUsersBtn) prevUsersBtn.addEventListener('click', () => {
            if (adminUsersPage > 1) {
                adminUsersPage--;
                loadAdminUsers();
            }
        });
        if (nextUsersBtn) nextUsersBtn.addEventListener('click', () => {
            adminUsersPage++;
            loadAdminUsers();
        });
        if (prevLinksBtn) prevLinksBtn.addEventListener('click', () => {
            if (adminLinksPage > 1) {
                adminLinksPage--;
                loadAdminLinks();
            }
        });
        if (nextLinksBtn) nextLinksBtn.addEventListener('click', () => {
            adminLinksPage++;
            loadAdminLinks();
        });
    }
    
    // При загрузке страницы обновляем информацию о пользователе
    if (authToken) {
        updateUserInfo();
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

// Обновление превью короткой ссылки
function updatePreview() {
    const code = customInput.value.trim();
    if (code) {
        previewCode.textContent = code;
        previewUrl.style.display = 'block';
    } else {
        previewUrl.style.display = 'none';
    }
}

// Создание короткой ссылки
async function shortenUrl() {
    const url = urlInput.value.trim();
    const title = titleInput.value.trim();
    const customCode = customInput.value.trim();
    
    if (!url) {
        showNotification('Enter URL', 'error');
        return;
    }
    
    if (!isValidUrl(url)) {
        showNotification('Enter a valid URL', 'error');
        return;
    }
    
    shortenBtn.disabled = true;
    shortenBtn.innerHTML = '<i class="fas fa-spinner fa-spin"></i> Processing...';
    
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
            throw new Error(error || 'Server error');
        }
        
        const data = await response.json();
        
        // Показываем результат
        shortUrlInput.value = data.short_url;
        resultDiv.style.display = 'block';
        
        // Прокручиваем к результату
        resultDiv.scrollIntoView({ behavior: 'smooth' });
        
        showNotification('Link created successfully!');
        
        // Обновляем список ссылок и статистику
        loadUserLinks();
        loadStats();
        
        // Очищаем поля
        urlInput.value = '';
        titleInput.value = '';
        customInput.value = '';
        
    } catch (error) {
        console.error('Error shortening URL:', error);
        showNotification(`Error: ${error.message}`, 'error');
    } finally {
        shortenBtn.disabled = false;
        shortenBtn.innerHTML = '<i class="fas fa-magic"></i> Shorten';
    }
}

// Копирование ссылки в буфер обмена
function copyToClipboard() {
    shortUrlInput.select();
    shortUrlInput.setSelectionRange(0, 99999); // Для мобильных устройств
    
    try {
        navigator.clipboard.writeText(shortUrlInput.value)
            .then(() => {
                showNotification('Link copied to clipboard!');
            })
            .catch(err => {
                // Fallback для старых браузеров
                document.execCommand('copy');
                showNotification('Link copied!');
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
                <p>Your shortened links will appear here</p>
            </div>
        `;
        return;
    }
    
    const html = links.map(link => `
        <div class="link-item">
            <div class="link-info">
                <h4>${link.title || 'No title'}</h4>
                <p>${link.original_url.substring(0, 50)}...</p>
                <small>Created: ${new Date(link.created_at).toLocaleDateString('en-US')}</small>
            </div>
            <div class="link-stats">
                <div class="clicks">${link.click_count || 0} clicks</div>
                <a href="${link.short_url}" target="_blank" class="btn-secondary" style="padding: 5px 10px; font-size: 0.8rem; margin-top: 5px;">
                    <i class="fas fa-external-link-alt"></i> Open
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
            title: 'Example link',
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
            // Если ошибка авторизации, используем демо-данные
            if (response.status === 401) {
                showDemoStats();
                return;
            }
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
        showDemoStats();
    }
}

function showDemoStats() {
    // Демо-статистика с реалистичными числами
    totalLinksEl.textContent = '5';
    totalClicksEl.textContent = '120';
    todayClicksEl.textContent = '12';
    // Не показываем уведомление, чтобы не раздражать пользователя
}

// Инициализация
async function init() {
    // Применяем сохранённый язык
    applyLanguage(currentLang);
    initLanguageSwitcher();
    
    // Показываем уведомление о загрузке
    showNotification('Loading...', 'info');
    
    try {
        // Аутентификация
        await ensureAuth();
        
        // Инициализация аутентификационного UI
        initAuth();
        
        // Загрузка начальных данных
        await loadUserLinks();
        await loadStats();
        
        showNotification('Ready to work!');
    } catch (error) {
        console.error('Ошибка инициализации:', error);
        showNotification('Demo mode: data may be limited', 'info');
    }
    
    // Обработчики событий
    shortenBtn.addEventListener('click', shortenUrl);
    copyBtn.addEventListener('click', copyToClipboard);
    refreshBtn.addEventListener('click', () => {
        loadUserLinks();
        loadStats();
        showNotification('Data updated');
    });
    
    // Превью пользовательского кода
    if (customInput) {
        customInput.addEventListener('input', updatePreview);
        // Инициализируем превью, если уже есть значение
        updatePreview();
    }
    
    // Отправка формы по Enter
    urlInput.addEventListener('keypress', (e) => {
        if (e.key === 'Enter') shortenUrl();
    });
}

// Запуск приложения
document.addEventListener('DOMContentLoaded', init);