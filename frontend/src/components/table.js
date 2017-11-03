
import React, {Component} from 'react'

// -- colsgroup --

const ColsGroup = ({checkbox, header}) => (
    <colgroup>
        {checkbox &&
            <col
                style={{width: 10}}
            />
        }
        {header.map((item, index) => (
            <col 
                key={index}
                style={{width: item.width}} 
            />
        ))}
    </colgroup>
)

// -- header arrow --

const HeaderImage = ({image, onChange}) => (
    <a
        href="javascript:void(0)"
        onClick={onChange}
    >
        <img src={image} />
    </a>
)

const HeaderArrow = ({header, sort, order, onChange}) => {
    
    if (header.name != sort) {
        return (
            <HeaderImage
                image={'/images/arrow_down_grey.png'}
                onChange={() => {onChange(header.name, 'asc')}}
            />
        )
    }

    if (order == 'asc') { 
        return (
            <HeaderImage
                image={'/images/arrow_down_black.png'}
                onChange={() => {onChange(header.name, 'desc')}}
            />
        )
    }

    return (
        <HeaderImage
            image={'/images/arrow_up_black.png'}
            onChange={() => onChange()}
        />
    )
}

// -- header --

const Header = ({header, sort, order, onChange}) => (
    <th>
        {header.label === undefined ? header.name : header.label}
        {header.sort &&
            <HeaderArrow header={header} sort={sort} order={order} onChange={onChange} />
        }
    </th>
)

// -- headers --

const Headers = ({checkbox, allSelected, header, sort, order, onSelectedAll, onChange}) => (
    <thead>
        <tr>
            {checkbox &&
                <th>
                    <input
                        type="checkbox"
                        checked={allSelected}
                        onChange={onSelectedAll}
                    />
                </th>
            }
            {header.map((item, index) => (
                <Header
                    key={index}
                    header={item}
                    sort={sort}
                    order={order}
                    onChange={onChange}
                />
            ))}
        </tr>
    </thead>
)

// -- plain item --

const PlainItem = ({field, name, link, onClick}) => {
    let text = (
        <div>{name || ""}</div>
    )

    if (link) {
        return (
            <a
                href="javascript:void(0)"
                onClick={() => {
                    onClick(field, name)
                }}
            >
                {text}
            </a>
        )
    }
    
    return text
}

// -- field --

const Field = ({onClick, item, value, header, onAction}) => {
    let {field} = header;

    if (field === false) {
        if (header.component) {
            return <header.component onAction={onAction} onClick={onClick} item={item} />
        }
    }
    
    if (header.component) {
        return <header.component onClick={onClick} value={value} />
    }

    return <PlainItem field={header.name} onClick={onClick} name={value} link={header.link} />
}

// -- item --

class Item extends Component {
    render() {
        let {onAction, checkbox, isSelected, item, onChecked, header, onClick} = this.props;

        return (
            <tr>
                {checkbox &&
                    <th>
                        <input
                            type="checkbox"
                            checked={isSelected}
                            onChange={(e) => {
                                onChecked(item.id, !isSelected)
                            }}
                        />
                    </th>
                }
                {header.map((h, index) => (
                    <th key={index}>
                        <Field onAction={onAction} onClick={onClick} item={item} value={item[h.name]} header={h} />
                    </th>
                    )
                )}
            </tr>
        )
    }
}

const Loading = ({num}) => (
    <tbody>
        <tr>
            <td colSpan={num}>
                Loading..
            </td>
        </tr>
    </tbody>
)

// -- body --

const Expanded = ({item, ExpComp, numfields}) => (
    <tr>
        <td colSpan={numfields} >
            <ExpComp item={item} />
        </td>
    </tr>
)

class Body extends Component {
    render() {
        let {items, header, selected, checkbox, isLoading, loaded, expanded, expComp} = this.props;

        if (isLoading && !loaded) {
            return <Loading 
                num={header.length}
            />
        }

        let expIndex = -1;
        let expItem = undefined;

        let children = items.map((item, index) => {
            if (expComp != undefined) {
                if (item.id == expanded) {
                    expIndex = index
                    expItem = item
                }
            }

            console.log(item.id)

            return (
                <Item
                    onChecked={this.props.onChecked}
                    checkbox={checkbox} 
                
                    onAction={this.props.onAction}

                    key={item.id} 
                    isSelected={selected.includes(item.id)}
                    onClick={this.props.onClick} 
                    item={item} 
                    header={header} 
                />
            )
        })

        console.log(expIndex)

        // add the expanded item. TODO. Using this is impossible to apply 'persiana' effect
        if (expIndex > -1) {
            children.splice(expIndex+1, 0, <Expanded item={expItem} ExpComp={expComp} numfields={header.length} key={'expanded'} />);
        }

        return (
            <tbody>
                {children}
            </tbody>
        )
    }
}

// -- table --

export class Table extends Component {
    render() {
        let {items, hasBack, hasNext} = this.props;
        let {header, order, sort} = this.props;
        let {checkbox, selected, onSelect, allSelected} = this.props;
        let {onClick, onSort, onSelectedAll} = this.props;
        let {isLoading, loaded} = this.props;
        let {expanded, expComp} = this.props;
        
        return (
            <div>
                <div className="table-responsive">
                    <table className="table">
                        <ColsGroup
                            header={header}
                            checkbox={checkbox}
                        />
                        <Headers
                            header={header}
                            order={order}
                            sort={sort}
                            checkbox={checkbox}
                            onChange={onSort}
                            allSelected={allSelected}
                            onSelectedAll={onSelectedAll}
                        />
                        <Body
                            header={header}
                            items={items}
                            isLoading={isLoading}
                            loaded={loaded}
                            selected={selected}
                            checkbox={checkbox}
                            onAction={this.props.onAction}
                            onClick={onClick}
                            onChecked={onSelect}

                            expanded={expanded}
                            expComp={expComp}
                        />
                    </table>
                </div>
                {hasBack &&
                    <a
                        href="javascript:void(0)"
                        onClick={this.props.back}
                    >
                        Back
                    </a>
                }
                {hasNext &&
                    <a
                        href="javascript:void(0)"
                        onClick={this.props.next}
                    >
                    Next
                    </a>
                }
            </div>
        )
    }
}
