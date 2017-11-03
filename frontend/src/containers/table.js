
import React, {Component} from 'react'
import {connect} from 'react-redux'

import {Table as TableDumb} from '../components/table'

import {actions} from '../actions/app'

const SECONDS = 1000
const RefreshInterval = 10 * SECONDS

// encadena todos los elementos para hacer el query a bd
function queryToString(fields) {
    if (fields.length === 0) {
        return '?'
    }

    let composed = Object.keys(fields).map(k => k + '=' + fields[k])
    return '?' + composed.join('&')
}

// clear query elimina el page si estaba
function clearQuery(oldquery) {
    let query = Object.assign({}, oldquery, {})
    
    let page = 0
    if (query['page'] !== undefined) {
        page = query['page']
        delete query['page']
    }

    return ({page: page, query, str: queryToString(query)})
}

// Objeto de tipo generico para que sea usado
// con los eventos de la ventana (page.. sort..)
// o de forma local en una ventana modal donde no tiene
// una referencia de objeto query tan clara

// solo le pasas dos cosas, el endpoint y el query.
// el endpoint lo sabes porque si quieres devices sabes que vas 
// al /devices y el query lo puedes sacar del url o construirlo si hace falta

class Table extends Component {
    constructor(props) {
        super(props)

        this.loadData = this.loadData.bind(this)
        this.fetchData = this.fetchData.bind(this)

        this.onSelectedAll = this.onSelectedAll.bind(this)
        this.onSelect = this.onSelect.bind(this)
        this.changeSort = this.changeSort.bind(this)
        this.changePage = this.changePage.bind(this)
        this.onExpand = this.onExpand.bind(this)
        this.onAction = this.onAction.bind(this)

        this.count = 0
        this.interval = undefined

        this.state = {
            expanded: undefined
        }
    }

    onAction(action, info) {
        console.log("-- action --")
        console.log(action)
        console.log(info)

        if (action == 'expand') {
            this.onExpand(info)
        } else {
            this.props.onAction(action, info)
        }
    }

    onExpand(id) {
        console.log("-- expand --")
        console.log(id)

        this.setState({
            expanded: (this.state.expanded == id ? undefined : id)
        })
    }

    onSelectedAll() {
        let {allSelected, items} = this.props;
        if (allSelected) {
            this.props.onSelect([])
        } else {
            this.props.onSelect(items.map(i => i.id))
        }
    }

    onSelect(item, checked) {
        console.log("-- item --")
        console.log(checked)

        let {selected} = this.props;
        if (!checked) {
            // eliminar
            this.props.onSelect(selected.filter(function(i) {return i != item}))
        } else {
            // annadir
            selected.push(item)
            this.props.onSelect(selected)
        }
    }

    changeSort(sortIndex, order) {
        if (sortIndex == undefined) {
            let query = this.props.query;
            delete query['sort']
            delete query['order']

            this.props.onQueryChange(query, queryToString(query))
        } else {
            let query = Object.assign({}, this.props.query, {sort: sortIndex, order:order})
            this.props.onQueryChange(query, queryToString(query))
        }
    }

    // -- cambiar la pagina --
    changePage(move) {
        console.log("-- change page --")
        
        let {page, query} = this.props;
        let newpage = parseInt(page) + move;

        // si es cero no se envia nada
        if (newpage != 0) {
            query = Object.assign({}, this.props.query, {page: newpage})
        }

        this.props.onQueryChange(query, queryToString(query))
    }

    // carga los datos seguro
    loadData(props) {
        let {queryStr, page, isLoading, fetchfn} = props
        this.props.loadData(fetchfn, {fields: queryStr, page})
    }

    startInterval() {
        if (this.interval != undefined) {
            clearInterval(this.interval)
        }

        this.interval = setInterval(() => {
            this.loadData(this.props)
        }, RefreshInterval)
    }
    
    // Comprobar si esta en fetching o no
    fetchData(props, justmounted=false) {
        if (props.notfound || justmounted) {
            if (this.count < 5) {
                this.loadData(props)        
                
                //  ... TODO ...
                //  se ha llamado porque los props se han cambiado, start interval again
                this.startInterval()
            } else {
                console.error("-- llamada mayor de cinco -- posible bucle infinito")
            }
            
            this.count++
        }
    }
    
    componentDidMount() {
        console.log("-- montado ahora --")
        this.fetchData(this.props, true)
    }

    componentWillUnmount() {
        if (this.interval != undefined) {
            clearInterval(this.interval)
        }
    }

    componentWillReceiveProps(nextProps) {
        this.fetchData(nextProps)
    }

    render() {
        let {header, items, hasBack, query} = this.props;
        let {checkbox, selected=[], onSelect, onAction, allSelected} = this.props;
        let {onClick} = this.props;
        let {isLoading, loaded} = this.props;
        
        let {order='asc', sort=''} = query;

        let {expComp} = this.props;
        let {expanded} = this.state;

        return (
            <div>
                <TableDumb
                    header={header}
                    order={order}
                    sort={sort}
                    onSort={this.changeSort}

                    isLoading={isLoading}
                    loaded={loaded}

                    checkbox={checkbox}
                    selected={selected}
                    allSelected={allSelected}
                    onSelect={this.onSelect}
                    onSelectedAll={this.onSelectedAll}

                    items={items}

                    hasNext={true}
                    next={() => {this.changePage(1)}}

                    hasBack={hasBack}
                    back={() => {this.changePage(-1)}}

                    onClick={onClick}
                    onAction={this.onAction}

                    expanded={expanded}
                    expComp={expComp}
                />
            </div>
        )
    }
}

// parece que query no llega aqui. Se debe pasar por props.
// el query puede tener el page.
const mapStateToProps = (state, ownProps) => {
    let {reducer='app', endpoint, checkbox, fetchfn} = ownProps;
    let {query, page, str} = clearQuery(ownProps.query)
    
    if (!fetchfn) {
        console.error("No hay funcion fetchfn")
    }

    // Comprobar el checkbox, si esta en true, debe haber una lista en selected
    // ademas, deben estar activos los eventos onSelected y onSelectedAll

    if (checkbox) {
        let {onSelect, selected} = ownProps;
        if (selected == undefined) {
            console.error("Debe existir selected si se quiere usar todo")
        }
        if (onSelect == undefined) {
            console.error("Debe existir el evento onSelect si se quieren usar checkboxs")
        }
    }

    // Get pagination and entities
    let pagination = state[reducer].pagination[endpoint]
    let entities = state[reducer].entities[endpoint]

    // Sacar el endpoint especifico (parte de esto puede no existir)
    let paginationItems = pagination[str] || {pages: {}};
    let paginationItemsPage = paginationItems.pages[page] || {loaded: false, isLoading: false, notfound: true, items: []};
    
    let items = paginationItemsPage.items.map(i => Object.assign({}, entities[i], {id: i}));

    // Comprobar si estan todos seleccionados o no (selected con paginationItemsPage.items)
    let selected = ownProps.selected || [];
    let allSelected = !(paginationItemsPage.items.map(i => selected.includes(i))).includes(false) && items.length != 0

    return {
        query,
        queryStr: str,
        page,
        items,
        hasBack: page > 0,
        isLoading: paginationItemsPage.isLoading,
        loaded: paginationItemsPage.loaded,
        notfound: paginationItemsPage.notfound,
        allSelected
    }
}

const mapDispatchToProps = (dispatch, getState) => {
    return {
        loadData: (fetchfn, query) => {
            dispatch(fetchfn(query))
        }
    }
}

export default connect(mapStateToProps, mapDispatchToProps)(Table)
