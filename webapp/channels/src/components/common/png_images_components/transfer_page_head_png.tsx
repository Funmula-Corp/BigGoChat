// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React from 'react';

import Img from './image/transfer_page_head.png';

type PngProps = {
    width?: number;
    height?: number;
    className?: string;
}

const Png = (props: PngProps) => (
    <img
        {...props}
        width={props.width ?? '96'}
        height={props.height ?? '96'}
        src={Img}
    />
);

export default Png;
