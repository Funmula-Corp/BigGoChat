// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React from 'react';
import type {CSSProperties} from 'react';
import {useIntl} from 'react-intl';

export default function BigGoChatLogo(props: React.HTMLAttributes<HTMLSpanElement>) {
    const {formatMessage} = useIntl();
    return (
        <span {...props}>
            <svg
                version='1.1'
                x='0px'
                y='0px'
                viewBox='0 0 64 64'
                enableBackground='new 0 0 64 64'
                role='img'
                aria-label={formatMessage({id: 'generic_icons.biggochat', defaultMessage: 'BigGo Chat Logo'})}
            >
                <g clipPath='url(#clip0_1299_16440)'>
                    <path
                        d='M47.5476 64H16.4538C5.952 64 0 58.048 0 47.5476V16.4524C0 5.952 5.952 0 16.4538 0H47.5476C58.048 0 64 5.952 64 16.4524V47.5476C64 58.048 58.048 64 47.5476 64Z'
                        fill='#488EFF'
                    />
                    <mask
                        id='mask0_1299_16440'
                        style={{maskType: 'luminance'}}
                        maskUnits='userSpaceOnUse'
                        x='0'
                        y='0'
                        width='64'
                        height='64'
                    >
                        <path
                            d='M47.5476 64H16.4538C5.952 64 0 58.048 0 47.5476V16.4524C0 5.952 5.952 0 16.4538 0H47.5476C58.048 0 64 5.952 64 16.4524V47.5476C64 58.048 58.048 64 47.5476 64Z'
                            fill='white'
                        />
                    </mask>
                    <g mask='url(#mask0_1299_16440)'>
                        <path
                            style={style}
                            fillRule='evenodd'
                            clipRule='evenodd'
                            d='M0 47.1594V44.9208L5.568 38.1528H7.45891L0 47.1594Z'
                            fill='black'
                        />
                        <path
                            style={style}
                            fillRule='evenodd'
                            clipRule='evenodd'
                            d='M64.0004 45.5753C60.1211 47.1229 55.3865 47.936 49.9814 47.936C46.1094 47.936 42.5851 47.5186 39.4694 46.7127C39.0273 46.6007 38.6607 46.7607 38.2374 46.9295C37.7909 47.1069 33.7894 48.9382 31.2585 50.0989C30.5298 50.4335 29.7807 49.6742 30.1065 48.9309C30.8862 47.1506 31.9058 44.9004 32.0411 44.5862C32.2433 44.1164 32.0789 43.5898 31.6993 43.3396C27.0331 40.1949 24.3945 35.5171 24.3945 29.5549C24.3945 17.2175 34.4571 11.1738 49.9814 11.1738C55.3036 11.1738 59.9669 11.8735 63.804 13.2844C63.9349 14.288 64.0004 15.3455 64.0004 16.4524V45.5753Z'
                            fill='white'
                        />
                        <path
                            style={style}
                            d='M10.6279 29.5818V29.0931C10.6279 24.7629 12.8621 21.48 15.0279 19.6865H18.2643C16.3079 21.9222 14.6076 25.1571 14.6076 28.9069V29.6283C14.6076 33.4698 16.3079 36.776 18.2643 39.0349H15.0279C12.8621 37.2182 10.6279 33.9353 10.6279 29.5818Z'
                            fill='black'
                        />
                        <path
                            style={style}
                            d='M22.242 20.8511H25.4551L28.8544 26.6023H25.6646L23.8479 23.9245L22.0093 26.6023H18.8428L22.242 20.8511Z'
                            fill='black'
                        />
                        <path
                            style={style}
                            d='M31.9707 22.9811C31.9707 21.7738 32.9496 20.7949 34.1569 20.7949H38.2325V22.4982H35.6565V37.2182H38.2325V38.9244H34.1569C32.9496 38.9244 31.9707 37.9455 31.9707 36.7367V22.9811Z'
                            fill='black'
                        />
                        <path
                            style={style}
                            d='M39.3652 37.2187H41.918V22.4987H39.3652V20.7939H43.4176C44.6263 20.7939 45.6052 21.7743 45.6052 22.9816V36.7372C45.6052 37.9459 44.6263 38.9249 43.4176 38.9249H39.3652V37.2187Z'
                            fill='black'
                        />
                        <path
                            style={style}
                            d='M52.1483 20.8511H55.3614L58.7607 26.6023H55.5708L53.7541 23.9245L51.9156 26.6023H48.749L52.1483 20.8511Z'
                            fill='black'
                        />
                    </g>
                </g>
                <defs>
                    <clipPath id='clip0_1299_16440'>
                        <rect
                            width='64'
                            height='64'
                            fill='white'
                        />
                    </clipPath>
                </defs>
            </svg>
        </span>
    );
}

const style: CSSProperties = {
    fillRule: 'evenodd',
    clipRule: 'evenodd',
};
