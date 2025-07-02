describe('Enter Nickname Test', () => {
  it('Should display enter nickname page correctly', () => {
    cy.visit('/enter-nickname');
    cy.get('button').should('exist').should('be.visible');
  });
});