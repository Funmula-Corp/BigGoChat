// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React, {useRef} from 'react';
import {useDispatch, useSelector} from 'react-redux';
import styled from 'styled-components';

import IconButton from '@mattermost/compass-components/components/icon-button'; // eslint-disable-line no-restricted-imports
import glyphMap from '@mattermost/compass-icons/components';

import {getInt} from 'mattermost-redux/selectors/entities/preferences';

import {setProductMenuSwitcherOpen} from 'actions/views/product_menu';
import {isSwitcherOpen} from 'selectors/views/product_menu';

import {
    GenericTaskSteps,
    OnboardingTaskCategory,
    OnboardingTasksName,
} from 'components/onboarding_tasks';
import {FINISHED} from 'components/tours';

import type {GlobalState} from 'types/store';

import {useClickOutsideRef} from '../../hooks';

export const ProductMenuContainer = styled.nav`
    display: flex;
    align-items: center;
    cursor: pointer;

    > * + * {
        margin-left: 12px;
    }
`;

export const ProductMenuButton = styled(IconButton).attrs(() => ({
    id: 'product_switch_menu',
    icon: 'products',
    size: 'sm',

    // we currently need this, since not passing a onClick handler is disabling the IconButton
    // this is a known issue and is being tracked by UI platform team
    // TODO@UI: remove the onClick, when it is not a mandatory prop anymore
    onClick: () => {},
    inverted: true,
    compact: true,
}))`
    > i::before {
        font-size: 20px;
        letter-spacing: 20px;
    }
`;

const ProductMenu = (): JSX.Element => {
    const dispatch = useDispatch();
    const switcherOpen = useSelector(isSwitcherOpen);
    const menuRef = useRef<HTMLDivElement>(null);

    const triggerStep = useSelector((state: GlobalState) => getInt(state, OnboardingTaskCategory, OnboardingTasksName.EXPLORE_OTHER_TOOLS, FINISHED));
    const exploreToolsTourTriggered = triggerStep === GenericTaskSteps.STARTED;

    useClickOutsideRef(menuRef, () => {
        if (exploreToolsTourTriggered || !switcherOpen) {
            return;
        }
        dispatch(setProductMenuSwitcherOpen(false));
    });

    const MenuItem = styled.div`
        && {
            text-decoration: none;
            color: inherit;
        }

        height: 40px;
        width: 270px;
        padding-left: 16px;
        padding-right: 20px;
        display: flex;
        align-items: center;
        position: relative;

        button {
            padding: 0 6px;
        }
    `;
    const ProductIcon = glyphMap['product-channels'];
    const MenuItemTextContainer = styled.div`
        margin-left: 8px;
        flex-grow: 1;
        font-weight: 600;
        font-size: 14px;
        line-height: 20px;
    `;

    return (
        <MenuItem>
            <ProductIcon
                size={24}
            />
            <MenuItemTextContainer>
                {'Channels'}
            </MenuItemTextContainer>
        </MenuItem>
    );
};

export default ProductMenu;
