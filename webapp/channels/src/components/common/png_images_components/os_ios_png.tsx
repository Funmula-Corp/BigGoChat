// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React from 'react';

import OSIos from './image/os_ios.png';

export function OSIosPNG(props: React.HTMLAttributes<HTMLSpanElement>) {
    return (
        <span {...props}>
            <img
                width={48}
                height={48}
                src={OSIos}
            />
        </span>
    );
}
