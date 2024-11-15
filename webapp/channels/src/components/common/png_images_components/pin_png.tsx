// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React from 'react';

import Pin from './image/pin.png';

export function PinPNG(props: React.HTMLAttributes<HTMLSpanElement>) {
    return (
        <span {...props}>
            <img src={Pin}/>
        </span>
    );
}
