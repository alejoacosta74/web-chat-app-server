import logo from './logo.svg';
import './App.css';
import React, { Component } from "react";
import { connect, sendMsg } from "./api";
import Header from './components/Header/Header';
import ChatHistory from './components/ChatHistory/ChatHistory';
import ChatInput from './components/ChatInput/ChatInput';
import {
  Button,
  Card,
  CardBody,
  CardGroup,
  CardTitle,
  CardImg,
  Row,
  CardImgOverlay,
  CardText,
  Col,
  Container,
  CardHeader
} from "reactstrap";
import 'bootstrap/dist/css/bootstrap.min.css';

class App extends Component{
  constructor(props){
    super(props);
    this.state = {
      chatHistory: []
    }
  }

  componentDidMount() {
    // within the componentDidMount life cycle, call `connect` and pass a call back function
    // that will update the state when a new msg is received via websocket
    connect((msg) => {
      console.log("New Message")
      this.setState(prevState => ({
        chatHistory: [...this.state.chatHistory, msg]
      }))
      console.log("(connect call back funtcion) => state: ", this.state);
    });
  }

  // By passing in this event, we’ll be able to query if the key pressed was the Enter key, 
  // if it is, we’ll be able to send the value of our <input/> field to our WebSocket endpoint and then subsequently clear that <input/>

  send(event) {
    if(event.keyCode === 13) {
      sendMsg(event.target.value);
      event.target.value = "";
    }
  }

  render() {
    return (
      <div className='App'>
        <header className='App-header'>
          <Header/>
        </header>
		  <Container>
			  <Row className='App-row'>
				<Col
						className="bg-light border"
						xs="5"
					>
						<Card
						inverse
						body
						color="danger"
						>
							<CardHeader tag="h2">
								Input text
							</CardHeader>
							<CardBody>
								<ChatInput send={this.send} />
								<br/>
							</CardBody>
						</Card>
				</Col>

				<Col
						className="bg-light border"
						xs="5"
					>

				<Card 
							body
							inverse
							style={{
								backgroundColor: '#333',
								borderColor: '#333'
							}}
							>
							<CardHeader tag="h2">
							Chat History
							</CardHeader>
							<CardBody>
								<ChatHistory chatHistory={this.state.chatHistory} />
							</CardBody>
						</Card>


				</Col>

			  </Row>
		  </Container>
      </div>
    );
  }

}

// function App() {
//   return (
//     <div className="App">
//       <header className="App-header">
//         <img src={logo} className="App-logo" alt="logo" />
//         <p>
//           Edit <code>src/App.js</code> and save to reload.
//         </p>
//         <a
//           className="App-link"
//           href="https://reactjs.org"
//           target="_blank"
//           rel="noopener noreferrer"
//         >
//           Learn React
//         </a>
//       </header>
//     </div>
//   );
// }

export default App;
