// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

Cypress.Commands.add('uiLogin', (user = {}) => {
    // BigGo dosen't have password login, so we use apiLogin to login
    cy.apiLogin(user).reload()

    // cy.url().should('include', '/login');
    // // # Type email and password, then Sign in
    // cy.get('#input_loginId').should('be.visible').type(user.email);
    // cy.get('#input_password-input').should('be.visible').type(user.password);
    // cy.get('#saveSetting').should('not.be.disabled').click();
});
