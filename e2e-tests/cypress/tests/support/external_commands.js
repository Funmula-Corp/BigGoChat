// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import {getAdminAccount} from './env';

Cypress.Commands.add('externalActivateUser', (userId, active = true) => {
    const baseUrl = Cypress.config('baseUrl');
    const admin = getAdminAccount();

    cy.externalRequest({user: admin, method: 'put', baseUrl, path: `users/${userId}/active`, data: {active}});
});

Cypress.Commands.add('externalSendInvite', (email, teamId) => {
    const baseUrl = Cypress.config('baseUrl');
    const admin = getAdminAccount();

    cy.externalRequest({user: admin, method: 'post', baseUrl, path: `teams/${teamId}/invite/email?graceful=false`, data: {'emails': [email]}});
});

Cypress.Commands.add('externalPatchUserRoles', (userId, roleNames) => {
    const baseUrl = Cypress.config('baseUrl');
    const admin = getAdminAccount();

    cy.externalRequest({user: admin, method: 'put', baseUrl, path: `users/${userId}/roles`, data: {roles: roleNames = ['system_user'].join(' ')},});
});
