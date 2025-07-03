describe('Authentication Test', () => {
  it('Should display Authentication page correctly', () => {
    cy.visit('/login');
    cy.get('button').should('exist').should("be.visible"); // проверка, что кнопка есть
  });
});