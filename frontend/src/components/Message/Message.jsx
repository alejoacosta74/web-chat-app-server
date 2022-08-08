// Message component is going to take in the message it needs to display through a prop. 
// It’ll then parse this prop called 'message' and store it in the components state which we can then use within our render function.

import React, { Component } from "react";
import "./Message.scss";

class Message extends Component {
  constructor(props) {
    super(props);
    let temp = JSON.parse(this.props.message);
    this.state = {
      message: temp
    };
  }

  render() {
    return <div className="Message">{this.state.message.body}</div>;
  }
}

export default Message;