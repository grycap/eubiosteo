
import React, {Component} from 'react'

const Header = ({onLogout, onHostChange, onProjectsBack}) => (
    <header className={"hg-header"}>
        <img
            className="logoindigo"
            src="/static/images/indigo.png"
            width="150"
            height="100"
        />
        <ul className={'user-info-menu right-links list-inline list-unstyled'}>
            <li>
                <a
                    href="javascript:void(0)"
                    onClick={onHostChange}
                >
                    <span>
                        Host
                    </span>
                </a>
            </li>
            <li>
                <a>
                    <span>
                        version 1.0
                    </span>
                </a>
            </li>
            <li>
                <a
                    href="javascript:void(0)"
                    onClick={onLogout}
                >
                    <span>
                        Logout
                    </span>
                </a>
            </li>
        </ul>
    </header>
)

export default Header
