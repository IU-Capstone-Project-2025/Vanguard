describe('Homepage Test', () => {
  it('Should display homepage correctly', () => {
    cy.visit('/');
    cy.get('button').should('exist').should("be.visible"); // проверка, что кнопка есть

  });
});