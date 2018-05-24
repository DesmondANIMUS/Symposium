import React, {Component} from 'react';
import PropTypes from 'prop-types';

class ChannelForm extends Component {
    onSubmit(e) {
        e.preventDefault();
        const node = this.refs.channel;
        const chanName = node.value;
        this.props.addChannel(chanName);
        node.value = '';
    }

    render() {
        return (            
            <form onSubmit={this.onSubmit.bind(this)}>
                <div className='form-group'> 
                    <input type='text' ref='channel' className='form-control' placeholder='Add Channel'/>
                </div>        
            </form>
        )
    }
}

ChannelForm.propTypes = {
    addChannel: PropTypes.func.isRequired
}

export default ChannelForm