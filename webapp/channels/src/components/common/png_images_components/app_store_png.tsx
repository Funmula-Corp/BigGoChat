// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React from 'react';

import AppStore from './image/app_store.png';

export function AppStorePNG(props: React.HTMLAttributes<HTMLSpanElement>) {
    return (
        <span {...props}>
            <img
                width={110}
                height={36}
                src={AppStore}
            />
        </span>
    );
}
