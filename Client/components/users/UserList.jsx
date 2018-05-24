import React, {Component} from 'react';
import PropTypes from 'prop-types';
import User from './User.jsx';

class UserList extends Component {
    render() {
        return (
            <ul>{
                this.props.users.map(person => {
                    return (<User user = {person} key = {person.id} />)
                })
            }</ul>
        )
    }
}

UserList.propTypes = {    
    users: PropTypes.array.isRequired
}

export default UserList