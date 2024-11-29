// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React, {PureComponent} from 'react';
import {FormattedMessage} from 'react-intl';

import BrowserStore from 'stores/browser_store';

import AppStorePNG from 'components/common/png_images_components/app_store_png';
import OSAndroidPNG from 'components/common/png_images_components/os_android_png';
import OSIosPNG from 'components/common/png_images_components/os_ios_png';
import OSLinuxPNG from 'components/common/png_images_components/os_linux_png';
import OSMacPNG from 'components/common/png_images_components/os_mac_png';
import OSWindowsPNG from 'components/common/png_images_components/os_windows_png';
import PlayStorePNG from 'components/common/png_images_components/play_store_png';
import BiggoChatLogoSVG from 'components/common/svg_images_components/biggo_chat_logo_svg';
import ExternalLink from 'components/external_link';

import './download.scss';

type Props = {
    window64Link: string;
    macARMLink: string;
    macIntelLink: string;
    linuxDEBLink: string;
    linuxRPMLink: string;
    linuxAPPIMAGELink: string;
    linuxTARGZLink: string;
    iosAppLink: string;
    androidAppLink: string;
}

type State = {
    rememberChecked: boolean;
    redirectPage: boolean;
    location: string;
    nativeLocation: string;
    brandImageError: boolean;
    navigating: boolean;
    desktopVersion: string;
    appVersion: string;
}

export default class DownloadPage extends PureComponent<Props, State> {
    constructor(props: Props) {
        super(props);

        const location = window.location.href.replace('/download', '');

        this.state = {
            rememberChecked: false,
            redirectPage: false,
            location,
            nativeLocation: location.replace(/^(https|http)/, 'mattermost'),
            brandImageError: false,
            navigating: false,
            desktopVersion: '5.9.0',
            appVersion: '5.9.0',
        };

        if (!BrowserStore.hasSeenLandingPage()) {
            BrowserStore.setLandingPageSeen(true);
        }
    }

    getIosLink = () => {
        return (
            <div className='get-app-card'>
                <div className='get-app-card-desc'>
                    <OSIosPNG/>
                    <div className='get-app-card-desc-title-wrap'>
                        <span className='get-app-card-desc-title'>
                            {'iOS'}
                        </span>
                        <span className='get-app-card-desc-subtitle'>
                            {`Supported on iOS 10+ v${this.state.appVersion}`}
                        </span>
                    </div>
                </div>
                <ExternalLink
                    className='get-app-link'
                    location='app_download_modal'
                    href={this.props.iosAppLink}
                >
                    <AppStorePNG/>
                </ExternalLink>
            </div>
        );
    };

    getAndroidLink = () => {
        return (
            <div className='get-app-card'>
                <div className='get-app-card-desc'>
                    <OSAndroidPNG/>
                    <div
                        className='get-app-card-desc-title-wrap'
                        style={{height: '69px'}}
                    >
                        <span className='get-app-card-desc-title'>
                            {'Android'}
                        </span>
                        <span className='get-app-card-desc-subtitle'>
                            {`Supported on Android 10+ v${this.state.appVersion}`}
                        </span>
                    </div>
                </div>
                <ExternalLink
                    className='get-app-link'
                    location='app_download_modal'
                    href={this.props.androidAppLink}
                >
                    <PlayStorePNG/>
                </ExternalLink>
            </div>
        );
    };

    getMacLink = () => {
        return (
            <div className='get-app-card'>
                <div className='get-app-card-desc'>
                    <OSMacPNG/>
                    <div className='get-app-card-desc-title-wrap'>
                        <span className='get-app-card-desc-title'>
                            {'MacOS'}
                        </span>
                        <span className='get-app-card-desc-subtitle'>
                            {`Supported on 10+ v${this.state.desktopVersion}`}
                        </span>
                    </div>
                </div>
                <div className='get-app-link-wrap'>
                    <div className='get-app-link-title'>
                        <FormattedMessage
                            id='download.download'
                            defaultMessage='Download'
                        />
                        {':'}
                    </div>
                    <div className='get-app-link-btn-wrap'>
                        <ExternalLink
                            className='get-app-link'
                            location='app_download_modal'
                            href={this.props.macARMLink}
                        >
                            <div className='btn get-app-link-btn'>
                                {'ARM'}
                            </div>
                        </ExternalLink>
                        <ExternalLink
                            className='get-app-link'
                            location='app_download_modal'
                            href={this.props.macIntelLink}
                        >
                            <div className='btn get-app-link-btn'>
                                {'Intel'}
                            </div>
                        </ExternalLink>
                    </div>
                </div>
            </div>
        );
    };

    getLinuxLink = () => {
        return (
            <div className='get-app-card'>
                <div className='get-app-card-desc'>
                    <OSLinuxPNG/>
                    <div
                        className='get-app-card-desc-title-wrap'
                        style={{height: '69px'}}
                    >
                        <span className='get-app-card-desc-title'>
                            {'Linux'}
                        </span>
                        <span className='get-app-card-desc-subtitle'>
                            {`Supported on Ubuntu 18.04+ v${this.state.desktopVersion}`}
                        </span>
                    </div>
                </div>
                <div className='get-app-link-wrap'>
                    <div className='get-app-link-title'>
                        <FormattedMessage
                            id='download.download'
                            defaultMessage='Download'
                        />
                        {':'}
                    </div>
                    <div className='get-app-link-btn-wrap'>
                        <ExternalLink
                            className='get-app-link'
                            location='app_download_modal'
                            href={this.props.linuxDEBLink}
                        >
                            <div className='btn get-app-link-btn'>
                                {'.DEB'}
                            </div>
                        </ExternalLink>
                        <ExternalLink
                            className='get-app-link'
                            location='app_download_modal'
                            href={this.props.linuxRPMLink}
                        >
                            <div className='btn get-app-link-btn'>
                                {'.RPM'}
                            </div>
                        </ExternalLink>
                        <ExternalLink
                            className='get-app-link'
                            location='app_download_modal'
                            href={this.props.linuxAPPIMAGELink}
                        >
                            <div className='btn get-app-link-btn'>
                                {'.Appimage'}
                            </div>
                        </ExternalLink>
                        <ExternalLink
                            className='get-app-link'
                            location='app_download_modal'
                            href={this.props.linuxTARGZLink}
                        >
                            <div className='btn get-app-link-btn'>
                                {'TAR.GZ'}
                            </div>
                        </ExternalLink>
                    </div>
                </div>
            </div>
        );
    };

    getWindowLink = () => {
        return (
            <div className='get-app-card'>
                <div className='get-app-card-desc'>
                    <OSWindowsPNG/>
                    <div className='get-app-card-desc-title-wrap'>
                        <span className='get-app-card-desc-title'>
                            {'Window'}
                        </span>
                        <span className='get-app-card-desc-subtitle'>
                            {`Supported on 12+ v${this.state.desktopVersion}`}
                        </span>
                    </div>
                </div>
                <div className='get-app-link-wrap'>
                    <div className='get-app-link-title'>
                        <FormattedMessage
                            id='download.download'
                            defaultMessage='Download'
                        />
                        {':'}
                    </div>
                    <div className='get-app-link-btn-wrap'>
                        <ExternalLink
                            className='get-app-link'
                            location='app_download_modal'
                            href={this.props.window64Link}
                        >
                            <div className='btn get-app-link-btn'>
                                {'64bit/ARM 64bit'}
                            </div>
                        </ExternalLink>
                    </div>
                </div>
            </div>
        );
    };

    handleBrandImageError = () => {
        this.setState({brandImageError: true});
    };

    renderDialogBody = () => {
        return (
            <div className='get-app-body'>
                <section className='get-app-top'>
                    <div className='get-app-title'>
                        <BiggoChatLogoSVG
                            fill='white'
                            height={52}
                            width={271}
                        />
                        <FormattedMessage
                            id='get_app.description'
                            defaultMessage='Mobile and Desktop Apps'
                        />
                    </div>
                </section>
                <section className='get-app-bottom'>
                    <div className='get-app-card-wrap'>
                        {this.getIosLink()}
                        {this.getAndroidLink()}
                        {this.getWindowLink()}
                        {this.getMacLink()}
                        {this.getLinuxLink()}
                    </div>
                </section>
            </div>
        );
    };

    render() {
        // const isMobile = UserAgent.isMobile();

        return (
            <div className='get-app'>
                {this.renderDialogBody()}
            </div>
        );
    }
}
