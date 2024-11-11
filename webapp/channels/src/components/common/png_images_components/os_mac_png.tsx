// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React from 'react';

import OSMac from './image/os_mac.png';

type PngProps = {
    width?: number;
    height?: number;
}

const Png = (props: PngProps) => (
    <img
        width={props.width ?? '48'}
        height={props.height ?? '48'}
        src={OSMac}
    />
);

export default Png;
