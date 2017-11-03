
import React from 'react' 

export default class Tools extends React.Component {
    render() {
        let {tools, onClick, pathname} = this.props;

        return (
            <div id="tools">
                <ul>
                    {tools.map((item, index) => {
                        let {name, endpoint, dest} = item;

                        let className = '';
                        if (pathname.indexOf(endpoint) !== -1) {
                            className += ' toolsel'
                        }

                        return (
                            <li
                                className={className}
                                key={index}
                                onClick={() => {
                                    let destiny = (dest == undefined) ? endpoint : dest;
                                    onClick(destiny)
                                }}
                            >
                                <div>
                                    {name}
                                </div>
                            </li>
                        )
                    })}
                </ul>
            </div>
        )
    }
}
