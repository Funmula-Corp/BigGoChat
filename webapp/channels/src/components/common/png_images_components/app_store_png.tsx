// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React from 'react';

import AppStore from './image/app_store.png';

const PNG = (props: Partial<React.ImgHTMLAttributes<HTMLImageElement>>) => (
    <img
        {...props}
        width={props.width ?? '95'}
        height={props.height ?? '32'}
        src={AppStore}
    />
);

export default PNG;
