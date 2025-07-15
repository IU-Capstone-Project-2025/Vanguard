describe('admin:home-create-ask-wait; user: play-code-nick-wait', () => {
  let sessionCode; // Будет доступно для всех тестов

  it('Admin creates session', () => {
    // Действия админа
    cy.visit('/');
    cy.get('.create-button').click();
    cy.wait(2000);
    cy.contains('.quiz-item', 'Basic Python Knowledge',{timeout:10000}).click();
    cy.get('.quiz-item.selected').should('exist'); // Явная проверка выбора
    cy.get('.play-button').click();

    // Получаем код сессии
    cy.url().should('match', /\/[A-Z0-9]{6}$/); // Проверяем формат URL
    cy.url().then((url) => {
      sessionCode = url.split('/').pop();
      cy.log('Admin created session with code:', sessionCode);
      expect(sessionCode).to.match(/^[A-Z0-9]{6}$/); // Проверка формата кода
    });
    cy.get('.play-button').click();
    cy.get('.player-box span').should('contain','Admin');
  });

  it('User joins session', () => {
    // Проверяем, что код сессии существует
    if (!sessionCode) {
      throw new Error('Session code was not created by admin!');
    }

    // Действия пользователя
    cy.visit('/');
    cy.get('.play-button').click();
    cy.get('input[type="text"]').type(sessionCode);
    cy.get('.play-button').click();
    cy.get('input[type="text"]').type("Bob");
    cy.get('.play-button').click();
    cy.get('.player-box span').should('contain','Admin').and('contain','Bob');

  });

});