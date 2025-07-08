const { defineConfig } = require('cypress')
const { BASE_URL } = require('../src/constants/api');

module.exports = defineConfig({
  e2e: {
    baseUrl: BASE_URL,
    supportFile: '../tests/cypress/support/e2e.js',
    specPattern: '../tests/cypress/e2e/**/*.cy.{js,jsx,ts,tsx}',
    fixturesFolder: '../tests/cypress/fixtures'
  }
})
