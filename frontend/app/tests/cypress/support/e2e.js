import './commands';

// Глобальные настройки
beforeEach(() => {
    // Установите базовый URL вашего приложения
    cy.visit('http://localhost:3000'); // Или ваш URL
});