describe('Register Test', () => {
  it('Should display register correctly', () => {
    cy.visit('/register');
    cy.get('button').should('exist').should('be.visible');
  });
});