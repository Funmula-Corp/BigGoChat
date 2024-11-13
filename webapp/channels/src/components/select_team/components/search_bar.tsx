// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React, {useRef, useState} from 'react';
import type {ChangeEvent} from 'react';
import {FormattedMessage, useIntl} from 'react-intl';

type Props = {
    handleChange: (e: ChangeEvent<HTMLInputElement>) => void;
    handleClear?: () => void;
    isError?: boolean;
}

const SearchBar: React.FunctionComponent<Props> = (props: Props): JSX.Element => {
    const [hasInput, setHasInput] = useState(false);
    const intl = useIntl();
    const {formatMessage} = intl;
    const inputRef = useRef<HTMLInputElement>(null);

    const handleInputClear = () => {
        if (inputRef.current) {
            inputRef.current.value = '';
        }
        props.handleClear?.();
        setHasInput(false);
    };

    const handleChange = (e: ChangeEvent<HTMLInputElement>) => {
        setHasInput(e.target.value.length > 0);
        props.handleChange?.(e);
    };

    return (
        <>
            <div
                id={'searchFormContainer'}
                className={`search-form__container search__form ${(hasInput && props.isError) ? 'error' : ''}`}
            >
                <input
                    ref={inputRef}
                    id='searchTeamInput'
                    className='search_team_input'
                    placeholder={formatMessage({id: 'signup_team.join_search-input-placeholder', defaultMessage: 'Enter team ID...'})}
                    onChange={handleChange}
                />
                <div className={`input-clear ${hasInput ? 'visible' : ''}`}>
                    <span className='input-clear-x'>
                        <i
                            className='icon icon-close-circle'
                            onClick={handleInputClear}
                        />
                    </span>
                </div>
                {(hasInput && props.isError) &&
                    <div className='search-form__error_text'>
                        <FormattedMessage
                            id='signup_team.join_search-input-error'
                            defaultMessage={'Team not found!'}
                        />
                    </div>
                }
            </div>
            {!hasInput &&
                <FormattedMessage
                    id='signup_team.join_search-input'
                    defaultMessage={'Enter team ID or paste URL to search'}
                />
            }
        </>
    );
};

export default SearchBar;
