describe('Create Session Test', () => {
  it('Should display create session correctly', () => {
    cy.visit('/create');
    cy.get('button').should('exist').should('be.visible');
  });
});