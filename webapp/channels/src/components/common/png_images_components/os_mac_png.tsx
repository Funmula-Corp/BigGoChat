// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React from 'react';

import OSMac from './image/os_mac.png';

export function OSMacPNG(props: React.HTMLAttributes<HTMLSpanElement>) {
    return (
        <span {...props}>
            <img
                width={48}
                height={48}
                src={OSMac}
            />
        </span>
    );
}
