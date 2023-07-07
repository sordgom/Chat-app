# Building a Chat Application in Golang, Fiber, PostgreSQL, Redis and  ReactJS

## Achievements

- [x] Create a basic frontend to host my services
- [x] Create a chat app service using Websocket
- [x] Create JWT Auhentication service
- [x] Everything is dockerized

## Issues

- If JWT Token is valid, user gets redirected to /home, the problem is the user itself doesnt get picked up by the frontend, so the user is not logged in. I need to find a way to pass the user to the frontend.

## Todo

- Create a video chat app service using WebRTC
- Make the services scalable

[![Run in Postman](https://run.pstmn.io/button.svg)](https://app.getpostman.com/run-collection/14390200-fb88bc8e-5710-4dce-9619-0f379111aa39?action=collection%2Ffork&source=rip_markdown&collection-url=entityId%3D14390200-fb88bc8e-5710-4dce-9619-0f379111aa39%26entityType%3Dcollection%26workspaceId%3D8e7cf06c-8991-43d1-8068-311d94c52000)

## UI

#### Landing page

![image](https://gcdnb.pbrd.co/images/fzFoZEpqQXcS.png?o=1) 

#### Login page

![image](https://gcdnb.pbrd.co/images/q2jdjOQXvQ1T.png?o=1)

#### Chat page

![image](https://gcdnb.pbrd.co/images/CUsExP9GrPRG.png?o=1)

#### Video Chat page

![image](https://gcdnb.pbrd.co/images/pWlutEV6rbam.png?o=1)

Credits to:

- https://medium.com/@ramezemadaiesec/from-zero-to-fully-functional-video-conference-app-using-go-and-webrtc-7d073c9287da

- https://codevoweb.com/how-to-properly-use-jwt-for-authentication-in-golang/