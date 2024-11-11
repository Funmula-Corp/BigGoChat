// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React from 'react';

import PlayStore from './image/play_store.png';

export function PlayStorePNG(props: React.HTMLAttributes<HTMLSpanElement>) {
    return (
        <span {...props}>
            <img
                width={127}
                height={38}
                src={PlayStore}
            />
        </span>
    );
}
