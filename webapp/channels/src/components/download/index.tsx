// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import {connect} from 'react-redux';

import {Client4} from 'mattermost-redux/client';
import {getConfig} from 'mattermost-redux/selectors/entities/general';
import {getTheme} from 'mattermost-redux/selectors/entities/preferences';

import type {GlobalState} from 'types/store';

import DownloadPage from './download';

function mapStateToProps(state: GlobalState) {
    const config = getConfig(state);

    return {
        window64Link: 'https://img.bgo.one/biggo-chat-download/biggo-chat-setup-5.10.0-win.exe',
        macARMLink: 'https://img.bgo.one/biggo-chat-download/biggo-chat-5.10.0-mac-arm64.dmg',
        macIntelLink: 'https://img.bgo.one/biggo-chat-download/biggo-chat-5.10.0-mac-x64.dmg',
        linuxDEBLink: 'https://img.bgo.one/biggo-chat-download/biggo-chat_5.10.0_amd64.deb',
        linuxRPMLink: 'https://img.bgo.one/biggo-chat-download/biggo-chat-5.10.0-linux-x86_64.rpm',
        linuxAPPIMAGELink: 'https://img.bgo.one/biggo-chat-download/biggo-chat-5.10.0-linux-x86_64.AppImage',
        linuxTARGZLink: 'https://img.bgo.one/biggo-chat-download/biggo-chat-5.10.0-linux-x64.tar.gz',
        iosAppLink: 'https://apps.apple.com/tw/app/biggo-chat/id6503929831',
        androidAppLink: 'https://play.google.com/store/apps/details?id=com.funmula.biggo.chat',
        defaultTheme: getTheme(state),
        siteUrl: config.SiteURL,
        siteName: config.SiteName,
        brandImageUrl: Client4.getBrandImageUrl('0'),
        enableCustomBrand: config.EnableCustomBrand === 'true',
    };
}

export default connect(mapStateToProps)(DownloadPage);
