import React, {Component} from 'react';
import PropTypes from 'prop-types';

class UserForm extends Component {
    onSubmit(e) {
        e.preventDefault();
        const node = this.refs.user;
        const useName = node.value;
        this.props.addUser(useName);
        node.value = '';
    }

    render() {
        return (
            <form onSubmit={this.onSubmit.bind(this)}>
                <div className='form-group'>
                    <input type='text' ref='user' className='form-control' placeholder='Add Yourself'/>
                </div>
            </form>    
        )        
    }
}

UserForm.propTypes = {
    addUser: PropTypes.func.isRequired
}

export default UserForm