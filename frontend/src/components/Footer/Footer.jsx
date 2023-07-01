import React from 'react';
import './Footer.css';
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";
import {faPhone, faEnvelope} from '@fortawesome/free-solid-svg-icons';
import {faFacebook, faTwitter, faInstagram} from '@fortawesome/free-brands-svg-icons';
import {Box, Heading, Center, Container} from '@chakra-ui/react';
function Footer() {
    return (
      <Box padding={8}>
        <Center>
          <Heading size="sm"></Heading>
        </Center>
      </Box>
    );
  }
  
export default Footer;
