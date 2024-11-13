// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React from 'react';
import {FormattedMessage, useIntl} from 'react-intl';

import ExternalLink from 'components/external_link';
import BaseSettingItem from 'components/widgets/modals/components/base_setting_item';
import CheckboxSettingItem from 'components/widgets/modals/components/checkbox_setting_item';

type Props = {
    allowOpenInvite: boolean;
    isGroupConstrained?: boolean;
    setAllowOpenInvite: (value: boolean) => void;
};

const OpenInvite = ({isGroupConstrained, allowOpenInvite, setAllowOpenInvite}: Props) => {
    const {formatMessage} = useIntl();
    if (isGroupConstrained) {
        const groupConstrainedContent = (
            <p id='groupConstrainedContent' >{
                formatMessage({
                    id: 'team_settings.openInviteDescription.groupConstrained',
                    defaultMessage: 'Members of this team are added and removed by linked groups. <link>Learn More</link>',
                }, {
                    link: (msg: React.ReactNode) => (
                        <ExternalLink
                            href='https://mattermost.com/pl/default-ldap-group-constrained-team-channel.html'
                            location='open_invite'
                        >
                            {msg}
                        </ExternalLink>
                    ),
                })}
            </p>
        );
        return (
            <BaseSettingItem
                className='access-invite-domains-section'
                title={formatMessage({
                    id: 'general_tab.openInviteText',
                    defaultMessage: 'Any user can join the team',
                })}
                description={formatMessage({
                    id: 'general_tab.openInviteDesc',
                    defaultMessage: 'When checked, any user can join the team on their own. If unchecked, it will become a private team requiring an invitation code or invite link to join. Previous invitation codes and links will become invalid after the update and cannot be used.',
                })}
                descriptionAboveContent={true}
                content={groupConstrainedContent}
            />
        );
    }

    return (
        <CheckboxSettingItem
            className='access-invite-domains-section'
            inputFieldTitle={
                <FormattedMessage
                    id='general_tab.openInviteTitle'
                    defaultMessage='Allow any user to join the team'
                />
            }
            inputFieldData={{name: 'name'}}
            inputFieldValue={allowOpenInvite}
            handleChange={setAllowOpenInvite}
            title={formatMessage({
                id: 'general_tab.openInviteText',
                defaultMessage: 'Any user can join the team',
            })}
            description={formatMessage({
                id: 'general_tab.openInviteDesc',
                defaultMessage: 'When checked, any user can join the team on their own. If unchecked, it will become a private team requiring an invitation code or invite link to join. Previous invitation codes and links will become invalid after the update and cannot be used.',
            })}
            descriptionAboveContent={true}
        />
    );
};

export default OpenInvite;
