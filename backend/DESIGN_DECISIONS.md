# Technical documentation

## User Authentication 

### Technologies used

* [Redis](https://redis.io/) - In-memory data structure store
* [JWT](https://jwt.io/) - JSON Web Tokens
* [bcrypt](https://www.npmjs.com/package/bcrypt) - A library to help you hash passwords
* [cookie-parser](https://www.npmjs.com/package/cookie-parser) - Parse HTTP request cookies
* [postgreSQL](https://www.postgresql.org/) - Open source relational database
* [Golang](https://golang.org/) - Open source programming language that makes it easy to build simple, reliable, and efficient software
* [Fiber](https://gofiber.io/) - Express inspired web framework written in Go
* [GORM](https://gorm.io/) - The fantastic ORM library for Golang
* [Docker](https://www.docker.com/) - Docker is a set of platform as a service products that use OS-level virtualization to deliver software in packages called containers
  
### JWT Access & Refresh Token Design Decisions

* Secure Storage: Storing token metadata in a Redis database allows for centralized and secure storage. Redis is known for its performance and reliability, and it can be configured to store data securely.

* Separation of Concerns: Storing token metadata separately from the actual token enhances security by separating sensitive information from the token itself. By storing metadata in a separate database, it minimizes the risk of exposing sensitive data if the token is compromised.

* HTTP-only Cookies: Using HTTP-only cookies to transmit tokens provides additional security measures. HTTP-only cookies cannot be accessed or manipulated by client-side Jav1aScript, reducing the risk of cross-site scripting (XSS) attacks. This helps protect the token from being stolen through client-side vulnerabilities.

* Authorization Header: Including the access token in the Authorization header as a Bearer token follows the industry-standard practice for token-based authentication. This allows for secure transmission of the token in API requests, ensuring that only authenticated users with valid tokens can access protected resources.

## Chat Feature

### Technologies used

* [Websockets](https://websockets.readthedocs.io/) is an advanced technology that makes it possible to open an interactive communication session between the user's browser and a server. With this API, you can send messages to a server and receive event-driven responses without having to poll the server for a reply.
* [Redis](https://redis.io/) - Chats and user contact list are stored on Redis memory
  
### Chat Feature Design Decisions

* [Websockets](https://websockets.readthedocs.io/): This allows for a more seamless chat experience, as users can see messages as they are sent and received.
