// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React from 'react';

import OSAndroid from './image/os_android.png';

type PngProps = {
    width?: number;
    height?: number;
    className?: string;
}

const Png = (props: PngProps) => (
    <img
        {...props}
        width={props.width ?? '48'}
        height={props.height ?? '48'}
        src={OSAndroid}
    />
);

export default Png;
