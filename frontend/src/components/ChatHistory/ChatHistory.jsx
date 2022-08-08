import React, { Component } from "react";
import "./ChatHistory.scss";
import Message from "../Message/Message";

class ChatHistory extends Component {
  render() {
    console.log(this.props.chatHistory);
    // .map function returns a <Message /> component with the message prop set to our msg.data. 
    // This will subsequently pass in the JSON string to every message component and it will then be able to parse and render that, as it wishes.
    const messages = this.props.chatHistory.map((msg, idx, history) => {
      if (idx > (history.length - 10)) {
        return <Message message={msg.data} />
      } else {
        // eslint-disable-next-line array-callback-return
        return
      } 
    })

    return (
      <div className="ChatHistory">
        {messages}
      </div>
    );
  }
}

export default ChatHistory;