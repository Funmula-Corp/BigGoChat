// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React from 'react';
import type {ReactNode, MouseEvent} from 'react';
import {injectIntl, type WrappedComponentProps} from 'react-intl';

import type {Team} from '@mattermost/types/teams';

import {Client4} from 'mattermost-redux/client';

import OverlayTrigger from 'components/overlay_trigger';
import Tooltip from 'components/tooltip';
import TeamInfoIcon from 'components/widgets/icons/team_info_icon';

import * as Utils from 'utils/utils';

interface Props extends WrappedComponentProps {
    teamIcon?: string | null;
    team: Team;
    teamJoined: boolean;
    inviteId: string;
    isInviteIdValid: boolean;
    onTeamClick: (team: Team, inviteId: string) => void;
    loading: boolean;
    canJoinPublicTeams: boolean;
    canJoinPrivateTeams: boolean;
}

interface State {
    inviteId: string;
    inviteIdError: boolean; // 邀請碼是否符合搜尋到的Team
}

export class SelectTeamItem extends React.PureComponent<Props, State> {
    state = {
        inviteId: this.props.inviteId, // 搜尋使用邀請碼直接帶入
        inviteIdError: false,
    };

    handleInviteIdChange = async (e: React.ChangeEvent<HTMLInputElement>) => {
        this.setState({inviteId: e.target.value});
        if (e.target.value.length === 0) {
            this.setState({inviteIdError: false});
            return;
        }
        try {
            const team = await Client4.getTeamInviteInfo(e.target.value);
            if (!team || team.id !== this.props.team.id) {
                this.setState({inviteIdError: true});
            } else {
                this.setState({inviteIdError: false});
            }
        } catch (error) {
            this.setState({inviteIdError: true});
        }
    };

    handleTeamClick = (e: MouseEvent): void => {
        e.preventDefault();
        this.props.onTeamClick(this.props.team, (this.state.inviteId ?? ''));
    };

    renderDescriptionTooltip = (): ReactNode => {
        const team = this.props.team;
        if (!team.description) {
            return null;
        }

        const descriptionTooltip = (
            <Tooltip id='team-description__tooltip'>
                {team.description}
            </Tooltip>
        );

        return (
            <OverlayTrigger
                delayShow={1000}
                placement='top'
                overlay={descriptionTooltip}
                rootClose={true}
                container={this}
            >
                <TeamInfoIcon className='icon icon--info'/>
            </OverlayTrigger>
        );
    };

    render() {
        const {canJoinPublicTeams, canJoinPrivateTeams, loading, team} = this.props;
        const canJoinPublic = team.allow_open_invite && canJoinPublicTeams;
        const canJoinPrivate = !team.allow_open_invite && canJoinPrivateTeams;
        const inviteIdVerified = !this.state.inviteIdError && this.state.inviteId.length > 0;
        const canJoin = this.props.teamJoined || canJoinPublic || canJoinPrivate || inviteIdVerified;

        return (
            <div className='signup-team-dir'>
                {/* {this.renderDescriptionTooltip()} */}
                <div className='signup-team-dir__content'>
                    {this.props.teamIcon ? (
                        <img
                            className='team-picture-section__team-icon'
                            src={this.props.teamIcon}
                        />
                    ) : (
                        <div className='team-picture-section__team-icon'>
                            <span
                                id='teamIconInitial'
                                className='team-picture-section__team-name'
                            >{team.display_name.charAt(0).toUpperCase() + team.display_name.charAt(1)}</span>
                        </div>
                    )}
                    <div className='signup-team-dir__info'>
                        <span className='signup-team-dir__name'>{team.display_name}</span>
                        {team.description &&
                            <span className='signup-team-desc'>{team.description}</span>
                        }
                        {!team.allow_open_invite &&
                            <span className='signup-team-text private'>
                                <i className='icon icon-lock-outline'/>
                                {this.props.intl.formatMessage({id: 'select_team.private.icon', defaultMessage: 'Private team'})}
                            </span>
                        }
                        {this.props.teamJoined &&
                            <span className='signup-team-text joined'>
                                <i className='icon icon-check'/>
                                {this.props.intl.formatMessage({id: 'select_team.joined', defaultMessage: 'You are already a team member'})}
                            </span>
                        }
                        {/* invite code input -> private team && not use invite code search && not joined */}
                        {!team.allow_open_invite && !this.props.isInviteIdValid && !this.props.teamJoined && (
                            <div className={`signup-team-invite-code ${this.state.inviteIdError ? 'error' : ''}`}>
                                <input
                                    placeholder={this.props.intl.formatMessage({id: 'select_team.invite_code', defaultMessage: 'Enter invitation code...'})}
                                    value={this.state.inviteId}
                                    onChange={this.handleInviteIdChange}
                                />
                                {this.state.inviteIdError &&
                                    <span>
                                        {this.props.intl.formatMessage({id: 'select_team.invite_code_error', defaultMessage: 'Incorrect invitation code, please re-enter'})}
                                    </span>
                                }
                            </div>
                        )
                        }
                    </div>
                </div>
                {/* join or go to button -> use invite code search (no need to re-enter invite code can join directly) || not joined || already joined (go to) */}
                {(this.props.isInviteIdValid || canJoin) && (
                    <a
                        href='#'
                        id={Utils.createSafeId(team.display_name)}
                        onClick={canJoin ? this.handleTeamClick : undefined}
                        className={'signup-team-dir__link'}
                    >
                        {this.props.teamJoined ? this.props.intl.formatMessage({id: 'select_team.forward', defaultMessage: 'Go'}) : this.props.intl.formatMessage({id: 'select_team.join', defaultMessage: 'Join'})}
                    </a>
                )}
            </div>
        );
    }
}

export default injectIntl(SelectTeamItem);
