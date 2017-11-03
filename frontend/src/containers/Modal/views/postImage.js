
import React, { Component, PropTypes } from 'react'
import { connect } from 'react-redux'

import {actions as CreateActions} from '../../../actions/create'

import {Text, Select} from '../../../components/Inputs'

// Formatos

const formatsNames = [
    'PNG',
    'JPEG',
    'NIFTI',
    'other'
]

class PostImage extends Component {
    constructor(props) {
        super(props)

        this.file = undefined;

        this.state = {
            name: '',
            format: '',
            fileLoaded: false,
        }

        this.onChange = this.onChange.bind(this)
        this.uploadImage = this.uploadImage.bind(this)
        this.onFileChange = this.onFileChange.bind(this)
    }

    onChange(name, value) {
        this.setState(Object.assign({}, this.state, {[name]: value}))
    }

    onFileChange(e) {
        const file = e.target.files[0]
        this.file = file;

        this.setState({
            fileLoaded: true
        })
    }

    uploadImage() {
        const {name, format, fileLoaded} = this.state;
        if (name == "" || format == "" || !fileLoaded) {
            console.error("-- fallo al upload --")
            return
        }

        console.log("-- upload things --")
        console.log(name)

        this.props.commitImage(name, format, this.file)
    }
    
    render() {
        const {name, format, fileLoaded} = this.state;

        console.log("-- name --")
        console.log(name)

        return (
            <div>
                {'Post image'}
                <div>
                    <Text
                        value={name}
                        onChange={(value) => {
                            this.onChange('name', value)
                        }}
                    />
                    <Select
                        value={format}
                        values={formatsNames}
                        onChange={(value) => {
                            this.onChange('format', value)
                        }}
                    />
                    <input
                        type="file"
                        onChange={this.onFileChange}
                    />
                    <a
                        href="javascript:void(0)"
                        onClick={this.uploadImage}
                    >
                        {'Enviar'}
                    </a>
                </div>
            </div>
        )
    }
}

const mapStateToProps = (state, ownProps) => {
    return ownProps
}

const mapDispathToProps = (dispatch) => {
    return {
        commitImage: (name, format, image) => {
            console.log("- upload -")
            console.log(name)
            console.log(image)
            console.log(format)
            
            dispatch(CreateActions.postImage(name, format, image))
        }
    }
}

export default connect(
    mapStateToProps, 
    mapDispathToProps
)(PostImage)
