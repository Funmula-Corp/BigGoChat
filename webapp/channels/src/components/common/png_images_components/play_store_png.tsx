// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React from 'react';

import PlayStore from './image/play_store.png';

type PngProps = {
    width?: number;
    height?: number;
}

const Png = (props: PngProps) => (
    <img
        width={props.width ?? '107'}
        height={props.height ?? '32'}
        src={PlayStore}
    />
);

export default Png;
