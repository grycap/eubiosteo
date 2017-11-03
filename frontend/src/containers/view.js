
import React from 'react'

import { browserHistory } from 'react-router';

// con el boton de back
const CreateUpdate = ({name, onBack}) => (
    <div className={'navbar'}>
        <div
            className={'title'}
        >
            {name}
        </div>
        <a
            href="javascript:void(0)"
            onClick={onBack}
        >
            {'Back'}
        </a>
    </div>
)

// la normal con los items de los views
const Read = ({name, pathname, hasViews, views}) => (
    <div className={'navbar'}>
        <div
            className={'title'}
        >
            {name}
        </div>
        {hasViews &&
            <div className={'tabs'}>
                {views.map((local, index) => {
                    let {name, endpoint} = local;

                    let className = 'tab'
                    if (pathname.indexOf(endpoint) !== -1) {
                        className += ' tabsel'
                    }

                    return (
                        <div
                            key={index}
                            className={className}
                             onClick={() => {
                                browserHistory.push(endpoint)
                            }}
                        >
                            {name}
                        </div>
                    )
                })}
            </div>
        }
    </div>
)

const getPathname = (props) => {
    try {
        return props.location.pathname
    } catch(err) {
        return undefined
    }
}

export default function withHeader(view) {
    return class extends React.Component {
        render() {
            let {name, views, routes} = view;
            let hasViews = (views != undefined && Object.keys(views).length != 0)
            
            let pathname = getPathname(this.props)
            if (pathname == undefined) {
                console.error("El pathname es undefined")
            }

            routes = (routes == undefined) ? {} : routes
            let isRoute = routes[pathname] != undefined

            let View = undefined
            if (isRoute) {
                let {name, back} = routes[pathname];

                View = <CreateUpdate 
                    name={name}
                    onBack={() => {
                        browserHistory.push(back)
                    }}
                />
            } else {
                View = <Read name={name} pathname={pathname} hasViews={hasViews} views={views} />
            }
            
            // mirar si el children es undeifned entonces puedes redirigir al bueno
            // eso en otra propiedad

            // estableceer cuales son las llamadas de crear o actualizar
            // si se pueden establecer params o algo asi.
            
            return (
                <div>
                    <div
                        style={{backgroundColor: 'cyan'}}
                    >
                        {View}
                    </div>
                    <div className={'panel'}>
                        {this.props.children}
                    </div>
                </div>
            )
        }
    }
}
