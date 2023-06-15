import React from "react";
import { Box, Flex, Heading } from "@chakra-ui/react";
import { ColorModeSwitcher } from "./ColorModeSwitcher";
import { Link } from 'react-router-dom';

const Header = () => (
  <Flex className="header" alignItems="center" justifyContent="space-between">
    <Link to="/" style={{ textDecoration: 'none', color: 'inherit' }}>
      <Heading size="2xl" font-size={"xl"}>
        Real time chat application
      </Heading>
    </Link>
    <Box>
      <ColorModeSwitcher />
    </Box>
  </Flex>
);

export default Header;
