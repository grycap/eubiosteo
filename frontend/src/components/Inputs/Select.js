
import React from 'react'

const Select = ({value, values, onChange, blank=true}) => (
    <select value={value} onChange={(e) => {onChange(e.target.value)}}>
        {blank &&
            <option value={''}></option>
        }
        {values.map((item, index) => (
            <option key={index} value={item}>{item}</option>
        ))}
    </select>
)

export default Select
