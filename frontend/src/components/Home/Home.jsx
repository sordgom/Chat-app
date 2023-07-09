import { Box, Button, Container, Link, VStack } from '@chakra-ui/react';
import React from 'react';
import { Link as RouterLink } from 'react-router-dom';

const Home = () => {
  return (
    <Container maxW="xl" centerContent>
      <VStack spacing={4} align="stretch">
        <Box borderWidth="1px" borderRadius="lg" overflow="hidden">
          <Box p="6">
            <Link as={RouterLink} to={"/videochat"} style={{ display: 'block', width: '100%' }}>
              <Button colorScheme="cyan" variant="solid" w="full">Go to Video Chat</Button>
            </Link>
          </Box>
        </Box>

        <Box borderWidth="1px" borderRadius="lg" overflow="hidden">
          <Box p="6">
            <Link as={RouterLink} to={"/chat"} style={{ display: 'block', width: '100%' }}>
              <Button colorScheme="cyan" variant="solid" w="full">Go to Chat</Button>
            </Link>
          </Box>
        </Box>
      </VStack>
    </Container>
  );
};

export default Home;
