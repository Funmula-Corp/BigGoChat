// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React from 'react';

import OSWindows from './image/os_windows.png';

const PNG = (props: Partial<React.ImgHTMLAttributes<HTMLImageElement>>) => (
    <img
        {...props}
        width={props.width ?? '48'}
        height={props.height ?? '48'}
        src={OSWindows}
    />
);

export default PNG;
