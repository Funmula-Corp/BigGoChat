// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React from 'react';

import Img from './image/transfer_page_head.png';

const PNG = (props: Partial<React.ImgHTMLAttributes<HTMLImageElement>>) => (
    <img
        {...props}
        width={props.width ?? '96'}
        height={props.height ?? '96'}
        src={Img}
    />
);

export default PNG;
