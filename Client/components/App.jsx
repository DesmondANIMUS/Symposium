import React, {Component} from 'react';
import PropTypes from 'prop-types';
import ChannelSection from './channels/ChannelSection.jsx';
import UserSection from './users/UserSection.jsx';
import MessageSection from './messages/MessageSection.jsx';
import Socket from '../Socket.js';

class App extends Component {
    constructor(props) {
        super(props);
        this.state = {
            channels: [],
            users: [],
            messages: [],
            activeChannel: {},
            connected: false
        };                
    }

    componentDidMount() {
        let ws = new WebSocket('ws://localhost:8888')
        let sock = this.sock = new Socket(ws);
        sock.on('connect', this.onConnect.bind(this));
        sock.on('disconnect', this.onDisconnect.bind(this));
        sock.on('channel add', this.onAddChannel.bind(this));

        sock.on('user add', this.onAddUser.bind(this));
        sock.on('user edit', this.onEditUser.bind(this));
        sock.on('user remove', this.onRemoveUser.bind(this));

        sock.on('message add', this.onMessageAdd.bind(this));
    }    

    onConnect() {
        this.setState({connected: true});
        this.sock.emit('channel subscribe');
        this.sock.emit('user subscribe');
    }
    onDisconnect() {                
        this.setState({connected: false});
    }


    onAddChannel(channel) {
        let {channels} = this.state;
        channels.push(channel);
        this.setState({channels});
    }
    addChannel(name) {
        this.sock.emit('channel add', {name});
    }
    setChannel(activeChannel) {
        this.setState({activeChannel});
        this.sock.emit('message unsubscribe');
        this.setState({messages: []});
        this.sock.emit('message subscribe', { channelId: activeChannel.id });
    }

    
    onRemoveUser(removeUser) {
        console.log("REMOVE USER CALLED")

        let {users} = this.state;
        users = users.filter(user => {
            return user.id !== removeUser.id;
        });

        this.setState({users});
    }
    onEditUser(editUser) {        
        console.log("EDIT USER CALLED")

        let {users} = this.state;
        users = users.map(user => {
            if (editUser.id === user.id) {
                return editUser;
            }

            return user;
        });

        this.setState({users});
    }
    onAddUser(user) {
        let {users} = this.state;
        users.push(user);
        this.setState({users});
    }
    addUser(name) {
        this.sock.emit('user edit', {name});
    }


    onMessageAdd(message) {
        let {messages} = this.state;
        messages.push(message);
        this.setState({messages});
    }
    addMessage(body) {
        let {activeChannel} = this.state;
        this.sock.emit('message add', {channelId: activeChannel.id, body});
    }


    render() {
        return(
            <div className='app'> 
                <div className='nav'>
                    <ChannelSection 
                        {...this.state}
                        addChannel = {this.addChannel.bind(this)}
                        setChannel = {this.setChannel.bind(this)}
                    />
                    <UserSection 
                        {...this.state}
                        addUser = {this.addUser.bind(this)}                        
                    />                    
                </div>
                <MessageSection 
                    {...this.state}
                    addMessage = {this.addMessage.bind(this)}
                />                               
            </div>
        )
    }
}

export default App