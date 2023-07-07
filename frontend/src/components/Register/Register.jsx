import axios from 'axios';
import React, { Component } from 'react';

import {
  Box,
  Button,
  Container,
  FormControl,
  FormErrorMessage,
  FormHelperText,
  FormLabel,
  Input,
  Stack,
} from '@chakra-ui/react';

import { EditIcon } from '@chakra-ui/icons';
import { withCookies } from 'react-cookie';
import { Navigate } from 'react-router-dom';

class Register extends Component {
  constructor(props) {
    super(props);
    this.state = {
      username: '',
      email: '',
      password: '',
      confirmPassword: '',
      message: '',
      isInvalid: '',
      endpoint: 'http://localhost:8080/api/auth/register',
      redirect: false,
      redirectTo: '/home',
      tokenValidated: false, // Flag to indicate if token validation has been performed
    };
  }


  decodeJWT = (token) => {
    const base64Url = token.split('.')[1];
    const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/');
    const jsonPayload = window.atob(base64);
    return JSON.parse(jsonPayload);
  }  
  
  // on change of input, set the value to the message state
  onChange = event => {
    this.setState({ [event.target.name]: event.target.value });
  };
  
  onSubmit = async e => {
    e.preventDefault();

    const { username, email, password, confirmPassword } = this.state;

    if (password !== confirmPassword) {
      this.setState({
        message: 'Passwords do not match',
        isInvalid: true,
      });
      return;
    }
    try {
      const res = await axios.post(this.state.endpoint, {
        name: this.state.username,
        email: this.state.email,
        password: this.state.password,
        passwordConfirm: this.state.confirmPassword,
      });

      console.log(res.data.status)
      if (res.data.status==='success') {
        const redirectTo = this.state.redirectTo ;
        this.setState({ redirect: true, redirectTo });
      } else {
        // on failed
        this.setState({ message: res.data.message, isInvalid: true });
      }
    } catch (error) {
      console.log(error.response.data);
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
                <FormLabel>Username</FormLabel>
                <Input
                  type="text"
                  placeholder="Username"
                  name="username"
                  value={this.state.username}
                  onChange={this.onChange}
                />

                {!this.state.isInvalid ? (
                  <></>
                ) : (
                  <FormErrorMessage>{this.state.message}</FormErrorMessage>
                )}
                {/* <FormHelperText>use a unique username</FormHelperText> */}
              </FormControl>
              <FormControl isInvalid={this.state.isInvalid}>
                <FormLabel>Email</FormLabel>
                <Input
                  type="text"
                  placeholder="Email"
                  name="email"
                  value={this.state.email}
                  onChange={this.onChange}
                />

                {!this.state.isInvalid ? (
                  <></>
                ) : (
                  <FormErrorMessage>{this.state.message}</FormErrorMessage>
                )}
                {/* <FormHelperText>use a unique email</FormHelperText> */}
              </FormControl>
              <FormControl>
                <FormLabel>Password</FormLabel>
                <Input
                  type="password"
                  placeholder="Password"
                  name="password"
                  value={this.state.password}
                  onChange={this.onChange}
                />
                <FormHelperText>Password must be more than 8 characters</FormHelperText>
              </FormControl>
              <FormControl isInvalid={this.state.isInvalid}>
                <FormLabel>Confirm Password</FormLabel>
                <Input
                  type="password"
                  placeholder="Confirm Password"
                  name="confirmPassword"
                  value={this.state.confirmPassword}
                  onChange={this.onChange}
                />
                {!this.state.isInvalid ? (
                  <></>
                ) : (
                  <FormErrorMessage>{this.state.message}</FormErrorMessage>
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
                Register
              </Button>
            </Stack>
          </Box>
        </Container>
      </div>
    );
  }
}

export default withCookies(Register);
