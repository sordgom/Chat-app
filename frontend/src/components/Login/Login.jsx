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
      email: 'admin@admin.com',
      password: 'password123',
      message: '',
      isInvalid: false,
      endpoint: 'http://localhost:8080/api/auth/login',
      redirect: false,
      redirectTo: '/home',
      tokenValidated: false, // Flag to indicate if token validation has been performed
    };
  }

  // on change of input, set the value to the message state
  onChange = event => {
    this.setState({ [event.target.name]: event.target.value });
  };  

  handleAuthentication = async () => {
    try {
      const response = await axios.post(this.state.endpoint, {
        email: this.state.email,
        password: this.state.password,
      },{
        withCredentials: true, // include cookies in the request
      });
      return response.data;
    } catch (error) {
      console.log(error);
      throw new Error('Authentication failed');
    }
  };

  onSubmit = async e => {
    e.preventDefault();

    try {
      const authResponse = await this.handleAuthentication();
      const { cookies } = this.props;
      if (authResponse.status === 'success' || cookies.loggedIn === 'true') {       
        const redirectTo = this.state.redirectTo ;
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
    const { redirect, redirectTo } = this.state;

    if (redirect) {
      return <Navigate to={redirectTo} replace={true} />;
    }
    return (
      <div>
        <Container marginBlockStart={10} textAlign={'left'} maxW="2xl">
          <Box borderRadius="lg" padding={10} borderWidth="2px">
            <Stack spacing={5}>
              <FormControl isInvalid={this.state.isInvalid}>
                <FormLabel>Email</FormLabel>
                <Input
                  type="text"
                  placeholder="Email"
                  name="email"
                  value={this.state.email}
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
                    invalid email or password
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
