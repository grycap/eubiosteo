
import React from 'react'

const Text = ({value, onChange}) => (
    <input
        type="text"
        value={value}
        onChange={(e) => {onChange(e.target.value)}}
    />
)

export default Text
