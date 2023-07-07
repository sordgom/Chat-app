import { Box, Button, Container, Link, VStack } from '@chakra-ui/react';
import React from 'react';
import { Link as RouterLink } from 'react-router-dom';

const Home = () => {
  const queryParams = new URLSearchParams(window.location.search);
  const user = queryParams.get('u');

  return (
    <Container maxW="xl" centerContent>
      <VStack spacing={4} align="stretch">
        <Box borderWidth="1px" borderRadius="lg" overflow="hidden">
          <Box p="6">
            <Link as={RouterLink} to={"/videochat?meetingId=07927fc8-af0a-11ea-b338-064f26a5f90a&userId="+user} style={{ display: 'block', width: '100%' }}>
              <Button colorScheme="teal" variant="solid" w="full">Go to Video Chat</Button>
            </Link>
          </Box>
        </Box>

        <Box borderWidth="1px" borderRadius="lg" overflow="hidden">
          <Box p="6">
            <Link as={RouterLink} to={"/chat"} style={{ display: 'block', width: '100%' }}>
              <Button colorScheme="teal" variant="solid" w="full">Go to Chat</Button>
            </Link>
          </Box>
        </Box>
      </VStack>
    </Container>
  );
};

export default Home;
