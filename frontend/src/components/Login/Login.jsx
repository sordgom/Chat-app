import axios from 'axios';
import React, { Component } from 'react';

import {
  Box,
  Button,
  Container,
  FormControl,
  FormErrorMessage,
  FormLabel,
  Input,
  Stack,
  Text,
} from '@chakra-ui/react';

import { EditIcon } from '@chakra-ui/icons';
import { withCookies } from 'react-cookie';
import { Navigate } from 'react-router-dom';

class Login extends Component {
  constructor(props) {
    super(props);
    this.state = {
      username: '',
      password: '',
      message: '',
      isInvalid: false,
      endpoint: 'http://localhost:8080/signin',
      redirect: false,
      redirectTo: '/home?u=',
      tokenValidated: false, // Flag to indicate if token validation has been performed
    };
  }

  // on change of input, set the value to the message state
  onChange = event => {
    this.setState({ [event.target.name]: event.target.value });
  };  
  
  componentDidMount() {
    this.performTokenValidation();
  }

  performTokenValidation = async () => {
    const { cookies } = this.props;
    const jwtToken = cookies.get('jwtToken') || '';
    try {
      const res = await axios.get("http://localhost:8080/check-status", {
        headers: {
          'Authorization': `Bearer ${jwtToken}`
        }
      });
      const user = this.decodeJWT(jwtToken).Usr;
      if (res.data) {
        const redirectTo = this.state.redirectTo + user;
        this.setState({ redirect: true, redirectTo });
      }
      this.setState({ tokenValidated: true }); // Set the flag to indicate token validation is done
    } catch (error) {
      console.log(error);
      this.setState({ tokenValidated: true }); // Set the flag to indicate token validation is done
    }
  };

  decodeJWT = (token) => {
    const base64Url = token.split('.')[1];
    const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/');
    const jsonPayload = window.atob(base64);
    return JSON.parse(jsonPayload);
  }  

  handleAuthentication = async () => {
    const { cookies } = this.props;
    const jwtToken = cookies.get('jwtToken') || '';

    try {
      const res = await axios.post(this.state.endpoint, {
        username: this.state.username,
        password: this.state.password,
      }, {
        headers: {
          'Authorization': `Bearer ${jwtToken}`
        }
      });

      return res.data;
    } catch (error) {
      console.log(error);
      throw new Error('Authentication failed');
    }
  };

  onSubmit = async e => {
    e.preventDefault();

    try {
      
      const authResponse = await this.handleAuthentication();
      
      if (authResponse.status) {       
        console.log(authResponse);
        const { cookies } = this.props;
        cookies.set('jwtToken', authResponse.token);
        const redirectTo = this.state.redirectTo + this.state.username;
        this.setState({ redirect: true, redirectTo });
      } else {
        this.setState({ message: authResponse.message, isInvalid: true });
      }
    } catch (error) {
      console.log(error);
      this.setState({ message: 'something went wrong', isInvalid: true });
    }
};

  render() {
    const { redirect, redirectTo, tokenValidated } = this.state;

    if (!tokenValidated) {
      return null; // Render nothing while token validation is in progress
    }
    if (redirect) {
      return <Navigate to={redirectTo} replace={true} />;
    }
    return (
      <div>
        {this.state.redirect && (
          <Navigate to={this.state.redirectTo} replace={true}></Navigate>
        )}

        <Container marginBlockStart={10} textAlign={'left'} maxW="2xl">
          <Box borderRadius="lg" padding={10} borderWidth="2px">
            <Stack spacing={5}>
              <FormControl isInvalid={this.state.isInvalid}>
                <FormLabel>Username</FormLabel>
                <Input
                  type="text"
                  placeholder="Username"
                  name="username"
                  value={this.state.username}
                  onChange={this.onChange}
                />
              </FormControl>
              <FormControl isInvalid={this.state.isInvalid}>
                <FormLabel>Password</FormLabel>
                <Input
                  type="password"
                  placeholder="Password"
                  name="password"
                  value={this.state.password}
                  onChange={this.onChange}
                />
                {!this.state.isInvalid ? (
                  ''
                ) : (
                  <FormErrorMessage>
                    invalid username or password
                  </FormErrorMessage>
                )}
              </FormControl>
              <Button
                size="lg"
                leftIcon={<EditIcon />}
                colorScheme="cyan"
                variant="solid"
                type="submit"
                onClick={this.onSubmit}
              >
                Login
              </Button>
            </Stack>
            <Box paddingTop={3}>
              <Text as="i" fontSize={'lg'} color={'red'}>
                {this.state.message}
              </Text>
            </Box>
          </Box>
        </Container>
      </div>
    );
  }
}

export default withCookies(Login);
