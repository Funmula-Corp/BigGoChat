// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React from 'react';
import {Modal} from 'react-bootstrap';
import {FormattedMessage, injectIntl} from 'react-intl';
import type {IntlShape} from 'react-intl';

import {General} from 'mattermost-redux/constants';

import {trackEvent} from 'actions/telemetry_actions.jsx';

import Constants from 'utils/constants';
import Markdown from 'components/markdown';

type Props = {
    channelDisplayName: string;
    channelId: string;
    intl: IntlShape;

    /**
     * Function injected by ModalController to be called when the modal can be unmounted
     */
    onExited: () => void;

    actions: {
        updateChannelPrivacy: (channelId: string, privacy: string) => void;
    };
}

type State = {
    show: boolean;
    confirmDisplayName: string;
}

export class ConvertChannelModal extends React.PureComponent<Props, State> {
    constructor(props: Props) {
        super(props);

        this.state = {
            show: true,
            confirmDisplayName: '',
        };
    }

    canConvert = () => {
        return this.state.confirmDisplayName === this.props.channelDisplayName;
    };

    handleConvert = () => {
        const {actions, channelId} = this.props;
        if (channelId.length !== Constants.CHANNEL_ID_LENGTH) {
            return;
        }

        actions.updateChannelPrivacy(channelId, General.PRIVATE_CHANNEL);
        trackEvent('actions', 'convert_to_private_channel', {channel_id: channelId});
        this.onHide();
    };

    onUpdateConfirmName = (event: React.ChangeEvent<HTMLInputElement>) => {
        this.setState({
            confirmDisplayName: event.target.value,
        });
    };

    onHide = () => {
        this.setState({show: false});
    };

    render() {
        const formatMessage = this.props.intl.formatMessage;
        const {
            channelDisplayName,
            onExited,
        } = this.props;

        return (
            <Modal
                dialogClassName='a11y__modal'
                show={this.state.show}
                onHide={this.onHide}
                onExited={onExited}
                role='dialog'
                aria-labelledby='convertChannelModalLabel'
            >
                <Modal.Header closeButton={true}>
                    <Modal.Title
                        componentClass='h1'
                        id='convertChannelModalLabel'
                    >
                        <FormattedMessage
                            id='admin.team_channel_settings.convertConfirmModal.toPrivateTitle'
                            defaultMessage='Confirm changing {displayName} to a private channel?'
                            values={{
                                displayName: channelDisplayName,
                            }}
                        />
                    </Modal.Title>
                </Modal.Header>
                <Modal.Body>
                    <p>
                        <Markdown
                            message={formatMessage({
                                id: 'change_to_private_channel_modal.desc',
                                defaultMessage: 'When you convert **{display_name}** to a private channel, history and membership are preserved. Publicly shared files remain accessible to anyone with the link. Membership in a private channel is by invitation only.',
                            }, {
                                display_name: channelDisplayName,
                            })}
                        />
                    </p>
                    <p style={{marginTop: '25px'}}>
                        <div className='Input_wrapper'>
                            <input
                                className='Input form-control medium new-channel-modal-name-input channel-name-input-field'
                                placeholder={formatMessage({id: 'change_to_private_channel_modal.input.placeholder', defaultMessage: 'Enter channel name'})}
                                onChange={this.onUpdateConfirmName}
                                autoFocus={true}
                            />
                        </div>
                    </p>
                    <p style={{fontSize: '12px', color: 'rgba(63, 67 89, 0.75)'}}>
                        <FormattedMessage
                            id='change_to_private_channel_modal.input.hint'
                            defaultMessage="Please enter this channel's name to confirm the change"
                        />
                    </p>
                </Modal.Body>
                <Modal.Footer>
                    <button
                        type='button'
                        className='btn btn-tertiary'
                        onClick={this.onHide}
                    >
                        <FormattedMessage
                            id='convert_channel.cancel'
                            defaultMessage='Cancel'
                        />
                    </button>
                    <button
                        type='button'
                        className='btn btn-primary'
                        data-dismiss='modal'
                        onClick={this.handleConvert}
                        autoFocus={true}
                        disabled={!this.canConvert()}
                        data-testid='convertChannelConfirm'
                    >
                        <FormattedMessage
                            id='convert_channel.confirm'
                            defaultMessage='Confirm changing'
                        />
                    </button>
                </Modal.Footer>
            </Modal>
        );
    }
}

export default injectIntl(ConvertChannelModal);
