describe('Join Game Test', () => {
  it('Should display join game page correctly', () => {
    cy.visit('/join');
    cy.get('button').should('exist').should('be.visible');
  });
});